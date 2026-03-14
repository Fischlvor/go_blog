<template>
  <div class="article-update-form">
    <el-form
        :model="articleUpdateFormData"
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

          <el-image v-if="articleUpdateFormData.cover" :src="articleUpdateFormData.cover" alt=""/>

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
            v-model="articleUpdateFormData.cover"
            size="large"
            disabled
        />
      </el-form-item>
      <el-form-item label="文章标题" prop="title">
        <el-input
            v-model="articleUpdateFormData.title"
            size="large"
            placeholder="请输入文章标题"
        />
      </el-form-item>
      <el-form-item label="文章类别" prop="category_id">
        <div style="display: flex; gap: 10px; width: 100%;">
          <el-select v-model="articleUpdateFormData.category_id" placeholder="请选择分类" size="large" style="flex: 1;">
            <el-option v-for="cat in categories" :key="cat.id" :label="cat.name" :value="cat.id" />
          </el-select>
          <el-button type="primary" size="large" @click="showNewCategoryDialog = true">新建</el-button>
        </div>
      </el-form-item>
      <el-form-item label="文章标签" prop="tag_ids">
        <div style="display: flex; gap: 10px; width: 100%; align-items: flex-start;">
          <div class="tag-checkbox-container" ref="tagContainerRef" style="flex: 1;" @wheel.prevent="handleTagScroll">
            <el-checkbox-group v-model="articleUpdateFormData.tag_ids">
              <el-checkbox v-for="tag in tags" :key="tag.id" :value="tag.id">{{ tag.name }}</el-checkbox>
            </el-checkbox-group>
          </div>
          <el-button type="primary" size="large" @click="showNewTagDialog = true">新建</el-button>
        </div>
      </el-form-item>
      <el-form-item label="文章简介" prop="abstract">
        <el-input
            v-model="articleUpdateFormData.abstract"
            type="textarea"
            placeholder="请输入文章简介"
        />
      </el-form-item>
      <el-form-item label="文章可见性" prop="visibility">
        <el-radio-group v-model="articleUpdateFormData.visibility">
          <el-radio value="public">公开</el-radio>
          <el-radio value="private">私有（仅自己可见）</el-radio>
        </el-radio-group>
      </el-form-item>

      <el-form-item label="文章内容" prop="content">
        <el-button @click="drawer = true" icon="EditPen">编辑内容</el-button>
        <el-drawer v-model="drawer" :direction="direction" size="80%">
          <template #header>
            编辑内容
          </template>
          <template #default>
            <!-- 添加 flex 布局容器 -->
            <div class="editor-container">
              <MdEditor
                  v-model="articleUpdateFormData.content"
                  @onUploadImg="onUploadImg"
                  class="full-height-editor"
              />
            </div>
          </template>
<!--          <template #footer>-->
<!--            <el-text>点击上方X或外部任意区域即可退出编辑</el-text>-->
<!--          </template>-->
        </el-drawer>
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
              @click="layoutStore.state.articleUpdateVisible = false"
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
import {type DrawerProps, ElMessage} from "element-plus";
import {type Article, articleUpdate, articleCategory, articleTags, createTag, createCategory, type ArticleUpdateRequest, type CategoryDetail, type TagDetail} from "@/api/article";
import type {ApiResponse} from "@/utils/request";
import type {ImageUploadResponse} from "@/api/image";
import {useUserStore} from "@/stores/user";
import {useLayoutStore} from "@/stores/layout";
import {MdEditor} from 'md-editor-v3';
import axios, {type AxiosResponse} from "axios";

const props = defineProps<{
  article: Article;
}>();

const userStore = useUserStore()
const layoutStore = useLayoutStore()

const path = ref(import.meta.env.VITE_BASE_API)

// 表单验证规则
const formRules = {
  category_id: [{ required: true, message: '请选择文章分类', trigger: 'change' }],
}

// 表单数据
const articleUpdateFormData = reactive({
  cover: props.article.featured_image || '',  // 完整 URL
  title: props.article.title,
  category_id: props.article.category?.id || 0,
  tag_ids: props.article.tags?.map(t => t.id) || [],
  abstract: props.article.excerpt || '',
  content: props.article.content || '',
  visibility: props.article.visibility || 'public',
})

// 保存原始数据用于提交
const originalArticle = props.article

// 分类和标签列表
const categories = ref<CategoryDetail[]>([])
const tags = ref<TagDetail[]>([])

// 标签滚动控制
const tagContainerRef = ref<HTMLElement | null>(null)
const ROW_HEIGHT = 40 // 每行高度
const isScrolling = ref(false)
const handleTagScroll = (e: WheelEvent) => {
  if (tagContainerRef.value && !isScrolling.value) {
    isScrolling.value = true
    const direction = e.deltaY > 0 ? 1 : -1
    tagContainerRef.value.scrollTo({
      top: tagContainerRef.value.scrollTop + direction * ROW_HEIGHT,
      behavior: 'smooth'
    })
    // 滚动动画大约 150ms，设置延迟防止连续触发
    setTimeout(() => {
      isScrolling.value = false
    }, 200)
  }
}

// 新建分类对话框
const showNewCategoryDialog = ref(false)
const newCategoryForm = reactive({
  name: '',
  slug: ''
})

// 创建分类
const handleCreateCategory = async () => {
  if (!newCategoryForm.name || !newCategoryForm.slug) {
    ElMessage.warning('请填写完整的分类信息')
    return
  }
  const res = await createCategory(newCategoryForm)
  if (res.code === '0000') {
    ElMessage.success('分类创建成功')
    showNewCategoryDialog.value = false
    newCategoryForm.name = ''
    newCategoryForm.slug = ''
    // 重新加载分类列表
    const catRes = await articleCategory()
    if (catRes.code === '0000') categories.value = catRes.data
  } else {
    ElMessage.error(res.message || '创建失败')
  }
}

// 新建标签对话框
const showNewTagDialog = ref(false)
const newTagForm = reactive({
  name: '',
  slug: ''
})

// 创建标签
const handleCreateTag = async () => {
  if (!newTagForm.name || !newTagForm.slug) {
    ElMessage.warning('请填写完整的标签信息')
    return
  }
  const res = await createTag(newTagForm)
  if (res.code === '0000') {
    ElMessage.success('标签创建成功')
    showNewTagDialog.value = false
    newTagForm.name = ''
    newTagForm.slug = ''
    // 重新加载标签列表
    const tagRes = await articleTags()
    if (tagRes.code === '0000') tags.value = tagRes.data
  } else {
    ElMessage.error(res.message || '创建失败')
  }
}

// 加载分类和标签
const loadCategoriesAndTags = async () => {
  const [catRes, tagRes] = await Promise.all([articleCategory(), articleTags()])
  if (catRes.code === '0000') categories.value = catRes.data
  if (tagRes.code === '0000') tags.value = tagRes.data
}
loadCategoriesAndTags()

const handleSuccess = (res: ApiResponse<ImageUploadResponse>) => {
  if (res.code === "0000") {
    articleUpdateFormData.cover = res.data.url
    ElMessage.success(res.message)
  }
}

// 构建提交数据
const buildSubmitData = (): ArticleUpdateRequest => {
  return {
    slug: originalArticle.slug,
    title: articleUpdateFormData.title,
    content: articleUpdateFormData.content,
    excerpt: articleUpdateFormData.abstract,
    featured_image: articleUpdateFormData.cover,
    category_id: articleUpdateFormData.category_id,
    tag_ids: articleUpdateFormData.tag_ids,
    status: originalArticle.status || 'published',
    visibility: articleUpdateFormData.visibility,
    is_featured: originalArticle.is_featured || false,
  }
}

const drawer = ref(false)
const direction = ref<DrawerProps['direction']>('rtl')

const onUploadImg = async (files: File[], callback: (urls: string[]) => void): Promise<void> => {
  const res = await Promise.all(
      files.map((file) => {
        return new Promise<AxiosResponse<ApiResponse<ImageUploadResponse>>>((resolve, reject) => {
          const form = new FormData();
          form.append('image', file);

          axios
              .post('/api/image/upload', form, {
                headers: {
                  'Content-Type': 'multipart/form-data',
                },
                withCredentials: true,
              })
              .then((response) => resolve(response))
              .catch((error) => reject(error));
        });
      })
  );

  callback(res.map((item) => item.data.data.url));
};

const submitForm = async () => {
  const submitData = buildSubmitData()
  const res = await articleUpdate(submitData)
  if (res.code === "0000") {
    ElMessage.success(res.message)
    layoutStore.state.articleUpdateVisible = false
    layoutStore.state.shouldRefreshArticleTable = true
  }
  // 错误提示已在 request.ts 拦截器中处理
}
</script>

<style scoped lang="scss">
.article-update-form {
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

      .editor-container {
        height: 100%;

        .full-height-editor {
          height: 100%;
        }
      }

      .tag-checkbox-container {
        max-height: 40px;
        overflow-y: auto;
        overflow-x: hidden;
        width: 100%;
        padding: 0;
        border: 1px solid var(--el-border-color);
        border-radius: 4px;
        scrollbar-width: none; // Firefox
        -ms-overflow-style: none; // IE/Edge
        
        &::-webkit-scrollbar {
          display: none; // Chrome/Safari
        }
        
        .el-checkbox-group {
          display: flex;
          flex-wrap: wrap;
          
          .el-checkbox {
            flex: 0 0 33.33%;
            margin: 0;
            height: 40px;
            line-height: 40px;
            padding-left: 10px;
          }
        }
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

.el-drawer{
  .md-editor .md-editor-toolbar-wrapper .md-editor-toolbar svg.md-editor-icon {
    height: 24px;
    width: 24px;
  }
}
</style>
