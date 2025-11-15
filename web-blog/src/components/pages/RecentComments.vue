<template>
  <el-card class="recent-comments">
    <el-row class="title">最新评论</el-row>
    <comment-item :comments="comments"/>
  </el-card>
</template>

<script setup lang="ts">
import CommentItem from "@/components/common/CommentItem.vue";
import {ref, watch} from "vue";
import {type Comment, commentNew} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import { parseEmojis, renderTextWithEmojisForText } from '@/utils/emojiParser'

const comments=ref<Comment[]>([])

// 解析并渲染评论内容中的 emoji 标记
const renderContentWithEmojis = async (text: string): Promise<string> => {
  if (!text) return text
  const parsed = await parseEmojis(text)
  return await renderTextWithEmojisForText(parsed)
}

// 递归处理最新评论的内容
const transformCommentsWithEmojis = async (list: Comment[]): Promise<Comment[]> => {
  const result: Comment[] = []
  for (const item of list) {
    const transformed: Comment = { ...item }
    transformed.content = await renderContentWithEmojis(item.content || '')
    if (item.children && item.children.length) {
      transformed.children = await transformCommentsWithEmojis(item.children)
    }
    result.push(transformed)
  }
  return result
}

const getRecentCommentInfo = async ()=>{
  const res = await commentNew()
  if (res.code===0){
    comments.value = await transformCommentsWithEmojis(res.data)
  }
}

getRecentCommentInfo()

const layoutStore = useLayoutStore()

watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getRecentCommentInfo()
  }
})
</script>

<style scoped lang="scss">
.recent-comments {
  margin-bottom: 20px;

  .title {
    font-size: 24px;
    margin-bottom: 20px;
  }
}
</style>
