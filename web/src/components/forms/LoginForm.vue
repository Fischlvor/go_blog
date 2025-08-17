<template>
  <div class="login-dialog-container">
    <div class="login-form">
    <el-form
        ref="loginForm"
        :model="loginFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
        @keyup.enter="submitForm"
        label-width="60px"
    >
      <el-form-item label="邮箱" prop="email">
        <el-input
            v-model="loginFormData.email"
            size="large"
            placeholder="请输入邮箱"
            suffix-icon="user"
        />
      </el-form-item>
      <el-form-item label="密码" prop="password">
        <el-input
            v-model="loginFormData.password"
            show-password
            size="large"
            type="password"
            placeholder="请输入密码"
        />
      </el-form-item>
      <el-form-item label="验证码" prop="captcha">
        <div class="captcha">
          <el-input
              v-model="loginFormData.captcha"
              placeholder="请输入验证码"
              size="large"
          />
          <el-image :src="picPath" alt="" @click="loginVerify"/>
        </div>
      </el-form-item>
      <el-form-item>
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
            :disabled="!isFormValid"
        >登 录
        </el-button>
      </el-form-item>
    </el-form>
  </div>
  </div>
</template>

<script setup lang="ts">
import {computed, reactive, ref, watch,} from "vue";
import type {LoginRequest} from "@/api/user";
import {useUserStore} from "@/stores/user";
import {captcha} from "@/api/base";
import type {FormInstance, FormRules} from 'element-plus';
import {useLayoutStore} from "@/stores/layout";

const loginForm = ref<FormInstance>()

const loginFormData = reactive<LoginRequest>({
  email: "",
  password: "",
  captcha: "",
  captcha_id: "",
})

const rules = reactive<FormRules<LoginRequest>>({
  email: [{
    required: true,
    type: 'email',
    trigger: 'blur',
    message: '请输入正确的邮箱格式'
  }],
  password: [{
    required: true,
    min: 8,
    max: 20,
    trigger: 'change',
    message: '密码的长度应为8~20位'
  }],
  captcha: [{
    required: true,
    len: 6,
    trigger: 'blur',
    message: '请输入6位的验证码'
  }],
})

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const picPath = ref('')

const loginVerify = () => {
  captcha().then(async (res) => {
    picPath.value = res.data.pic_path
    loginFormData.captcha_id = res.data.captcha_id
  })
  loginFormData.captcha = ""
}
loginVerify()


const isEmailValid = ref(false); // 邮箱有效性
const isPasswordValid = ref(false); // 密码有效性
const isCaptchaValid = ref(false); // 验证码有效性
// 验证邮箱格式
const validateEmail = (email: string) => {
  const emailReg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  isEmailValid.value = emailReg.test(email.trim()) && email.trim() !== "";
};
// 验证密码格式
const validatePassword = (password: string) => {
  isPasswordValid.value = password.length >= 8 && password.length <= 20;
};
// 验证验证码格式
const validateCaptcha = (captcha: string) => {
  isCaptchaValid.value = captcha.length === 6;
};
// 监听邮箱变化
watch(() => loginFormData.email, (val) => validateEmail(val), { immediate: true });
// 监听密码变化
watch(() => loginFormData.password, (val) => validatePassword(val), { immediate: true });
// 监听验证码变化
watch(() => loginFormData.captcha, (val) => validateCaptcha(val), { immediate: true });
// 组合验证状态
const isFormValid = computed(() =>
    isEmailValid.value &&
    isPasswordValid.value &&
    isCaptchaValid.value
);

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    loginForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const res = await userStore.loginIn(loginFormData)

    if (!res.flag) {
      ElMessage.error(res.message)
      loginVerify()
    } else {
      ElMessage.success(res.message)
      layoutStore.state.loginVisible = false
      layoutStore.state.popoverVisible = false
    }
  }
  loginVerify()
  return false
}
</script>


<style scoped lang="scss">
.login-form {
  //background-color: rgba(255, 255, 255, 0.5); /* 半透明白色背景 */
}
.login-form {
  // 主容器：使用 flex 布局实现整体居中
  display: flex;
  justify-content: center;
  align-items: center;


  .el-form {

    .captcha {
      display: flex;
      gap: 12px; // 验证码输入框和图片之间的间距

      .el-input {
        flex: 1; // 输入框占满剩余宽度
      }

      .el-image {
        width: 120px; // 验证码图片宽度
        height: 40px; // 验证码图片高度
        cursor: pointer;
        transition: transform 0.2s;

        &:hover {
          transform: scale(1.20); // 悬停时微放大
        }
      }
    }

    // 按钮样式优化
    .el-button--primary {
      transition: all 0.2s;

      &:hover {
        transform: translateY(-1px); // 悬停时上移
        box-shadow: 0 4px 12px rgba(22, 93, 255, 0.2); // 添加阴影
      }

      &:active {
        transform: translateY(0); // 点击时恢复原位
      }
    }


  }
}
</style>


