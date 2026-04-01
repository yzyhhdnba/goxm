<script>
import Header from "@/components/header/HeaderFix.vue";
import Seuser from "@/components/search/seuser.vue";
export default {
  components: {
    Header,
    Seuser,
  },
  created() {},
};
</script>

<template>
  <div>
    <Header />
    <Seuser />
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
