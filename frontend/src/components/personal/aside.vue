<template>
  <div class="personal-aside">
    <!--
    personal-aside类名
    -->
    <el-space wrap>
      <el-card class="enter-ceration">
        <el-container>
          <router-link to="/upload" class="img-creation">
            <img class="creation" width="200" :src="imgcreation" />
          </router-link>
        </el-container>
      </el-card>
    </el-space>

    <el-card class="attention-boxcard">
      <el-col>
        <span class="person-master" @click="to('/')"><b>个人资料</b></span>
        <span class="info-command">pid</span>
        <span class="info-value">{{ pid }}</span>
      </el-col>
      <el-col>
        <span class="info-command">播放数</span>
        <span class="info-value">{{ numplay }}</span>
      </el-col>
      <!-- <el-col>
        <div class="right">
          <div class="info">
            <a
              href="/member/feeds/following"
              class=""
              log-item='{"click":"CLICK_PROFILE","position":"关注","module":"PROFILE"}'
              ><per>{{ numattention }}</per>
              <p>关注</p></a
            >
            <a
              href="/member/feeds/fans"
              class=""
              log-item='{"click":"CLICK_PROFILE","position":"粉丝","module":"PROFILE"}'
              ><per>{{ numfan }}</per>
              <p>粉丝</p></a
            >
            <a
              rel="noopener"
              log-item='{"click":"CLICK_PROFILE","position":"播放数","module":"PROFILE"}'
              target="_blank"
              href="//member.acfun.cn/video-history"
              ><per>{{ numplay }}</per>
              <p>播放数</p></a
            >
          </div>
          <div> -->
      <el-button
        text
        @click="opendarker"
        class="sign-in-btn ac-button ac-button-default ac-button-normal"
      >
        <span class="ac-icon" style="margin-right: 4px"
          ><img
            src="//ali-imgs.acfun.cn/kos/nlav10360/static/img/sign.cf4133ff.svg"
            class="icon-img"
        /></span>
        <span class=""> 签到得闪电 </span>
        <span class="signin-show-banana" style="display: none"
          ><span class="ac-icon"
            ><img
              src="//ali-imgs.acfun.cn/kos/nlav10360/static/img/sign_banana.3bbae706.svg"
              class="icon-img"
          /></span>
          +3
        </span>
      </el-button>

      <!---->
      <!-- </div>
        </div> -->
      <!-- </el-col> -->
    </el-card>

    <el-card class="enter-pilipili">
      <el-container>
        <div class="sectione" style="">
          <div class="elec-action">
            <div
              text
              @click="dialogVisible = true"
              type="button"
              class="elec-trigger"
            >
              <span class="elec-trigger-icon"></span>为TA充电
            </div>
            <div class="elec-map">
              <div class="elec-status">
                共 <span class="elec-count">{{ numpep }}</span
                >人为TA充电
              </div>
            </div>
          </div>
        </div>
      </el-container>
    </el-card>
    <el-dialog
      v-model="dialogVisible"
      title="充电UP"
      width="30%"
    >
      <img width="216" height="336" :src="imgpay" />
      <template #footer>
        <span class="dialog-footer">
          <!-- <el-button @click="dialogVisible = false">Cancel</el-button> -->
          <el-button type="primary" @click="dialogVisible = false"
            >( ͡° ͜ʖ ͡°)✧充电完成！</el-button
          >
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessageBox } from "element-plus";

export default {
  name: "PersonalAside",
  props: {
    dashboard: {
      type: Object,
      default: null,
    },
  },
  setup(props) {
    let router = useRouter();
    const to = (url) => {
      router.push(url);
    };
    const opendarker = () => {
      ElMessageBox.alert("签到成功✧*｡٩(ˊᗜˋ*)و✧*", "来自世界最高城的提示", {
        confirmButtonText: "ＯＫ(ゝω・´★)",
      });
    };
    const dialogVisible = ref(false);
    const user = computed(() => props.dashboard?.user || JSON.parse(localStorage.getItem("userInfo") || "null") || {});
    const stats = computed(() => props.dashboard?.stats || {});

    return {
      opendarker,
      dialogVisible,
      to,
      imgcreation:
        "https://s1.hdslb.com/bfs/static/jinkela/space/assets/icon_createCenters.png",
      pid: computed(() => user.value.id || "未登录"),
      numpep: computed(() => user.value.follower_count || 0),
      numplay: computed(() => stats.value.total_view_count || 0),
      imgpay: "https://s3.bmp.ovh/imgs/2022/06/28/7cf90dc29615a197.jpg",
    };
  },
};
</script>

<style scoped>
.enter-ceration {
  width: 290px;
  height: 110px;
  margin-bottom: 5px;
  margin-top: 5px;
  border-radius: 30px;
}

.attention-boxcard {
  width: 290px;
  height: 200px;
  border-radius: 30px;
}

.enter-pilipili {
  width: 290px;
  height: 200px;
  border-radius: 30px;
}

.img-creation {
  margin: 0;
  padding-left: 30px;
}
.img-creation:hover,
.img-creation:focus {
  color: #66ccff;
  text-decoration: underline;
}

.person-master {
  font-size: 20px;
  font-weight: 700;
  padding-right: 150px;
  padding-bottom: 30px;
}

.info-command {
  display: inline-block;
  min-width: 24px;
  font-size: 18px;
  font-family: Microsoft YaHei;
  line-height: 16px;
  color: #9499a0;
  margin-right: 6px;
}

.info-value {
  color: rgb(51, 51, 51);
  font-size: 18px;
}

/* .person-informantion {
  min-width: 24px;
  font-size: 16px;
  font-family: Microsoft YaHei;
  line-height: 30px;
  color: #9499a0;
  margin-right: 6px;
}

.right {
  color: #666;
  -webkit-font-smoothing: antialiased;
  margin: 0;
  padding: 0;
  border: 0;
  font: inherit;
  font-size: 100%;
  box-sizing: border-box;
  flex: none;
}

.info {
  color: #666;
  -webkit-font-smoothing: antialiased;
  width: 250px;
  font: inherit;
  font-size: 100%;
  box-sizing: border-box;
  display: flex;
  margin-top: 20px;
} */

/* a {
  -webkit-font-smoothing: antialiased;
  margin: 0;
  border: 0;
  font: inherit;
  font-size: 100%;
  outline: none;
  text-decoration: none;
  color: #333;
  box-sizing: border-box;
  width: 80px;
  padding: 3px 0 40px;
}

per:hover,
per:focus {
  color: #66ccff;
}

p,
per {
  -webkit-font-smoothing: antialiased;
  color: #333;
  margin: 0;
  padding: 0;
  border: 0;
  font: inherit;
  box-sizing: border-box;
  font-family: PingFangSC, PingFangSC-Regular;
  font-weight: 500;
  text-align: center;
  height: 20px;
  width: 80px;
  font-size: 18px;
  line-height: 18px;
} */

.sign-in-btn ac-button ac-button-default ac-button-normal {
  -webkit-font-smoothing: antialiased;
  padding: 3px 10px;
  font-size: 14px;
  line-height: 14px;
  cursor: pointer;
  position: relative;
  font-weight: 400;
  white-space: nowrap;
  text-align: center;
  color: #999;
  box-sizing: border-box;
  width: 210px;
  height: 52px;
  border-radius: 5px;
  outline: none;
  border: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #f5f5f5;
  background-color: #e5e5e5;
}

span {
  -webkit-font-smoothing: antialiased;
  cursor: pointer;
  white-space: nowrap;
  margin: 0;
  padding: 0;
  border: 0;
  font: inherit;
  box-sizing: border-box;
  font-size: 16px;
  font-family: PingFangSC, PingFangSC-Regular;
  font-weight: 400;
  text-align: left;
  color: #333;
  line-height: 24px;
}

.sectione {
  color: #222;
  margin: center;
  width: 290px;
  font-size: 12px;
  line-height: 1.7em;
  font-family: Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif;
  margin: 0;
  vertical-align: baseline;
  word-break: break-word;
  position: relative;
  background: #fff;
  border: 1px solid #eee;
  border-radius: 4px;
  padding: 15px 20px 18px;
  margin-bottom: 10px;
  padding-top: 15px;
}

.elec-trigger {
  font-family: Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif;
  margin: 0;
  border: 0;
  vertical-align: baseline;
  word-break: break-word;
  background: #f25d8e;
  border-radius: 4px;
  box-shadow: 0 4px 4px rgba(255, 112, 159, 0.3);
  color: #fff;
  cursor: pointer;
  display: inline-block;
  font-size: 18px;
  line-height: 50px;
  padding: 0 24px;
}

.elec-trigger-icon {
  font-family: Hiragino Sans GB, Microsoft YaHei, Arial, sans-serif;
  color: #fff;
  cursor: pointer;
  font-size: 18px;
  line-height: 50px;
  margin: 0;
  padding: 0;
  border: 0;
  word-break: break-word;
  display: inline-block;
  background-image: url(//s1.hdslb.com/bfs/static/jinkela/space/assets/icons.png);
  background-repeat: no-repeat;
  background-position: -278px -918px;
  width: 20px;
  height: 24px;
  vertical-align: middle;
  margin-right: 12px;
}
</style>
*/
