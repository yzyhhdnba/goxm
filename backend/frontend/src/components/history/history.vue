<template>
  <div class="history">
    <div class="b-head-c" _v-0a0880cb="">
      <i class="b-icon b-icon-history" _v-0a0880cb=""></i>
      <span class="b-head-t" _v-0a0880cb="">历史记录</span>
    </div>
    <div class="swiper_right_content" v-if="flag">
      <div v-for="(item, i) in dataList" :key="i" class="s_r_c_item">
        <div @click="jumpPath(item.title)">
          <one-card :item="item" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Card from "@/components/common/card/Card";
import OneCard from "@/components/common/card/Historycard";
import axios from "axios";
export default {
  created() {
    const { userId } = JSON.parse(localStorage.getItem('userInfo'))
    axios.get('http://172.20.10.6:8081/histories', {
      params: {
        page: 1,
        userId
      }
    }).then(res => {
      console.log(res);
      this.dataList=res.data;
      this.flag=true;
    })
  },
  props: {},
  components: {
    Card,
    OneCard,
  },
  data() {
    return {
      flag:false,
      dataList: [],
    };
  },
  methods: {
    jumpPath() { },
  },
};
</script>

<style scoped>
.history {
  width: 1200px;
  margin: auto;
  padding-top: 70px;
}

.b-head-c {
  font: 12px Helvetica Neue, Helvetica, Arial, Microsoft Yahei, Hiragino Sans GB,
    Heiti SC, WenQuanYi Micro Hei, sans-serif;
  font-size: 12px;
  margin: 0;
  padding: 0;
  float: left;
}

.b-icon,
.b-icon-history {
  font: 12px Helvetica Neue, Helvetica, Arial, Microsoft Yahei, Hiragino Sans GB,
    Heiti SC, WenQuanYi Micro Hei, sans-serif;
  font-size: 12px;
  margin: 0;
  padding: 0;
  font-style: normal;
  font-weight: 400;
  background: url(//s1.hdslb.com/bfs/static/history-record/./img/icons.png) -83px -850px no-repeat;
  width: 28px;
  height: 28px;
  display: inline-block;
  margin-right: 8px;
  vertical-align: middle;
}

.b-head-t {
  font: 12px Helvetica Neue, Helvetica, Arial, Microsoft Yahei, Hiragino Sans GB,
    Heiti SC, WenQuanYi Micro Hei, sans-serif;
  margin: 0;
  padding: 0;
  vertical-align: middle;
  display: inline-block;
  font-size: 18px;
  line-height: 24px;
  color: #222;
}

.s_r_c_item {
  position: relative;
  width: 98%;
  height: 48%;
  border-radius: 6px;
  overflow: hidden;
  margin-top: 20px;
  margin-bottom: 25px;
  margin-left: 20px;
}

.swiper_right_content {
  margin-top: 30px;
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  /* justify-content: space-between; */
  align-content: space-between;
  margin-left: 1vw;
}
</style>
