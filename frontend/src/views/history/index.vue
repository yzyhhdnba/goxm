<template>
  <div style="width: 100%; height: 100%">
    <HeaderFix></HeaderFix>
    <History></History>
  </div>
</template>

<script>
import History from "@/components/history/history.vue";
import HeaderFix from "@/components/header/HeaderFix.vue";
export default {
  components: {
    History,
    HeaderFix,
  },
  beforeCreate() {
    if (!localStorage.getItem("isLogin")) {
      alert("您还未登录！");
      this.$router.replace("/");
    }
  },
};
</script>

<style></style>
