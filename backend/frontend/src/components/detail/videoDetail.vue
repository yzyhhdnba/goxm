<template>
  <div class="box">
    <div style="height: 64px"></div>
    <div class="container">
      <el-container>
        <el-main>
          <div class="main-title">
            <h1 class="title-headline">{{ videoDetail.videoTitle }}</h1>
            <div class="title-data">
              <span class="title-item">播放量:{{ videoDetail.videoHits }}</span>
              <span class="title-item">评论数:{{ videoDetail.videoComments }}</span>
              <span class="title-item">发布时间:{{ videoDetail.videoContribute }}</span>
              <span class="title-item">未经作者授权，禁止转载</span>
            </div>
          </div>
          <!-- 视频位置 -->
          <div style="width: 770px; height: 490px;margin-bottom: 10px;">
            <VideoPlayer> </VideoPlayer>
          </div>
          <!--  -->
          <div class="main-info">
            {{ videoDetail.videoContent }}
          </div>
          <div class="main-comments">
            <div class="reco-title">
              评论
              <!-- <span style="font-size: 14px; color: rgb(148, 153, 160)">{{videoDetail}}</span> -->
            </div>
            <div class="comments-body">
              <div class="comments-mine">
                <div class="my-avatar">
                  <el-avatar :size="48" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
                </div>
                <div class="comments-inputs">
                  <textarea v-model="commentContent" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
                </div>
                <div class="comments-botton">
                  <el-button color="#00aeec" :dark="isDark" style="height: 100%; width: 100%; color: #ffffff"
                    @click="submit">发布
                  </el-button>
                </div>
              </div>
              <div class="comments-main" v-if="flag">
                <Floor v-for="item in comments" :item="item" />
              </div>
            </div>
          </div>
        </el-main>
        <el-aside width="350px">
          <div class="aside">
            <div class="up">
              <div class="up-avatar">
                <el-avatar :size="48" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
              </div>
              <div class="up-info">
                <div class="up-name">{{ videoUser }}</div>
                <div class="up-intro">reahbera</div>
                <el-row>
                  <el-col :span="8">
                    <el-button style="width: 80%">充电</el-button>
                  </el-col>
                  <el-col :span="16">
                    <el-button color="#00aeec" :dark="isDark" style="width: 80%; color: #ffffff" @click="follow"
                      v-if="!isFollow">关注</el-button>
                    <el-button type="danger" style="width: 80%;" @click="unfollow" v-if="isFollow">已关注</el-button>
                  </el-col>
                </el-row>
              </div>
            </div>
            <div class="aside-reco">
              <div class="reco-title">推荐视频</div>
              <Recommend v-for="item in posts" :item="item" />
            </div>
          </div>
        </el-aside>
      </el-container>
    </div>
  </div>
</template>

<script>
import Recommend from "./recommend";
import Floor from "./floor";
import VideoPlayer from "@/components/video/videoMore.vue";
import axios from 'axios'
import { reactive, ref, h, toRaw } from "vue";
import { useRouter, useRoute } from 'vue-router'
import { ElNotification } from "element-plus";
export default {
  name: "VideoDetail",
  props: {},
  components: {
    Recommend,
    Floor,
    VideoPlayer,
  },
  setup() {
    let isFollow = ref(false)
    let signal = ref('')
    let comments = ref([])
    let flag = ref(false)
    let router = useRouter();
    let route = useRoute();
    let { videoId } = route.query;
    let commentContent = ref('')
    console.log(videoId);
    let userId = localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).userId : 'null';
    let followId = ''
    let videoUser = ref('')
    axios.get('http://172.20.10.6:8081/video', {
      params: {
        videoId: route.query.videoId,
        userId
      }
    }).then(res => {
      console.log(res);
      videoUser.value = res.data.userName;
      res.data.videoContribute = res.data.videoContribute.slice(0, 10)
      videoDetail.value = res.data;
      followId = res.data.userId;
      axios.get('http://172.20.10.6:8081/followed', {
        params: {
          followId,
          userId
        }
      }).then(res => {
        console.log(res);
        if (res.data.status === 0) {
          isFollow.value = false;
        }
        else if (res.data.status === 1) {
          isFollow.value = true;
        }
      })
    }).catch(err => {

    })
    axios.get('http://172.20.10.6:8081/comment', {
      params: {
        videoId,
        page: 1
      }
    }).then(res => {
      console.log(res);
      for (let i = 0; i < res.data.length; i++) {
        axios.get('http://172.20.10.6:8081/comment/liked', {
          params: {
            userId,
            commentId: res.data[i].commentId
          }
        }).then(ress => {
          console.log('liked', ress);
          res.data[i].isLike = ress.data.status === 1 ? true : false
          if (i === res.data.length - 1) {
            flag.value = true;
          }
        })
      }
      comments.value = res.data;

    })
    let videoDetail = ref({})

    let posts = ref([{
      pic: "https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917531190280.jpg",
      title:
        '【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国 PV2【MCE汉化组】" class="title">【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国PV2【MCE汉化组】',
      up: "夏日幻听MCE",
      hits: "1144",
      comments: "5144",
      duration: "19:19",
    },
    {
      pic: "https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg",
      title:
        '【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国 PV2【MCE汉化组】" class="title">【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国PV2【MCE汉化组】',
      up: "夏日幻听MCE",
      hits: "1144",
      comments: "5144",
      duration: "19:19",
    },
    {
      pic: "https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg",
      title:
        '【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国 PV2【MCE汉化组】" class="title">【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国PV2【MCE汉化组】',
      up: "夏日幻听MCE",
      hits: "1144",
      comments: "5144",
      duration: "19:19",
    },
    {
      pic: "https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg",
      title:
        '【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国 PV2【MCE汉化组】" class="title">【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国PV2【MCE汉化组】',
      up: "夏日幻听MCE",
      hits: "1144",
      comments: "5144",
      duration: "19:19",
    },
    {
      pic: "https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg",
      title:
        '【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国 PV2【MCE汉化组】" class="title">【7月/虚渊玄/早见沙织/中字首发】RWBY 冰雪帝国PV2【MCE汉化组】',
      up: "夏日幻听MCE",
      hits: "1144",
      comments: "5144",
      duration: "19:19",
    },])
    const submit = () => {
      if (!JSON.parse(localStorage.getItem('userInfo'))) {
        ElNotification({
          title: "错误",
          message: h("p", { style: "color: red" }, "请先登录再发评论哦！"),
          type: "error",
        });
        return;
      }
      let { userId } = JSON.parse(localStorage.getItem('userInfo'))
      axios.post('http://172.20.10.6:8081/comment/make', {
        videoId,
        userId,
        commentContent: commentContent.value
      }).then(res => {
        console.log(res);
        if (res.data.status === 1) {
          ElNotification({
            title: "成功",
            message: h("p", { style: "color: green" }, "发布评论成功"),
            type: "success",
          });
          location.reload()
        }
      }).catch(err => {

      })
    }

    const unfollow = () => {
      axios.get('http://172.20.10.6:8081/unfollow', {
        params: {
          followId,
          userId
        }
      }).then(res => {
        console.log(res);
        if (res.data.status === 1) {
          ElNotification({
            title: "成功",
            message: h("p", { style: "color: green" }, "取消关注成功"),
            type: "success",
          });
          isFollow.value = false;
        }
      })
    }
    const follow = () => {
      axios.get('http://172.20.10.6:8081/follow', {
        params: {
          followId,
          userId
        }
      }).then(res => {
        console.log(res);
        if (res.data.status === 1) {
          ElNotification({
            title: "成功",
            message: h("p", { style: "color: green" }, "关注成功"),
            type: "success",
          });
          isFollow.value = true;
        }
      })
    }
    axios.post('http://172.20.10.6:8081/videoList', {
      areaId: 1,
      count: 6
    }).then(res => {
      console.log(res);
      posts.value = res.data;
    }).catch(err => {
      console.log(err);
    })
    return {
      flag,
      videoDetail,
      posts,
      submit,
      comments,
      follow,
      unfollow,
      isFollow,
      videoUser,
      commentContent,
      // comments: [
      //   {
      //     avatar:
      //       "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //     name: "折木-南",
      //     text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //     time: "2022-06-06 20:20",
      //     like: "336",
      //     within: [
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //         withinMore:[
      //           {
      //             avatar:
      //               "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //             name: "折木-南",
      //             text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //             time: "2022-06-06 20:20",
      //             like: "336",
      //           },
      //         ],
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //     ],
      //   },
      //   {
      //     avatar:
      //       "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //     name: "折木-南",
      //     text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //     time: "2022-06-06 20:20",
      //     like: "336",
      //     within: [
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //       {
      //         avatar:
      //           "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png",
      //         name: "折木-南",
      //         text: "up花了很多篇幅提及左翼学运，来反驳部分人的观点——《冰菓》的右翼倾向和反学运倾向。事实上，这完全是不攻自破的（我看来）。真正想要抹去一段历史，很简单的方法就是不再去提及，冰菓在我眼里没有和其他校园四霸乃至无数校园番趋同，就在于冰菓敢于直面一段真实的历史，暗含了作者大量对历史和社会问题的反思，并以大众化的动漫题材去回顾。就像最后千反田说的，我们数十年后或许就忘记了，而从这个角度来说，打破了所谓“不能说的历史”叙事环境，已经是需要非常大的勇气了，也就已经是潜在的承认并且支持那段历史。",
      //         time: "2022-06-06 20:20",
      //         like: "336",
      //       },
      //     ],
      //   },
      // ],
      textarea1: "",
    };
  },
  methods: {},
};
</script>

<style scoped>
.box {
  font-size: 12px;
  line-height: 12px;
  color: #666;
}

.container {
  margin: 0 auto;
  max-width: 1150px;
  min-width: 1080px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  box-sizing: content-box;
}

.el-main {
  --el-main-padding: 0px !important;
}

.main-title {
  height: 106px;
  box-sizing: border-box;
  padding-top: 24px;
}

.title-headline {
  line-height: 28px;
  font-size: 20px;
  font-weight: 500;
  color: #18191c;
  margin-bottom: 6px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.title-data {
  font-size: 13px;
  color: rgb(148, 153, 160);
  text-overflow: ellipsis;
  white-space: nowrap;
  display: flex;
  align-items: center;
  height: 24px;
  line-height: 18px;
}

.title-item {
  display: inline-flex;
  align-items: center;
  margin-right: 12px;
}

.main-info {
  margin: 16px 0;
  position: relative;
  white-space: pre-line;
  font-size: 15px;
  color: #18191c;
  letter-spacing: 0;
  line-height: 24px;
  overflow: hidden;
  padding-bottom: 16px;
  padding-top: 10px;
  border-bottom: 1px solid #e3e5e7;
}

.main-comments {
  margin-top: 30px;
  z-index: 0;
  position: relative;
}

.comments-mine {
  display: flex;
  flex-direction: row;
  height: 50px;
}

.comments-main {
  margin-top: 14px;
  padding-bottom: 100px;
}

.my-avatar {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 80px;
  height: 50px;
}

.comments-inputs {
  position: relative;
  flex: 1;
}

.comments-botton {
  display: flex;
  justify-content: center;
  align-items: center;
  position: relative;
  flex-basis: 80px;
  margin-left: 10px;
  border-radius: 4px;
  cursor: pointer;
}

.inputs1 {
  width: 100%;
  height: 100%;
  padding: 0;
  border: 1px solid #f1f2f3;
  border-radius: 6px;
  font-family: inherit;
  line-height: 38px;
  color: #18191c;
  resize: none;
  outline: none;
}

.aside {
  margin-left: 30px;
  padding-bottom: 20px;
}

.up {
  box-sizing: border-box;
  min-height: 58px;
  margin: 18px 0 10px 0;
  display: flex;
  align-items: center;
}

.up-avatar {
  align-self: flex-start;
  width: 48px;
  height: 48px;
  position: relative;
  flex-shrink: 0;
}

.up-info {
  margin-left: 12px;
  flex: 1;
  overflow: auto;
}

.up-name {
  color: #fb7299;
  height: 22px;
  line-height: 22px;
  font-size: 15px;
  align-items: center;
  font-weight: 500;
}

.up-intro {
  margin: 2px 0 10px 0;
  font-size: 13px;
  line-height: 16px;
  height: 16px;
  color: rgb(148, 153, 160);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.aside-reco {
  margin-top: 20px;
}

.reco-title {
  font-weight: 500;
  font-size: 20px;
  color: #18191c;
  margin-bottom: 15px;
}
</style>