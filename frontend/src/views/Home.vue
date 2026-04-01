<template>
  <div ref="shao"></div>
  <Login :isopen="isopen" @loginSuccess="loginSuccess"></Login>
  <div class="video-layer">
    <video loop src="../assets/video.webm" autoplay muted data-height="180" data-width="1920"></video>
  </div>
  <transition>
    <div ref="header1" v-show="isheader1">
      <Header @openChange="openChange" :islogin='islogin' />
    </div>
  </transition>

  <transition>
    <div ref="header2" v-show="!isheader1">
      <HeaderFix @openChange="openChange"></HeaderFix>
    </div>
  </transition>
  <div style="height:178px" v-show="!isheader1"></div>
  <suspense>
    <Body/>
  </suspense>
  <!-- <button @click="qq">aa</button> -->
  <div style="margin:40px;"></div>
</template>

<script>
import Login from '../components/header/login.vue'
import Video from "../components/video/video.vue";
import Header from "../components/header/Header.vue";
import HeaderFix from "@/components/header/HeaderFix.vue";
import Body from '@/components/body/Body.vue';
export default {
  created(){
    // location.reload()
  },
  mounted() {
    window.addEventListener('scroll', this.scrollHeader, true);
  },
  unmounted() {
    window.removeEventListener('scroll', this.scrollHeader, true);
  },
  components: {
    Video,
    Header,
    HeaderFix,
    Login,Body
  },
  methods: {
    loginSuccess(v) {
      this.islogin.islogin = v;
    },
    openChange(value) {

      this.isopen.dialogFormVisible = value;
      console.log(this.isopen);
    },
    scrollHeader() {
      // console.log(1,this.$refs.header1.getBoundingClientRect().top);
      // console.log(2,this.$refs.header2.getBoundingClientRect().top);
      if (this.isheader1) {
        if (this.$refs.shao.getBoundingClientRect().top < 0) {
          this.isheader1 = false;
        }
      }
      else {
        if (this.$refs.shao.getBoundingClientRect().top == 0) {
          this.isheader1 = true;
        }
      }

      // else if(this.$refs.header1.getBoundingClientRect().top >= 0){
      //   this.isheader1=true;
      // }
    },
    change(e) {
      this.childList[e].status = !this.childList[e].status;
      this.sanlianNum = 0;
      // console.log(this.childList);
    },
    sanlian() {
      clearInterval(this.timer);
      this.timeStart = new Date().getTime();
      this.timeEnd;
      let that = this;
      this.timer = setInterval(function () {
        that.timeEnd = new Date().getTime();
        if (that.sanlianNum >= 100 && that.timeEnd - that.timeStart >= 3000) {
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
      islogin: { islogin: false },
      isopen: {
        dialogFormVisible: false
      },
      isheader1: true,
      timer: null,
      childList: [
        {
          num: 0,
          src1: require("../assets/image/good.svg"),
          src2: require("../assets/image/good-fill.svg"),
          status: false,
          id: 1,
        },
        {
          num: 0,
          src1: require("../assets/image/collection.svg"),
          src2: require("../assets/image/collection-fill.svg"),
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
.header1 {
  display: none;
  transition: all 0.5s;
}

.video-footer {
  display: flex;
  width: 736px;
  margin-left: 20px;
}

.video-layer {
  position: absolute;
  top: -40%;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.video-footer-child {
  display: flex;
  padding: 12px 0;
  align-items: center;
  margin-right: 30px;

  img:hover {
    cursor: pointer;
    color: rgb(0, 174, 236);
    transition: all 0.5s
  }

  img {
    transition: all 0.5s;
    height: 40px;
    width: 40px;
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