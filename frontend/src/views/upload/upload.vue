<template>
  <div class="header">
    <div class="header-box">
      <div class="header-left">
        <div class="title">创作者中心</div>
        <div class="tv"></div>
        <div @click="to('/')" class="to-index">主站</div>
      </div>
      <div class="header-right">

      </div>
    </div>
  </div>
  <el-divider />

  <div class="tac">
    <el-row class='main'>
      <el-col :span="4" style="height:1000px;">
        <el-menu default-active="/upload" class="el-menu-vertical-demo" :router="true">
          <el-menu-item index="/upload">
            <span>视频上传</span>
          </el-menu-item>
          <el-menu-item index="/upload/pass">
            <span>已通过稿件</span>
          </el-menu-item>
          <el-menu-item index="/upload/toBePass">
            <span>待审核稿件</span>
          </el-menu-item>
          <el-menu-item index="/upload/notPass">
            <span>未通过稿件</span>
          </el-menu-item>
        </el-menu>
      </el-col>
      <el-col :span="18" class="col-right">
        <router-view></router-view>
      </el-col>
    </el-row>
  </div>

</template>

<script>



import { useRouter } from 'vue-router';
export default {

  setup() {
    let defaultAction = '/upload'
    let router = useRouter();
    let to = (url) => {
      router.push(url)
    }
    return {
      to, defaultAction
    }
  }
}
</script>

<style scoped lang="less">
.header {
  height: 40px;
  width: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 20px;

  .header-box {
    width: 90%;
    display: flex;
    justify-content: space-between;
  }

  .header-left {
    display: flex;
    align-items: center;

    .title {
      font-weight: 600;
      font-size: 36px;
      color: rgb(0, 174, 236);
    }

    .tv {
      margin-left: 30px;
      height: 20px;
      width: 20px;
      background-image: url('../../assets/image/icons-upload.png');
      background-position: 298px -212px;
    }

    .to-index {
      cursor: pointer;
      color: black;
      transition: all 0.2s;
    }

    .to-index:hover {
      transition: all 0.2s;
      color: rgb(0, 174, 236);
    }
  }
}

.tac {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;

  .main {
    width: 95%;
  }

  .col-right {
    display: flex;
    width: 100%;
    position: relative;
    left: 40px;
    // align-items: center;
    justify-content: center;
    margin-left: 20px;
  }
}
</style>