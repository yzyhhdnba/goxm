<template>
  <div>
    <!-- 图片 -->
    <div class="item-pic">
      <img class="img-pic" :src="item.pic" alt="" />
      <div class="count">
        <el-row>
          <el-col :span="1" class="icon">
            <img
              src="https://s3.bmp.ovh/imgs/2022/06/25/74d27e4d317e41f1.png"
              alt=""
              style="width: 16px"
            />
          </el-col>
          <el-col :span="3">
            {{ item.view }}
          </el-col>
          <el-col :span="1" class="icon">
            <img
              src="https://s3.bmp.ovh/imgs/2022/06/25/da0d57c5fdacbfb1.png"
              alt=""
              style="width: 18px"
            />
          </el-col>
          <el-col :span="3">
            {{ item.comments }}
          </el-col>
        </el-row>
      </div>
      <div class="duration">
        {{ item.time }}
      </div>
    </div>
    <!-- 标题 -->
    <el-row class="item-title">
      <a href="#">{{ item.title }}</a>
    </el-row>
    <!-- up主 -->
    <el-row class="item-up">
      <el-row>
        <el-col :span="4" class="icon1">
          <img
            src="https://s3.bmp.ovh/imgs/2022/06/28/7b004d8bb8f75253.png"
            alt=""
            style="width: 17px"
          />
        </el-col>
        <el-col :span="20" style="float: left">
          <a href="#" style="text-decoration: none">
            {{ item.uploadtime }}
          </a>
        </el-col>
      </el-row>
    </el-row>
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
.item-pic {
  position: relative;
  width: 100%;
  max-height: 150px;
  border-radius: 5px;
  overflow: hidden;
}

.img-pic {
  width: 238.43px;
  height: 134.36px;
}

.item-pic img {
  display: block;
  width: 100%;
}

.item-pic .count {
  position: absolute;
  left: 0;
  bottom: 0;
  font-size: 12px;
  padding: 3px -5px;
  color: #fff;
  background-image: linear-gradient(
    to top,
    rgba(0, 0, 0, 0.6),
    rgba(119, 119, 119, 0.1)
  );
  width: 100%;
  /* 一行显示 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  overflow: hidden;
}

.icon {
  margin-top: 10px;
}

.icon1 {
  margin-top: 5px;
}
.duration {
  position: absolute;
  right: 0;
  bottom: 0;
  font-size: 12px;
  padding-right: 5px;
  color: #fff;
  background-image: linear-gradient(
    to left bottom,
    rgba(0, 0, 0, 0.6),
    rgba(255, 255, 255, 0.1)
  );
}

.item-title {
  height: 21%;
  margin-top: 5px;
  margin-bottom: 0px;
  font-size: 14px;
  /* 两行显示 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  color: #4d4d4d;
}

.item-up {
  padding-left: 5px;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.item-up a {
  font-size: 12px;
  color: #999;
}
</style>
