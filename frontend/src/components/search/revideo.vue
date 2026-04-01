<template>
  <div class="swiper_right_content" v-if="flag">
    <div v-for="(item, i) in dataList.list" :key="i" class="s_r_c_item">
      <div @click="jumpPath(item)">
        <one-card :item="item" />
      </div>
    </div>
  </div>
</template>

<script>
import OneCard from "@/components/common/card/OneCard.vue";
import { SearchAPI } from "@/api/index";
export default {
  props: {},
  components: {
    OneCard,
  },
  data() {
    return {
      dataList: { list: [] },
      flag: false,
    };
  },
  watch: {
    "$route.query.text": {
      immediate: true,
      handler() {
        this.loadVideos();
      },
    },
  },
  methods: {
    loadVideos() {
      const keyword = this.$route.query.text;
      if (!keyword) {
        this.dataList = { list: [] };
        this.flag = true;
        return;
      }
      SearchAPI.searchVideos({ keyword, page: 1, page_size: 20 }).then((res) => {
        this.dataList = res || { list: [] };
        this.flag = true;
      });
    },
    jumpPath(item) {
      const videoId = item.id || item.videoId;
      if (!videoId) {
        return;
      }
      this.$router.push({
        path: "/video",
        query: {
          videoId,
        },
      });
    },
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
