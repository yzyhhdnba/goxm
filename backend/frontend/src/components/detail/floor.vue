<template>
  <div class="floor">
    <div class="floor-avatar">
      <el-avatar :size="48" :src="require('@/assets/head.jpg')" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">{{ item.userName }}</div>
        <div class="user-content">
          <span>{{ item.commentContent }}</span>
        </div>
        <div class="content-info">
          <span>{{ item.commentTime.slice(0, 10) }}</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-if="!item.isLike"><img src="@/assets/image/good.svg"></span>
            <span>{{ item.commentLikes }}</span>
          </div>

          <el-button @click="showToggle" type="text" v-show="isShow">隐藏</el-button>
          <el-button @click="showToggle" type="text" v-show="!isShow">点击显示{{ item.commentReplies }}条评论</el-button>
          <el-button @click="open" type="text" v-show="isShow">回复</el-button>
        </div>
        <div class="content-within" v-show="isShow">
          <el-scrollbar max-height="400px">
            <FloorWithin v-for="item in within" :item="item" v-if='flag' />
          </el-scrollbar>
          <div class="comments-mine" v-show="isShowMine">
            <div class="my-avatar">
              <el-avatar :size="48" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
            </div>
            <div class="comments-inputs">
              <textarea v-model="replyText" placeholder="请输入一条友善的评论" class="inputs1"></textarea>
            </div>
            <div class="comments-botton">
              <el-button color="#00aeec" :dark="isDark" style="height: 100%; width: 100%; color: #ffffff"
                @click="submitComment(item.commentId)">发布
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
import axios from "axios";
import { h } from 'vue';
export default {
  mounted() {

  },
  created() {
    let { userId } = JSON.parse(localStorage.getItem('userInfo'))
    console.log('item', this.item);
    axios.get('http://172.20.10.6:8081/reply', {
      params: {
        commentId: this.item.commentId,
        page: 1,
      }
    }).then(res => {
      console.log(res);
      for (let i = 0; i < res.data.length; i++) {
        axios.get('http://172.20.10.6:8081/reply/liked', {
          params: {
            userId,
            replyId: res.data[i].replyId
          }
        }).then(ress => {
          console.log('reply liked', ress);
          res.data[i].isLike = ress.data.status === 1 ? true : false
          if (i === res.data.length - 1) {
            this.flag = true;
          }
        })
      }
      this.within = res.data
    })
  },
  name: "Floor",
  props: {
    item: { type: Object, },
  },
  components: {
    FloorWithin,
  },
  data() {
    return {
      within: [],
      isShow: false,
      isShowMine: false,
      flag: false,
      xin: false,
      replyText: ''
    }
  },
  methods: {
    open() {
      this.isShowMine = !this.isShowMine
    },
    submitComment(commentId) {
      console.log(commentId);
      let { userId } = JSON.parse(localStorage.getItem('userInfo'))
      axios.post('http://172.20.10.6:8081/reply/make', {
        commentId: this.item.commentId,
        userId,
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
    dianzan() {
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
        axios.get('http://172.20.10.6:8081/comment/like', {
          params: {
            commentId: this.item.commentId,
            userId
          }
        }).then(res => {
          if (res.data.status === 1) {
            ElNotification({
              title: "成功",
              message: h("p", { style: "color: green" }, "点赞成功！"),
              type: "success",
            });
          }
        })
        this.item.commentLikes++;
        this.item.isLike = !this.item.isLike
      }
      else {
        axios.get('http://172.20.10.6:8081/comment/dislike', {
          params: {
            commentId: this.item.commentId,
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
        })
        this.item.commentLikes--;
        this.item.isLike = !this.item.isLike
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