<template>
  <div class="user-comment">
    <el-row class="title">我的评论</el-row>
    <div class="table-data" v-for="item in userCommentTableData">
      <div class="link">
        <el-link :href="'/article/'+item.article_id">{{ item.article_id }}</el-link>
      </div>
      <div class="item">
        <comment-item :comments="[item]"/>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import {ref, watch} from "vue";
import CommentItem from "@/components/common/CommentItem.vue";
import {type Comment, commentInfo} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";


const userCommentTableData = ref<Comment[]>()
const getUserCommentTableData = async () => {
  const table = await commentInfo()

  if (table.code === 0) {
    userCommentTableData.value = table.data
  }
}
getUserCommentTableData()

const layoutStore = useLayoutStore()

watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getUserCommentTableData()
    //console.log("watch "+layoutStore.state.shouldRefreshCommentList)
    // 重置 shouldRefreshCommentList 状态
    //layoutStore.state.shouldRefreshCommentList = false
    layoutStore.hide("shouldRefreshCommentList")
  }
})
</script>

<style scoped lang="scss">
.user-comment {
  background-color: rgba(255,255,255,0.7);
  padding: 10px;
  min-height: calc(100vh - 178px); /* 计算高度：视口高度 - 头部80px - 顶部间距100px（需根据实际调整） */


  .title {
    margin-bottom: 20px;
    font-size: 24px;
  }

  .table-data {
    display: flex;
    border: 1px solid #DCDFE6;
    background-color: rgba(255,255,255,0.3);

    .link{
      width: 200px;
      display: flex;
      padding: 5px;
      border: 1px solid #000000;
      .el-link{
        color: black;
        text-align: center;
      }
    }
    .item{
      width: 100%;

    }
  }
}
</style>
