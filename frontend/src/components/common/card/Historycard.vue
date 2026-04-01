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
        <el-col :span="2" class="history-time">{{ displayTime }}</el-col>
        <el-col :span="19">
          <el-row>
            <el-col :span="7" @click="to(videoId)">
              <el-image
                style="width: 244px; height: 137.5px"
                class="video-img"
                :src="coverUrl"
              />
            </el-col>
            <el-col :span="17" class="video-history">
              <el-row class="video-title" @click="to(videoId)"> {{ videoTitle }} </el-row>
              <el-row class="up-name">
                {{ authorName }}&nbsp;丨&nbsp;{{ areaName }}
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
import { computed } from 'vue'
export default {
  props: {
    item: { type: Object, default: () => ({}) },
  },
  components: {},
  setup(props) {
    const router= useRouter();
    const watchedAt = computed(() => props.item.watched_at || props.item.latestTime || '');
    const displayTime = computed(() => {
      if (!watchedAt.value) {
        return '';
      }
      return watchedAt.value.slice(0, 10)+' '+ watchedAt.value.slice(11, 19);
    });
    const videoId = computed(() => props.item.video_id || props.item.videoId || 0);
    const videoTitle = computed(() => props.item.video_title || props.item.videoTitle || '');
    const coverUrl = computed(() => props.item.cover_url || '');
    const authorName = computed(() => props.item.author_name || props.item.userName || '');
    const areaName = computed(() => props.item.area_name || props.item.areaName || '');
    const to=(videoId)=>{
       router.push({
        path:'/video',
        query:{
          videoId
        }
      })
    }
    return {
      to,
      displayTime,
      videoId,
      videoTitle,
      coverUrl,
      authorName,
      areaName,
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
