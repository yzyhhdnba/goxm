<template>
  <div class="card-list">

    <head class="card-title">
      <div class="c-t-left">
        <div style="float:left;">
          <img src="@/assets/image/bilibili.svg" alt="" style="width: 32px;margin-right:5px" />
        </div>
        <span>{{ mdnameItem }}</span>

      </div>
      <div class="c-t-right">
        <div class="btn-change">
          <el-button text :icon="RefreshLeft" @click="change">换一换</el-button>
        </div>
        <div class="more">
          <a href="#">
            <el-button text :icon="CaretRight">更多</el-button>
          </a>
        </div>
      </div>
    </head>

    <div class="card-content">
      <div class="c-c-item" v-for="(item, i) in videoList" :key="i" @click="jumpPath(item.title)">
        <one-card :item="item" />
      </div>
    </div>
  </div>
</template>

<script>
import OneCard from "./OneCard";
import { RefreshLeft, CaretRight } from '@element-plus/icons-vue';
import { ref } from 'vue'
export default {
  name: "CardList",
  props: {
    maindataItem: {
      type: Array,
      default() {
        return [];
      },
    },
    mdnameItem: {
      type: String,
      default: "",
    },

  },
  components: {
    OneCard,
  },
  setup(props) {
    let videoList=ref(props.maindataItem.slice(0, 5))
    const change = async () => {
      try {
        const { VideoAPI } = require('@/api/index');
        const res = await VideoAPI.getRecommend({ limit: 5, cursor: Math.random().toString() });
        if (res && res.items) {
          videoList.value = res.items;
        }
      } catch (err) {
        console.error(err);
      }
    }
    return {
      RefreshLeft,
      CaretRight,
      videoList,
      change
    }
  },
  methods: {

  }
}
</script>

<style scoped>
a {
  text-decoration: none;
  outline: none;
  color: #000;
}


.card-list {
  width: 100%;
  font: 12px Helvetica Neue, Helvetica, Arial, Microsoft Yahei, Hiragino Sans GB,
    Heiti SC, WenQuanYi Micro Hei, sans-serif;
}

.card-title {
  padding: 10px;
  display: flex;
  justify-content: space-between;
}

.c-t-left {
  display: flex;
  align-items: center;
  font-size: 20px;
  cursor: pointer;
}

.c-t-right {
  display: flex;
  justify-content: space-around;
  align-items: center;
  text-align: center;
  width: 10%;
  font-size: 12px;
  margin-right: 23px;
}

.c-t-right .btn-change {
  flex: 1;
  margin-right: 5px;
  padding: 2px 5px;
  border-radius: 5px;
  border: 1px solid #ccc;
  cursor: pointer;
}

.c-t-right .more {
  padding: 2px 5px;
  border-radius: 5px;
  border: 1px solid #ccc;
}

.btn-change:hover {
  background: #f5f5f5;
}

.card-content {
  display: flex;
  justify-content: space-between;
  flex-wrap: wrap;
}

.c-c-item {
  width: 19%;
  margin-bottom: 20px;
}
</style>
