import axios from "axios";
import type {AxiosRequestConfig, AxiosResponse, AxiosError, InternalAxiosRequestConfig} from "axios";
import {useUserStore} from '@/stores/user';
import router from '@/router/index';
import {useLayoutStore} from "@/stores/layout";

const service = axios.create({
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 10000,
})

export interface ApiResponse<T> {
    code: number;
    msg: string;
    data: T;
}

service.interceptors.request.use(
    (config: AxiosRequestConfig) => {
        const userStore = useUserStore();
        config.headers = {
            'Content-Type': 'application/json',
            'x-access-token': userStore.state.accessToken,
            ...config.headers,
        }
        return config as InternalAxiosRequestConfig
    },
    (error: AxiosError) => {
        ElMessage.error({
            showClose: true,
            message: error.message,
            type: 'error',
        })
        return Promise.reject(error)
    }
)

service.interceptors.response.use(
    (response: AxiosResponse) => {
        const userStore = useUserStore()
        if (response.headers['new-access-token']) {
            userStore.state.accessToken = (response.headers['new-access-token'])
        }
        if (response.data.code !== 0) {
            ElMessage.error(response.data.msg)

            if (response.data.data && response.data.data.reload) {
                userStore.reset()
                const layoutStore = useLayoutStore()
                localStorage.clear()
                router.push({name: 'index', replace: true}).then(() => {
                    layoutStore.state.popoverVisible = true
                    layoutStore.state.loginVisible = true
                })
            }
        }
        return response.data
    },
    (error: AxiosError) => {
        if (!error.response) {
            ElMessageBox.confirm(`
        <p>检测到请求错误</p>
        <p>${error.message}</p>
      `, '请求报错', {
                dangerouslyUseHTMLString: true,
                distinguishCancelAndClose: true,
                confirmButtonText: '稍后重试',
                cancelButtonText: '取消',
            }).then()
            return Promise.reject(error)
        }

        switch (error.response.status) {
            case 500:
                return handleSpecificError(500, error)
            case 404:
                return handleSpecificError(404, error)
            case 403:
                return handleSpecificError(403, error)
        }
        return Promise.reject(error)
    }
);

// 处理具体错误状态
const handleSpecificError = (status: number, error: AxiosError) => {
    const errorMessages: { [key: number]: string } = {
        500: `
            <p>检测到接口错误: ${error.message}</p>
            <p>错误码：<span style="color:red">500</span></p>
            <p>此类错误通常由后台服务器发生不可预料的错误（如panic）引起。请先查看后台日志以获取更多信息。</p>
            <p>如果此错误影响您的正常使用，建议您清理缓存并重新登录。</p>
        `,
        404: `
            <p>检测到接口错误: ${error.message}</p>
            <p>错误码：<span style="color:red">404</span></p>
            <p>此错误通常表示请求的接口未注册（或服务未重启）或请求路径（方法）与API路径（方法）不符。</p>
            <p>请检查您请求的URL和方法，确保它们正确无误。</p>
        `,
        403: `
            <p>检测到权限错误: ${error.message}</p>
            <p>错误码：<span style="color:red">403</span></p>
            <p>您没有权限访问此路由（admin）。请确认您的用户角色是否具备访问该页面的权限。</p>
            <p>如果您认为这是一个错误，请联系系统管理员获取帮助。</p>
        `,
    }

    ElMessageBox.confirm(errorMessages[status], '接口报错', {
        dangerouslyUseHTMLString: true,
        distinguishCancelAndClose: true,
        confirmButtonText: '清理缓存',
        cancelButtonText: '取消',
    }).then(() => {
        const userStore = useUserStore()
        userStore.$reset()
        const layoutStore = useLayoutStore()
        localStorage.clear()
        router.push({name: 'index', replace: true}).then(() => {
            layoutStore.state.popoverVisible = true
            layoutStore.state.loginVisible = true
        });
    });

    return Promise.reject(error)
};

export default service

// 流式请求配置接口
export interface StreamRequestConfig {
    url: string;
    method?: string;
    data?: any;
    headers?: Record<string, string>;
    timeout?: number; // 超时时间（毫秒）
}

// 流式响应数据接口
export interface StreamResponseData {
    content?: string;
    is_complete?: boolean;
    [key: string]: any;
}

// 流式请求回调接口
export interface StreamCallbacks {
    onData?: (data: StreamResponseData) => void;
    onComplete?: (data: StreamResponseData) => void;
    onError?: (error: Error) => void;
}

// 通用流式请求方法 - 支持拦截器
export const streamRequest = (
    config: StreamRequestConfig,
    callbacks: StreamCallbacks
): Promise<void> => {
    return new Promise((resolve, reject) => {
        // 获取token（模拟拦截器行为）
        const userStore = useUserStore();
        const token = userStore.state.accessToken;
        
        if (!token) {
            const error = new Error('未登录');
            callbacks.onError?.(error);
            reject(error);
            return;
        }

        // 创建XMLHttpRequest来处理流式数据
        const xhr = new XMLHttpRequest();
        // 使用与axios相同的baseURL
        const baseURL = import.meta.env.VITE_BASE_API;
        const fullURL = config.url.startsWith('http') ? config.url : `${baseURL}${config.url}`;
        xhr.open(config.method?.toUpperCase() || 'POST', fullURL, true);
        
        // 设置默认头部（模拟拦截器行为）
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.setRequestHeader('x-access-token', token);
        
        // 设置自定义头部
        if (config.headers) {
            Object.entries(config.headers).forEach(([key, value]) => {
                xhr.setRequestHeader(key, value);
            });
        }

        // 设置超时
        if (config.timeout) {
            xhr.timeout = config.timeout;
        }

        let buffer = '';
        let timeoutId: number | null = null;
        let isAborted = false; // 添加标记，避免主动abort后触发错误

        // 设置超时处理
        if (config.timeout) {
            timeoutId = setTimeout(() => {
                isAborted = true;
                xhr.abort();
                const error = new Error(`请求超时 (${config.timeout}ms)`);
                callbacks.onError?.(error);
                reject(error);
            }, config.timeout);
        }

        xhr.onreadystatechange = function() {
            if (xhr.readyState === 3) { // 接收到部分数据
                // 重置超时计时器
                if (timeoutId) {
                    clearTimeout(timeoutId);
                    timeoutId = setTimeout(() => {
                        isAborted = true;
                        xhr.abort();
                        const error = new Error(`请求超时 (${config.timeout}ms)`);
                        callbacks.onError?.(error);
                        reject(error);
                    }, config.timeout!);
                }

                const newData = xhr.responseText.substring(buffer.length);
                buffer = xhr.responseText;
                
                // 处理SSE数据
                const lines = newData.split('\n');
                
                for (const line of lines) {
                    if (line.startsWith('data: ')) {
                        const data = line.slice(6); // 移除 'data: ' 前缀
                        
                        try {
                            const parsed = JSON.parse(data);
                            
                            // 检查是否是完成事件
                            if (parsed.event_id === 2 || parsed.is_complete) {
                                if (timeoutId) clearTimeout(timeoutId);
                                isAborted = true; // 标记为主动abort
                                callbacks.onComplete?.(parsed);
                                xhr.abort(); // 主动结束请求
                                resolve();
                                return;
                            } else {
                                callbacks.onData?.(parsed);
                            }
                        } catch (e) {
                            // 忽略解析错误，继续处理下一行
                            console.warn('解析流式数据失败:', e, 'data:', data);
                        }
                    }
                }
            } else if (xhr.readyState === 4) { // 请求完成
                if (timeoutId) clearTimeout(timeoutId);
                
                // 如果是主动abort的，不触发错误
                if (isAborted) {
                    return;
                }
                
                if (xhr.status === 200) {
                    resolve();
                } else {
                    const error = new Error(`HTTP error! status: ${xhr.status}`);
                    callbacks.onError?.(error);
                    reject(error);
                }
            }
        };

        xhr.ontimeout = function() {
            isAborted = true;
            const error = new Error(`请求超时 (${config.timeout}ms)`);
            callbacks.onError?.(error);
            reject(error);
        };

        xhr.onerror = function() {
            if (timeoutId) clearTimeout(timeoutId);
            
            // 如果是主动abort的，不触发错误
            if (isAborted) {
                return;
            }
            
            const error = new Error('网络请求失败');
            callbacks.onError?.(error);
            reject(error);
        };

        // 发送请求
        if (config.data) {
            xhr.send(JSON.stringify(config.data));
        } else {
            xhr.send();
        }
    });
};