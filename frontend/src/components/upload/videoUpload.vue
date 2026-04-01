<template>
  <!-- <img src="../../assets/image/tiao.jpg"/> -->
  <el-upload class="upload-demo" drag action="#" :before-upload="show()"
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
import { h, defineEmits } from 'vue'
import { useStore } from 'vuex'
import { ElNotification } from "element-plus";
import { UploadAPI } from '@/api/index'
const emit=defineEmits(['open'])
const store = useStore();
// submitUpload 是投稿两段式流程的第二步：读取 Store 中的 videoId，再上传源文件。
const submitUpload = (content: any) => {
  const videoId = store.state.videoId;
  if (!videoId) {
    ElNotification({
      title: "上传失败",
      message: h("p", { style: "color: red" }, "请先填写稿件信息"),
      type: "error",
    });
    return;
  }

  UploadAPI.uploadSource(videoId, content.file).then(() => {
    ElNotification({
      title: "上传成功",
      message: h("p", { style: "color: green" }, "视频上传成功"),
      type: "success",
    });
    emit('open',true)
  }).catch(() => {
      ElNotification({
        title: "上传失败",
        message: h("p", { style: "color: red" }, "视频上传错误"),
        type: "error",
      });
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
