<template>
  <el-container class="personal-center">
    <!--
    personal-center
    -->
    <el-row span="8">
      <el-aside width="300px"><Aside /></el-aside>
    </el-row>
    <el-row span="16">
      <el-main>
        <el-tabs
          :tab-position="tabPosition"
          style="height: 600px; width: 1050px"
          class="demo-tabs"
        >
          <el-tab-pane label="❀&nbsp;视频"> <Revideo /> </el-tab-pane>
          <el-tab-pane label="➹&nbsp;收藏"> <College /> </el-tab-pane>
          <el-tab-pane label="✧&nbsp;关注"> <Attention /> </el-tab-pane>
          <el-tab-pane label="✲&nbsp;设置"> <Install /> </el-tab-pane>
        </el-tabs>
      </el-main>
    </el-row>
  </el-container>
</template>
<script>
import Aside from "@/components/personal/aside";
import Revideo from "@/components/personal/revideo";
import College from "@/components/personal/college";
import Attention from "@/components/personal/attention";
import Install from "@/components/personal/install";
import { ref } from "vue";

const tabPosition = ref("left");
export default {
  name: "PersonalCenter",
  components: {
    Revideo,
    Aside,
    College,
    tabPosition,
    Attention,
    Install,
  },
  return() {
    return {
      addressName: "admin",
      uid: "10000000",
      sizeList: ["large"],
      tableData: [],
    };
  },
  methods: {},
};
</script>

<style scoped>
.personal-center {
  margin: auto;
  width: 1500px;
  opacity: 0.825;
}

.el-header {
  height: 200px;
  line-height: 60px;
  text-align: center;
  background-color: #b3c0d1;
  background-image: url(https://s3.bmp.ovh/imgs/2022/04/16/97d72f38964e925c.jpg);
}

.el-aside {
  line-height: 44px;
  text-align: center;
  background-color: rgb(255, 255, 255);

  color: Transparent;
}

.el-main {
  line-height: 36px;
  text-align: center;
  width: 1200px;
  background-color: rgb(255, 255, 255);
}

.message {
  text-decoration: none;
  color: #000;
}

.message:hover,
.message:focus {
  color: #66ccff;
  text-decoration: underline;
}

.demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}

.el-tabs--right .el-tabs__content,
.el-tabs--left .el-tabs__content {
  height: 100%;
}
</style>
