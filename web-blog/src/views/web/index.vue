<template>
  <div class="web">
    <el-container class="app-container">
      <el-header class="app-header" :class="{ 'collapsed-header': !showFullHeader }">
        <web-navbar :class="{ 'fixed-nav': isNavFixed || !showFullHeader }"/>
        <div v-if="showFullHeader" class="site-info">
          <h1 class="typing-text">Welcome to my Blog!✨✨✨</h1>
        </div>
        <el-button v-if="showFullHeader" class="scroll-down-btn" @click="scrollDown" circle>
          <el-icon class="arrow-icon"><ArrowDown /></el-icon>
        </el-button>
      </el-header>
      <el-main class="app-main">
        <router-view/>
      </el-main>
      <el-footer class="app-footer">
        <web-footer/>
      </el-footer>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ArrowDown } from '@element-plus/icons-vue'
import WebNavbar from "@/components/layout/WebNavbar.vue";
import WebFooter from '@/components/layout/WebFooter.vue';
import { onMounted, onBeforeUnmount, ref, watch, nextTick } from 'vue';
import { useRoute } from 'vue-router';

const route = useRoute();
const isNavFixed = ref(false);
const showFullHeader = ref(route.path === '/');
let typingInterval: number | null = null;

// 滚动逻辑：每次滚动距离为 1个视口高度 - 当前滚动高度
const scrollDown = () => {
  const currentScrollY = window.scrollY;
  const viewportHeight = window.innerHeight;

  // 计算目标滚动位置：当前位置 + (视口高度 - 当前位置)，确保至少滚动一些距离
  const targetScrollY = Math.max(currentScrollY + 100, currentScrollY + (viewportHeight - currentScrollY));

  window.scrollTo({
    top: targetScrollY,
    behavior: 'smooth'
  });
};

// 封装打字动画逻辑为独立函数
const initTypingAnimation = () => {
  // 清除之前的定时器
  if (typingInterval) {
    clearInterval(typingInterval);
    typingInterval = null;
  }

  const text = "Welcome to my Blog!✨✨✨";
  const typingElement = document.querySelector('.typing-text');
  if (!typingElement) return;

  // 重置文本内容
  typingElement.textContent = '';

  let isDeleting = false;
  let i = 0;

  const type = () => {
    if (typingElement) {
      const currentText = isDeleting
          ? text.substring(0, i - 1)
          : text.substring(0, i + 1);

      typingElement.textContent = currentText;

      if (!isDeleting && i === text.length) {
        isDeleting = true;
        clearInterval(typingInterval as number);
        setTimeout(() => {
          typingInterval = setInterval(type, 100);
        }, 500);
      } else if (isDeleting && i === 0) {
        isDeleting = false;
        clearInterval(typingInterval as number);
        setTimeout(() => {
          typingInterval = setInterval(type, 200);
        }, 200);
      }

      isDeleting ? i-- : i++;
    }
  };

  typingInterval = setInterval(type, 200);
};

// 监听路由变化，控制header显示内容并重新初始化动画
watch(
    () => route.path,
    async (newPath) => {
      showFullHeader.value = newPath === '/';

      // 等待DOM更新完成后再初始化动画
      if (showFullHeader.value) {
        await nextTick();
        initTypingAnimation();
      } else if (typingInterval) {
        clearInterval(typingInterval);
        typingInterval = null;
      }
    },
    { immediate: true } // 初始执行一次
);

// 添加滚动事件监听
const handleScroll = () => {
  isNavFixed.value = window.scrollY > 50 && showFullHeader.value;
};

onMounted(() => {
  window.addEventListener('scroll', handleScroll);

  // 初始加载时初始化动画
  if (showFullHeader.value) {
    initTypingAnimation();
  }
});

// 组件卸载时清除定时器和事件监听
onBeforeUnmount(() => {
  if (typingInterval) {
    clearInterval(typingInterval);
    typingInterval = null;
  }
  window.removeEventListener('scroll', handleScroll);
});

</script>

<style scoped lang="scss">
/* 保持原有样式不变 */
.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh; // 关键：确保容器至少占满整个视口高度

  .app-header {
    padding: 0px;
    background-color: rgba(128,128,128,0.2);
    height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    position: relative;
    transition: all 0.3s ease;

    .fixed-nav {
      position: fixed !important;
      top: 0;
      left: 0;
      right: 0;
      z-index: 1000;
      background-color: rgba(0, 0, 0, 0.5) !important;
      box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
      transition: all 0.3s ease;
    }

    .web-navbar {
      background-color: rgba(255,255,255,0.0);
      position: absolute;
      top: 0;
      width: 100%;
      transition: all 0.3s ease;
    }

    .site-info {
      text-align: center;
      color: white;
      font-size: clamp(1.5rem, 5vw, 2.5rem);
      font-weight: bold;
      text-shadow: 0 0 10px rgba(255,255,255,0.5);

      .typing-text {
        display: inline-block;
        position: relative;
        &::after {
          content: '|';
          right: -10px;
          animation: blink 0.7s infinite;
        }
      }
    }

    .scroll-down-btn {
      position: absolute;
      bottom: 20px;
      left: 50%;
      transform: translateX(-50%);
      background: transparent !important;
      border: none !important;
      animation: bounce 2s infinite;

      .arrow-icon {
        color: #ffffff;
        font-size: 48px;
        filter: drop-shadow(0 0 1px currentColor);
        transform: scale(1.2);
      }
    }
  }

  .collapsed-header {
    height: auto;
    min-height: 100px;
    background-color: transparent;

  }

  .app-main {
    flex: 1;
  }

  .app-footer {
    padding: 0;
    margin: 0;
    height: auto;

    .web-footer{
      background-color: rgba(0, 0, 0, 0.5);
    }

  }
}

@keyframes bounce {
  0%, 100% {
    transform: translate(-50%, -10px);
  }
  50% {
    transform: translate(-50%, 0);
  }
}

@keyframes blink {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}
</style>