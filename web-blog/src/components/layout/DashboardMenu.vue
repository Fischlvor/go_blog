<template>
  <div class="dashboard-menu">
        <div class="collapse-button">
          <el-button @click="handleCollapse">
            <el-icon>
              <component :is="icon"></component>
            </el-icon>
          </el-button>
        </div>
        <el-menu :collapse="isCollapse" :collapse-transition="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item v-if="!item.subItems" :index="generatePathForSingleItem(item)">
              <el-icon>
                <component :is="item.icon"></component>
              </el-icon>
              <span>{{ item.title }}</span>
            </el-menu-item>
            <el-sub-menu v-else-if="!item.admin_role || userStore.isAdmin" :index="item.name">
              <template #title>
                <el-icon>
                  <component :is="item.icon"></component>
                </el-icon>
                <span>{{ item.title }}</span>
              </template>
              <el-menu-item v-for="subItem in item.subItems" :index="generatePathForSubItem(item,subItem)"
                            @click="handleClick(subItem)">
                <el-icon>
                  <component :is="subItem.icon"></component>
                </el-icon>
                <span>{{ subItem.title }}</span>
              </el-menu-item>
            </el-sub-menu>
          </template>
        </el-menu>
  </div>
</template>

<script setup lang="ts">
import type {Tag} from "@/stores/tag";
import {computed} from "vue";
import {useLayoutStore} from "@/stores/layout";
import {useTagStore} from "@/stores/tag";
import {useUserStore} from "@/stores/user";

interface MenuItem {
  title: string;
  name: string;
  icon: string;
  subItems?: MenuItem[];
  admin_role?:boolean;
}

const menuList: MenuItem[] = [
  {
    title: "主页",
    name: "",
    icon: "House"
  },
  {
    title: "个人中心",
    name: "user-center",
    icon: "Monitor",
    subItems: [
      {
        title: "我的信息",
        name: "user-info",
        icon: "Postcard"
      },
      {
        title: "我的收藏",
        name: "user-star",
        icon: "Star"
      },
      {
        title: "我的评论",
        name: "user-comment",
        icon: "ChatDotRound"
      },
      {
        title: "我的反馈",
        name: "user-feedback",
        icon: "Message"
      }
    ]
  },
  {
    title: "用户管理",
    name: "users",
    icon: "User",
    admin_role:true,
    subItems: [
      {
        title: "用户列表",
        name: "user-list",
        icon: "SetUp"
      }
    ]
  },
  {
    title: "文章管理",
    name: "articles",
    icon: "Document",
    admin_role:true,
    subItems: [
      {
        title: "发布文章",
        name: "article-publish",
        icon: "Collection"
      },
      {
        title: "评论列表",
        name: "comment-list",
        icon: "ChatLineRound"
      },
      {
        title: "文章列表",
        name: "article-list",
        icon: "DocumentCopy"
      }
    ]
  },
  {
    title: "图片管理",
    name: "images",
    icon: "Picture",
    admin_role:true,
    subItems: [
      {
        title: "图片列表",
        name: "image-list",
        icon: "PictureRounded"
      }
    ]
  },
  {
    title: "AI对话管理",
    name: "ai-management",
    icon: "ChatDotRound",
    admin_role:true,
    subItems: [
      {
        title: "模型管理",
        name: "models",
        icon: "Setting"
      },
      {
        title: "会话管理",
        name: "sessions",
        icon: "ChatLineRound"
      },
      {
        title: "消息管理",
        name: "messages",
        icon: "Message"
      }
    ]
  },
  {
    title: "表情管理",
    name: "emoji",
    icon: "Sunny",
    admin_role:true,
    subItems: [
      {
        title: "表情列表",
        name: "emoji-list",
        icon: "Picture"
      },
      {
        title: "表情组管理",
        name: "emoji-groups",
        icon: "Collection"
      },
      {
        title: "雪碧图管理",
        name: "emoji-sprites",
        icon: "Grid"
      }
    ]
  },
  {
    title: "系统管理",
    name: "system",
    icon: "Coin",
    admin_role:true,
    subItems: [
      {
        title: "反馈列表",
        name: "feedback-list",
        icon: "Position"
      },
      {
        title: "广告列表",
        name: "advertisement-list",
        icon: "Connection"
      },
      {
        title: "友链列表",
        name: "friend-link-list",
        icon: "Link"
      },
      {
        title: "登录日志",
        name: "login-logs",
        icon: "memo"
      },
      {
        title: "应用配置",
        name: "app-config",
        icon: "setting"
      }
    ]
  }
]

const layoutStore = useLayoutStore()
const isCollapse = computed(() => layoutStore.state.isCollapse)
const tagStore = useTagStore()
const icon = computed(() => layoutStore.state.isCollapse ? "Expand" : "Fold")
const userStore = useUserStore()

const handleCollapse = () => {
  layoutStore.stateChange("isCollapse")
}

function generatePathForSingleItem(item: MenuItem): string {
  return '/dashboard/' + item.name;
}

function generatePathForSubItem(parentItem: MenuItem, subItem: MenuItem): string {
  return '/dashboard/' + parentItem.name + '/' + subItem.name;
}

function handleClick(subItem: MenuItem) {
  const newTag: Tag = {
    title: subItem.title,
    name: subItem.name
  }
  const exists = tagStore.state.tags.some(tag => tag.name === newTag.name);
  if (exists) {
    return;
  }
  tagStore.addTag(newTag);
}

</script>

<style scoped lang="scss">
.dashboard-menu {
  background: rgba(255, 255, 255, 0.0); // 或使用 transparent
  width: v-bind('isCollapse ? "64px" : "240px"'); // 关键修改
  transition: width 0.3s ease; // 添加过渡效果

  .collapse-button {
    display: flex;
    justify-content: center;  // Horizontally center the button
    align-items: center;      // Vertically center the button
    background: transparent;
    padding: 10px;

    .el-button {
      margin-left: auto;
      border: none;
      background: transparent;
      color: #fff;
      font-size: 18px;
      transition: all 0.3s ease;

      &:hover {
        color: var(--el-color-primary);
        background: rgba(255, 255, 255, 0.1);
      }
    }
  }

  .el-menu {
    border-right: none;
    --el-menu-item-height: 52px;
    --el-menu-sub-item-height: 44px;
    background: transparent;
    transition: all 0.3s ease;

    :deep(.el-menu) {
      background: transparent !important;
    }

    // 菜单项通用样式
    .el-menu-item, :deep(.el-sub-menu__title) {
      color: #ffffff;
      font-size: 20px;
      font-weight: bold;
      transition: all 0.3s ease;

      &:hover {
        color: #fff;
        background: rgba(255, 255, 255, 0.20) !important;
      }

      &.is-active {
        color: gold;
        font-weight: bold;
        //background: rgba(255, 255, 255, 0.20) !important;

        &::after {
          content: '';
          position: absolute;
          right: 0;
          top: 0;
          bottom: 0;
          width: 8px;
          background: var(--el-color-primary);
        }
      }
    }

    .el-sub-menu{
      background: rgba(255, 255, 0, 0.0);

      .el-menu-item{
        background: rgba(255, 255, 0, 0.0);
        font-size: 16px;
      }

      &.is-opened {

        .el-menu-item{
          background: rgba(255, 255, 255, 0.1);
        }
      }
    }
  }
}
</style>

