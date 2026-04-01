<template>
  <div class="reco-item" @click="to(item.videoId)">
    <div class="card-box" @click="to(item.videoId)">
      <div class="pic-box">
        <div class="pic-pre">
          <a>
            <img :src="'http://172.20.10.6:8081/cover/cover-' + item.videoId + '.jpeg'" class="pic"
              @click="to(item.videoId)" />
          </a>
        </div>
      </div>
      <div class="info">
        <a>
          <p class="info-title" @click="to(item.videoId)">{{ item.videoTitle }}</p>
        </a>
        <div class="info-up">
          <a class="up-link"><span style="overflow: hidden;text-overflow: ellipsis;" @click="to(item.videoId)">{{
              item.userName
          }}</span></a>
        </div>
        <div class="info-play">
          播放量:{{ item.videoHits }} 评论:{{ item.videoComments }}
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { useRouter } from "vue-router"
export default {
  name: "Recommend",
  props: {
    item: { type: Object },
  },
  setup() {
    const router = useRouter();
    const to = (videoId) => {
      router.push({
        path: '/video',
        query: { videoId }
      })
      location.reload()
    }
    return {
      to
    }
  },
  methods: {

  }
}
</script>

<style scoped>
.reco-item {
  margin-bottom: 12px;
}

.card-box {
  display: flex;
}

.pic-box {
  position: relative;
  width: 141px;
  height: 80px;
  border-radius: 6px;
  background: #C9CCD0;
  flex: 0 0 auto;
}

.pic-pre {
  height: calc(100% + 10px);
  position: absolute;

  width: 100%;
  bottom: 0;
  border-radius: 6px;
}

.pic-pre a {
  display: block;
  position: absolute;
  top: 10px;
  text-decoration: none;
  outline: none;
}

.pic {
  width: 141px;
  height: 78.2px;
  object-fit: cover;
  display: block;
  margin: 0 auto;
  border-radius: 6px;
}

.info {
  margin-left: 10px;
  flex: 1;
  font-size: 13px;
  line-height: 15px;
}

.info a {
  text-decoration: none;
  outline: none;
}

.info-title {
  color: #18191C;
  font-size: 15px;
  line-height: 19px;
  display: -webkit-box;
  overflow: hidden;
  -webkit-box-orient: vertical;
  text-overflow: ellipsis;
  word-break: break-word;
  -webkit-line-clamp: 2;
  font-weight: 500;
  -webkit-font-smoothing: antialiased;
}

.info-up {
  cursor: pointer;
  margin: 2px 0;
  height: 20px;
}

.up-link {
  color: #9499A0;
  display: inline-flex;
  align-items: center;
  width: 100%;
  height: 100%;
  text-decoration: none;
  outline: none;
}

.info-play {
  color: #9499A0;
  display: inline-flex;
  align-items: center;
}
</style>