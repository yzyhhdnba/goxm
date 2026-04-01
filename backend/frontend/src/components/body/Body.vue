<template>
  <div class="body">
    <!-- 导引栏内容 -->
    <div class="b-wrap first-content">
      <div class="space-between">
        <!-- 轮播图  -->
        <el-carousel height="300px" indicator-position="none" trigger="click">
          <el-carousel-item v-for="(item, i) in swiperList" :key="i">
            <div @click="jumpPath(item.title)">
              <img class="swiper_img" :src="item.pic" />
              <div class="swiper_title">
                <div>{{ item.title }}</div>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
        <!-- 轮播图右推荐列表内容 -->
        <div class="swiper_right_content" v-if="swiperList_r && swiperList_rFlag">
          <div v-for="(item, i) in swiperList_r" :key="i" class="s_r_c_item">
            <div @click="jumpPath(item.title)">
              <one-card :item="item" />
            </div>
          </div>

        </div>
        <div class="swiper_right_content" v-if="!swiperList_rFlag">
          <div class="right-skeleton" v-for="(item, index) in ['', '', '', '', '', '']" :key="index">
            <el-skeleton :rows="3" animated />
          </div>
        </div>
      </div>
    </div>

    <div class="b-wrap">
      <!-- 推荐 -->
      <el-skeleton :rows="5" v-if="!animationListFlag" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag" animated />
      <div class="various-content" v-if="animationList && animationListFlag">

        <!-- 卡片区域 -->
        <card class="various-c-card" :maindata="animationList" :mdname="'推荐'"></card>

      </div>
      <!-- 番剧 -->
      <div class="various-content" v-if="animationList && animationListFlag">
        <!-- 卡片区域 -->
        <card class="various-c-card" :maindata="animationList" :mdname="'番剧'"></card>

      </div>
    </div>

  </div>
  <el-footer class="main-footer">
    <Foot />
  </el-footer>
</template>

<script>
import Card from "@/components/common/card/Card";
import OneCard from "@/components/common/card/OneCard";
import { onMounted, reactive, ref } from 'vue'
import Foot from "@/components/foot/index";
import axios from 'axios'
export default {
  props: {

  },
  components: {
    Card,
    OneCard,
    Foot,
  },
  data() {
    return {
      swiperList: [{ title: '苏渊默', pic: 'http://101.35.142.191:8081/tiao.jpg' },
      { title: '开发者名单', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg' },
      { title: 'icy', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917531190280.jpg' },
      { title: 'hhd', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917534678902.jpg' },
      { title: '苏渊默', pic: 'https://uploadstatic.mihoyo.com/contentweb/20200329/2020032917002825368.jpg' }],
      swiperList_r: [],
      animationList: [],
      swiperList_rFlag: false,
      animationListFlag: false
    }
  },
  created() {
    axios.post('http://172.20.10.6:8081/videoList', {
      areaId: 1,
      count: 6
    }).then(res => {
      console.log(res);
      this.swiperList_r = res.data;
      this.swiperList_rFlag = true;
    }).catch(err => {
      console.log(err);
    })
    axios.post('http://172.20.10.6:8081/videoList', {
      areaId: 1,
      count: 8
    }).then(res => {
      console.log(res);
      this.animationList = res.data;
      this.animationListFlag = true
    }).catch(err => {
      console.log(err);
    })

  },
  methods: {

  }
}
</script>

<style scoped lang="less">
.body {
  margin-top: 20px;
}

.b-wrap {
  width: 90vw;
  margin: 0 auto;
  /* 最小宽度 */
  min-width: 1320px;
}

/* 第一内容 */
.first-content {
  height: 450px;
  min-height: 220px;
}


.space-between {
  display: flex;
  justify-content: space-between;
  height: 100%;
}

/* 轮播图 */
.el-carousel {
  width: 550px;
  height: 380px !important;
  border-radius: 8px;

}

.el-carousel__container {
  height: 380px !important;
}

.el-carousel__item {
  /* position: relative; */
  width: 550px;
  height: 380px !important;
}

.swiper_img {
  display: inline-block;
  width: 100%;
  height: 380px;

}

.swiper_title {
  display: flex;
  width: 100%;
  height: 60px;
  position: absolute;
  left: 0px;
  bottom: 0px;
  padding-left: 10px;
  color: #fff;
  font-size: 17 px;
  font-weight: 400;
  background-color: rgba(123, 122, 122, 0.7);
  align-items: center;
  /* 一行显示 */
  overflow: hidden;

}

/* 轮播图右边内容 */
.swiper_right_content {
  flex: 1;
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-content: space-between;
  margin-left: 1vw;

  .right-skeleton {
    width: 32%;
    height: 48%;
  }
}

.s_r_c_item {
  position: relative;
  width: 32%;
  height: 48%;
  border-radius: 6px;
  overflow: hidden;
  margin-bottom: 20px;

}

.s_r_c_item img {
  width: 100%;
  display: block;
}

.s_r_c_title {
  height: 10%;
  margin-top: 5px;
  left: 15px;
  font-size: 14px;
  /* 两行显示 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  color: #4d4d4d;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.s_r_c_author {
  padding-left: 5px;
}

.s_r_c_author a {
  font-size: 12px;
  color: #999;
}


.various-content {
  display: flex;
  justify-content: space-around;
  margin-top: 30px;
}

.various-c-card {
  width: 100%;
}
</style>