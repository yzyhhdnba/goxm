<template>
  <div>
    <div style="width: 99vw; height: 140px">
      <img src="https://s3.bmp.ovh/imgs/2022/06/23/a628d7fade3c0f61.jpg" style="width: 100%" />
    </div>
    <div class="card-box">
      <el-row justify="center">
        <el-col class="title" :span="3">Message</el-col>
      </el-row>
      <br />
      <el-row justify="center">
        <el-col :span="1" :offset="2">
          <el-icon :size="25">
            <Check />
          </el-icon>
        </el-col>
        <el-col :span="16">
          <span>一段话一段话一段话一段话一段话一段话一段话一段话一段话<a herf="#">>>>查看更多</a></span>
        </el-col>

      </el-row>
    </div>
  </div>
</template>

<script>
import { Check, } from '@element-plus/icons-vue';
  export default {
    props: {

    },
    components: {

    },
    setup(){
      return {
        Check: Check,
      }
    },
    methods: {

    }
  }
</script>

<style scoped>

.card-box {
  width:75vw;
  height:200px;
  border: 5px solid #adc056;
  border-radius: 10px;
  background-color:#d1f8cf;
  margin: 0 auto;
}
.title {
  font-size:18px;
}
</style>