<template>
  <div>
    <div class="selecter">
      <el-select v-model="value" class="m-2" placeholder="排序方式">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </div>
    <!-- <el-divider /> -->
    <!-- {{dataLists}} -->
    <div v-for="item in dataLists" :key="item.videoId">
      <div class="container" >
        <div class="left">
          <img :src="'http://172.20.10.6:8081/cover/cover-' + item.videoId + '.jpeg'" />
        </div>
        <div class="middle">
          <div class="title">{{ item.videoTitle }}</div>
          <div class="time">{{ item.videoApprove.slice(0, 10) }}</div>
          <div class="buttom">
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/play-fill.svg" />
              </div>
              <div>1</div>
            </div>
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/good.svg" />
              </div>
              <div>{{ item.videoHits }}</div>
            </div>
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/collection.svg" />
              </div>
              <div>{{ item.videoLikes }}</div>
            </div>
          </div>
        </div>
        <div class="right">
          <div>
            <el-button @click="editOpen">编辑视频</el-button>
          </div>
        </div>
      </div>
      <el-divider />
    </div>
  </div>
  <ChangeInfo :isopen="isEditOpen"></ChangeInfo>
</template>

<script>
import { ref, reactive } from 'vue';
// import qs from 'qs';
import ChangeInfo from '@/components/upload/changeInfo.vue';
import axios from 'axios'
export default {
  components: {
    ChangeInfo
  },
  setup() {
    let dataLists = ref([])
    let { userId } = JSON.parse(localStorage.getItem('userInfo'))
    axios({
      url: 'http://172.20.10.6:8081/creator/ToBeApproved',
      method: 'post',
      data: {
        pageNo: 1,
        userId
      },
    }).then(res => {
      console.log(res);
      dataLists.value = res.data;
      console.log(dataLists.value);
    })
    const options = [
      {
        value: "1",
        label: "按发布时间排序",
      },
      {
        value: "2",
        label: "按视频播放量排序",
      },
      {
        value: "3",
        label: "按点赞排序",
      },
      {
        value: "4",
        label: "按收藏排序",
      }
    ]
    let isEditOpen = reactive({
      isopen: false
    });
    const editOpen = () => {
      console.log('open');
      isEditOpen.isopen = true
    }
    let value = ref('')
    return {
      options, value, editOpen, isEditOpen, dataLists
    };
  },
};
</script>

<style lang="less" scoped>
* {
  color: #505050;
}

.selecter {
  position: relative;
  // left: 500px;
  margin-bottom: 20px;
}

.container {
  display: flex;
  width: 1150px;
  align-items: center;

  .left {
    height: 100px;
    width: 155px;
    border-radius: 5px;

    img {
      width: 150px;
      height: 100px;
      overflow: hidden;
      border-radius: 5px;
    }
  }

  .middle {
    margin-left: 40px;
    display: flex;
    justify-content: center;
    flex-direction: column;

    div {
      margin: 6px 0;
    }

    .title {
      font-size: 20px;
      margin-top: 3px;
    }

    .time {
      margin: 0;
      margin-top: 3px;
    }

    .buttom {
      margin: 0;
      display: flex;
      align-items: center;

      .item {
        display: flex;
        margin-right: 30px;
      }

      .icon {
        display: flex;
        align-items: center;
        height: 25px;
        width: 25px;
        margin-right: 5px;

        img {
          height: 25px;
          width: 25px;
        }
      }
    }
  }

  .right {
    flex: 5;
    display: flex;
    justify-content: flex-end;
  }
}
</style>