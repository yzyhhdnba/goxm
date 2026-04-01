<template>
  <div class="search-user">
    <!--
    personal-center
    -->
    <!-- <Header /> -->
    <!-- <Top /> -->
    <Ser />
  </div>
</template>
<script>
import Aside from "@/components/personal/aside";
import Reuser from "@/components/search/reuser";
import Top from "@/components/personal/top.vue";
import Ser from "@/components/personal/themain.vue";
import Header from "@/components/header/HeaderFix.vue";
import { ref } from "vue";

const tabPosition = ref("left");
export default {
  name: "PersonalCenter",
  components: {
    Reuser,
    Aside,
    tabPosition,
    Ser,
    Top,
    Header,
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
