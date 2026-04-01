<template>
  <div>
    <div class="selecter">
      <el-select v-model="value" class="m-2" placeholder="排序方式">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </div>
    <div v-if="flag">
      <div v-for="item in videos" :key="item.id">
        <div class="container">
          <div class="left">
            <img :src="item.cover_url" @click="watchVideo(item.play_url)" style="cursor: pointer" />
          </div>
          <div class="middle">
            <div class="title">{{ item.title }}</div>
            <div>{{ item.author_username }}</div>
            <div class="time">{{ formatDate(item.created_at) }}</div>
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
    <Video :url="url" />
  </el-dialog>
  <el-dialog v-model="isOpenDef" title="">
    <el-input v-model="submitForm.reason" :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" placeholder="请输入拒绝意见" />
    <el-button @click="notPass()">提交</el-button>
  </el-dialog>
</template>

<script>
import { ref, reactive, h, onMounted } from "vue";
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
    let submitForm = reactive({
      videoId: '',
      reason: ''
    });
    let isOpenDef = ref(false);
    let url = ref('');
    const options = [
      { value: "1", label: "按发布时间排序" },
      { value: "2", label: "按视频播放量排序" },
      { value: "3", label: "按点赞排序" },
      { value: "4", label: "按收藏排序" },
    ];
    // loadVideos 对应文档“管理后台：审核列表与统计看板”。
    // 这是审核页加载待审稿件列表的入口，也是 Mermaid 投稿与审核流里的 AdminPage 起点。
    const loadVideos = () => {
      AdminAPI.listPendingVideos({
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
    // pass 对应后台审核“通过”动作，最终会进入后端 admin.Repository.Review 事务。
    const pass = (videoId) => {
      AdminAPI.approveVideo(videoId).then(() => {
        ElNotification({
          title: "通过",
          message: h("p", { style: "color: green" }, "审核通过"),
          type: 'success'
        });
        loadVideos();
      });
    };
    // notPass 对应后台审核“驳回”动作，会把拒绝原因一起提交到后端。
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
    let value = ref("");
    return {
      submitForm,
      isOpenDef,
      openDef,
      options,
      value,
      pass,
      notPass,
      watchVideo,
      url,
      dialogVisible,
      videos,
      flag,
      formatDate,
    };
  },
};
</script>

<style lang="less" scoped>
* {
  color: #505050;
}

.selecter {
  position: relative;
  margin-bottom: 20px;
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
