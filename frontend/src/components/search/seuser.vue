<template>
  <el-container class="search-user">
    <!--
    personal-center
    -->
    <el-row span="8"> </el-row>
    <el-row span="16">
      <el-main>
        <el-tabs :tab-position="tabPosition" style="height: 600px; width: 1400px" class="demo-tabs">
          <el-tab-pane label="❀&nbsp;用户">
            <Reuser />
          </el-tab-pane>
        </el-tabs>
      </el-main>
    </el-row>
  </el-container>
</template>
<script>
import Aside from "@/components/personal/aside";
import Reuser from "@/components/search/reuser";
import { ref } from "vue";

const tabPosition = ref("left");
export default {
  name: "PersonalCenter",
  components: {
    Reuser,
    Aside,
    tabPosition,
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
.search-user {
  margin: auto;
  padding-top: 45px;
  opacity: 0.825;
}

.demo-tabs>.el-tabs__content {
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
