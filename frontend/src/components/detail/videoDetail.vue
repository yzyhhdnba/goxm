<template>
  <div class="box">
    <div style="height: 64px"></div>
    <div class="container">
      <el-container>
        <el-main>
          <div class="main-title">
            <h1 class="title-headline">{{ videoDetail.title || videoDetail.videoTitle }}</h1>
            <div class="title-data">
              <span class="title-item">播放量:{{ videoDetail.view_count !== undefined ? videoDetail.view_count : videoDetail.videoHits }}</span>
              <span class="title-item">评论数:{{ videoDetail.comment_count !== undefined ? videoDetail.comment_count : videoDetail.videoComments }}</span>
              <span class="title-item">发布时间:{{ videoDetail.published_at || videoDetail.videoContribute }}</span>
              <span class="title-item">未经作者授权，禁止转载</span>
            </div>
          </div>
          <!-- 视频位置 -->
          <div style="width: 770px; height: 490px;margin-bottom: 10px;">
            <VideoPlayer> </VideoPlayer>
          </div>
          <!--  -->
          <div class="main-info">
            {{ videoDetail.description || videoDetail.videoContent }}
          </div>
          <div class="main-comments">
            <div class="reco-title">
              评论
              <!-- <span style="font-size: 14px; color: rgb(148, 153, 160)">{{videoDetail}}</span> -->
            </div>
            <div class="comments-body">
              <div class="comments-mine">
                <div class="my-avatar">
                  <el-avatar :size="48" :src="currentUserAvatar" />
                </div>
                <div class="comments-inputs">
                  <textarea v-model="commentContent" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
                </div>
                <div class="comments-botton">
                  <el-button color="#00aeec" style="height: 100%; width: 100%; color: #ffffff"
                    @click="submit">发布
                  </el-button>
                </div>
              </div>
              <div class="comments-main" v-if="flag">
                <Floor v-for="item in comments" :key="item.id || item.commentId" :item="item" />
              </div>
            </div>
          </div>
        </el-main>
        <el-aside width="350px">
          <div class="aside">
            <div class="up">
              <div class="up-avatar">
                <el-avatar :size="48" :src="authorAvatar" />
              </div>
              <div class="up-info">
                <div class="up-name">{{ (videoDetail.author && videoDetail.author.username) || videoUser }}</div>
                <div class="up-intro">reahbera</div>
                <el-row>
                  <el-col :span="8">
                    <el-button style="width: 80%">充电</el-button>
                  </el-col>
                  <el-col :span="16">
                    <el-button color="#00aeec" style="width: 80%; color: #ffffff" @click="follow"
                      v-if="!isFollow">关注</el-button>
                    <el-button type="danger" style="width: 80%;" @click="unfollow" v-if="isFollow">已关注</el-button>
                  </el-col>
                </el-row>
              </div>
            </div>
            <div class="aside-reco">
              <div class="reco-title">推荐视频</div>
              <Recommend v-for="item in posts" :key="item.id || item.videoId" :item="item" />
            </div>
          </div>
        </el-aside>
      </el-container>
    </div>
  </div>
</template>

<script>
import { computed, h, onMounted, ref } from "vue";
import { useRoute } from "vue-router";
import { ElNotification } from "element-plus";
import { InteractAPI, VideoAPI } from "@/api/index";
import Recommend from "./recommend";
import Floor from "./floor";
import VideoPlayer from "@/components/video/videoMore.vue";

const defaultAvatar = require("@/assets/head.jpg");

export default {
  name: "VideoDetail",
  components: {
    Recommend,
    Floor,
    VideoPlayer,
  },
  setup() {
    const route = useRoute();
    const videoId = route.query.videoId;
    const isFollow = ref(false);
    const comments = ref([]);
    const flag = ref(false);
    const posts = ref([]);
    const commentContent = ref("");
    const videoUser = ref("");
    const videoDetail = ref({});
    const followId = ref(null);
    const currentUserAvatar = ref(defaultAvatar);

    const authorAvatar = computed(() => {
      return videoDetail.value?.author?.avatar_url || defaultAvatar;
    });

    const hasAccessToken = () => Boolean(localStorage.getItem("access_token"));

    const notify = (title, message, type) => {
      ElNotification({
        title,
        message: h("p", { style: `color: ${type === "error" ? "red" : "green"}` }, message),
        type,
      });
    };

    const formatDate = (value) => {
      return value ? String(value).slice(0, 10) : "";
    };

    const decorateComment = (item) => {
      return {
        ...item,
        isLike: item?.viewer_state ? item.viewer_state.liked : Boolean(item?.isLike),
      };
    };

    const syncCurrentUserAvatar = () => {
      const raw = localStorage.getItem("userInfo");
      if (!raw) {
        currentUserAvatar.value = defaultAvatar;
        return;
      }
      try {
        const userInfo = JSON.parse(raw);
        currentUserAvatar.value = userInfo.avatar_url || userInfo.avatar || defaultAvatar;
      } catch (error) {
        currentUserAvatar.value = defaultAvatar;
      }
    };

    // loadVideoDetail 对应文档“视频详情页：一页串起详情、评论、推荐、关注”。
    // 核心是优先消费后端返回的 viewer_state，只有缺失时才回退到单独查关注状态。
    const loadVideoDetail = async () => {
      if (!videoId) {
        return;
      }
      const data = await VideoAPI.getVideoDetail(videoId);
      if (!data) {
        return;
      }

      videoUser.value = data?.author?.username || data.userName || "";
      videoDetail.value = {
        ...data,
        published_at: formatDate(data.published_at || data.videoContribute),
      };
      followId.value = data.userId || data?.author?.id || null;

      if (data.viewer_state && typeof data.viewer_state.followed !== "undefined") {
        isFollow.value = Boolean(data.viewer_state.followed);
        return;
      }
      if (!followId.value || !hasAccessToken()) {
        isFollow.value = false;
        return;
      }

      try {
        const followState = await InteractAPI.getFollowStatus(followId.value);
        isFollow.value = Boolean(followState?.is_following || followState?.followed);
      } catch (error) {
        isFollow.value = false;
      }
    };

    // loadComments 负责把评论列表与后端的 viewer_state.liked 对齐，供楼层组件直接渲染。
    const loadComments = async () => {
      if (!videoId) {
        flag.value = true;
        return;
      }
      flag.value = false;
      try {
        const data = await InteractAPI.getComments(videoId, { page: 1, page_size: 20 });
        comments.value = (data?.list || []).map(decorateComment);
      } catch (error) {
        comments.value = [];
      } finally {
        flag.value = true;
      }
    };

    // loadRecommend 负责补齐详情页右侧推荐区域。
    // 它和详情、评论一起在 onMounted 中并发加载，减少页面串行等待。
    const loadRecommend = async () => {
      try {
        const data = await VideoAPI.getRecommend({ limit: 6 });
        posts.value = data?.items || data?.list || [];
      } catch (error) {
        posts.value = [];
      }
    };

    // submit 对应详情页评论写路径：发评论成功后，前端直接更新本地列表和评论数，而不是整页刷新。
    const submit = async () => {
      if (!hasAccessToken()) {
        notify("错误", "请先登录再发评论哦！", "error");
        return;
      }

      const content = commentContent.value.trim();
      if (!content) {
        notify("错误", "评论内容不能为空", "error");
        return;
      }

      const created = await InteractAPI.postComment(videoId, { content });
      comments.value = [decorateComment(created), ...comments.value];
      commentContent.value = "";
      flag.value = true;

      if (typeof videoDetail.value.comment_count === "number") {
        videoDetail.value.comment_count += 1;
      } else if (typeof videoDetail.value.videoComments === "number") {
        videoDetail.value.videoComments += 1;
      }

      notify("成功", "发布评论成功", "success");
    };

    // ensureFollowable 把关注前的前置校验收口，避免 follow / unfollow 各自重复写判空逻辑。
    const ensureFollowable = () => {
      if (!hasAccessToken()) {
        notify("错误", "请先登录后再操作关注", "error");
        return false;
      }
      if (!followId.value) {
        notify("错误", "未找到目标作者", "error");
        return false;
      }
      return true;
    };

    // follow / unfollow 对应详情页作者关系的写路径。
    // 前端在成功后直接更新本地 isFollow，和后端 viewer_state 形成闭环。
    const follow = async () => {
      if (!ensureFollowable()) {
        return;
      }
      await InteractAPI.followUser(followId.value);
      isFollow.value = true;
      notify("成功", "关注成功", "success");
    };

    const unfollow = async () => {
      if (!ensureFollowable()) {
        return;
      }
      await InteractAPI.unfollowUser(followId.value);
      isFollow.value = false;
      notify("成功", "取消关注成功", "success");
    };

    // onMounted 是详情页的总装入口：并发拉详情、评论、推荐，保证首屏主信息尽快就绪。
    onMounted(async () => {
      syncCurrentUserAvatar();
      await Promise.all([loadVideoDetail(), loadComments(), loadRecommend()]);
    });

    return {
      authorAvatar,
      commentContent,
      comments,
      currentUserAvatar,
      flag,
      follow,
      isFollow,
      posts,
      submit,
      unfollow,
      videoDetail,
      videoUser,
    };
  },
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
