<template>
  <div class="floor">
    <div class="floor-avatar">
      <el-avatar :size="48" :src="avatarSrc(item)" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">{{ (item.user && item.user.username) || item.userName }}</div>
        <div class="user-content">
          <span>{{ item.content || item.commentContent }}</span>
        </div>
        <div class="content-info">
          <span>{{ formatDate(item.created_at || item.commentTime) }}</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-if="!item.isLike"><img src="@/assets/image/good.svg"></span>
            <span>{{ item.like_count !== undefined ? item.like_count : item.commentLikes }}</span>
          </div>

          <el-button @click="showToggle" type="text" v-show="isShow">隐藏</el-button>
          <el-button @click="showToggle" type="text" v-show="!isShow">点击显示{{ item.reply_count !== undefined ? item.reply_count : item.commentReplies }}条评论</el-button>
          <el-button @click="open" type="text" v-show="isShow">回复</el-button>
        </div>
        <div class="content-within" v-show="isShow">
          <el-scrollbar max-height="400px">
            <FloorWithin
              v-for="reply in within"
              :key="reply.id || reply.replyId"
              :item="reply"
              v-if="flag"
              @reply-created="handleReplyCreated"
            />
          </el-scrollbar>
          <div class="comments-mine" v-show="isShowMine">
            <div class="my-avatar">
              <el-avatar :size="48" :src="currentUserAvatar" />
            </div>
            <div class="comments-inputs">
              <textarea v-model="replyText" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
            </div>
            <div class="comments-botton">
              <el-button color="#00aeec" style="height: 100%; width: 100%; color: #ffffff"
                @click="submitComment">发布
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import FloorWithin from "./floorWithin";
import { ElNotification } from "element-plus";
import { h } from 'vue';
import { InteractAPI } from "@/api/index";

const defaultAvatar = require("@/assets/head.jpg");

export default {
  name: "Floor",
  props: {
    item: { type: Object, },
  },
  components: {
    FloorWithin,
  },
  // created 对应文档“评论楼层组件：回复列表与评论点赞”。
  // 楼层组件初始化时会同步当前用户头像，并主动拉取该一级评论下的回复列表。
  created() {
    this.syncCurrentUserAvatar();
    this.fetchReplies();
  },
  data() {
    return {
      within: [],
      isShow: false,
      isShowMine: false,
      flag: false,
      xin: false,
      replyText: '',
      currentUserAvatar: defaultAvatar,
    }
  },
  watch: {
    item: {
      deep: true,
      handler() {
        this.fetchReplies();
      },
    },
  },
  methods: {
    hasAccessToken() {
      return Boolean(localStorage.getItem("access_token"));
    },
    syncCurrentUserAvatar() {
      const raw = localStorage.getItem("userInfo");
      if (!raw) {
        this.currentUserAvatar = defaultAvatar;
        return;
      }
      try {
        const userInfo = JSON.parse(raw);
        this.currentUserAvatar = userInfo.avatar_url || userInfo.avatar || defaultAvatar;
      } catch (error) {
        this.currentUserAvatar = defaultAvatar;
      }
    },
    avatarSrc(entry) {
      return entry?.user?.avatar_url || entry?.avatar || defaultAvatar;
    },
    formatDate(value) {
      return value ? String(value).slice(0, 10) : "";
    },
    decorateReply(item) {
      return {
        ...item,
        isLike: item?.viewer_state ? item.viewer_state.liked : Boolean(item?.isLike),
      };
    },
    // fetchReplies 是评论树二级回复的读取入口。
    // 它把后端返回的 viewer_state.liked 转成组件本地的 isLike，降低模板判断复杂度。
    async fetchReplies() {
      const parentId = this.item.id || this.item.commentId;
      this.flag = false;
      if (!parentId) {
        this.within = [];
        this.flag = true;
        return;
      }
      try {
        const res = await InteractAPI.getReplies(parentId, { page: 1, page_size: 10 });
        this.within = (res?.list || []).map(this.decorateReply);
      } catch (error) {
        this.within = [];
      } finally {
        this.flag = true;
      }
    },
    open() {
      this.isShowMine = !this.isShowMine
    },
    // submitComment 对应“回复写路径”：成功后直接把新回复插入本地列表，并同步 reply_count。
    async submitComment() {
      if (!this.hasAccessToken()) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
        return;
      }
      const content = this.replyText.trim();
      if (!content) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "回复内容不能为空"), type: "error" });
        return;
      }
      const targetId = this.item.id || this.item.commentId || this.item.replyId;
      const created = await InteractAPI.postReply(targetId, { content });
      this.within = [...this.within, this.decorateReply(created)];
      this.replyText = "";
      this.isShow = true;
      this.isShowMine = false;
      if (typeof this.item.reply_count === "number") {
        this.item.reply_count += 1;
      } else if (typeof this.item.commentReplies === "number") {
        this.item.commentReplies += 1;
      }
      ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "回复成功"), type: "success" });
    },
    handleReplyCreated(created) {
      this.within = [...this.within, this.decorateReply(created)];
      if (typeof this.item.reply_count === "number") {
        this.item.reply_count += 1;
      } else if (typeof this.item.commentReplies === "number") {
        this.item.commentReplies += 1;
      }
      this.isShow = true;
    },
    showToggle: function () {
      this.isShow = !this.isShow
    },
    showMyComment: function () {
      this.isShowMine = !this.isShowMine
    },
    // dianzan 是评论点赞/取消点赞的本地乐观切换入口。
    // 这里先调接口，再直接更新当前楼层状态，避免整个评论区重新拉取。
    async dianzan() {
    if (!this.hasAccessToken()) {
      ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
      return;
    }
    const targetId = this.item.id || this.item.commentId || this.item.replyId;
      if (!this.item.isLike) {
        await InteractAPI.likeComment(targetId);
        ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "点赞成功！"), type: "success" });
        if (this.item.like_count !== undefined) this.item.like_count++;
        else if (this.item.commentLikes !== undefined) this.item.commentLikes++;
        else if (this.item.replyLikes !== undefined) this.item.replyLikes++;
        this.item.isLike = true;
      } else {
        await InteractAPI.unlikeComment(targetId);
        ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "取消点赞成功！"), type: "success" });
        if (this.item.like_count !== undefined) this.item.like_count--;
        else if (this.item.commentLikes !== undefined) this.item.commentLikes--;
        else if (this.item.replyLikes !== undefined) this.item.replyLikes--;
        this.item.isLike = false;
      }
  },

  }
}
</script>

<style scoped lang="less">
.floor {
  position: relative;
  padding: 22px 0 0 80px;
}

.floor-avatar {
  display: flex;
  justify-content: center;
  position: absolute;
  left: 0;
  width: 80px;
  cursor: pointer;
}

.content {
  flex: 1;
  position: relative;
  padding-bottom: 14px;
  border-bottom: 1px solid #E3E5E7;
}

.user-info {
  align-items: center;
}

.user-name {
  font-weight: 500;
  font-size: 14px;
  color: #FB7299;
  cursor: pointer;
  margin: 15px 5px 15px 0;
}

.user-content {
  position: relative;
  font-size: 15px;
  line-height: 24px;
  color: #000000;
}

.content-info {
  display: flex;
  align-items: center;
  position: relative;
  font-size: 13px;
  color: #9499A0;
  margin-top: 15px;

  img {
    height: 14px;
    width: 14px;
    margin-right: -10px;
  }

  ;

  a {
    margin-right: 20px;
    text-decoration: none;
    outline: none;

  }
}

.content-info span {
  margin-right: 20px;
}

.content-within {
  margin-top: 15px;
}

.comments-mine {
  display: flex;
  flex-direction: row;
  height: 50px;
  margin-top: 15px;
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
</style>
