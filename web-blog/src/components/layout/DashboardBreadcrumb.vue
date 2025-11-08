<template>
  <div class="breadcrumb">
    <el-breadcrumb separator="/">
      <template v-for="item in route.matched">
        <el-breadcrumb-item v-if="item.components" :to="{ path: item.path }"> {{ item.meta.title }}</el-breadcrumb-item>
        <el-breadcrumb-item v-else>{{ item.meta.title }}</el-breadcrumb-item>
      </template>
    </el-breadcrumb>
  </div>
</template>

<script setup lang="ts">
import {useRoute} from "vue-router";

const route = useRoute()
</script>

<style scoped lang="scss">
.breadcrumb {
  display: flex;          /* 启用flex布局 */
  align-items: flex-end;  /* 内容下对齐 */
  //height: 40px;
  //background-color: rgba(255, 255, 255, 0.1);

  .el-breadcrumb {
    font-size: 18px;
    //padding-left: 15px;
    //padding-right: 15px;
    //padding-top: 15px;
    //padding-bottom: 5px;


    // 面包屑项文字颜色
    :deep(.el-breadcrumb__inner),
    :deep(.el-breadcrumb__inner a) {
      color: white;
      font-weight: normal;
    }

    // 分隔符颜色
    :deep(.el-breadcrumb__separator) {
      color: rgba(255, 255, 255, 0.8);
    }

    // 鼠标悬停效果
    :deep(.el-breadcrumb__inner:hover),
    :deep(.el-breadcrumb__inner a:hover) {
      color: var(--el-color-primary-light-3);
    }
  }
}
</style>
