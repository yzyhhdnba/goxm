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
import axios from 'axios';
import { h } from 'vue'
import { ElNotification } from "element-plus";
export default {
  created() {
    this.videoUrl = 'http://172.20.10.6:8081/video/processed/video-' + this.$route.query.videoId + '/ts/index.m3u8';
    console.log(this.videoUrl);
    console.log(111, this.$route.query);
    if (!localStorage.getItem('userInfo')) {
      return;
    }
    console.log(localStorage.getItem('userInfo'));
    let { userId } = JSON.parse(localStorage.getItem('userInfo'))
    this.videoId = this.$route.query.videoId;
     axios.get('http://172.20.10.6:8081/video', {
      params: {
        videoId: this.$route.query.videoId,
        userId
      }
    }).then(res => {
      console.log(res);
      this.childList[0].num = res.data.videoLikes;
      this.childList[1].num = res.data.videoCollects;
    })
    axios.get('http://172.20.10.6:8081/video/liked', {
      params: {
        userId,
        videoId: this.$route.query.videoId
      }
    }).then(res => {
      console.log(res);
      if (res.data.status === 1) {
        this.childList[0].status = true;
      }
    })
    axios.get('http://172.20.10.6:8081/video/collected', {
      params: {
        userId,
        videoId: this.$route.query.videoId
      }
    }).then(res => {
      console.log(res);
      if (res.data.status === 1) {
        this.childList[1].status = true;
      }
    })
  },
  components: {
    Video,
  },
  methods: {
    qq() {
      this.sanlianNum += 1;
    },
    change(e) {
      if (!localStorage.getItem('userInfo')) {
        ElNotification({
          title: "错误",
          message: h("p", { style: "color: red" }, "请先登录再进行操作哦！"),
          type: "error",
        });
      }
      let { userId } = JSON.parse(localStorage.getItem('userInfo'))
      console.log(userId);
      if (e === 0) {
        if (this.childList[0].status) {
          axios.get('http://172.20.10.6:8081/video/dislike', {
            params: {
              videoId: this.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            if (res.data.status === 1) {
              this.childList[0].status = false;
              this.childList[0].num--;
            }
          })
        }
        else if (!this.childList[0].status) {
          axios.get('http://172.20.10.6:8081/video/like', {
            params: {
              videoId: this.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            if (res.data.status === 1) {
              this.childList[0].status = true;
              this.childList[0].num++;
            }
          })
        }
      }
      else if (e === 1) {
        if (this.childList[1].status) {
          axios.get('http://172.20.10.6:8081/video/cancel', {
            params: {
              videoId: this.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            if (res.data.status === 1) {
              this.childList[1].status = false;
              this.childList[1].num--;
            }
          })
        }
        else if (!this.childList[1].status) {
          axios.get('http://172.20.10.6:8081/video/collect', {
            params: {
              videoId: this.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            if (res.data.status === 1) {
              this.childList[1].status = true;
              this.childList[1].num++;
            }
          })
        }
      }
      this.childList[e].status = !this.childList[e].status;
      this.sanlianNum = 0;
      // console.log(this.childList);
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
          let { userId } = JSON.parse(localStorage.getItem('userInfo'))
          axios.get('http://172.20.10.6:8081/video/like', {
            params: {
              videoId: that.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            // alert(111)
            if (res.data.status === 1) {
              that.childList[0].status = true;
              that.childList[0].num++;
            }
          })
          axios.get('http://172.20.10.6:8081/video/collect', {
            params: {
              videoId: that.videoId,
              userId
            }
          }).then(res => {
            console.log(res);
            // alert(222)
            if (res.data.status === 1) {

              that.childList[1].status = true;
              that.childList[1].num++;
            }
          })
          for (let i of that.childList) {
            i.status = true;
          }
          that.sanlianNum = 0;
          clearInterval(that.timer);
          // alert('三连成功')
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