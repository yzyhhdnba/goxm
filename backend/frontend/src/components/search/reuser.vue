<template>
  <div class="swiper_right_content">
    <div v-for="(item, i) in dataList" :key="i" class="s_r_c_item">
      <div @click="jumpPath(item.title)">
        <one-card :item="item" />
      </div>
    </div>
  </div>
</template>

<script>
import Card from "@/components/common/card/Card";
import OneCard from "@/components/common/card/Usercard";
import { onMounted, reactive, ref } from "vue";
import axios from "axios";
export default {
  props: {},
  components: {
    Card,
    OneCard,
  },
  data() {
    return {
      dataList: [
        
      ],
    };
  },
  created() {
    axios.get('http://172.20.10.6:8081/search/user',{
      params: {
        key:this.$route.query.text,
        pageNo:1
      }
    }).then(res=>{
       console.log(res);
       this.dataList=res.data;
    })
  },

  methods: {
    jumpPath() {},
  },
};
</script>

<style scoped>
.s_r_c_item {
  position: relative;
  width: 45%;
  height: 48%;
  border-radius: 6px;
  overflow: hidden;
  margin-top: 20px;
  margin-bottom: 25px;
  margin-left: 20px;
}

.swiper_right_content {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  /* justify-content: space-between; */
  align-content: space-between;
  margin-left: 1vw;
}
</style>
