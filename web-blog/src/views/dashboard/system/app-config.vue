<template>
  <div class="system-config">
    <el-container>
      <el-aside width="160px">
        <el-menu :collapse-transition="false" :router="true" :default-active="$route.path">
          <template v-for="item in menuList">
            <el-menu-item :index="generatePathForSingleItem(item)">
              <span>{{ item.title }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view/>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
interface MenuItem {
  title: string;
  name: string;
}

const menuList: MenuItem[] = [
  {
    title: "网站配置",
    name: "site-config",
  },
  {
    title: "系统配置",
    name: "system-config",
  },
  {
    title: "邮箱配置",
    name: "email-config",
  },
  {
    title: "QQ登录配置",
    name: "qq-config",
  },
  {

    title: "七牛云配置",
    name: "qiniu-config",
  },
  {
    title: "jwt配置",
    name: "jwt-config",
  },
  {
    title: "高德配置",
    name: "gaode-config",
  }
]

function generatePathForSingleItem(item: MenuItem): string {
  return '/dashboard/system/app-config/' + item.name;
}
</script>

<style scoped lang="scss">
.system-config{
  background-color: rgba(255,255,255,0.7);
  padding: 10px;
  height: calc(100vh - 158px); /* 计算高度：视口高度 - 头部80px - 顶部间距100px（需根据实际调整） */


  .el-container{
    height: 100%;

    .el-aside{

      .el-menu{
        background-color: rgba(255,255,255,0.0);
        //height: 100vh;
      }
    }

    .el-main{
      //padding-bottom: 10px;
      margin-top: 0px; /* 为main-top留出空间 */
      overflow-y: auto; /* 垂直方向滚动 */
      scrollbar-width: none; /* 隐藏滚动条（可选） */
      padding: 5px;
      //height: 100vh;
    }
  }




}

</style>
