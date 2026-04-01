<template>
  <div style="width:100%;height:100%">
    <HeaderFix></HeaderFix>
    <Partition></Partition>
  </div>
</template>

<script>
import Partition from '@/components/body/partition.vue';
import HeaderFix from "@/components/header/HeaderFix.vue";
  export default {
    props: {

    },
    components: {
      HeaderFix,
      Partition
    },
    data(){
      return {
        
      }
    },
    methods: {}
  }
</script>

<style scoped>

</style>