<template>
  <div class="user-feedback">
    <el-row class="title">我的反馈</el-row>

    <el-table
        :data="userFeedbackTableData"
        :row-style="{height: '80px'}"
    >
      <el-table-column label="时间" width="150">
        <template #default="scope:{ row: Feedback, column: any, $index: number }">
          {{ getTime(scope.row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column prop="content" label="内容"/>
      <el-table-column prop="reply" label="回复"/>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type Feedback, feedbackInfo} from "@/api/feedback";

const getTime = (date: Date): string => {
  const time = new Date(date)
  return time.toLocaleString()
}

const userFeedbackTableData = ref<Feedback[]>()
const getUserFeedbackTableData = async () => {
  const table = await feedbackInfo()

  if (table.code === 0) {
    userFeedbackTableData.value = table.data
  }
}
getUserFeedbackTableData()
</script>

<style scoped lang="scss">
.user-feedback {
  background-color: rgba(255,255,255,0.7);
  padding: 10px;
  min-height: calc(100vh - 178px); /* 计算高度：视口高度 - 头部80px - 顶部间距100px（需根据实际调整） */

  .title {
    margin-bottom: 20px;
    font-size: 24px;
  }

  .el-table {
    border: 1px solid #DCDFE6;
  }
}
</style>
