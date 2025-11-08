import {defineStore} from 'pinia';
import {computed, ref, watch} from 'vue';
import {userInfo, login, logout,type RegisterRequest, register} from '@/api/user'
import type {User} from '@/api/user'
import type {LoginRequest} from '@/api/user'
import router from "@/router";


function initState() {
    const userInfo = ref<User>({
        id: 0,
        created_at: new Date(),
        updated_at: new Date(),
        uuid: '',
        username: '',
        email: '',
        openid: '',
        avatar: '',
        address: '',
        signature: '',
        role_id: 0,
        register: '',
        freeze: false,
    })
    const savedIsUserLoggedInBefore = localStorage.getItem('isUserLoggedInBefore');
    // ✅ SSO模式：从localStorage恢复accessToken
    const savedAccessToken = localStorage.getItem('accessToken') || '';
    return {
        userInfo,
        accessToken: savedAccessToken,
        userInfoInitialized: false,
        isUserLoggedInBefore: savedIsUserLoggedInBefore === 'true'
    }
}

export const useUserStore = defineStore('user', () => {
    const state = ref(initState());

    const reset = () => {
        state.value.userInfo = {
            id: 0,
            created_at: new Date(),
            updated_at: new Date(),
            uuid: '',
            username: '',
            email: '',
            openid: '',
            avatar: '',
            address: '',
            signature: '',
            role_id: 0,
            register: '',
            freeze: false,
        }
        state.value.accessToken = ''
        state.value.userInfoInitialized = false
        state.value.isUserLoggedInBefore = false
    }

    /* 登录*/
    const loginIn = async (loginInfo: LoginRequest) => {
        const res = await login(loginInfo)
        if (res.code === 0) {
            state.value.userInfo = res.data.user
            state.value.accessToken = res.data.access_token
            state.value.isUserLoggedInBefore = true
            return { message: res.msg ,flag: true };
        } else {
            return { message: res.msg ,flag: false }
        }
    }
    /* 注册*/
    const registerIn = async (registerInfo: RegisterRequest)=>{
        const res = await register(registerInfo)
        if (res.code===0){
            state.value.userInfo = res.data.user
            state.value.accessToken = res.data.access_token
            state.value.isUserLoggedInBefore = true
            return { message: res.msg ,flag: true };
        } else {
            return { message: res.msg ,flag: false };
        }
    }
    /* 登出*/
    const loginOut = async () => {
        await logout()
        const userStore = useUserStore()
        userStore.reset()
        localStorage.clear()
        router.push({name: 'index'}).then()
        ElMessage.success("Logout successfully.")
        return { message: "" ,flag: true };
    }

    watch(() => state.value.isUserLoggedInBefore, (newIsUserLoggedInBefore) => {
        localStorage.setItem('isUserLoggedInBefore', String(newIsUserLoggedInBefore));
    })

    // ✅ SSO模式：持久化accessToken
    watch(() => state.value.accessToken, (newAccessToken) => {
        if (newAccessToken) {
            localStorage.setItem('accessToken', newAccessToken);
        } else {
            localStorage.removeItem('accessToken');
        }
    })

    const initializeUserInfo = async () => {
        // ✅ SSO模式：必须同时满足 isUserLoggedInBefore 和 accessToken 存在
        if (state.value.isUserLoggedInBefore && state.value.accessToken && !state.value.userInfoInitialized) {
            // ✅ 标记为已初始化（无论成功与否，避免重复请求）
            state.value.userInfoInitialized = true
            
            try {
                const res = await userInfo();
                if (res.code === 0) {
                    state.value.userInfo = res.data;
                } else {
                    // 如果获取用户信息失败，清除登录状态
                    console.warn('获取用户信息失败，清除登录状态');
                    reset();
                }
            } catch (error) {
                console.error('获取用户信息异常:', error);
                reset();
            }
        }
    }

    // 添加验证方法：验证是否已登录
    const isLoggedIn = computed(() => {
        return state.value.userInfo.role_id !== 0;
    });

    // 添加验证方法：验证是否是管理员
    const isAdmin = computed(() => {
        return state.value.userInfo.role_id === 2;
    });

    return {
        state,
        reset,
        loginIn,
        registerIn,
        loginOut,
        initializeUserInfo,
        isLoggedIn,
        isAdmin
    };
});