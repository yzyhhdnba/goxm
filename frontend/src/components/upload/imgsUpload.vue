<template>
  <el-upload action="#" list-type="picture-card" :show-file-list="true"
    :http-request="submitUpload" ref="uploadRef" :on-success="success" accept=".jpg">
    <el-icon>
      <Plus />
    </el-icon>
    <template #file="{ file }">
      <div>
        <img class="el-upload-list__item-thumbnail" :src="file.url" alt="" />
        <span class="el-upload-list__item-actions">
          <span class="el-upload-list__item-preview" @click="handlePictureCardPreview(file)">
            <el-icon>
              <zoom-in />
            </el-icon>
          </span>
        </span>
      </div>
    </template>
  </el-upload>
  <el-dialog v-model="dialogVisible">
    <img w-full :src="dialogImageUrl" alt="Preview Image" />
  </el-dialog>
</template>
<script lang="ts">
import { ref, h } from 'vue'
import { Delete, Download, Plus, ZoomIn } from '@element-plus/icons-vue'
import type { UploadFile, UploadInstance } from 'element-plus'
import { useStore } from 'vuex'
import { ElNotification } from "element-plus";
import { UploadAPI } from '@/api/index'
export default {
  setup() {
    const store = useStore();
    const dialogImageUrl = ref('')
    const dialogVisible = ref(false)
    const disabled = ref(false)
    const uploadRef = ref<UploadInstance>()
    const handleRemove = (file: UploadFile) => {
      console.log(file)
    }
    const success = () => {
      console.log('okk');
    }
    const handlePictureCardPreview = (file: UploadFile) => {
      dialogImageUrl.value = file.url!
      dialogVisible.value = true;
    }
    const handleChange = (file: UploadFile) => {

    }
    const handleDownload = (file: UploadFile) => {
      console.log(file)
    }

    const submitUpload = (content: any) => {
      const videoId = store.state.videoId;
      if (!videoId) {
        ElNotification({
          title: "上传失败",
          message: h("p", { style: "color: red" }, "请先上传视频源文件"),
          type: "error",
        });
        return;
      }

      UploadAPI.uploadCover(videoId, content.file).then(() => {
        ElNotification({
          title: "上传成功",
          message: h("p", { style: "color: green" }, "图片上传成功"),
          type: "success",
        });
      }).catch(() => {
          ElNotification({
            title: "上传失败",
            message: h("p", { style: "color: red" }, "图片上传错误"),
            type: "error",
          });
      })
    }
    return {
      dialogImageUrl, dialogVisible, disabled, handleRemove, handlePictureCardPreview, handleDownload, submitUpload, uploadRef, success,
      Delete, Download, Plus, ZoomIn
    }
  }
}

</script>
