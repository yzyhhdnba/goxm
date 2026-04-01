<script>
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup

export default {
  mounted () {
    // 检测浏览器路由改变页面不刷新问题,hash模式的工作原理是hashchange事件
    // window.addEventListener('hashchange', () => {
    //   console.log(window.location.hash,'window.location.hash')
    //   let currentPath = window.location.hash.slice(1)
    //   console.log(currentPath,'currentPath')
    //   if (this.$route.path !== currentPath) {
    //     this.$router.push(currentPath)
    //   }
    // }, false)
  },

  components: {},
};
</script>

<template>
  <router-view></router-view>
</template>

<style>
* {
  padding: 0;
  margin: 0;
  font-family: PingFang SC, HarmonyOS_Regular, Helvetica Neue, Microsoft YaHei,
    sans-serif !important;
}
li {
  list-style: none;
}
input {
  border: none;
  outline: none;
}
input:focus {
  border: none;
}
</style>
