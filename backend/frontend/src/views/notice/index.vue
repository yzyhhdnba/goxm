<template>
  <div style="width:100%;height:100%">
    <HeaderFix></HeaderFix>
    <Notice></Notice>
  </div>
</template>

<script>
import Notice from '@/components/notice/notice.vue';
import HeaderFix from "@/components/header/HeaderFix.vue";
export default {
  components:{
    Notice,
    HeaderFix,
  },
  beforeCreate() {
    if(!localStorage.getItem('isLogin')){
      alert('您还未登录！')
      this.$router.replace('/')
    }
  }
}
</script>

<style>

</style>