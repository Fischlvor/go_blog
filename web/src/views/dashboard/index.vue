<template>
  <div class="dashboard">
    <el-container>
      <el-header>
          <logo/>
          <auth-popover/>
      </el-header>
      <el-container class="head-container">
        <el-aside :class="{collapsed: isCollapsed}">
          <div class="aside-bottom">
            <dashboard-menu/>
          </div>
        </el-aside>
        <el-main>
          <div class="main-top">
            <div class="breadcrumb-bg">
              <breadcrumb/>
            </div>
            <div class="dashboard-tag-bg">
              <dashboard-tag/>
            </div>
          </div>
          <div class="main-bottom">
            <div class="route-view-bg">
              <router-view/>
            </div>
          </div>

        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import Logo from "@/components/widgets/Logo.vue";
import Breadcrumb from "@/components/layout/DashboardBreadcrumb.vue";
import DashboardMenu from "@/components/layout/DashboardMenu.vue";
import {computed} from "vue";
import DashboardTag from "@/components/layout/DashboardTag.vue";
import AuthPopover from "@/components/common/AuthPopover.vue";
import {useLayoutStore} from "@/stores/layout";

const store = useLayoutStore()

const isCollapsed = computed(() => store.state.isCollapse)

</script>

<style scoped lang="scss">
.dashboard {
  display: flex;
  background-color: rgba(255, 255, 255, 0.1);

  .el-header{
    width: 100%;                /* 撑满父容器宽度 */
    height: 80px;
    display: flex;              /* 子容器再次开启Flex布局 */
    justify-content: space-between; /* 左右元素两端对齐 */
    padding: 0px;

    .logo{

    }

    .auth-popover{
      //display: flex; /* 开启 Flex 布局 */
      align-items: center; /* 子元素垂直居中 */
      padding: 20px;
    }
  }


  .el-aside {
    width: 240px;
    height: calc(100vh - 80px);
    display: flex;
    flex-direction: column;
    overflow: auto; /* 允许滚动但隐藏滚动条 */

    /* 隐藏滚动条但保留功能 */
    scrollbar-width: none; /* Firefox */
    -ms-overflow-style: none; /* IE and Edge */

    ::-webkit-scrollbar {
      display: none;
    }

    .aside-bottom{
      flex: 1; /* 自动填充剩余空间 */
      min-height: 0; /* 允许 flex 子项收缩到最小内容高度以下 */
      overflow-y: auto; /* 如果内容超出，显示滚动条 */
      background-color: rgba(255, 255, 255, 0.1);

      .dashboard-menu{
        //border-right: 1px solid #DCDFE6;
        flex-shrink: 0; /* 防止菜单被压缩 */
      }
    }
  }

  .el-main{
    background-color: rgba(255,255,255,0.1);
    //width: 70vh;
    height: calc(100vh - 80px); /* 保持整体高度为视口高度 */
    padding: 10px;


    .main-top{
      position: sticky; /* 粘性定位，固定在顶部 */
      top: 0px; /* 与头部高度一致，避免遮挡 */
      z-index: 5; /* 确保在滚动内容上方 */
      padding: 5px;
      //background-color: rgba(255,255,255,0.8);

      .breadcrumb-bg{
        padding-bottom: 5px;
        border-bottom: 1px solid #DCDFE6
      }
    }

    .main-bottom{
      height: calc(100vh - 168px); /* 计算高度：视口高度 - 头部80px - 顶部间距100px（需根据实际调整） */
      margin-top: 0px; /* 为main-top留出空间 */
      overflow-y: auto; /* 垂直方向滚动 */
      scrollbar-width: none; /* 隐藏滚动条（可选） */
      padding: 5px;


      &::-webkit-scrollbar {
        display: none; /* 隐藏WebKit滚动条 */
      }

      .route-view-bg {
        height: 100%; /* 撑满容器高度 */
      }
    }


  }
}
</style>
