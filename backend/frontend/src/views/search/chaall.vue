<script>
import Header from "@/components/header/HeaderFix.vue";
import Seall from "@/components/search/seall.vue";
export default {
  components: {
    Header,
    Seall,
  },
  created() {},
};
</script>

<template>
  <div>
    <Header></Header>
    <Seall />
  </div>
</template>

<style scoped></style>
<style>
* {
  padding: 0;
  margin: 0;
  font-family: PingFang SC, HarmonyOS_Regular, Helvetica Neue, Microsoft YaHei,
    sans-serif !important;
  font-weight: 400;
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
