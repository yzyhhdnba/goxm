<template>
  <div class="main-player" ref="mainPlayer">
    <div id="dplayer" class="play-root"></div>
    <slot name="footer"></slot>
  </div>
</template>

<script>
import Hls from 'hls.js'
import DPlayer from "dplayer";

export default {
  name: "optimizeVideo",
  props:{
    height:String,
    width:String,
    url:String,
  },
  data() {
    return {
      dp: null,
    };
  },
  methods: {
    onPlay() {
      this.dp.play();
    },
  },
  // 一定要在mounted中创建
  mounted() {
    this.$refs.mainPlayer.style.height=this.height;
    this.$refs.mainPlayer.style.width=this.width;
    let url=this.url
    const dp = new DPlayer({
      // 配置参数
      container: document.getElementById("dplayer"),
      
      autoplay: false,
      theme: "#FADFA3",
      loop: true,
      screenshot: true, // 是否允许截图（按钮），点击可以自动将截图下载到本地
      hotkey: true,
      lang: "zh-cn",
      preload: "auto",
      // logo: 'logo.png',
      volume: 0.7,
      video: {
        url: url,
        // url: "http://172.20.10.6:49150/20220621/2-55fca9a44e124af5babb048147bfb8aa/ts/index.m3u8",
        // pic: "dplayer.png",
        type: 'customHls',
        customType: {
            customHls: function (video, player) {
                const hls = new Hls();
                hls.loadSource(video.src);
                hls.attachMedia(video);
            },
        },
      },
    });

    // 禁止右键下载视频
    document.oncontextmenu = new Function("event.returnValue=false;");
    document.onselectstart = new Function("event.returnValue=false;");

    // 修改循环播放显示
    document
      .getElementsByClassName("dplayer-setting-item dplayer-setting-loop")[0]
      .getElementsByClassName("dplayer-label")[0].innerText = "循环播放";
    // 修改倍速播放显示
    document
      .getElementsByClassName("dplayer-setting-item dplayer-setting-speed")[0]
      .getElementsByClassName("dplayer-label")[0].innerText = "播放倍速";
    dp.on("play", function () {
        document.getElementById("mseVideoUrl") .classList.add("dplayer-hide-controller");
    });
    
  },
};
</script>

<style scoped>
.play-root {
  width: 100%;
  height: 450px;
  background-color: coral;
  margin: 0 auto;
}
/* .main-player{
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
} */
/* 在浏览器找到对应的class名称。然后通过穿透对样式进行更改 */

/* 禁止循环播放显示 */
/* .paly-root >>> .dplayer-setting-loop {
    background-color: cyan;
    display: none;
  } */

/* 禁止出现快进多少秒提示 */
/* .play-root /deep/ .dplayer-notice {
    display: none;
  } */

/* 禁止右键自定义列表 */
/* .play-root /deep/ .dplayer-menu-show{
     display: none;
  } */
</style>
