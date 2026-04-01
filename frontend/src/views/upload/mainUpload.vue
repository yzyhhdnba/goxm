<template>
  <div class="upload">
    <div class="video-layer">
      <video loop src="../../assets/video.webm" autoplay muted data-height="180" data-width="1920"></video>
    </div>

    <!-- <div style="margin:20px 0;">
          <Imageupload></Imageupload>
        </div> -->
    <div style="margin:10px 0;"></div>
    <Form v-if="!finished" @finish="finish"></Form>
    <videoUpload v-if="finished && isVideo" @open="open"></videoUpload>
    <imgsUpload v-if="finished && !isVideo"></imgsUpload>
  </div>
</template>

<script>
import videoUpload from '../../components/upload/videoUpload.vue';
import imgsUpload from '../../components/upload/imgsUpload.vue'
import Form from './form.vue'
export default {
  components: {
    videoUpload, Form, imgsUpload
  },
  data() {
    return {
      isVideo: true,
      finished: false,
    }
  },
  methods: {
    finish(val) {
      this.finished = val
    },
    open(val) {
      this.isVideo = false;
    }
  }
}
</script>

<style scoped lang="less">
.upload {
  width: 80%;

}

.video-layer {
  // position: absolute;
  // top: 0;
  // left: 0;
  width: 1100px;
  height: 150px;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  margin-bottom: 20px;
}
</style>