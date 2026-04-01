<template>
  <div class="swiper_right_content">
    <div v-for="(item, i) in dataList" :key="i" class="s_r_c_item">
      <one-card :item="item" />
    </div>
  </div>
</template>

<script>
import OneCard from "@/components/common/card/Perusercard";
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
      return this.dashboard?.following_users || [];
    },
  },
};
</script>

<style scoped>
.s_r_c_item {
  position: relative;
  width: 99%;
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
