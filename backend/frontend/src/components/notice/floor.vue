<template>
  <div class="floor">
    <div class="read">
      <a href="#" style="color:#498ee7;text-decoration: none;outline: none;" @click="read">标记为已读</a>
    </div>
    <div class="card-box">
      <div class="card-top">
        <span class="title">{{ item.noticeTitle }}</span>
        <span class="time">{{ item.noticeTime.slice(0, 10) }}</span>
      </div>
      <div class="card-body">
        <span>{{ item.noticeContent }}</span>
        <!-- <a>{{ item.link }}</a> -->
      </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
export default {
  name: "Floor",
  props: {
    item: { type: Object },
  },
  components: {

  },
  setup(props) {
    let userInfo=JSON.parse(localStorage.getItem('userInfo'));
    let noticeId=props.item.noticeId
    const read = () => {
      axios.get('http://172.20.10.6:8081/user/read',{
        params: {
          userId:userInfo.userId,
          noticeId
        }
      }).then(res=>{
        console.log(res);
      })
    }
    return {
      read
    }
  },
  methods: {

  }
}
</script>

<style scoped>
.floor {
  margin-bottom: 10px;
}

.read {
  float: right;
  margin: 10px 20px 0 0;
}

.card-box {
  padding: 24px 16px;
  background-color: #fff;
  box-shadow: 0 2px 4px 0 rgb(121 146 185 / 54%);
  margin-bottom: 10px;
  border-radius: 4px;
}

.title {
  color: #333;
  font-weight: 700;
}

.time {
  color: #999;
  font-size: 12px;
  line-height: 22px;
  margin: 0 10px;
}

.card-body {
  color: #666;
  padding-left: 8px;
}
</style>