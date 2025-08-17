<template>
  <div id="app">
    <el-config-provider :locale="zhCn">
      <router-view/>
    </el-config-provider>
  </div>
</template>

<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'
</script>

<style>
/* 全局样式，可放在 App.vue 或 index.css */
body {
  margin: 0;
  padding: 0;
}

#app {
  /* 背景图片设置 */
  background-image: url('/image/bg.png');
  background-size: cover;       /* 覆盖整个容器 */
  background-position: center;  /* 图片居中 */
  background-repeat: no-repeat; /* 不重复 */
  background-attachment: fixed; /* 背景固定，内容可滚动 */

  /* 确保容器足够高以显示滚动内容 */
  min-height: 100vh;

  /* 其他样式 */
  width: 100%;
  margin: 0;
  padding: 0;
}

/* 当窗口宽高比 >= 图片宽高比时，改为宽度填满，高度按比例缩放 */
@media (min-aspect-ratio: 16/9) { /* 16/9 替换为你的图片实际宽高比 */
  #app {
    background-size: 100% auto;
  }
}
</style>
