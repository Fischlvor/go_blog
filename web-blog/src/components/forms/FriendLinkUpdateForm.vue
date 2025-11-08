<template>
  <div class="friend-link-update-form">
    <el-form
        :model="friendLinkUpdateFormData"
        :validate-on-rule-change="false"
    >


      <el-form-item label="logo图片" prop="logo">
        <el-upload
            :action="`${path}/image/upload`"
            drag
            with-credentials
            :headers="{'Authorization': `Bearer ${userStore.state.accessToken}`}"
            :show-file-list="false"
            :on-success="handleSuccess"
            :on-error="handleSuccess"
            name="image"
        >

          <el-image v-if="friendLinkUpdateFormData.logo" :src="friendLinkUpdateFormData.logo" alt=""/>
          <el-image v-else :src="friendLinkUpdateFormData.oldLogo" alt=""/>

          <template #tip>
            <div class="el-upload__tip">
              jpg/png/jpeg/ico/tiff/gif/svg/webp files with a size less than 20MB.
            </div>
          </template>
        </el-upload>
      </el-form-item>
      <el-form-item label="友链链接" prop="link">
        <el-input
            v-model="friendLinkUpdateFormData.link"
            size="large"
            placeholder="请输入友链链接"
        />
      </el-form-item>
      <el-form-item label="友链名称" prop="name">
        <el-input
            v-model="friendLinkUpdateFormData.name"
            size="large"
            placeholder="请输入友链名称"
        />
      </el-form-item>
      <el-form-item label="友链描述" prop="description">
        <el-input
            v-model="friendLinkUpdateFormData.description"
            size="large"
            placeholder="请输入友链描述"
        />
      </el-form-item>
      <el-form-item>
        <div class="button-group">
          <el-button
              type="primary"
              size="large"
              @click="submitForm"
          >确定
          </el-button>
          <el-button
              size="large"
              @click="layoutStore.state.friendLinkUpdateVisible = false"
          >取消
          </el-button>
        </div>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import {defineProps, reactive} from 'vue';
import {ElMessage} from "element-plus";
import {type FriendLink, friendLinkUpdate, type FriendLinkUpdateRequest} from "@/api/friend-link";
import {useLayoutStore} from "@/stores/layout";
import {type ApiResponse} from "@/utils/request";
import {type ImageUploadResponse} from "@/api/image";
import {ref} from "vue";
import {useUserStore} from "@/stores/user";

const layoutStore = useLayoutStore()
const userStore = useUserStore()
const path = ref(import.meta.env.VITE_BASE_API)

const props = defineProps<{
  friendLink: FriendLink;
}>();

const friendLinkUpdateFormData = reactive<FriendLinkUpdateRequest>({
  id: props.friendLink.id,
  logo: "",
  oldLogo: props.friendLink.logo,
  link: props.friendLink.link,
  name: props.friendLink.name,
  description: props.friendLink.description,
})

const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === 0) {
    friendLinkUpdateFormData.logo = res.data.url
    ElMessage.success(res.msg)
  }
}

const submitForm = async () => {
  if (friendLinkUpdateFormData.logo === friendLinkUpdateFormData.oldLogo) {
    delete friendLinkUpdateFormData.logo
    delete friendLinkUpdateFormData.oldLogo
  }
  const res = await friendLinkUpdate(friendLinkUpdateFormData)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    layoutStore.state.shouldRefreshFriendLinkTable = true
    layoutStore.state.friendLinkUpdateVisible = false
  }
}
</script>

<style scoped lang="scss">
.friend-link-update-form {
  .el-form {
    .el-form-item {
      .el-image {
        height: 120px;
      }

      .button-group {
        margin-left: auto;
      }
    }
  }
}
</style>
