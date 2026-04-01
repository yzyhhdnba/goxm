<template>
  <div class="swiper_right_content" v-if="flag">
    {{}}
    <div v-for="(item, i) in dataList.list" :key="i" class="s_r_c_item">
      <div @click="jumpPath(item.title)">
        <one-card :item="item" />
      </div>
    </div>
  </div>
</template>

<script>
import Card from "@/components/common/card/Card";
import OneCard from "@/components/common/card/OneCard.vue";
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
      dataList: {
        list: [
          {
            videoId: 4,
            videoTitle: "哥们名叫丁真",
            userId: 1,
            userName: "yanshsf",
            videoContribute: "2022-06-26T12:52:43.000+00:00",
            videoApprove: "2022-06-27T16:42:11.000+00:00",
            videoComments: 0,
            videoHits: 0,
            videoLikes: 0,
            videoDuration: "09:99",
            videoContent: "L!T!C!",
            areaId: 1,
            denyReason: null,
          },
        ],
      },
      flag: false,
    };
  },
  created() {
    console.log(this.$route.query.text);
    axios
      .get("http://172.20.10.6:8081/search/video", {
        params: {
          key: this.$route.query.text,
          pageNo: 1,
        },
      })
      .then((res) => {
        console.log(res);
        this.dataList.list = res.data;
        this.flag = true;
      });
  },
  methods: {
    jumpPath() {},
  },
};
</script>

<style scoped>
.s_r_c_item {
  position: relative;
  width: 40%;
  height: 55%;
  /* width:18%;
  height:48%; */
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
