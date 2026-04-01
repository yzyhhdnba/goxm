<template>
  <div class="container">
    <div class="box">
      <el-result
        icon="success"
        title="注册成功！"
        sub-title="请等待3秒或点击下面的按钮完成注册"
      >
        <template #extra>
          <el-button type="primary">完成注册</el-button>
        </template>
      </el-result>
    </div>
  </div>
</template>

<script>
export default {
  mounted() {
    let that = this;
    setTimeout(function(){
      that.$router.replace('/')
    },3000)
  }
}
</script>

<style lang="less" scoped>
.container {
  width: 100%;
  height: 900px;
  display: flex;
  justify-content: center;
  align-items: center;
  .box{
    width:80%;
    border:2px solid #030303;
    position: relative;
    top:-50px;
    display: flex;
    justify-content: center;
    align-items:center;
    flex-direction: column;
    img{
      // translate:scale(2);
      margin-bottom: -20px;
      height:150px;
    }
    .title{
      margin-bottom: 20px;
    }
  }
}
</style>