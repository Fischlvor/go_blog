<template>
  <div class="user-popover-content">
    <div class="title">
      <el-row>欢迎回来，{{ userStore.state.userInfo.username }}！</el-row>
    </div>
    <div class="user-info">
      <el-avatar :size="50" :src="userStore.state.userInfo.avatar"/>
      <el-row>{{ userStore.state.userInfo.username }}</el-row>
      <el-row>{{ userStore.state.userInfo.address }}</el-row>
    </div>
    <el-row class="uuid">uuid：{{ userStore.state.userInfo.uuid }}</el-row>
    <div class="container">
      <el-row>签名：{{ userStore.state.userInfo.signature }}</el-row>
      <div class="action-button">
        <el-button @click="logout">退出登录</el-button>
        <el-button @click="goIndexOrToDashboard">{{ label }}</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from "@/stores/user";
import { useRoute } from "vue-router";
import router from "@/router";
import {useLayoutStore} from "@/stores/layout";

const userStore = useUserStore();
const layoutStore = useLayoutStore();
const route = useRoute();
const isInDashboard = route.matched.some(record => record.name === 'dashboard');
const label = isInDashboard ? '返回首页' : '进入后台';

const goIndexOrToDashboard = (() => {
  console.log("in")
  if (isInDashboard) {
    router.push({name: 'index'});
  } else {
    router.push({name: 'home'});
  }
});

const logout = (() => {
  userStore.loginOut()
  layoutStore.hide('popoverVisible')
});
</script>

<style scoped lang="scss">
.user-popover-content {
  .title {
    display: flex;
    height: 60px;

    .el-row {
      margin: auto auto;
    }
  }

  .user-info {
    text-align: center;

    .el-row {
      display: flex;
      justify-content: center;
      margin: 5px;
    }
  }

  .uuid {
    margin: 10px 20px;
  }

  .container {
    background-color: white;

    .el-row {
      padding: 20px;
    }

    .action-button {
      display: flex;
      justify-content: center;

      .el-button {
        border: none;
        background-color: transparent;
        margin-bottom: 10px;
      }
    }
  }
}
</style>