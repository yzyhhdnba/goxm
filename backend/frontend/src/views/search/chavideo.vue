<script>
import Header from "@/components/header/HeaderFix.vue";
import Sevideo from "@/components/search/sevideo.vue";
export default {
  components: {
    Header,
    Sevideo,
  },
  created() {},
};
</script>

<template>
  <div>
    <Header />
    <Sevideo />
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
