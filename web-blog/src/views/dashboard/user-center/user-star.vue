<template>
  <div class="user-star">
    <el-row class="title">我的收藏</el-row>

    <el-table
        :data="articleLikesListData"
    >
      <el-table-column label="封面" width="100">
        <template #default="scope:{ row: Article, column: any, $index: number }">
          <el-image :src="scope.row.featured_image" alt=""/>
        </template>
      </el-table-column>
      <el-table-column prop="title" label="标题" width="120"/>
      <el-table-column prop="category.name" label="类别" width="80"/>
      <el-table-column label="标签" width="120">
        <template #default="scope:{ row: Article, column: any, $index: number }">
          <el-tag v-for="tag in scope.row.tags" :key="tag.id">{{ tag.name }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="简介">
        <template #default="scope:{ row: Article, column: any, $index: number }">
          <el-text line-clamp="5">{{ scope.row.excerpt }}</el-text>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="发布时间" width="102"/>
      <el-table-column label="文章id" width="200">
        <template #default="scope:{ row: Article, column: any, $index: number }">
          <el-link :href="'/article/'+scope.row.slug">{{ scope.row.slug }}</el-link>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
        :current-page="page"
        :page-size="page_size"
        :page-sizes="[10, 30, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
    />
  </div>
</template>

<script setup lang="ts">
import {nextTick, onMounted, reactive, ref, watch} from "vue";
import {type Article} from "@/api/article";
import {useRoute, useRouter} from "vue-router";
import type {Hit, PageInfo} from "@/api/common";
import {articleLikesList} from "@/api/article";


const articleLikesListData = ref<Article[]>()
const page = ref(1)
const page_size = ref(10)
const total = ref(0)

const articleLikesListRequest = reactive<PageInfo>({
  page: 1,
  page_size: 10,
})

const route = useRoute()
const router = useRouter()

onMounted(() => {
  page.value = Number(route.query.page) || 1
  page_size.value = Number(route.query.page_size) || 10
})

const getArticleLikesListData = async () => {
  articleLikesListRequest.page = page.value
  articleLikesListRequest.page_size = page_size.value

  const table = await articleLikesList(articleLikesListRequest)

  if (table.code === "0000") {
    articleLikesListData.value = table.data.list
    total.value = table.data.total_items

    await router.push({
      path: router.currentRoute.value.path,
      query: {
        page: articleLikesListRequest.page,
        page_size: articleLikesListRequest.page_size,
      },
    })
  }
}

watch(() => route.query, (newQuery) => {
  articleLikesListRequest.page = Number(newQuery.page) || 1
  articleLikesListRequest.page_size = Number(newQuery.page_size) || 10
}, {immediate: true})

nextTick(() => {
  getArticleLikesListData()
})

const handleSizeChange = (val: number) => {
  page_size.value = val
  getArticleLikesListData()
}

const handleCurrentChange = (val: number) => {
  page.value = val
  getArticleLikesListData()
}
</script>

<style scoped lang="scss">
.user-star {
  background-color: rgba(255,255,255,0.7);
  padding: 10px;
  min-height: calc(100vh - 178px); /* 计算高度：视口高度 - 头部80px - 顶部间距100px（需根据实际调整） */


  .title {
    margin-bottom: 20px;
    font-size: 24px;
  }

  .el-table {
    border: 1px solid #DCDFE6;

    .el-image {
      height: 48px;
    }
  }

  .el-pagination {
    display: flex;
    justify-content: center;
  }
}
</style>
