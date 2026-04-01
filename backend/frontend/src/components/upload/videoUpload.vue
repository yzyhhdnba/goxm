<template>
  <!-- <img src="../../assets/image/tiao.jpg"/> -->
  <el-upload class="upload-demo" drag action="http://172.20.10.6:8081/upload/Video" :before-upload="show()"
    accept=".mp4" :http-request="submitUpload">
    <el-icon class="el-icon--upload">
      <upload-filled />
    </el-icon>
    <div class="el-upload__text">
      拖入视频到此处或 <em>点此上传</em>
    </div>
    <template #tip>
      <div class="el-upload__tip">
        视频最大限制为200Mb
      </div>
    </template>
  </el-upload>
</template>

<script setup lang="ts">
import { UploadFilled } from '@element-plus/icons-vue'
import { ref, h,toRaw,defineEmits } from 'vue'
import { useStore } from 'vuex'
import { ElNotification } from "element-plus";
import axios from "axios"
const emit=defineEmits(['open'])
let isshow = ref(false);
const store = useStore();
const submitUpload = (content: any) => {
  let videoInfo = toRaw(store.state.videoInfo);
  delete videoInfo.videoId;
  console.log(111,videoInfo);
  
  console.log(content);
  const formData = new FormData()
  formData.append('file', content.file)
  formData.append('videoInfo', new Blob([JSON.stringify({
    videoTitle:videoInfo.videoTitle,
    userId:videoInfo.userId,
    areaId:videoInfo.areaId,
    videoContent:videoInfo.videoContent
    })], { type: "application/json; charset=utf-8" }))
  // uploadRef.value!.submit()
  axios({
    method: 'POST',
    url: 'http://172.20.10.6:8081/upload/Video',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  }).then(res => {
    console.log(res);
    if (res.data.status === 0) {
      ElNotification({
        title: "上传失败",
        message: h("p", { style: "color: red" }, "视频上传错误"),
        type: "error",
      });
    }
    else  {
      ElNotification({
        title: "上传成功",
        message: h("p", { style: "color: green" }, "视频上传成功"),
        type: "success",
      });
      store.commit('addVideoId',res.data.status)
      emit('open',true)
    }

  })
}
let show = () => {

}
</script>
<style lang="less" scoped>
img {
  width: 100%;
  height: 250px;
}
</style>