<template>
  <div style="display: flex;justify-content: center;">
    <div class="nav">
      <div class="logo" @click="toIndex()">
        <img src="../../assets/image/logo.svg" style='transform: scale(2.5)'/>
      </div>
      <div class="right-menu">
        <div class="search">
          <div class="search-input">
            <input class="header-input" placeholder="PiliPili爱你哦" />
          </div>
          <div class="search-btn">
            <svg
              width="17"
              height="17"
              viewBox="0 0 17 17"
              fill="none"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill-rule="evenodd"
                clip-rule="evenodd"
                d="M16.3451 15.2003C16.6377 15.4915 16.4752 15.772 16.1934 16.0632C16.15 16.1279 16.0958 16.1818 16.0525 16.2249C15.7707 16.473 15.4456 16.624 15.1854 16.3652L11.6848 12.8815C10.4709 13.8198 8.97529 14.3267 7.44714 14.3267C3.62134 14.3267 0.5 11.2314 0.5 7.41337C0.5 3.60616 3.6105 0.5 7.44714 0.5C11.2729 0.5 14.3943 3.59538 14.3943 7.41337C14.3943 8.98802 13.8524 10.5087 12.8661 11.7383L16.3451 15.2003ZM2.13647 7.4026C2.13647 10.3146 4.52083 12.6766 7.43624 12.6766C10.3517 12.6766 12.736 10.3146 12.736 7.4026C12.736 4.49058 10.3517 2.1286 7.43624 2.1286C4.50999 2.1286 2.13647 4.50136 2.13647 7.4026Z"
                fill="currentColor"
              />
            </svg>
          </div>
        </div>
        <div v-for="(item, index) in navMenu" :key="index" class="menu">{{ item.text }}</div>
        <el-avatar :size="50" :src="av" fit="fill" class="avater" @click="toUser()"></el-avatar>
        <el-popconfirm
          title="确定要退出吗？"
          confirmButtonText="确认"
          cancelButtonText="取消"
          @confirm="loginOut"
        >
          <template #reference>
            <a class="menu" style="margin-left: 20px;color:orange;">退出</a>
          </template>
        </el-popconfirm>
      </div>
    </div>
  </div>
  <el-divider></el-divider>
</template>

<script lang="ts">
import { useRouter } from 'vue-router'
import { useStore } from 'vuex';
import { AuthAPI } from '@/api/index';
import defaultAvatar from '@/assets/head.jpg';
export default {
  setup() {
    const navMenu = [{
      text: '主站',
      to: '/'
    }, {
      text: '创作者中心',
      to: ''
    }, {
      text: '个人中心',
      to: ''
    }, {
      text: '联系我们',
      to: ''
    }]
    const router = useRouter();
    const store = useStore();
    const toUser = () => {
      router.push('/personal')
    }
    const av = defaultAvatar;
    function toIndex() {
      router.push('/')
    }
    async function loginOut() {
      try {
        if (localStorage.getItem('access_token')) {
          await AuthAPI.logout();
        }
      } catch (error) {
        console.error('Logout error', error);
      } finally {
        localStorage.removeItem('isLogin')
        localStorage.removeItem('access_token')
        localStorage.removeItem('refresh_token')
        localStorage.removeItem('userInfo')
        router.replace('/managementLogin')
      }
    }
    return {
      navMenu,
      av,
      toIndex,
      loginOut,
      toUser
    }
  }
}

</script>

<style scoped lang="less">
input {
  border: transparent;
  outline: none;
}
.nav {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  top: 10px;
  width: 95%;
  margin: auto 0;
  .logo {
    width: 128px;
    height: 28px;
    margin-left: 20px;
    object-fit: center;
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
    img {
      width: 130px;
      height: 55px;
    }
  }
  .right-menu {
    display: flex;
    align-items: center;
    margin-right: 30px;
    .menu {
      cursor: pointer;
    }
    div {
      margin-left: 30px;
    }
    .avater {
      margin-left: 30px;
    }
  }
}
.search-btn {
  position: relative;
  color: #000000;
  cursor: pointer;
  /* visibility: hidden; */
  z-index: 1999;
  border-radius: 6px;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: 0.5s;
}
.search-btn:hover {
  background-color: #e3e5e7;
  transition: 0.5s;
  cursor: pointer;
}
.search-input {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: center;
  align-items: center;
  border-radius: 6px;
  background-color: #f1f2f3;
  transition: 0.5s;
}
.search {
  display: flex;
  height: 40px;
  border: 1px;
  background-color: #f1f2f3;
  /* padding: 0 48px; */
  border-radius: 8px;
  justify-content: space-between;
  align-items: center;
  margin-right: 30px;
  padding-right: 5px;
  flex: 1;
  transition: 0.5s;
  opacity: 0.8;
}
.search:hover,
.search-input:hover {
  transition: 0.5s;
}
.header-input {
  line-height: 20px;
  width: 90%;
  transition: 0.5s;
  border-radius: 6px;
  margin-left: 5px;
  padding: 5px 0;
  padding-right: 30px;
  margin-right: 3px;
  background-color: #f1f2f3;
}
.avater {
  cursor: pointer;
}
</style>
