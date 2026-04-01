<script>
import Header from "@/components/header/HeaderFix.vue";
import Thop from "@/components/personal/top.vue";
import Ser from "@/components/personal/themain.vue";
import axios from "axios";
export default {
  components: {
    Header,
    Thop,
    Ser,
  },
  created() {
    console.log("created");
    
  },
};
</script>

<template>
  <div>
    <Header />
    <Thop />
    <Ser />
  </div>
</template>

<style scoped></style>
<style>
* {
  padding: 0;
  margin: 0;
  font-family: PingFang SC, HarmonyOS_Regular, Helvetica Neue, Microsoft YaHei,
    sans-serif !important;
}

li {
  list-style: none;
}

input {
  border: none;
  outline: none;
}

input:focus {
  border: none;
}
</style>
