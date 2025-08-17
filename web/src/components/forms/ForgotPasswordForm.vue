<template>
  <div class="forgot-password-form">
    <el-form
        ref="forgotPasswordForm"
        :model="forgotPasswordFormData"
        :rules="rules"
        :validate-on-rule-change="false"
        hide-required-asterisk
        @keyup.enter="submitForm"
        label-width="90px"
    >

      <el-form-item label="邮箱" prop="email">
        <el-input
            v-model="forgotPasswordFormData.email"
            size="large"
            placeholder="请输入邮箱"
            @blur="isFormValid"
            @input="isFormValid"
        />
      </el-form-item>
      <el-form-item prop="captcha">
        <div class="captcha">
          <el-input
              v-model="emailRequest.captcha"
              placeholder="请输入图片验证码"
              size="large"
              maxlength="6"
              minlength="6"
          />
          <el-image :src="picPath" alt="" @click="emailVerify"/>
          <el-button
              @click="sendCode"
              :disabled="!isFormValid"
              :class="{ 'is-disabled': !isFormValid }"
          >发送验证码</el-button>
        </div>
      </el-form-item>
      <el-form-item label="邮箱验证码" prop="verification_code">
        <el-input
            v-model="forgotPasswordFormData.verification_code"
            size="large"
            placeholder="请输入邮箱验证码"
        />
      </el-form-item>
      <el-form-item label="密码" prop="new_password">
        <el-input
            v-model="forgotPasswordFormData.new_password"
            show-password
            size="large"
            type="password"
            placeholder="请输入新密码"
        />
      </el-form-item>
      <el-form-item label="确认密码">
        <el-input
            v-model="repeatPassword"
            show-password
            size="large"
            type="password"
            placeholder="请再次输入新密码"
        />
        <el-text v-if="forgotPasswordFormData.new_password!==repeatPassword">两次密码不一致！</el-text>
      </el-form-item>
      <el-form-item>
        <el-button
            type="primary"
            size="large"
            @click="submitForm"
        >确定
        </el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {computed, reactive, ref, watch,} from "vue";
import {forgotPassword, type ForgotPasswordRequest} from "@/api/user";
import {captcha, type EmailRequest, sendEmailVerificationCode} from "@/api/base";
import type {FormInstance, FormRules} from 'element-plus';
import {useLayoutStore} from "@/stores/layout";

const forgotPasswordForm = ref<FormInstance>()

const forgotPasswordFormData = reactive<ForgotPasswordRequest>({
  email: "",
  verification_code: "",
  new_password:"",
})

const repeatPassword = ref('')

const rules = reactive<FormRules<ForgotPasswordRequest>>({
  email: [{
    required: true,
    type:'email',
    trigger: 'blur',
    message: '请输入正确的邮箱格式'
  }],
  verification_code: [{
    required: true,
    len: 6,
    trigger: 'blur',
    message: '请输入6位的验证码'
  }],
  new_password:[{
    required:true,
    min:8,
    max:20,
    trigger:'change',
    message:'密码的长度应为8~20位'
  }]
})

const layoutStore = useLayoutStore()

const emailRequest = reactive<EmailRequest>({
  email: "",
  captcha: "",
  captcha_id: "",
})
const picPath = ref('')
const emailVerify = () => {
  captcha().then(async (res) => {
    picPath.value = res.data.pic_path
    emailRequest.captcha_id = res.data.captcha_id
  })
  emailRequest.captcha = ""
}
emailVerify()

const isEmailValid = ref(false); // 邮箱验证状态
const isCaptchaValid = ref(false); // 验证码验证状态
const emailPattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/; // 邮箱正则表达式
// 验证邮箱格式
const validateEmail = () => {
  isEmailValid.value = emailPattern.test(forgotPasswordFormData.email.trim());
};
// 验证验证码长度
const validateCaptcha = () => {
  isCaptchaValid.value = emailRequest.captcha.length === 6;
};
// 监听邮箱变化
watch(() => forgotPasswordFormData.email, validateEmail, { immediate: true });
// 监听验证码变化
watch(() => emailRequest.captcha, validateCaptcha, { immediate: true });
// 组合验证状态（如果需要同时依赖两者）
const isFormValid = computed(() => isEmailValid.value && isCaptchaValid.value);

const sendCode = async () => {
  emailRequest.email = forgotPasswordFormData.email
  const res = await sendEmailVerificationCode(emailRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
  } else {
    ElMessage.error(res.msg)
    emailVerify()
  }
}

const submitForm = async () => {
  const isValid: boolean = await new Promise((resolve) => {
    forgotPasswordForm.value?.validate((valid: boolean) => {
      resolve(valid)
    })
  })

  if (isValid) {
    const res = await forgotPassword(forgotPasswordFormData)

    if (res.code===0) {
      ElMessage.success(res.msg)
      layoutStore.state.forgotPasswordVisible = false
      layoutStore.state.popoverVisible = false
    }
  }
  return false
}
</script>


<style scoped lang="scss">
.forgot-password-form {
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

