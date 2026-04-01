<template>
  <div class="swiper_right_content">
    <div v-for="(item, i) in dataList" :key="i" class="s_r_c_item">
      <div @click="jumpPath(item)">
        <one-card :item="item" />
      </div>
    </div>
  </div>
</template>

<script>
import OneCard from "@/components/common/card/Usercard";
import { SearchAPI } from "@/api/index";
export default {
  props: {},
  components: {
    OneCard,
  },
  data() {
    return {
      dataList: [],
    };
  },
  watch: {
    "$route.query.text": {
      immediate: true,
      handler() {
        this.loadUsers();
      },
    },
  },
  methods: {
    loadUsers() {
      const keyword = this.$route.query.text;
      if (!keyword) {
        this.dataList = [];
        return;
      }
      SearchAPI.searchUsers({ keyword, page: 1, page_size: 20 }).then((res) => {
        this.dataList = res?.list || [];
      });
    },
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
