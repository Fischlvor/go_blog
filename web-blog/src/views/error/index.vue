<template>
  <el-container class="error">
    <el-header>
      <web-navbar :noScroll="true"/>
    </el-header>
    <el-main>
      <el-image src="/image/404.jpg" alt=""/>
    </el-main>
    <el-footer>
      <web-footer/>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">

import WebNavbar from "@/components/layout/WebNavbar.vue";
import WebFooter from "@/components/layout/WebFooter.vue";
</script>

<style scoped lang="scss">
.error{
  display: flex;
  flex-direction: column;
  min-height: 100vh; // 关键：确保容器至少占满整个视口高度

  .el-header{
    padding: 0;
    margin: 0;
    height: auto;

    .web-navbar{
      background-color: rgba(0, 0, 0, 0.5) !important;
    }
  }


  .el-main{
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 0;
    //background-color: #f5f7fa;
    .el-image{
      height: 50vh;
    }
  }

  .el-footer{
    padding: 0;
    margin: 0;
    height: auto;

    .web-footer{
      background-color: rgba(0, 0, 0, 0.5);
    }

  }
}



</style>
