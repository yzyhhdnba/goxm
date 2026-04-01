<template>
  <div>
    <div>
      <div class="search">
        <div class="search-input">
          <input v-model="searchKeyword" class="header-input" placeholder="搜索稿件标题或作者" />
        </div>
        <div class="search-btn">
          <svg
            width="17"
            height="17"
            viewBox="0 0 17 17"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              clip-rule="evenodd"
              d="M16.3451 15.2003C16.6377 15.4915 16.4752 15.772 16.1934 16.0632C16.15 16.1279 16.0958 16.1818 16.0525 16.2249C15.7707 16.473 15.4456 16.624 15.1854 16.3652L11.6848 12.8815C10.4709 13.8198 8.97529 14.3267 7.44714 14.3267C3.62134 14.3267 0.5 11.2314 0.5 7.41337C0.5 3.60616 3.6105 0.5 7.44714 0.5C11.2729 0.5 14.3943 3.59538 14.3943 7.41337C14.3943 8.98802 13.8524 10.5087 12.8661 11.7383L16.3451 15.2003ZM2.13647 7.4026C2.13647 10.3146 4.52083 12.6766 7.43624 12.6766C10.3517 12.6766 12.736 10.3146 12.736 7.4026C12.736 4.49058 10.3517 2.1286 7.43624 2.1286C4.50999 2.1286 2.13647 4.50136 2.13647 7.4026Z"
              fill="currentColor"
            />
          </svg>
        </div>
      </div>
    </div>
    <div v-if="flag">
      <div v-for="item in displayVideos" :key="item.id">
        <div class="container">
          <div class="left">
            <img :src="item.cover_url" @click="watchVideo(item.play_url)" style="cursor: pointer" />
          </div>
          <div class="middle">
            <div class="title">{{ item.title }}</div>
            <div>{{ item.author_username }}</div>
            <div class="time">{{ formatDate(item.updated_at || item.created_at) }}</div>
            <div class="time" v-if="item.review_reason">审核备注：{{ item.review_reason }}</div>
          </div>
          <div class="right">
            <div>
              <el-button @click="pass(item.id)">审核通过</el-button>
              <el-button @click="openDef(item.id)" type="danger" style="color: white">不通过</el-button>
            </div>
          </div>
        </div>
        <el-divider />
      </div>
    </div>
  </div>
  <el-dialog v-model="dialogVisible" title="视频预览">
    <Video :url="url"></Video>
  </el-dialog>
  <el-dialog v-model="isOpenDef" title="">
    <el-input v-model="submitForm.reason" :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" placeholder="请输入拒绝意见" />
    <el-button @click="notPass()">提交</el-button>
  </el-dialog>
</template>

<script>
import { ref, reactive, h, onMounted, computed } from "vue";
import Video from "@/components/video/video.vue";
import { AdminAPI } from "@/api/index";
import { ElNotification } from "element-plus";
export default {
  components: {
    Video,
  },
  setup() {
    let flag = ref(false);
    let videos = ref([]);
    let dialogVisible = ref(false);
    let searchKeyword = ref('');
    let submitForm = reactive({
      videoId: '',
      reason: ''
    });
    let isOpenDef = ref(false);
    let url = ref('');
    const loadVideos = () => {
      AdminAPI.listVideos({
        review_status: 'reviewed',
        page: 1,
        page_size: 20,
      }).then((res) => {
        videos.value = res?.list || [];
        flag.value = true;
      });
    };
    onMounted(() => {
      loadVideos();
    });
    const openDef = (id) => {
      isOpenDef.value = true;
      submitForm.videoId = id;
      submitForm.reason = '';
    };
    const pass = (videoId) => {
      AdminAPI.approveVideo(videoId).then(() => {
        ElNotification({
          title: "通过",
          message: h("p", { style: "color: green" }, "审核通过"),
          type: 'success'
        });
        loadVideos();
      })
    };
    const notPass = () => {
      AdminAPI.rejectVideo(submitForm.videoId, { reason: submitForm.reason }).then(() => {
        isOpenDef.value = false;
        ElNotification({
          title: "提示",
          message: h("p", { style: "color: green" }, "驳回成功"),
          type: 'success'
        });
        loadVideos();
      })
    };
    const watchVideo = (playUrl) => {
      if (!playUrl) {
        return;
      }
      url.value = playUrl;
      dialogVisible.value = true;
    };
    const formatDate = (value) => {
      return value ? value.slice(0, 10) : '';
    };
    const displayVideos = computed(() => {
      const keyword = searchKeyword.value.trim().toLowerCase();
      if (!keyword) {
        return videos.value;
      }
      return videos.value.filter((item) => {
        const title = item.title?.toLowerCase() || '';
        const author = item.author_username?.toLowerCase() || '';
        return title.includes(keyword) || author.includes(keyword);
      });
    });
    return {
      submitForm,
      isOpenDef,
      openDef,
      pass,
      notPass,
      watchVideo,
      url,
      dialogVisible,
      videos,
      flag,
      displayVideos,
      searchKeyword,
      formatDate,
    };
  },
};
</script>

<style lang="less" scoped>
.search-btn {
  position: relative;
  color: #000000;
  cursor: pointer;
  z-index: 1999;
  border-radius: 6px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: 0.5s;
}
.search {
  margin-bottom:20px;
  display: flex;
  height: 40px;
  border: 1px;
  background-color: #f1f2f3;
  border-radius: 8px;
  justify-content: space-between;
  align-items: center;
  margin-right: 30px;
  padding-right: 5px;
  flex: 1;
  transition: 0.5s;
  opacity: 0.8;
}
.header-input {
  line-height: 20px;
  width: 90%;
  transition: 0.5s;
  border-radius: 6px;
  margin-left: 5px;
  padding: 5px 0;
  padding-right: 30px;
  margin-right: 3px;
  background-color: #f1f2f3;
}
* {
  color: #505050;
}

.container {
  display: flex;
  width: 1200px;
  align-items: center;

  .left {
    height: 100px;
    width: 155px;
    border-radius: 5px;

    img {
      width: 150px;
      height: 100px;
      overflow: hidden;
      border-radius: 5px;
      object-fit: cover;
    }
  }

  .middle {
    margin-left: 40px;
    display: flex;
    justify-content: space-around;
    flex-direction: column;
    height: 100px;

    div {
      margin: 6px 0;
    }

    .title {
      font-size: 20px;
      margin-top: 3px;
    }

    .time {
      margin: 0;
      margin-top: 3px;
    }
  }

  .right {
    flex: 5;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
