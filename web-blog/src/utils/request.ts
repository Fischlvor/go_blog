import axios from "axios";
import type {AxiosRequestConfig, AxiosResponse, AxiosError, InternalAxiosRequestConfig} from "axios";
import {useUserStore} from '@/stores/user';
import router from '@/router/index';
import {useLayoutStore} from "@/stores/layout";
import {getDeviceId} from './deviceId';

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
            // ✅ SSO模式：使用 Authorization: Bearer <token>
            'Authorization': userStore.state.accessToken ? `Bearer ${userStore.state.accessToken}` : '',
            // ✅ 添加设备ID用于限流
            'X-Device-Id': getDeviceId(),
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
        // ✅ SSO自动刷新token：后端在响应头返回新token
        const newAccessToken = response.headers['x-new-access-token'] || response.headers['X-New-Access-Token']
        if (newAccessToken) {
            userStore.state.accessToken = newAccessToken
            console.log('✓ Token已自动刷新')
        }
        if (response.data.code !== 0) {
            // ✅ 只有非静默接口才显示错误提示
            if (!response.config.url?.includes('/user/info')) {
                ElMessage.error(response.data.msg)
            }

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
                // ✅ SSO模式：403可能是token无效，静默处理/user/info接口的403
                if (error.config?.url?.includes('/user/info')) {
                    const userStore = useUserStore()
                    userStore.reset()
                    console.warn('Token无效，已清除登录状态')
                    return Promise.reject(error)
                }
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
    autoContentType?: boolean; // 是否自动设置Content-Type，默认true
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
    onSSEEvent?: (event: string, data: any) => void; // SSE事件回调
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
        // 检查是否需要自动设置Content-Type
        const isFormData = config.data instanceof FormData;
        const shouldSetContentType = config.autoContentType !== false && !isFormData;
        
        if (shouldSetContentType) {
            xhr.setRequestHeader('Content-Type', 'application/json');
        }
        
        // ✅ SSO模式：使用 Authorization: Bearer <token>
        xhr.setRequestHeader('Authorization', `Bearer ${token}`);
        
        // ✅ 添加设备ID用于限流
        xhr.setRequestHeader('X-Device-Id', getDeviceId());
        
        // 设置自定义头部
        if (config.headers) {
            Object.entries(config.headers).forEach(([key, value]) => {
                // 如果是FormData且用户没有显式设置Content-Type，跳过设置
                if (isFormData && key.toLowerCase() === 'content-type' && !value) {
                    return;
                }
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
                
                // 处理流式数据 - 支持两种格式
                const lines = newData.split('\n');
                let currentEvent = '';
                
                for (const line of lines) {
                    // SSE格式处理
                    if (callbacks.onSSEEvent) {
                        if (line.startsWith('event:')) {
                            currentEvent = line.substring(6).trim();
                        } else if (line.startsWith('data:')) {
                            const jsonStr = line.substring(5).trim();
                            if (!jsonStr) continue; // 跳过空数据行
                            
                            try {
                                const data = JSON.parse(jsonStr);
                                callbacks.onSSEEvent(currentEvent, data);
                                
                                // 检查是否是完成事件 - 只标记完成，不立即abort
                                if (currentEvent === 'complete') {
                                    if (timeoutId) clearTimeout(timeoutId);
                                    callbacks.onComplete?.({ ...data, is_complete: true });
                                    // 延迟abort，让complete事件在应用层处理完成
                                    setTimeout(() => {
                                        if (!isAborted) {
                                            isAborted = true;
                                            xhr.abort();
                                            resolve();
                                        }
                                    }, 2000); // 增加延迟时间到2秒
                                }
                            } catch (e) {
                                console.warn('解析SSE数据失败:', e, 'event:', currentEvent, 'data:', jsonStr);
                                // 如果是complete事件解析失败，仍然尝试处理
                                if (currentEvent === 'complete') {
                                    console.log('complete事件解析失败，尝试强制完成');
                                    if (timeoutId) clearTimeout(timeoutId);
                                    isAborted = true;
                                    callbacks.onComplete?.({ is_complete: true });
                                    setTimeout(() => {
                                        xhr.abort();
                                        resolve();
                                    }, 100);
                                    return;
                                }
                            }
                        }
                    }
                    // 传统格式处理（向后兼容）
                    else if (line.startsWith('data: ')) {
                        const data = line.slice(6); // 移除 'data: ' 前缀
                        
                        try {
                            const parsed = JSON.parse(data);
                            
                            // 检查是否是完成事件
                            if (parsed.event_id === 2 || parsed.is_complete) {
                                if (timeoutId) clearTimeout(timeoutId);
                                callbacks.onComplete?.(parsed);
                                // 延迟abort，确保完成事件处理完成
                                setTimeout(() => {
                                    if (!isAborted) {
                                        isAborted = true; // 标记为主动abort
                                        xhr.abort(); // 主动结束请求
                                        resolve();
                                    }
                                }, 2000);
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
            // 如果是FormData，直接发送；否则JSON序列化
            if (config.data instanceof FormData) {
                xhr.send(config.data);
            } else {
                xhr.send(JSON.stringify(config.data));
            }
        } else {
            xhr.send();
        }
    });
};