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
import Card from "@/components/common/card/OneCard";
import OneCard from "@/components/common/card/OneCard";
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
      dataList: []
    };
  },
  created() {
    let userInfo = JSON.parse(localStorage.getItem("userInfo"));
    console.log(userInfo);
    axios
      .get("http://172.20.10.6:8081/getMyIndexInfo", {
        params: { userId: userInfo.userId },
      })
      .then((res) => {
        console.log(res);
        this.dataList = res.data.myVideos;
      });
  },

  methods: {
    jumpPath() { },
  },
};
</script>

<style scoped>
.s_r_c_item {
  position: relative;
  width: 23%;
  height: 48%;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 20px;
  margin-left: 5px;
  margin-right: 14px;
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
