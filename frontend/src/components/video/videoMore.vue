<template>
  <Video :height="'500px'" :width="'736px'" :url="videoUrl">
    <template v-slot:footer>
      <div class="video-footer">
        <div class="video-footer-child" v-for="(item, index) in childList" :key="item.id">
          <el-progress type="circle" :percentage="0" v-if="item.status" :stroke-width="2" :width="40">
            <img :src="item.src2" @click="change(index)" />
            <!-- <div>{{ item.num }}</div> -->
          </el-progress>
          <el-progress type="circle" :percentage="sanlianNum" v-if="!item.status" :stroke-width="2" :width="40"
            :indeterminate="true" @click.stop="change(index)">
            <img :src="item.src1" @mousedown.stop="sanlian" @mouseup.stop="cancel" />
            <!-- <div>{{ item.num }}</div> -->
          </el-progress>
          <!-- {{item}} -->
          <div style="margin-left:10px;">{{ item.num }}</div>
        </div>
      </div>
    </template>
  </Video>
  <!-- <button @click="qq">aa</button> -->
</template>

<script>
import Video from "@/components/video/video.vue";
import { h } from 'vue'
import { ElNotification } from "element-plus";
import { HistoryAPI, InteractAPI, VideoAPI } from "@/api/index";
export default {
  // created 是播放区互动组件的初始化入口，对应文档“播放区互动：点赞、收藏与历史上报”。
  // 这里会拉详情回填 viewer_state，并在已登录时主动上报一次观看历史。
  created() {
    this.videoId = this.$route.query.videoId;
    VideoAPI.getVideoDetail(this.videoId).then(data => {
      if (data) {
        this.videoUrl = data.play_url || '';
        this.childList[0].num = data.like_count !== undefined ? data.like_count : (data.videoLikes || 0);
        this.childList[1].num = data.favorite_count !== undefined ? data.favorite_count : (data.videoCollects || 0);
        if (data.viewer_state) {
          this.childList[0].status = data.viewer_state.liked;
          this.childList[1].status = data.viewer_state.favorited;
        }
        if (localStorage.getItem('access_token')) {
          HistoryAPI.report({
            video_id: Number(this.videoId),
            progress_seconds: 0,
          }).catch(() => {});
        }
      }
    });
  },
  components: {
    Video,
  },
  methods: {
    qq() {
      this.sanlianNum += 1;
    },
    // change 是点赞/收藏按钮的统一切换入口。
    // 阅读时可以和后端 video.Service.GetDetail 返回的 viewer_state 一起理解前后端状态如何闭环。
    change(e) {
      if (!localStorage.getItem('userInfo')) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录再进行操作哦！"), type: "error" });
        return;
      }
      if (e === 0) {
        if (this.childList[0].status) {
          InteractAPI.unlikeVideo(this.videoId).then(res => {
            console.log(res);
            if (res) {
              this.childList[0].status = false;
              if (this.childList[0].num > 0) {
                this.childList[0].num--;
              }
            }
          })
        }
        else if (!this.childList[0].status) {
          InteractAPI.likeVideo(this.videoId).then(res => {
            console.log(res);
            if (res) {
              this.childList[0].status = true;
              this.childList[0].num++;
            }
          })
        }
      }
      else if (e === 1) {
        if (this.childList[1].status) {
          InteractAPI.uncollectVideo(this.videoId).then(res => {
            console.log(res);
            if (res) {
              this.childList[1].status = false;
              this.childList[1].num--;
            }
          })
        }
        else if (!this.childList[1].status) {
          InteractAPI.collectVideo(this.videoId).then(res => {
            console.log(res);
            if (res) {
              this.childList[1].status = true;
              this.childList[1].num++;
            }
          })
        }
      }
      this.sanlianNum = 0;
    },
    sanlian() {
      if (!localStorage.getItem('userInfo')) {
        ElNotification({
          title: "错误",
          message: h("p", { style: "color: red" }, "请先登录再进行操作哦！"),
          type: "error",
        });
        return;
      }
      clearInterval(this.timer);
      this.timeStart = new Date().getTime();
      this.timeEnd;
      let that = this;
      this.timer = setInterval(function () {
        that.timeEnd = new Date().getTime();
        if (that.sanlianNum >= 100 && that.timeEnd - that.timeStart >= 3000) {
          InteractAPI.likeVideo(that.videoId).then(res => {
            console.log(res);
            // alert(111)
            if (res) {
              that.childList[0].status = true;
              that.childList[0].num++;
            }
          })
          InteractAPI.collectVideo(that.videoId).then(res => {
            console.log(res);
            // alert(222)
            if (res) {

              that.childList[1].status = true;
              that.childList[1].num++;
            }
          })
          for (let i of that.childList) {
            i.status = true;
          }
          that.sanlianNum = 0;
          clearInterval(that.timer);
        }
        else {
          if (that.sanlianNum < 100)
            that.sanlianNum += 1.2;
        }
      }, 30);
    },
    cancel() {
      clearInterval(this.timer);
      let that = this;
      this.timer = setInterval(function () {
        that.timeEnd = new Date().getTime();
        if (that.sanlianNum > 0) {
          that.sanlianNum -= 1;
        }
        else {
          clearInterval(that.timer);
          that.sanlianNum = 0;
          // for(let i of that.childList){
          //   i.status=false;
          // }
          this.sanlianNum = 0;
        }
      }, 30);

    }
  },
  data() {
    return {
      videoId: '',
      videoUrl: '',
      timer: null,
      childList: [
        {
          num: 0,
          src1: require("@/assets/image/good.svg"),
          src2: require("@/assets/image/good-fill.svg"),
          status: false,
          id: 1,
        },
        {
          num: 0,
          src1: require("@/assets/image/collection.svg"),
          src2: require("@/assets/image/collection-fill.svg"),
          status: false,
          id: 2,
        },
      ],
      sanlianNum: 0,
      timeStart: null,
      timeEnd: null
    };
  },
};
</script>
<style scoped lang="less">
.video-footer {
  display: flex;
  width: 736px;
  margin-left: 20px;
}

.video-footer-child {
  display: flex;
  padding: 10px 5px;
  align-items: center;
  margin-right: 30px;

  img:hover {
    cursor: pointer;
    color: rgb(0, 174, 236);
    transition: all 0.5s
  }

  img {
    transition: all 0.5s;
    height: 30px;
    width: 30px;
    margin-right: 10px;
    // margin-right: 10px;
  }

  div {
    font-size: 18px;
    margin-right: 15px;
  }

  .circle {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
