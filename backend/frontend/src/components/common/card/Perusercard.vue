<template>
  <div>
    <!-- 图片 -->
    <div class="item-user">
      <el-row>
        <el-col :span="3">
          <el-avatar :size="80" :src="require('@/assets/head.jpg')" />
        </el-col>
        <el-col :span="10">
          <el-row class="item-username">
            <a href="#">{{ item.userName }}</a>
          </el-row>
          <el-row class="item-introduce">
            <p>{{ item.userSignature }}</p>
          </el-row>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    item: { type: Array, default: [] },
  },
  components: {},
  setup() {
    return {};
  },
  methods: {},
};
</script>

<style scoped>
.item-user {
  border-bottom-style: solid;
  border-bottom-color: #999;
}

.item-username {
  height: 21%;

  margin-top: 3px;
  margin-bottom: 20px;
  font-size: 17px;
  font-weight: 700;
  /* 两行显示 */
  color: #000000;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.item-introduce {
  font-size: 15px;
  color: #999;
}
</style>
