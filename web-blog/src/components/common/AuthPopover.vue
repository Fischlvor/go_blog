<template>
  <div class="auth-popover" ref="popoverContainer">
    <!-- 已登录状态 -->
    <el-popover
        v-if="isLoggedIn"
        width="280"
        trigger="click"
        :visible="layoutStore.state.popoverVisible"
        @show="layoutStore.show('popoverVisible')"
        @hide="layoutStore.hide('popoverVisible')"
    >
      <template #reference>
        <div class="avatar-circle"> <!-- 新增外层容器 -->
          <el-avatar
              :size="42"
              :src="userStore.state.userInfo.avatar"
              class="clickable-avatar"
              @click="handleStateChangePopover"
          />
        </div>
      </template>
      <user-popover-content/>
    </el-popover>

    <!-- 未登录状态 -->
    <el-popover
        v-else
        width="280"
        trigger="click"
        :visible="layoutStore.state.popoverVisible"
        @show="layoutStore.show('popoverVisible')"
        @hide="layoutStore.hide('popoverVisible')"
    >
      <template #reference>
        <div class="avatar-circle"> <!-- 新增外层容器 -->
          <el-avatar
              :size="42"
              :icon="User"
              class="clickable-avatar"
              @click="handleStateChangePopover"
          />
        </div>
      </template>
      <auth-popover-content/>
    </el-popover>
  </div>
</template>

<script setup lang="ts">
import { User } from "@element-plus/icons-vue";
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useLayoutStore } from "@/stores/layout";
import { useUserStore } from "@/stores/user";
import UserPopoverContent from "@/components/common/popover/UserPopoverContent.vue";
import AuthPopoverContent from "@/components/common/popover/AuthPopoverContent.vue";

const userStore = useUserStore();
const layoutStore = useLayoutStore();
const popoverContainer = ref<HTMLElement | null>(null);
const isLoggedIn = computed(() => userStore.isLoggedIn);

const handleStateChangePopover = () => {
  layoutStore.stateChange('popoverVisible');
};

// 添加全局点击事件监听作为备用
const handleClickOutside = (event: MouseEvent) => {
  const popover = document.querySelector('.el-popover');
  const reference = document.querySelector('.clickable-avatar');

  if (popoverContainer.value &&
      !popover?.contains(event.target as Node) &&
      !reference?.contains(event.target as Node)) {
    layoutStore.hide('popoverVisible');
  }
};

onMounted(() => {
  document.addEventListener('click', handleClickOutside);
});

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside);
});
</script>

<style scoped lang="scss">
.auth-popover {
  display: flex;
  --el-text-color-disabled: transparent;
  //padding: 0px;

  .avatar-circle {
    display: inline-flex;       /* 使用 flex 布局确保完美居中 */
    align-items: center;        /* 垂直居中 */
    justify-content: center;    /* 水平居中 */
    border-radius: 50%;         /* 确保是正圆 */
    padding: 2px;               /* 圆圈与头像之间的间距 */
    border: 2px solid #e0e0e0; /* 浅灰色边框 (#e0e0e0) */
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.1);
    transition: all 0.2s ease;
    width: 48px;                /* 固定宽度，确保圆形 */
    height: 48px;               /* 固定高度，确保圆形 */
    box-sizing: border-box;     /* 防止 padding 影响尺寸 */

    &:hover {
      transform: scale(1.05); /* 悬停放大 */
      box-shadow: 0 0 10px rgba(0, 0, 0, 0.2); /* 悬停阴影加深 */
    }
  }

  .clickable-avatar {
    cursor: pointer;
    transition: transform 0.2s;
  }

  .el-avatar--icon {
    font-size: 36px;
  }
}
</style>

<style lang="scss">
.el-popover.el-popper {
  background-color: rgba(255, 255, 255, 0.75);
  padding: 0;

  .el-popper__arrow:before {
    background-color: transparent;
  }
}
</style>