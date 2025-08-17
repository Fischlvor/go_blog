<template>
  <div class="user-popover-content">
    <div class="title">
      <el-row>请登录以获取完整的功能体验。</el-row>
    </div>
    <div class="auth-button">
      <div class="button-item">
        <el-button class="login" icon="User" @click="layoutStore.show('loginVisible')"/>
        登录
        <el-dialog
            v-model="layoutStore.state.loginVisible"
            width="500"
            destroy-on-close
            align-center
            class="login-dialog"
        >
        <template #header>登录</template>
        <login-form/>
        </el-dialog>
      </div>

      <div class="button-item">
        <el-button class="register" icon="Edit" @click="layoutStore.show('registerVisible')"/>
        注册
        <el-dialog
            v-model="layoutStore.state.registerVisible"
            width="640"
            destroy-on-close
            align-center
            class="register-dialog"
        >
        <template #header>注册</template>
        <register-form/>
        </el-dialog>
      </div>

      <div class="button-item">
        <el-button class="forgot-password" icon="Unlock" @click="layoutStore.show('forgotPasswordVisible')"/>
        找回密码
        <el-dialog
            v-model="layoutStore.state.forgotPasswordVisible"
            width="570"
            destroy-on-close
            align-center
            class="forgot-password-dialog"
        >
        <template #header>找回密码</template>
        <forgot-password-form/>
        </el-dialog>
      </div>
    </div>
    <div class="container">
      <el-divider>快速登录</el-divider>
      <div class="oauth-login">
        <div class="login-item">
          <el-image class="qq-login" src="/image/qq_symbol.jpg" alt="" @click="qqLogin"/>
          QQ登录
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import LoginForm from "@/components/forms/LoginForm.vue";
import RegisterForm from "@/components/forms/RegisterForm.vue";
import ForgotPasswordForm from "@/components/forms/ForgotPasswordForm.vue";
import { useLayoutStore } from "@/stores/layout";
import { qqLoginURL } from "@/api/base";

const layoutStore = useLayoutStore();

const qqLogin = async () => {
  const res = await qqLoginURL();
  if (res.code === 0) {
    window.location.href = res.data;
  }
};
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

  .auth-button {
    display: flex;
    justify-content: center;

    .button-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      width: 80px;

      .el-button {
        border: none;
        --el-font-size-base: 32px;
        height: 62px;
        background-color: transparent;
      }
    }
  }

  .container {
    background-color: white;

    .el-divider {
      --el-bg-color: transparent;
    }

    .oauth-login {
      display: flex;
      justify-content: center;

      .login-item {
        display: flex;
        flex-direction: column;
        align-items: center;
        width: 80px;
        margin-top: auto;

        .el-image {
          height: 32px;
          margin-bottom: 8px;
        }
      }
    }
  }
}
</style>