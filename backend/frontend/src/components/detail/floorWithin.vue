<template>
  <div class="floor2">
    <div class="floor-avatar">
      <el-avatar :size="36" :src="require('@/assets/head.jpg')" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">
          <span>{{ item.userName }}</span>
          <!-- <el-button @click="showMyComment" type="text">回复</el-button> -->
          <span v-if="item.replyFor">回复{{'   '+item.replyForName}}</span>
        </div>
        <div class="user-content">
          <span>{{ item.replyContent }}</span>
        </div>
        <div class="content-info">
          <span>2020/6/29</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-else><img src="@/assets/image/good.svg"></span>
            <span>{{ item.replyLikes }}</span>
          </div>

          <el-button @click="showToggle" type="text" v-show="isShow">隐藏</el-button>
          <el-button @click="showToggle" type="text" v-show="!isShow">显示{{ 0 }}条评论</el-button>
          <el-button type="text" v-show="isShow" @click="open">回复</el-button>
        </div>
        <div class="content-within">
          <FloorWithinMore v-for="item in item.withinMore" :item="item" v-show="isShow" />
          <div class="comments-mine" v-show="isShowMine">
            <div class="my-avatar">
              <el-avatar :size="24" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
            </div>
            <div class="comments-inputs">
              <textarea v-model="replyText" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
            </div>
            <div class="comments-botton">
              <el-button color="#00aeec" :dark="isDark" style="height: 100%; width: 100%; color: #ffffff"
                @click="submit(item.replyId, item.commentId)">发布
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import FloorWithinMore from "./floorWithinMore";
import axios from "axios";
import { h } from 'vue';
import { ElNotification } from "element-plus";
export default {
  name: "FloorWithin",
  props: {
    item: { type: Object },
  },
  components: {
    FloorWithinMore
  },
  data() {
    return {
      isShow: false,
      isShowMine: false,
      xin: false,
      replyText: ''
    }
  },
  methods: {
    open() {
      this.isShowMine = !this.isShowMine;
    },
    submit(replyFor, commentId) {
      console.log(commentId);
      let { userId } = JSON.parse(localStorage.getItem('userInfo'))
      axios.post('http://172.20.10.6:8081/reply/make', {
        commentId,
        userId,
        replyFor,
        replyContent: this.replyText
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
      })
    },
    showToggle: function () {
      this.isShow = !this.isShow
    },
    showMyComment: function () {
      this.isShowMine = !this.isShowMine
    },
    dianzan: function () {
      // alert("点咯")
      if (!JSON.parse(localStorage.getItem('userInfo'))) {
        ElNotification({
          title: "错误",
          message: h("p", { style: "color: red" }, "请先登录再点赞哦！"),
          type: "error",
        });
        return;
      }
      let { userId } = JSON.parse(localStorage.getItem('userInfo'))

      if (!this.item.isLike) {
        axios.get('http://172.20.10.6:8081/reply/like', {
          params: {
            replyId: this.item.replyId,
            userId
          }
        }).then(res => {
          if (res.data.status === 1) {
            ElNotification({
              title: "成功",
              message: h("p", { style: "color: green" }, "点赞成功！"),
              type: "success",
            });
          } else {
            ElNotification({
              title: "失败",
              message: h("p", { style: "color: green" }, "未知错误！"),
              type: "success",
            });
          }
        })
        this.item.replyLikes++;
        this.item.isLike = !this.item.isLike
      }
      else {
        axios.get('http://172.20.10.6:8081/reply/dislike', {
          params: {
            replyId: this.item.replyId,
            userId
          }
        }).then(res => {
          if (res.data.status === 1) {
            ElNotification({
              title: "成功",
              message: h("p", { style: "color: green" }, "取消点赞成功！"),
              type: "success",
            });
          }
          else {
            ElNotification({
              title: "失败",
              message: h("p", { style: "color: green" }, "未知错误！"),
              type: "success",
            });
          }
        })
        this.item.replyLikes--;
        this.item.isLike = !this.item.isLike
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