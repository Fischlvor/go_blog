<template>
  <div id="app" :style="appStyle">
    <el-config-provider :locale="zhCn">
      <router-view/>
    </el-config-provider>
  </div>
</template>

<script setup lang="ts">
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import { cdn } from "@/utils/cdn"
import { computed } from "vue"

const appStyle = computed(() => ({
  backgroundImage: `url('${cdn("blog/5523c88dd347d1b7cc617f632b7efdb7-20251111150011.png")}')`,
  backgroundSize: "cover",
  backgroundPosition: "center",
  backgroundRepeat: "no-repeat",
  backgroundAttachment: "fixed",
  minHeight: "100vh",
  width: "100%",
  margin: "0",
  padding: "0"
}))
</script>

<style>
/* 全局样式，可放在 App.vue 或 index.css */
body {
  margin: 0;
  padding: 0;
}

/* 背景相关样式已移至内联绑定，便于使用环境变量域名 */
/* 当窗口宽高比 >= 图片宽高比时，改为宽度填满，高度按比例缩放 */
@media (min-aspect-ratio: 16/9) { /* 16/9 替换为你的图片实际宽高比 */
  #app {
    background-size: 100% auto;
  }
}
</style>
