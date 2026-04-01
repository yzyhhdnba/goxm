<template>
  <div style="width:100%;height:100%">
    <HeaderFix></HeaderFix>
    <VideoDetail></VideoDetail>
  </div>
</template>

<script>
import VideoDetail from '@/components/detail/videoDetail.vue';
import HeaderFix from "@/components/header/HeaderFix.vue";
  export default {
    props: {

    },
    components: {
      VideoDetail,
      HeaderFix,
    },
    setup(){
      return {
        
      }
    },
    methods: {

    }
  }
</script>

<style scoped>

</style>