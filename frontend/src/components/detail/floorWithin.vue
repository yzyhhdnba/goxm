<template>
  <div class="floor2">
    <div class="floor-avatar">
      <el-avatar :size="36" :src="avatarSrc(item)" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">
          <span>{{ (item.user && item.user.username) || item.userName }}</span>
        </div>
        <div class="user-content">
          <span>{{ item.content || item.replyContent }}</span>
        </div>
        <div class="content-info">
          <span>{{ formatDate(item.created_at || item.replyTime) }}</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-else><img src="@/assets/image/good.svg"></span>
            <span>{{ item.like_count !== undefined ? item.like_count : item.replyLikes }}</span>
          </div>
          <el-button type="text" @click="open">回复</el-button>
        </div>
        <div class="content-within">
          <div class="comments-mine" v-show="isShowMine">
            <div class="my-avatar">
              <el-avatar :size="24" :src="currentUserAvatar" />
            </div>
            <div class="comments-inputs">
              <textarea v-model="replyText" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
            </div>
            <div class="comments-botton">
              <el-button color="#00aeec" style="height: 100%; width: 100%; color: #ffffff"
                @click="submit">发布
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { h } from 'vue';
import { ElNotification } from "element-plus";
import { InteractAPI } from "@/api/index";

const defaultAvatar = require("@/assets/head.jpg");

export default {
  name: "FloorWithin",
  emits: ["reply-created"],
  props: {
    item: { type: Object },
  },
  created() {
    this.syncCurrentUserAvatar();
  },
  data() {
    return {
      isShowMine: false,
      xin: false,
      replyText: '',
      currentUserAvatar: defaultAvatar,
    }
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
    open() {
      this.isShowMine = !this.isShowMine;
    },
    async submit() {
      if (!this.hasAccessToken()) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
        return;
      }
      const content = this.replyText.trim();
      if (!content) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "回复内容不能为空"), type: "error" });
        return;
      }
      const targetId = this.item.id || this.item.replyId || this.item.commentId;
      const created = await InteractAPI.postReply(targetId, { content });
      this.replyText = "";
      this.isShowMine = false;
      this.$emit("reply-created", created);
      ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "发布评论成功"), type: "success" });
    },
    async dianzan() {
    if (!this.hasAccessToken()) {
      ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录再点赞哦！"), type: "error" });
      return;
    }
    const targetId = this.item.id || this.item.replyId || this.item.commentId;
      if (!this.item.isLike) {
        await InteractAPI.likeComment(targetId);
        ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "点赞成功！"), type: "success" });
        if (this.item.like_count !== undefined) this.item.like_count++;
        else if (this.item.replyLikes !== undefined) this.item.replyLikes++;
        this.item.isLike = true;
      } else {
        await InteractAPI.unlikeComment(targetId);
        ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "取消点赞成功！"), type: "success" });
        if (this.item.like_count !== undefined) this.item.like_count--;
        else if (this.item.replyLikes !== undefined) this.item.replyLikes--;
        this.item.isLike = false;
      }
  },
  }
}
</script>

<style scoped lang="less">
.floor2 {
  position: relative;
  padding: 22px 0 0 60px;
}

.floor-avatar {
  display: flex;
  justify-content: center;
  position: absolute;
  left: 0;
  width: 40px;
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
  margin: 5px 5px 10px 0;

  span {
    margin-right: 10px;
  }
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
    height: 12px;
    width: 12px;
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
  height: 40px;
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
