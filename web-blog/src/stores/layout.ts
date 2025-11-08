import {defineStore} from "pinia";
import {ref, watch} from 'vue';

function initState() {
    const savedIsCollapse = localStorage.getItem('isCollapse');
    return {
        isCollapse: savedIsCollapse === 'true',
        popoverVisible: false,
        userPopoverVisible: false,
        authPopoverVisible: false,
        loginVisible: false,
        registerVisible: false,
        forgotPasswordVisible: false,
        passwordResetVisible:false,
        shouldRefreshUserTable:false,
        shouldRefreshImageTable:false,
        articleCreateVisible:false,
        articleUpdateVisible:false,
        shouldRefreshArticleTable:false,
        friendLinkCreateVisible: false,
        friendLinkUpdateVisible: false,
        shouldRefreshFriendLinkTable: false,
        advertisementCreateVisible:false,
        advertisementUpdateVisible:false,
        shouldRefreshAdvertisementTable: false,
        feedbackReplyVisible:false,
        shouldRefreshFeedbackTable:false,
        shouldRefreshCommentTable:false,
        shouldRefreshCommentList:false,
        emojiPopoverVisible:false,
    };
}

export const useLayoutStore = defineStore('layout', () => {
    const state = ref(initState());

    // 监视 isCollapse 的变化并更新 localStorage
    watch(() => state.value.isCollapse, (newIsCollapse) => {
        localStorage.setItem('isCollapse', String(newIsCollapse));
    });

    const show = (field: keyof typeof state.value) => {
        state.value[field] = true;
    };

    const hide = (field: keyof typeof state.value) => {
        state.value[field] = false;
    };

    // 保持原来的 stateChange 方法
    const stateChange = (field: keyof typeof state.value) => {
        state.value[field] = !state.value[field];
    };

    return { state, show, hide, stateChange };
});