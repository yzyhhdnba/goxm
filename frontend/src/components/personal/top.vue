<template>
  <div class="personal-top">
    <el-col :span="3">
      <el-row class="myhead-portrait">
        <el-avatar :size="90" :src="avatarUrl"></el-avatar>
      </el-row>
    </el-col>
    <el-col :span="21">
      <el-row class="myname">
        <span>
          {{ username }}
        </span>
      </el-row>
      <el-input
        class="mywrite"
        :model-value="bio"
        readonly
        placeholder="这个人很神秘，还没有留下签名"
      ></el-input>
    </el-col>
  </div>
</template>
<script>
export default {
  name: "PersonalTop",
  props: {
    dashboard: {
      type: Object,
      default: null,
    },
  },
  computed: {
    userInfo() {
      return this.dashboard?.user || JSON.parse(localStorage.getItem("userInfo") || "null") || {};
    },
    username() {
      return this.userInfo.username || this.userInfo.userName || "PILIPILI Go";
    },
    bio() {
      return this.userInfo.bio || "这个人很神秘，还没有留下签名";
    },
    avatarUrl() {
      return this.userInfo.avatar_url || this.userInfo.avatarUrl || require("@/assets/head.jpg");
    },
  },
};
</script>
<style scoped>
.personal-top {
  height: 200px;
  width: 1200px;
  margin: auto;
  line-height: 60px;
  opacity: 0.825;
  text-align: center;
  background-color: #b3c0d1;
  background-image: url(https://s3.bmp.ovh/imgs/2022/04/16/97d72f38964e925c.jpg);
}
.demo-basic--circle {
  align-self: center;
  margin-top: 30px;
  margin-left: 150px;
}
.myname {
  margin-left: 25px;
  font-weight: bold;
  font-size: 160%;
  color: rgb(255, 255, 255);
  position: relative;
  float: left;
  padding-top: 90px;
  font-weight: bolder;
}

.myhead-portrait {
  position: relative;
  float: left;
  margin-left: 20px;
  padding-top: 95px;
}
.mywrite {
  position: relative;
  float: left;
  margin-left: 130px;
  margin-top: -40px;
}

.address-name {
  color: white;
  text-decoration: none;
  margin-left: 20px;
}

.address-name:hover,
.address-name:focus {
  color: #66ccff;
  text-decoration: underline;
}
</style>
