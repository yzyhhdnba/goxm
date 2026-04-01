<script>
import Header from "@/components/header/HeaderFix.vue";
import Thop from "@/components/personal/top.vue";
import Ser from "@/components/personal/themain.vue";
import { UserAPI } from "@/api/index";
export default {
  components: {
    Header,
    Thop,
    Ser,
  },
  data() {
    return {
      dashboard: null,
    };
  },
  async created() {
    if (!localStorage.getItem("access_token")) {
      return;
    }

    try {
      this.dashboard = await UserAPI.getDashboard();
    } catch (error) {
      console.error("load personal dashboard failed", error);
    }
  },
};
</script>

<template>
  <div>
    <Header />
    <Thop :dashboard="dashboard" />
    <Ser :dashboard="dashboard" />
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
