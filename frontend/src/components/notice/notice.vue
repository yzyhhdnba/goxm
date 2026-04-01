<template>
  <div class="message">
    <div style="height:64px"></div>
    <div class="container">
      <el-container>
        <el-aside width="140px">
          <div class="aside">
            <div class="aside-title">
              <div class="icon1">
              </div>
              <span>
                消息中心
              </span>
            </div>
            <div class="aside-list">
              系统消息
            </div>
          </div>
        </el-aside>
        <el-main class="main">
          <div class="main-top">
            <div class="main-top-title">
              系统通知
            </div>
          </div>

          <div class="main-body">
            <el-scrollbar max-height="580px">
              <Floor v-for="item in posts" :key="item.id" :item="item" @read="handleRead" />
            </el-scrollbar>
          </div>

        </el-main>
      </el-container>
    </div>
  </div>
</template>

<script>
import Floor from "./floor";
import { ref } from 'vue'
import { NoticeAPI } from '@/api/index'
export default {
  name: "Notice",
  components: {
    Floor,
  },
  setup() {
    const posts = ref([]);
    NoticeAPI.list({
      page: 1,
      page_size: 50,
    }).then((res) => {
      posts.value = res?.list || [];
    }).catch((error) => {
      console.error("load notices failed", error);
    });

    const handleRead = (noticeId) => {
      posts.value = posts.value.map((item) => {
        if (item.id !== noticeId) {
          return item;
        }
        return {
          ...item,
          read: true,
          read_at: new Date().toISOString(),
        };
      });
    };

    return {
      posts,
      handleRead,
    }
  },
}
</script>

<style scoped>
.message {
  background: url(//s1.hdslb.com/bfs/static/blive/blfe-message-web/static/img/infocenterbg.a1a0d152.jpg) top/cover no-repeat fixed;
  font-size: 12px;
  line-height: 12px;
  color: #666;
  padding: 0;
  margin: 0;
}

.container {
  height: calc(100vh - 64px);
  margin: 0 auto;
  max-width: 1143px;
}

.aside {
  height: calc(100vh - 64px);
  width: 140px;
  min-width: 140px;
  background-color: rgba(255, 255, 255, 0.8);
}

.aside-title {
  height: 62px;
  display: flex;
  -webkit-box-pack: center;
  justify-content: center;
  -webkit-box-align: center;
  align-items: center;
  color: #333;
  font-size: 14px;
  font-weight: 700;
}

.icon1 {
  width: 14px;
  height: 16px;
  margin-right: 10px;
  background: url(https://s1.hdslb.com/bfs/static/blive/blfe-message-web/static/img/plane.c9984cf0.svg) center/contain no-repeat;
}

.aside-list {
  margin: 0;
  padding-left: 40px;
  font-size: 14px;
}

.el-main {
  --el-main-padding: 0px !important;
}

.main {
  height: calc(100vh - 64px);
  background-color: rgba(255, 255, 255, 0.5);
}

.main-top {
  padding: 10px 10px 0;
}

.main-top-title {
  height: 42px;
  background-color: #fff;
  box-shadow: 0 2px 4px 0 rgb(121 146 185 / 54%);
  display: flex;
  -webkit-box-align: center;
  align-items: center;
  -webkit-box-pack: justify;
  justify-content: space-between;
  padding: 0 16px;
  font-size: 15px;
  color: #666;
  border-radius: 4px;
}

.main-body {
  padding: 10px;
}
</style>
