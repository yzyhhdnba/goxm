<template>
  <div>
    <div class="selecter">
      <el-select v-model="value" class="m-2" placeholder="排序方式">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </div>
    <!-- <el-divider /> -->
    <div v-if="flag">
      <div v-for="item in videos" :key="item.videoId">
        <div class="container">
          <div class="left">
            <img :src="'http://172.20.10.6:8081/cover/cover-' + item.videoId + '.jpeg'" @click="watchVideo(item.videoId)" style="cursor: pointer" />
          </div>
          <div class="middle">
            <div class="title">{{ item.videoTitle }}</div>
            <div>
              <span style="height: 16px; width: 16px"><img src="@/assets/image/bilibili.svg"
                  style="height: 20px; width: 20px" /></span>
              {{ item.userName }}
            </div>
            <div class="time">{{ item.videoContribute.slice(0, 10) }}</div>
          </div>
          <div class="right">
            <div>
              <el-button @click="pass(item.videoId)">审核通过</el-button>
              <el-button @click="openDef(item.videoId)" type="danger" style="color: white">不通过</el-button>
            </div>
          </div>
        </div>
        <el-divider />
      </div>
    </div>
  </div>
  <el-dialog v-model="dialogVisible" title="视频预览" v-if="videoflag">
    <Video :url="url" ></Video>
  </el-dialog>
  <el-dialog v-model="isOpenDef" title="">
    <el-input v-model="submitForm.denyReson" :autosize="{ minRows: 2, maxRows: 4 }" type="textarea" placeholder="请输入拒绝意见" />
    <el-button @click="notPass()">提交</el-button>
  </el-dialog>
</template>

<script>
import { ref, reactive, h } from "vue";
import Video from "@/components/video/video.vue";
import axios from "axios";
import { ElNotification } from "element-plus";
import { useRouter } from "vue-router";
export default {
  components: {
    Video,
  },
  setup(props) {
    let videoflag=ref(false)
    let flag = ref(false);
    let videos = ref([]);
    axios.get("http://172.20.10.6:8081/admin/video", {
        params: {
          page: 1,
        },
      })
      .then((res) => {
        console.log(res);
        videos.value = res.data;
        console.log(videos.value);
        flag.value = true;
      });
    let dialogVisible = ref(false);
    let submitForm=reactive({
      videoId:'',
      denyReson:''
    })
    let isOpenDef=ref(false)
    const openDef=(id)=>{
      isOpenDef.value=true;
      submitForm.videoId=id;
    }
    const options = [
      {
        value: "1",
        label: "按发布时间排序",
      },
      {
        value: "2",
        label: "按视频播放量排序",
      },
      {
        value: "3",
        label: "按点赞排序",
      },
      {
        value: "4",
        label: "按收藏排序",
      },
    ];
    const router = useRouter();
    const pass = (videoId) => {
      axios.get('http://172.20.10.6:8081/admin/approve', {
        params: {
          videoId
        }
      }).then(res => {
        if (res.data.status === 1) {
          ElNotification({
            title: "通过",
            message: h("p", { style: "color: green" }, "审核通过"),
            type: 'success'
          });
        }
        router.replace('/management')
      })
    };
    const notPass = () => {
      axios.post('http://172.20.10.6:8081/admin/deny', 
        submitForm
      ).then(res => {
        isOpenDef.value=false;
        if (res.data.status === 1) {
          ElNotification({
            title: "提示",
            message: h("p", { style: "color: green" }, "上传成功"),
            type: 'success'
          });
        }
        router.replace('/management')
        // location.reload();
      })
    };
    let url;
    const watchVideo = (id) => {
      url = 'http://172.20.10.6:8081/video/processed/video-' + id + '/ts/index.m3u8';
      videoflag.value=true;
      dialogVisible.value = true;
    };
    console.log(props);
    
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
      videoflag
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
  // left: 500px;
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
    }
  }

  .middle {
    margin-left: 40px;
    display: flex;
    justify-content: space-around;
    flex-direction: column;
    // align-items: space-between;
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

    .buttom {
      margin: 0;
      display: flex;
      align-items: center;

      .item {
        display: flex;
        margin-right: 30px;
      }

      .icon {
        display: flex;
        align-items: center;
        height: 25px;
        width: 25px;
        margin-right: 5px;

        img {
          height: 25px;
          width: 25px;
        }
      }
    }
  }

  .right {
    flex: 5;
    display: flex;
    justify-content: flex-end;
  }
}
</style>