<template>
  <div>
    <div class="selecter">
      <el-select v-model="value" class="m-2" placeholder="排序方式">
        <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value" />
      </el-select>
    </div>
    <!-- <el-divider /> -->
    <!-- {{dataLists}} -->
    <div v-for="item in dataLists" :key="item.id">
      <div class="container" >
        <div class="left">
          <img :src="item.cover_url" />
        </div>
        <div class="middle">
          <div class="title">{{ item.title }}</div>
          <div class="time">{{ formatDate(item.updated_at || item.created_at) }}</div>
          <div class="time" v-if="item.review_reason">驳回原因：{{ item.review_reason }}</div>
          <div class="buttom">
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/play-fill.svg" />
              </div>
              <div>{{ item.view_count }}</div>
            </div>
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/good.svg" />
              </div>
              <div>{{ item.like_count }}</div>
            </div>
            <div class="item">
              <div class="icon">
                <img src="../../assets/image/collection.svg" />
              </div>
              <div>{{ item.favorite_count }}</div>
            </div>
          </div>
        </div>
        <div class="right">
          <div>
            <el-button @click="editOpen(item)">编辑视频</el-button>
          </div>
        </div>
      </div>
      <el-divider />
    </div>
  </div>
  <ChangeInfo v-model:visible="isEditOpen" :video="currentVideo" @updated="loadVideos"></ChangeInfo>
</template>

<script>
import { ref } from 'vue';
import ChangeInfo from '@/components/upload/changeInfo.vue';
import { CreatorAPI } from '@/api/index'
export default {
  components: {
    ChangeInfo
  },
  setup() {
    let dataLists = ref([])
    const loadVideos = () => {
      CreatorAPI.listVideos({
        review_status: 'rejected',
        page: 1,
        page_size: 20,
      }).then(res => {
        dataLists.value = res?.list || [];
      })
    }
    loadVideos()
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
    const isEditOpen = ref(false)
    const currentVideo = ref(null)
    const editOpen = (item) => {
      currentVideo.value = item
      isEditOpen.value = true
    }
    const formatDate = (value) => {
      return value ? value.slice(0, 10) : '';
    }
    let value = ref('')
    return {
      options, value, editOpen, isEditOpen, dataLists, formatDate, currentVideo, loadVideos
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
