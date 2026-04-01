<template>
  <div>
    <!-- 图片 -->
    <div class="item-user">
      <el-row>
        <el-col :span="4" style="cursor: pointer">
          <el-avatar :size="100" :src="avatarUrl" @click="toUser(userId)" />
        </el-col>
        <el-col :span="10">
          <el-row class="item-username">
            <a href="#">{{ username }}</a>
          </el-row>
          <el-row class="item-introduce">
            <p>{{ followerCount }}粉丝&nbsp;·&nbsp;{{ videoCount }}个视频</p>
          </el-row>
          <el-row class="item-attention">
            <el-button color="#00aeec" :dark="isDark" style="width: 30%; color: #ffffff" @click="follow"
              v-if="!isFollow">关注</el-button>
            <el-button type="danger" style="width: 30%;" @click="unfollow" v-if="isFollow">已关注</el-button>
          </el-row>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script>
import { computed, ref } from 'vue'
export default {
  props: {
    item: { type: Object, default: () => ({}) },
  },
  components: {},
  setup(props) {
    let isFollow = ref(false)
    const userId = computed(() => props.item.id || props.item.userId || 0);
    const username = computed(() => props.item.username || props.item.userName || '');
    const avatarUrl = computed(() => props.item.avatar_url || props.item.avatarUrl || require('@/assets/head.jpg'));
    const followerCount = computed(() => props.item.follower_count ?? props.item.followerCount ?? 0);
    const videoCount = computed(() => props.item.video_count ?? props.item.videoCount ?? 0);
    return {
      isFollow,
      userId,
      username,
      avatarUrl,
      followerCount,
      videoCount,
    };
  },
  methods: {
    toUser(userId) {

    }
  },
};
</script>

<style scoped>
.item-username {
  height: 21%;
  margin-left: 10px;
  margin-top: 3px;
  margin-bottom: 10px;
  font-size: 17px;
  font-weight: 700;
  /* 两行显示 */
  color: #000000;
}

a {
  text-decoration: none;
  outline: none;
  color: #000;
}

.item-introduce {
  margin-left: 10px;
  font-size: 15px;
  color: #999;
}

.item-attention {
  margin-left: 10px;
  margin-top: 7px;
}
</style>
