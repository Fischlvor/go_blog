<template>
  <el-card class="tag-cloud">
    <el-row class="title">标签云</el-row>
    <el-tag v-for="item in tagCloudArray" :key="item.name" :type="item.type" size="large" effect="plain"
            @click="handleSearchJumps(item.name)">
      {{ item.name }} {{ item.article_count }}
    </el-tag>
  </el-card>
</template>

<script setup lang="ts">
import {ref} from "vue";
import {type TagDetail, articleTags} from "@/api/article";

const tagTypes = ["primary", "success", "info", "warning", "danger"]

interface TagCloudItem {
  name: string;
  article_count: number;
  type: string;
}

const tagCloudArray = ref<TagCloudItem[]>([])

const getTagCloudArray = async () => {
  const res = await articleTags()
  if (res.code === "0000") {
    const tagsArray: TagDetail[] = res.data
    for (let i = 0; i < tagsArray.length; i++) {
      const item = tagsArray[i];
      const tagCloud: TagCloudItem = {
        name: item.name,
        article_count: item.article_count,
        type: tagTypes[i % tagTypes.length]
      }
      tagCloudArray.value.push(tagCloud);
    }
  }
}

getTagCloudArray()

const handleSearchJumps = (tag: string) => {
  window.open("/search?tag=" + tag)
}
</script>

<style scoped lang="scss">
.tag-cloud {
  margin-bottom: 20px;

  .title {
    font-size: 24px;
    margin-bottom: 20px;
  }
}
</style>
