<template>
  <!-- // historytime: "2022-06-30 00:16",
          // videoimg: "https://s3.bmp.ovh/imgs/2022/06/21/89170619ddaad042.png",
          // videotitle: "丁真今晚吃电弧星",
          // upname: "狗三儿_Official",
          // areatip: "国创", -->
  <div>
    <el-main class="history-card">
      <el-row>
        <el-col :span="1" class="triangle-sign"
          ><el-icon class="triangle-right"><CaretRight /></el-icon
        ></el-col>
        <el-col :span="2" class="history-time">{{ item.latestTime.slice(0, 10)+' '+ item.latestTime.slice(11, 19)}}</el-col>
        <el-col :span="19">
          <el-row>
            <el-col :span="7" @click="to(item.videoId)">
              <el-image
                style="width: 244px; height: 137.5px"
                class="video-img"
                :src="'http://172.20.10.6:8081/cover/cover-' + item.videoId + '.jpeg'"
              />
            </el-col>
            <el-col :span="17" class="video-history">
              <el-row class="video-title" @click="to(item.videoId)"> {{ item.videoTitle }} </el-row>
              <el-row class="up-name">
                {{ item.userName }}&nbsp;丨&nbsp;{{ item.areaName }}
              </el-row>
            </el-col>
          </el-row>
        </el-col>

        <!-- <el-col :span="3">
          <el-image :size="80" :src="item.pic" />
        </el-col>
        <el-col :span="10">
          <el-row class="item-username">
            <a href="#">{{ item.username }}</a>
          </el-row>
          <el-row class="item-introduce">
            <p>{{ item.fans }}粉丝&nbsp;·&nbsp;{{ item.video }}个视频</p>
          </el-row>
        </el-col> -->
      </el-row>
    </el-main>
  </div>
</template>

<script>
import {useRouter} from 'vue-router'
export default {
  props: {
    item: { type: Array, default: [] },
  },
  components: {},
  setup() {
    const router= useRouter();
    const to=(videoId)=>{
       router.push({
        path:'/video',
        query:{
          videoId
        }
      })
    }
    return {
      to
    };
  },
};
</script>

<style scoped>
.item-username {
  height: 21%;
  margin-top: 3px;
  margin-bottom: 20px;
  font-size: 17px;
  font-weight: 700;
  /* 两行显示 */
  color: #000000;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.triangle-sign {
  border-left-style: solid;
  border-left-color: #999;
}
.triangle-right {
  padding-top: 57px;
  color: #999;
}
.history-time {
  padding-top: 55px;
  font-size: 15px;
  color: #999;
}

.video-img {
  border-radius: 5px;
  cursor: pointer;
}

.video-history {
  border-bottom-style: solid;

  border-bottom-color: #f1f1f1;
}

.video-title {
  font-size: 18px;
  font-weight: 700;
  padding-top: 10px;
  padding-bottom: 78px;
  cursor: pointer;
}

.up-name {
  font-size: 15px;
  color: #999;
}
</style>
