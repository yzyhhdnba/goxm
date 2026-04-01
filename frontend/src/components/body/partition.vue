<template>
  <div style="height: 64px"></div>
  <div class="partition">
    <div class="b-wrap" style="font-size:24px;line-height:34px;color:#18191C;height:64px">
      <span>番剧</span>
    </div>
    <div class="b-wrap">
      <div class="top">
        <el-carousel height="300px" trigger="click" type="card">
          <el-carousel-item v-for="(item, i) in swiperList" :key="i">
            <div @click="jumpPath(item.title)">
              <img class="swiper_img" :src="item.pic" />
              <div class="swiper_title">
                <div>{{ item.title }}</div>
              </div>
            </div>
          </el-carousel-item>
        </el-carousel>
      </div>
    </div>
    <div class="b-wrap">
      <el-skeleton :rows="5" v-if="!animationListFlag1" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag1" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag1" animated style='margin-bottom: 50px;' />
      <el-skeleton :rows="5" v-if="!animationListFlag1" animated />
      <div>
        <div class="various-content" v-if="animationList1 && animationListFlag1">
          <card class="various-c-card" :maindata="animationList1" :mdname="'前方高能'"></card>
        </div>
        <div class="various-content" v-if="animationList2 && animationListFlag2">
          <card class="various-c-card" :maindata="animationList2" :mdname="'鬼畜推荐'"></card>
        </div>
        <div class="various-content" v-if="animationList3 && animationListFlag3">
          <card class="various-c-card" :maindata="animationList3" :mdname="'你所热爱的，就是你的生活'"></card>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import Card from "@/components/common/card/Card";
import { VideoAPI } from "@/api/index";

export default {
  async created() {
    try {
      const res1 = await VideoAPI.getRecommend({ limit: 6 });
      if (res1 && res1.items) {
          this.animationList1 = res1.items;
      }
      this.animationListFlag1 = true;
      
      const res2 = await VideoAPI.getRecommend({ limit: 6, cursor: Math.random().toString() }); // Or similar to offset
      if (res2 && res2.items) {
          this.animationList2 = res2.items;
      }
      this.animationListFlag2 = true;

      const res3 = await VideoAPI.getRecommend({ limit: 6, cursor: Math.random().toString() });
      if (res3 && res3.items) {
          this.animationList3 = res3.items;
      }
      this.animationListFlag3 = true;
    } catch (err) {
      console.error(err);
      this.animationListFlag1 = true;
      this.animationListFlag2 = true;
      this.animationListFlag3 = true;
    }
  },
  name: "Partition",
  components: {
    Card,
  },
  data() {
    return {
      swiperList: [
        { title: '精选强档', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg' },
        { title: '番剧出击', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917531190280.jpg' },
        { title: '前方高能', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917534678902.jpg' }
      ],
      animationList1: [],
      animationListFlag1: false,
      animationList2: [],
      animationListFlag2: false,
      animationList3: [],
      animationListFlag3: false,
    }
  },
  methods: {
    jumpPath(title) {
      console.log('Jumping to', title);
    }
  }
}
</script>

<style scoped lang="less">
.partition {
  margin-top: 20px;
}

.b-wrap {
  width: 90vw;
  margin: 0 auto;
  /* 最小宽度 */
  min-width: 1320px;
}

.top {
  display: flex;
  justify-content: center;
  height: 100%;
}

.el-carousel {
  width: 100%;
  height: 100% !important;
  border-radius: 8px;
}

.el-carousel__container {
  height: 100% !important;
}

.el-carousel__item {
  /* position: relative; */
  width: 700px;
  height: 100% !important;
  box-shadow: 0px 0px 5px rgba(0, 0, 0, 0.3);
  border-radius: 8px;
}

.swiper_img {
  display: inline-block;
  width: 100%;
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

.various-content {
  display: flex;
  justify-content: space-around;
  margin-top: 30px;
}

.various-c-card {
  width: 100%;
}
</style>