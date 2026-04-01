<template>
  <div class="swiper_right_content">
    <div v-for="(item, i) in dataList" :key="i" class="s_r_c_item">
      <one-card :item="item" />
    </div>
  </div>
</template>

<script>
import OneCard from "@/components/common/card/OneCard";
export default {
  props: {
    dashboard: {
      type: Object,
      default: null,
    },
  },
  components: {
    OneCard,
  },
  computed: {
    dataList() {
      return this.dashboard?.favorite_videos || [];
    },
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
