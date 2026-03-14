<template>
  <div class="article-create-form">
    <el-form
        :model="articleCreateFormData"
        :rules="formRules"
        :validate-on-rule-change="false"
        label-width="90px"
    >
      <el-form-item label="文章封面" prop="cover">
        <el-upload
            :action="`${path}/image/upload`"
            drag
            with-credentials
            :headers="{'Authorization': `Bearer ${userStore.state.accessToken}`}"
            :show-file-list="false"
            :on-success="handleSuccess"
            :on-error="handleSuccess"
            name="file"
        >

          <el-image v-if="articleCreateFormData.cover" :src="articleCreateFormData.cover" alt=""/>

          <div v-else class="upload-content">
            <div class="container">
              <component is="UploadFilled" class="upload-filled"></component>
              <div class="el-upload__text">
                Drop file here or <em>click to upload</em>
              </div>
            </div>
          </div>

          <template #tip>
            <div class="el-upload__tip">
              jpg/png/jpeg/ico/tiff/gif/svg/webp files with a size less than 20MB.
            </div>
          </template>
        </el-upload>

        <el-input
            v-model="articleCreateFormData.cover"
            size="large"
            disabled
        />
      </el-form-item>
      <el-form-item label="文章标题" prop="title">
        <el-input
            v-model="articleCreateFormData.title"
            size="large"
            placeholder="请输入文章标题"
        />
      </el-form-item>
      <el-form-item label="文章类别" prop="category_id">
        <div style="display: flex; gap: 8px; width: 100%">
          <el-select v-model="articleCreateFormData.category_id" placeholder="请选择分类" size="large" style="flex: 1">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
          <el-button type="primary" size="large" @click="showNewCategoryDialog = true">新建</el-button>
        </div>
      </el-form-item>
      <el-form-item label="文章标签" prop="tag_ids">
        <div class="tag-checkbox-container">
          <el-checkbox-group v-model="articleCreateFormData.tag_ids">
            <el-checkbox v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</el-checkbox>
          </el-checkbox-group>
          <el-button type="primary" link @click="showNewTagDialog = true" style="margin-top: 8px">新建标签</el-button>
        </div>
      </el-form-item>
      <el-form-item label="文章简介" prop="abstract">
        <el-input
            v-model="articleCreateFormData.abstract"
            type="textarea"
            placeholder="请输入文章简介"
        />
      </el-form-item>
      <el-form-item label="文章可见性" prop="visibility">
        <el-radio-group v-model="articleCreateFormData.visibility">
          <el-radio value="public">公开</el-radio>
          <el-radio value="private">私有（仅自己可见）</el-radio>
        </el-radio-group>
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
              @click="layoutStore.state.articleCreateVisible = false"
          >取消
          </el-button>
        </div>
      </el-form-item>
    </el-form>

    <!-- 新建分类对话框 -->
    <el-dialog v-model="showNewCategoryDialog" title="新建分类" width="400px">
      <el-form :model="newCategoryForm" label-width="80px">
        <el-form-item label="分类名称" required>
          <el-input v-model="newCategoryForm.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="分类标识" required>
          <el-input v-model="newCategoryForm.slug" placeholder="请输入分类标识（英文）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNewCategoryDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateCategory">确定</el-button>
      </template>
    </el-dialog>

    <!-- 新建标签对话框 -->
    <el-dialog v-model="showNewTagDialog" title="新建标签" width="400px">
      <el-form :model="newTagForm" label-width="80px">
        <el-form-item label="标签名称" required>
          <el-input v-model="newTagForm.name" placeholder="请输入标签名称" />
        </el-form-item>
        <el-form-item label="标签标识" required>
          <el-input v-model="newTagForm.slug" placeholder="请输入标签标识（英文）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showNewTagDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreateTag">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {reactive, ref} from "vue";
import {ElMessage} from "element-plus";
import {articleCreate, articleCategory, articleTags, createCategory, createTag, type ArticleCreateRequest, type CategoryDetail, type TagDetail} from "@/api/article";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";

const props = defineProps<{
  title: string;
  content: string;
}>();

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const path = ref(import.meta.env.VITE_BASE_API)

// 表单验证规则
const formRules = {
  category_id: [{ required: true, message: '请选择文章分类', trigger: 'change' }],
}

// 表单数据
const articleCreateFormData = reactive({
  cover: '',      // 完整 URL（用于预览和提交）
  title: props.title,
  category_id: null as number | null, // 必填，默认为空
  tag_ids: [] as number[],
  abstract: '',
  content: props.content,
  visibility: 'public',
})

// 分类和标签列表
const categories = ref<CategoryDetail[]>([])
const tags = ref<TagDetail[]>([])

// 新建分类/标签对话框
const showNewCategoryDialog = ref(false)
const showNewTagDialog = ref(false)
const newCategoryForm = reactive({ name: '', slug: '' })
const newTagForm = reactive({ name: '', slug: '' })

// 加载分类和标签
const loadCategoriesAndTags = async () => {
  const [catRes, tagRes] = await Promise.all([articleCategory(), articleTags()])
  if (catRes.code === '0000') categories.value = catRes.data
  if (tagRes.code === '0000') tags.value = tagRes.data
}
loadCategoriesAndTags()

// 创建分类
const handleCreateCategory = async () => {
  if (!newCategoryForm.name || !newCategoryForm.slug) {
    ElMessage.warning('请填写完整信息')
    return
  }
  const res = await createCategory(newCategoryForm)
  if (res.code === '0000') {
    ElMessage.success('创建成功')
    showNewCategoryDialog.value = false
    articleCreateFormData.category_id = res.data.id
    newCategoryForm.name = ''
    newCategoryForm.slug = ''
    loadCategoriesAndTags()
  } else {
    ElMessage.error(res.message)
  }
}

// 创建标签
const handleCreateTag = async () => {
  if (!newTagForm.name || !newTagForm.slug) {
    ElMessage.warning('请填写完整信息')
    return
  }
  const res = await createTag(newTagForm)
  if (res.code === '0000') {
    ElMessage.success('创建成功')
    showNewTagDialog.value = false
    articleCreateFormData.tag_ids.push(res.data.id)
    newTagForm.name = ''
    newTagForm.slug = ''
    loadCategoriesAndTags()
  } else {
    ElMessage.error(res.message)
  }
}

// 构建提交数据
const buildSubmitData = (): ArticleCreateRequest => {
  return {
    title: articleCreateFormData.title,
    slug: '',
    content: articleCreateFormData.content,
    excerpt: articleCreateFormData.abstract,
    featured_image: articleCreateFormData.cover,
    category_id: articleCreateFormData.category_id || 0,
    tag_ids: articleCreateFormData.tag_ids,
    status: 'published',
    visibility: articleCreateFormData.visibility,
    is_featured: false,
  }
}

const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === "0000") {
    articleCreateFormData.cover = res.data.url
    ElMessage.success(res.message)
  }
}


const submitForm = async () => {
  // 前端验证
  if (!articleCreateFormData.content || articleCreateFormData.content.trim() === '') {
    ElMessage.warning('请输入文章内容')
    return
  }
  if (!articleCreateFormData.category_id) {
    ElMessage.warning('请选择文章分类')
    return
  }

  const submitData = buildSubmitData()
  const res = await articleCreate(submitData)
  if (res.code === "0000") {
    ElMessage.success(res.message)
    layoutStore.state.articleCreateVisible = false
    layoutStore.state.shouldRefreshArticleTable = true
  }
  // 错误提示已在 request.ts 拦截器中处理
}
</script>

<style scoped lang="scss">
.article-create-form {
  .el-form {
    .el-form-item {
      .el-image {
        height: 120px;
      }

      .upload-content {
        display: flex;
        height: 120px;

        .container {
          margin: auto;

          .upload-filled {
            height: 32px;
            width: 32px;
          }
        }
      }

      .tag-checkbox-container {
        max-height: 120px;
        overflow-y: auto;
        width: 100%;
        padding: 8px;
        border: 1px solid var(--el-border-color);
        border-radius: 4px;
      }

      .button-group {
        margin-left: auto;
      }
    }
  }
}
</style>

<style lang="scss">
.el-upload {
  --el-upload-dragger-padding-horizontal: 0px;
  --el-upload-dragger-padding-vertical: 0px;
  line-height: 0;
}
</style>
