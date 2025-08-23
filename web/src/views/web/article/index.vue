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
                      <div> <!-- 新增包裹容器 -->
                        <el-image
                            v-for="emoji in emojiList"
                            :src="'/emoji/'+emoji"
                            @click="insertEmoji(emoji)"
                        />
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

const openEmojiList = () => {
  layoutStore.show("emojiPopoverVisible") // 关闭表情弹框
  console.log(layoutStore.state.emojiPopoverVisible)
}

const insertEmoji = (emoji: string) => {
  content.value = content.value + '![](' + '/emoji/' + emoji + ')'
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
    articleInfo.value = res.data
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

const emojiList = ref<string[]>([]);
onMounted(async () => {
  try {
    const response = await fetch('/emoji/emoji_cache.txt');
    if (!response.ok) throw new Error('Failed to fetch emoji list');

    const text = await response.text();
    emojiList.value = text.split('\n').filter(line => line.trim() !== '');
  } catch (error) {
    console.error('Error loading emoji cache:', error);
    // Fallback to default numbers if needed
    //emojiList.value = Array.from({ length: 349 }, (_, i) => `s${i}.png`);
  }
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

const getArticleCommentsInfo = async () => {
  comments.value = []
  const res = await commentInfoByArticleID(articleID.value as string)
  if (res.code === 0) {
    comments.value = res.data
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