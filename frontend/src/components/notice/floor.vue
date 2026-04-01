<template>
  <div class="floor">
    <div class="read">
      <a
        href="#"
        style="color:#498ee7;text-decoration: none;outline: none;"
        @click.prevent="read"
      >{{ item.read ? '已读' : '标记为已读' }}</a>
    </div>
    <div class="card-box">
      <div class="card-top">
        <span class="title">{{ item.title }}</span>
        <span class="time">{{ formatDate(item.created_at) }}</span>
      </div>
      <div class="card-body">
        <span>{{ item.content }}</span>
      </div>
    </div>
  </div>
</template>

<script>
import { NoticeAPI } from '@/api/index';
export default {
  name: "Floor",
  props: {
    item: { type: Object },
  },
  emits: ["read"],
  methods: {
    async read() {
      if (this.item.read) {
        return;
      }

      try {
        await NoticeAPI.markRead(this.item.id);
        this.$emit("read", this.item.id);
      } catch (error) {
        console.error("mark notice read failed", error);
      }
    },
    formatDate(value) {
      return value ? value.slice(0, 10) : "";
    },
  },
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
