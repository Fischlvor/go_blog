<template>
  <div class="article">
    <el-container class="main-content">
      <div class="main-container">
        <el-main>
          <!-- 上部容器：包含信息和Markdown预览 -->
          <div class="upper-container">
            <div class="info">
              <el-row class="title">{{ articleInfo.title }}</el-row>
              <div class="time">
                <el-text>发布：{{ articleInfo.created_at }} 更新：{{ articleInfo.updated_at }}</el-text>
              </div>
              <el-row class="category">类别：{{ articleInfo.category }}</el-row>
              <el-row class="tags">标签：
                <el-tag v-for="item in articleInfo.tags" :key="item" effect="plain">{{ item }}</el-tag>
              </el-row>
              <div class="abstract">
                <el-text>{{ articleInfo.abstract }}</el-text>
              </div>
            </div>
            <MdPreview :id="mdID" :modelValue="articleInfo.content"/>
          </div>

<!--          &lt;!&ndash; 间距 &ndash;&gt;-->
<!--          <div class="container-spacing"></div>-->

          <!-- 下部容器：包含评论相关内容 -->
          <div class="lower-container">
            <div class="comment" id="comment">
              <el-input v-model="content" :autosize="{ minRows: 4, maxRows: 8 }" type="textarea"
                        placeholder="喜欢这篇文章吗？在这里与大家分享您的想法！" maxlength="320"/>
              <el-text>tip:请登录后再进行反馈!</el-text>
              <div class="operation">
                <div class="comment-tool">
                  <el-popover :visible="layoutStore.state.emojiPopoverVisible"
                              width="502"
                              trigger="click"
                              placement="right"
                  >
                    <template #reference>
                      <el-avatar src="/emoji/s14.png" @click="changeEmojiListState"/>
                    </template>
                    <template #default>
                      <div class="emoji-grid"> <!-- 新增包裹容器 -->
                        <div
                            v-for="emoji in visibleEmojis"
                            :key="emoji.newKey"
                            class="emoji-item"
                            :class="[
                              'emoji',
                              `emoji-sprite-${emoji.spriteGroup}`,
                              `emoji-${emoji.newKey}`
                            ]"
                            :title="`${emoji.oldKey} -> ${emoji.newKey}`"
                            @click="insertEmoji(emoji)"
                        ></div>
                        <div v-if="hasMoreEmojis" @click="loadMoreEmojis" class="load-more-btn">
                          加载更多...
                        </div>
                      </div>
                    </template>
                  </el-popover>
                </div>
                <div class="button-group">
                  <el-button size="large" type="primary" @click="submitComment">发表</el-button>
                  <el-button size="large" @click="content=''">取消</el-button>
                </div>
              </div>
            </div>
            <div class="comment-list">
              <el-row class="title">评论</el-row>
              <comment-item :comments="comments"/>
            </div>
          </div>
        </el-main>
      </div>
      <div class="aside-container">
        <el-aside>
          <div class="aside-content">
            <div class="catalog">
              <el-row class="title">目录</el-row>
              <MdCatalog :editorId="mdID" :scrollElement="scrollElement" :offsetTop="100" :scrollElementOffsetTop="80"/>
            </div>
            <div class="status">
              <el-icon size="24">
                <component is="View"/>
              </el-icon>
              {{ articleInfo.views }}
              <el-icon size="24">
                <component is="ChatDotRound"/>
              </el-icon>
              {{ articleInfo.comments }}
              <el-icon size="24" @click="handelLike">
                <component v-if="!isLike" is="Star"/>
                <component v-else is="StarFilled"/>
              </el-icon>
              {{ articleInfo.likes }}
            </div>
          </div>
        </el-aside>
      </div>
    </el-container>
  </div>
</template>

<script setup lang="ts">
// 此处代码与原代码相同，无需修改
import {useRoute} from "vue-router";
import {type Article, articleInfoByID} from "@/api/article";
import router from "@/router";
import {computed, onMounted, ref, watch} from "vue";
import {MdPreview, MdCatalog} from 'md-editor-v3';
import WebNavbar from "@/components/layout/WebNavbar.vue";
import CommentItem from "@/components/common/CommentItem.vue";
import {articleIsLike, articleLike, type ArticleLikeRequest} from "@/api/article";
import {type Comment, commentCreate, type CommentCreateRequest, commentInfoByArticleID} from "@/api/comment";
import {useLayoutStore} from "@/stores/layout";
import {useUserStore} from "@/stores/user";
import { parseEmojis, renderTextWithEmojisForText, getAllEmojis, type EmojiInfo } from '@/utils/emojiParser'

const mdID = "md-id"

const articleInfo = ref<Article>({
  created_at: '',
  updated_at: '',
  cover: '',
  title: '',
  keyword: '',
  category: '',
  tags: [],
  abstract: '',
  content: '',
  comments: 0,
  views: 0,
  likes: 0,
})

const scrollElement = document.documentElement

const route = useRoute()

const articleID = computed(() => route.params.id)

const layoutStore = useLayoutStore()

// 将内容中的 :emoji: 标记解析并渲染为雪碧图 span
const renderContentWithEmojis = async (text: string): Promise<string> => {
  if (!text) return text
  const parsed = await parseEmojis(text)
  return await renderTextWithEmojisForText(parsed)
}

const openEmojiList = () => {
  layoutStore.show("emojiPopoverVisible") // 关闭表情弹框
  console.log(layoutStore.state.emojiPopoverVisible)
}

const insertEmoji = (emoji: EmojiInfo) => {
  content.value = content.value + `:emoji:${emoji.newKey}:`
  layoutStore.hide("emojiPopoverVisible") // 关闭表情弹框
  console.log(layoutStore.state.emojiPopoverVisible)
}

const changeEmojiListState = () => {
  layoutStore.stateChange("emojiPopoverVisible") // 关闭表情弹框
  console.log(layoutStore.state.emojiPopoverVisible)
}

const getArticleInfo = async () => {
  const res = await articleInfoByID(articleID.value as string)
  if (res.code === 0) {
    const data = res.data
    data.content = await renderContentWithEmojis(data.content || '')
    articleInfo.value = data
  } else {
    await router.push({name: "404"})
  }
}

getArticleInfo()

const isLike = ref(false)

const getIsLikeInfo = async () => {
  const req: ArticleLikeRequest = {
    article_id: articleID.value as string
  }
  const res = await articleIsLike(req)
  if (res.code === 0) {
    isLike.value = res.data
  }
}

if (useUserStore().state.userInfo.role_id !== 0) {
  getIsLikeInfo()
}

const handelLike = async () => {
  const req: ArticleLikeRequest = {
    article_id: articleID.value as string
  }
  const res = await articleLike(req)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    articleInfo.value.likes += isLike.value ? -1 : 1
    isLike.value = !isLike.value
  }
}

const content = ref('')

const emojiList = ref<EmojiInfo[]>([]);
const visibleEmojis = ref<EmojiInfo[]>([]);
const currentPage = ref(0);
const pageSize = 48;

onMounted(async () => {
  try {
    // 使用新的emoji解析工具
    emojiList.value = await getAllEmojis();
    loadMoreEmojis();
  } catch (error) {
    console.error('Error loading emoji list:', error);
    // 降级到旧方式
    const response = await fetch('/emoji/emoji_cache.txt');
    if (response.ok) {
      const text = await response.text();
      const oldEmojiList = text.split('\n').filter(line => line.trim() !== '');
      // 转换为新格式（临时兼容）
      emojiList.value = oldEmojiList.map((filename, index) => ({
        oldKey: filename.replace('.png', ''),
        newKey: `e${index}`,
        spriteGroup: Math.floor(index / 128),
        index: index % 128
      }));
      loadMoreEmojis();
    }
  }
});

const loadMoreEmojis = () => {
  const startIndex = currentPage.value * pageSize;
  const endIndex = Math.min(startIndex + pageSize, emojiList.value.length);
  
  const newEmojis = emojiList.value.slice(startIndex, endIndex);
  visibleEmojis.value.push(...newEmojis);
  
  currentPage.value++;
};

const hasMoreEmojis = computed(() => {
  return visibleEmojis.value.length < emojiList.value.length;
});

const submitComment = async () => {
  const commentCreateRequest: CommentCreateRequest = {
    article_id: articleID.value as string,
    p_id: null,
    content: content.value,
  }
  const res = await commentCreate(commentCreateRequest)
  if (res.code === 0) {
    ElMessage.success(res.msg)
    content.value = ''
    layoutStore.state.shouldRefreshCommentList = true
  }
}

const comments = ref<Comment[]>([])

// 递归处理评论及其子评论中的 emoji 标记
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

const getArticleCommentsInfo = async () => {
  comments.value = []
  const res = await commentInfoByArticleID(articleID.value as string)
  if (res.code === 0) {
    comments.value = await transformCommentsWithEmojis(res.data)
  }
}

onMounted(() => {
  getArticleCommentsInfo()
})

watch(() => layoutStore.state.shouldRefreshCommentList, (newVal) => {
  if (newVal) {
    getArticleCommentsInfo()
    layoutStore.state.shouldRefreshCommentList = false
  }
})
</script>

<style scoped>
/* Emoji选择器样式 */
.emoji-grid {
  display: grid;
  grid-template-columns: repeat(12, 1fr);
  gap: 4px;
  padding: 8px;
  max-height: 360px;
  overflow-y: auto;
}

.emoji-item {
  width: 32px !important;
  height: 32px !important;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
  /* background-size继承全局CSS */
}

.emoji-item:hover {
  background-color: #f0f9ff;
  transform: scale(1.2);
}

.load-more-btn {
  grid-column: 1 / -1;
  text-align: center;
  padding: 8px;
  cursor: pointer;
  color: #409eff;
  border: 1px dashed #409eff;
  border-radius: 4px;
  margin-top: 8px;
}

.load-more-btn:hover {
  background-color: #f0f9ff;
}

/* Emoji样式已全局导入，这里只保留组件特定样式 */
</style>

<style scoped lang="scss">
.article {
    width: 1600px;
    display: flex;
    justify-content: center;
    align-items: center;
    margin: 0 auto;

  .main-content {
    display: flex;
    justify-content: center;
    gap: 10px;

    .main-container {
      width: 70%;
      padding: 0px;
      //background-color: white;

      .el-main{
        padding: 0px;
        display: flex;
        flex-direction: column;
        gap: 10px; /* 设置容器间的间距 */

        /* 上部容器样式 */
        .upper-container {
          padding: 20px;
          background-color: white;

          .info {
            border: 1px solid #DCDFE6;
            padding: 10px;

            .title {
              font-size: 24px;
              margin-bottom: 10px;
            }

            .time {
              margin-bottom: 10px;
            }

            .category {
              margin-bottom: 10px;
            }

            .tags {
              margin-bottom: 10px;
            }

            .abstract {
              margin-bottom: 10px;
            }
          }
        }

        /* 下部容器样式 */
        .lower-container {
          padding: 20px;
          background-color: white;

          .comment {
            border-top: 1px solid #DCDFE6;
            padding-top: 20px;

            .operation {
              margin-top: 20px;
              margin-bottom: 20px;
              display: flex;

              .comment-tool {
                margin-right: auto;

                .el-avatar {
                  background-color: unset;
                }
              }

              .button-group {
                margin-left: auto;
              }
            }
          }

          .comment-list {
            .title {
              font-size: 24px;
              margin-bottom: 10px;
            }
          }
        }
      }
    }

    .aside-container {
      width: 20%;
      padding:0px;
      //background-color: white;

      .el-aside{
        padding: 0px;
        display: flex;
        flex-direction: column;
        gap: 10px; /* 设置容器间的间距 */

        .aside-content {
          position: fixed;
          padding: 20px;
          background-color: white;

          .catalog {
            width: 100%;
            height: 50vh;
            overflow: auto;
            padding: 10px;
            border: 1px solid #DCDFE6;

            .title {
              font-size: 24px;
              margin-bottom: 10px;
            }
          }

          .status {
            justify-content: center;
            display: flex;
            padding: 20px;
            border-left: 1px solid #DCDFE6;
            border-right: 1px solid #DCDFE6;
            border-bottom: 1px solid #DCDFE6;

            .el-icon {
              margin-left: 20px;
              margin-right: 20px;
            }
          }
        }
      }


    }
  }
}

.el-popover {
  .el-image {
    height: 50px;
    width: 50px;
  }

  > div {
    max-height: 300px;
    overflow-y: auto;

    scrollbar-width: none;
    -ms-overflow-style: none;

    &::-webkit-scrollbar {
      display: none;
    }
  }
}
</style>