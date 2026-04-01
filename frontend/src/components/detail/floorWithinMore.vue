<template>
  <div class="floor">
    <div class="floor-avatar">
      <el-avatar :size="24" :src="avatarSrc(item)" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">
          <span>{{ (item.user && item.user.username) || item.name || item.userName }}</span>
          <el-button @click="showMyComment" type="text">回复</el-button>
        </div>
        <div class="user-content">
          <span>{{ item.content || item.text || item.replyContent }}</span>
        </div>
        <div class="content-info">
          <span>{{ formatDate(item.created_at || item.time || item.replyTime) }}</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-else><img src="@/assets/image/good.svg"></span>
            <span>{{ item.like_count !== undefined ? item.like_count : item.replyLikes }}</span>
          </div>
          <el-button @click="showMyComment" type="text">回复</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ElNotification } from "element-plus";
import { h } from 'vue';
import { InteractAPI } from "@/api/index";

const defaultAvatar = require("@/assets/head.jpg");

export default {
  name: "FloorWithinMore",
  props: {
    item: { type: Object }
  },
  methods: {
    avatarSrc(entry) {
      return entry?.user?.avatar_url || entry?.avatar || defaultAvatar;
    },
    formatDate(value) {
      return value ? String(value).slice(0, 10) : "";
    },
    showMyComment() {
      // display input logics if any
    },
    async dianzan() {
      if (!localStorage.getItem('access_token')) {
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
    }
  }
}
</script>

<style scoped lang="less">
.floor {
  position: relative;
  padding: 22px 0 0 40px;
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
  margin: 0 5px 10px 0;

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
</style>
