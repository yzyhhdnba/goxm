<template>
  <div class="top">
  </div>
  <div class="container">
    <div class="login-logo">
      <div class="img">
        <img src="../../assets/image/logo.svg" />
      </div>
      <div class="title">PiliPli后台系统</div>
    </div>
    <Login/>
  </div>
  <footer>
    <a v-for="(item, index) of footerData" :key="index" :href="index.url">
      <span>|</span>
      {{ item.text }}
    </a>
  </footer>
</template>

<script setup>
import Login from '@/components/management/login.vue'
import { useRouter } from 'vue-router'
import { ref, reactive } from 'vue'
import { onMounted, onUnmounted } from 'vue'
import managementBg from '@/assets/bg.jpg'
let isLoading = ref(false);
let phoneNumber = ref('');
let show = ref(false);
const open = () => {
  show.value = !show.value;
}
const footerData = reactive([{
  text: '隐私政策',
  url: ''
}, {
  text: '条款与条件',
  url: ''
}, {
  text: '版权所有 © 2022 PiliPili. 保留所有权利。',
  url: ''
},
])
onMounted(() => {
  document.body.style.backgroundImage = `url('${managementBg}')`;
  document.body.style.backgroundSize = 'cover';
  document.body.style.backgroundPosition = 'center center';
  document.body.style.margin = '0 0 0 0';
  document.body.style.padding = '0 0 0 0';
})
onUnmounted(() => {
  document.body.style.backgroundImage = '';
  document.body.style.backgroundSize = '';
  document.body.style.backgroundPosition = '';
})

const router = useRouter()
const toLogin=() => {
  router.replace('/')
}
const sendUser = () => {
  console.log(phoneNumber.value);
  if (phoneNumber.value.length < 3 || phoneNumber.value.length > 12) {
    alert('用户名为3-11个字符')
    return;
  }
  else {
    isLoading.value = true;
  }
}

</script>

<style scoped lang="less">
* {
  color: #494949;
  font-family: SF Pro Display, SF Pro Icons, Helvetica Neue, Helvetica, Arial,
    sans-serif;
}
.top {
  height: 50px;
  background-color: #323232;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  color: white !important;
  img {
    position: relative;
    top: -5px;
    left: 20px;
    width: 45px;
    height: 45px;
    cursor: pointer;
  }
  .top-right {
    position: relative;
    right: 20px;
  }
}
footer {
  display: flex;
  justify-content: center;
  position: absolute;
  bottom: 0px;
  height: 50px;
  width: 80%;
  left: 50%;
  transform: translate(-50%, 0);
  color: rgba(0, 0, 0, 0.3);
  line-height: 1.8;
  text-align: center;
  font-size: 14px;
  span {
    margin: 0 10px;
  }
}
.login-logo {
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
  img {
    width: 200px;
    height: 200px;
    overflow: hidden;
  }
  .img {
    width: 100px;
    height: 100px;
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative;
    left: 50px;
    background-color: #fff;
    border-radius: 20px;
    box-shadow: 0 0 8px 8px rgba(0, 0, 0, 0.1);
    margin-bottom: 10px;
  }
}
.find-password {
  margin-top: 20px;
  color: #0070c9;
  font-size: 14px;
  cursor: pointer;
  font-size: 18px;
}
.container {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 30%;
  flex-direction: column;
  margin: 0 auto;
  position: relative;
  top: 100px;
  padding-right: 50px;
  padding-top: 70px;
  padding-bottom: 70px;
  // filter:blur(50.6px);
}
.title {
  display: flex;
  justify-content: space-around;
  font-size: 25px;
  font-weight: bold;
  width: 65%;
  margin: 10px 0;
  padding-left: 100px;
  font-weight: 400;
  margin-bottom: 20px;
}
</style>
