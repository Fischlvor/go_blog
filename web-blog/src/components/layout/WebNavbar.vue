<template>
  <div :class="{'web-navbar': true,show: isShow}">
    <div class="container">
      <logo/>
      <div class="web-menu">
        <el-menu mode="horizontal" :ellipsis="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="item.name"><span>{{ item.title }}</span></el-menu-item>
          </template>
        </el-menu>
      </div>
      <music-player/>
      <auth-popover/>
    </div>
  </div>
</template>

<script setup lang="ts">
import AuthPopover from "@/components/common/AuthPopover.vue";
import Logo from "@/components/widgets/Logo.vue";
import {ref} from "vue";
import {onUnmounted} from "vue";


const isShow = ref(true)

const props = defineProps<{
  noScroll?:boolean
}>()

if (!props.noScroll){
  isShow.value = false
  window.addEventListener("scroll", scroll)
  scroll()
}

function scroll() {
  let top = document.documentElement.scrollTop
  isShow.value = top >= 100;
}

onUnmounted(() => {
  if (!props.noScroll) {
    window.removeEventListener("scroll", scroll)
  }
})

interface MenuItem {
  title: string;
  name: string;
}

const menuList: MenuItem[] = [
  {
    title: "首页",
    name: "/",
  },
  {
    title: "搜索",
    name: "/search",
  },
  {
    title: "新闻",
    name: "/news",
  },
  {
    title: "AI助手",
    name: "/ai-assistant",
  },
  {
    title: "友链",
    name: "/friend-link",
  },
  {
    title: "关于",
    name: "/about",
  }
]

</script>

<style scoped lang="scss">
.web-navbar {
  display: flex;
  justify-content: center;
  width: 100%;
  height: 100px;
  position: relative;
  z-index: 6;
  background-color: rgba(128, 128, 128, 0.4); /* 灰色，50%透明度 */


  //&.show {
  //  top: 0;
  //  background-color: white;
  //  color: black;
  //  --el-menu-text-color: black;
  //  --color: black;
  //
  //  .container {
  //    //margin-top: 8px;
  //  }
  //}

  .container {
    display: flex;
    max-width: 1800px;
    width: 100%;
    justify-content: space-between; /* 关键修改 */
    align-items: center; /* 垂直居中 */

    .logo {
      margin-left: 10px;
      height: 80px;
      /* 移除多余的 margin-top/bottom: auto */
    }

    .web-menu {
      margin-left: auto;
      margin-right: auto;

      .el-menu {
        background-color: transparent;
        border-bottom: none;
        --el-menu-item-font-size: 20px;

        .el-menu-item {
          color: white;
          border-bottom: none;
          background-color: transparent;
        }
      }
    }

    .music-player {
      margin-left: auto; /* 自动左外边距 */
      margin-right: auto; /* 自动右外边距 */
    }

    .auth-popover {
      padding-right: 10px;
    }
  }
}
</style>
