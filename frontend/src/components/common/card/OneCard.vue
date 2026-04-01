<template>
  <div>
    <div class="item-pic" @click="play(item.id || item.videoId)">
      <img :src="item.cover_url || item.coverUrl || require('@/assets/image/tiao.jpg')" />
      <div class="count">
        <div style="position:relative;top:14px;">
          <el-row>
            <el-col :span="1" class="icon">
              <img src="@/assets/image/player.svg" alt="" style="width: 16px" />
            </el-col>
            <el-col :span="3">
              {{ item.view_count !== undefined ? item.view_count : item.videoHits }}
            </el-col>
            <el-col :span="1" class="icon">
              <img src="@/assets/image/pinglun.svg" alt="" style="width: 18px" />
            </el-col>
            <el-col :span="3">
              {{ item.comment_count !== undefined ? item.comment_count : item.videoComments }}
            </el-col>
          </el-row>
        </div>
      </div>
      <div class="duration">
        {{ item.duration_seconds !== undefined ? formatDuration(item.duration_seconds) : item.videoDuration }}
      </div>
    </div>
    <!-- 标题 -->
    <div class="item-title">
      <a href="#">{{ item.title || item.videoTitle }}</a>
    </div>
    <!-- up主 -->
    <div class="item-up">
      <a href="#" style="text-decoration:none;">
        <el-row>
          <el-col :span="1" class="icon" style="margin-top:4px;">
            <img src="@/assets/image/up.png" alt="" style="width: 17px;" />
          </el-col>
          <el-col :span="3" style="margin-left:5px;font-size: 14px;">
            {{ (item.author && item.author.username) || item.userName }}
          </el-col>
        </el-row>
      </a>
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router'
export default {
  props: {
    item: { type: Object, default: () => ({}) },
  },
  setup(props) {
    const router = useRouter();
    const play = (id) => {
      router.push({
        path:'/video',
        query:{
          videoId: id
        }
      })
    }
    const formatDuration = (seconds) => {
      if (seconds === undefined || seconds === null) return '00:00';
      const m = Math.floor(seconds / 60).toString().padStart(2, '0');
      const s = (seconds % 60).toString().padStart(2, '0');
      return m + ':' + s;
    }
    return {
      props,
      play,
      formatDuration
    }
  },
}
</script>

<style scoped lang="less">
.item-pic {
  cursor: pointer;
  position: relative;
  width: 100%;
  max-height: 150px;
  border-radius: 5px;
  overflow: hidden;
}

.item-pic img {
  display: block;
  width: 100%;
}

.item-pic .count {
  position: absolute;
  left: 0;
  bottom: 0;
  font-size: 12px;
  padding: 3px 5px;
  color: #fff;
  width: 100%;
  height: 30px;
  /* 一行显示 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
  overflow: hidden;
  background-image: linear-gradient(180deg, rgba(0, 0, 0, 0) 0%, rgba(0, 0, 0, .8) 100%);
  align-items: center;
  /* overflow: hidden; */
}

.icon {
  margin-right: 5px;

}

.duration {
  position: absolute;
  right: 0;
  bottom: 0;
  font-size: 12px;
  padding: 3px 5px;
  color: #fff;
  background-image: linear-gradient(to left bottom,
      rgba(0, 0, 0, 0.6),
      rgba(255, 255, 255, 0.1));
  border-bottom-left-radius: 8px;
}

.item-title {
  height: 21%;
  margin: 12px 0;
  font-size: 18px;
  /* 两行显示 */
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  color: #4d4d4d;
}

.item-up {
  padding-left: 5px;
  font-size: 16px;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.item-up a {
  font-size: 12px;
  color: #999;
}
</style>
