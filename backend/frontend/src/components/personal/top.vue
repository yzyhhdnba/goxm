<template>
  <div class="personal-top">
    <!--
    personal-top类名
    说明这个是头部组件
    -->
    <el-col :span="3">
      <el-row class="myhead-portrait">
        <!--
        头像展示
        -->
        <el-avatar :size="90" :src="require('@/assets/head.jpg')"></el-avatar>
      </el-row>
    </el-col>
    <el-col :span="21">
      <el-row class="myname">
        <!--
        用户姓名
        -->
        <span>
          {{ userInfo.userName }}
        </span>
      </el-row>
      <!--
      个人签名
      -->
      <el-input
        class="mywrite"
        v-model="idiograph"
        placeholder="编辑个性签名"
      ></el-input>
    </el-col>
  </div>
</template>
<script>
export default {
  name: "PersonalTop",
  setup() {
    let userInfo=JSON.parse(localStorage.getItem('userInfo'));
    return {
      userInfo,
      idiograph: "俄而彩蝶萦丛绕,不识桃花倾于心",
      circleUrl: "",
      sizeList: ["large"],
      tableData: [{}],
    };
  },
  methods: {},
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
