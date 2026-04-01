<template>
  <div class="floor">
    <div class="floor-avatar">
      <el-avatar :size="24" :src="item.avatar" />
    </div>
    <div class="content">
      <div class="user-info">
        <div class="user-name">
          <span>{{ item.name }}</span>
          <el-button @click="showMyComment" type="text">回复</el-button>
        </div>
        <div class="user-content">
          <span>{{ item.text }}</span>
        </div>
        <div class="content-info">
          <span>{{ item.time }}</span>
          <div @click="dianzan" style="cursor: pointer;">
            <span v-if="item.isLike"><img src="@/assets/image/good-fill.svg"></span>
            <span v-else><img src="@/assets/image/good.svg"></span>
            <span>{{ item.replyLikes }}</span>
          </div>
          <el-button @click="showMyComment" type="text">回复</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ElNotification } from "element-plus";
import axios from "axios";
import { h } from 'vue';
export default {
  created() {
    console.log('22222',this.item);
  },
  name: "FloorWithinMore",
  props: {
    item: { type: Object },
  },
  components: {},
  data() {
    return {
  
    }
  },
  methods: {
    dianzan() {
      alert("点咯")
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