<template>
  <div class="header">
    <div class="header-video">
      <div class="header-top">
        <ul class="header-top-left">
          <li v-for="(item, index) in leftData" :key="index">
            <a class="entry-title" href="#" @click.prevent="navigate(item)">
              <span v-if="item.sag != ''" v-html="item.sag" class="svg"></span>
              <span class="left-ch">{{ item.text }}</span>
              <Headerch />
            </a>
          </li>
          <li class="download">
            <a class="entry-title" src="#" target="_blank">
              <svg width="17" height="16" viewBox="0 0 17 16" fill="none" xmlns="http://www.w3.org/2000/svg"
                class="download-entry__icon">
                <path
                  d="M14.9752 8.76432C15.6917 8.76432 16.2725 9.31186 16.2725 9.98725V14.7771C16.2725 15.4525 15.6917 16 14.9752 16H10.3806C9.6641 16 9.08327 15.4525 9.08327 14.7771V9.98725C9.08327 9.31186 9.6641 8.76432 10.3806 8.76432H14.9752ZM1.08329 11.465C1.53108 11.465 1.8941 11.8072 1.8941 12.2293V15.0828C1.8941 15.5049 1.53108 15.8471 1.08329 15.8471C0.635485 15.8471 0.272481 15.5049 0.272481 15.0828V12.2293C0.272481 11.8072 0.635485 11.465 1.08329 11.465ZM6.7049 11.465C7.15271 11.465 7.51571 11.8072 7.51571 12.2293V15.0828C7.51571 15.5049 7.15271 15.8471 6.7049 15.8471C6.25709 15.8471 5.89409 15.5049 5.89409 15.0828V12.2293C5.89409 11.8072 6.25709 11.465 6.7049 11.465ZM3.89409 12.1274C4.34189 12.1274 4.7049 12.4696 4.7049 12.8917V14.4204C4.7049 14.8425 4.34189 15.1847 3.89409 15.1847C3.4463 15.1847 3.08328 14.8425 3.08328 14.4204V12.8917C3.08328 12.4696 3.4463 12.1274 3.89409 12.1274ZM14.6508 10.293H10.7049V14.4713H14.6508V10.293ZM6.7049 8.76432C7.15271 8.76432 7.51571 9.10654 7.51571 9.52865C7.51571 9.95079 7.15271 10.293 6.7049 10.293H1.08327C0.635504 10.293 0.272461 9.95079 0.272461 9.52865C0.272461 9.10654 0.635504 8.76432 1.08327 8.76432H6.7049ZM6.16435 0C6.88084 0 7.46165 0.547516 7.46165 1.22293V6.01273C7.46165 6.68814 6.88084 7.23566 6.16435 7.23566H1.56978C0.853291 7.23566 0.272481 6.68814 0.272481 6.01273V1.22293C0.272481 0.547516 0.853291 0 1.56978 0H6.16435ZM14.9752 0C15.6917 0 16.2725 0.547516 16.2725 1.22293V6.01273C16.2725 6.68814 15.6917 7.23566 14.9752 7.23566H10.3806C9.6641 7.23566 9.08327 6.68814 9.08327 6.01273V1.22293C9.08327 0.547516 9.6641 0 10.3806 0H14.9752ZM5.84003 1.52866H1.8941V5.707H5.84003V1.52866ZM14.6508 1.52866H10.7049V5.707H14.6508V1.52866Z"
                  fill="currentColor"></path>
              </svg>
              <span class="left-ch">下载APP</span>
            </a>
          </li>
        </ul>
        <div class="search">
          <el-select v-model="inputSelect" style="margin-left: 5px; width: 100px">
            <el-option v-for="item in inputSelects" :key="item.value" :label="item.label" :value="item.value">
            </el-option>
          </el-select>
          <div class="search-input">
            <input class="header-input" placeholder="乌克兰局势" v-model="searchtext" />
          </div>
          <div class="search-btn" @click="searchBtn">
            <svg width="17" height="17" viewBox="0 0 17 17" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd" clip-rule="evenodd"
                d="M16.3451 15.2003C16.6377 15.4915 16.4752 15.772 16.1934 16.0632C16.15 16.1279 16.0958 16.1818 16.0525 16.2249C15.7707 16.473 15.4456 16.624 15.1854 16.3652L11.6848 12.8815C10.4709 13.8198 8.97529 14.3267 7.44714 14.3267C3.62134 14.3267 0.5 11.2314 0.5 7.41337C0.5 3.60616 3.6105 0.5 7.44714 0.5C11.2729 0.5 14.3943 3.59538 14.3943 7.41337C14.3943 8.98802 13.8524 10.5087 12.8661 11.7383L16.3451 15.2003ZM2.13647 7.4026C2.13647 10.3146 4.52083 12.6766 7.43624 12.6766C10.3517 12.6766 12.736 10.3146 12.736 7.4026C12.736 4.49058 10.3517 2.1286 7.43624 2.1286C4.50999 2.1286 2.13647 4.50136 2.13647 7.4026Z"
                fill="currentColor" />
            </svg>
          </div>
        </div>

        <ul class="header-top-right">
          <div class="av">
            <div @click="open" v-if="!isLogin">登录</div>
            <div>
              <el-popover :width="300" trigger="hover">
                <template #reference>
                  <img src="../../assets/head.jpg" v-if="isLogin" @click="to('/personal')" class="av-img" />
                </template>
                <template #default>
                  <div class="demo-rich-conent" style="display: flex; gap: 16px; flex-direction: row">
                    <div>
                      确认退出账号吗
                      <el-button @click="logOut">退出</el-button>
                    </div>
                  </div>
                </template>
              </el-popover>
            </div>
          </div>
          <li v-for="(item, index) in rightData" :key="index">
            <a class="entry-title entry-right" href="#" @click.prevent="navigate(item)">
              <div v-if="item.sag != ''" v-html="item.sag" class="right-sag"></div>
              <div class="right-ch">{{ item.text }}</div>
            </a>
          </li>
          <div class="upload" @click="to('/upload')">
            <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg"
              class="header-upload-entry__icon">
              <path
                d="M12.0824 10H14.1412C15.0508 10 15.7882 10.7374 15.7882 11.6471V12.8824C15.7882 13.792 15.0508 14.5294 14.1412 14.5294H3.84707C2.93743 14.5294 2.20001 13.792 2.20001 12.8824V11.6471C2.20001 10.7374 2.93743 10 3.84707 10H5.90589"
                stroke="currentColor" stroke-width="1.7" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M8.99413 11.2353L8.99413 3.82353" stroke="currentColor" stroke-width="1.7" stroke-linecap="round"
                stroke-linejoin="round" />
              <path d="M12.0823 6.29413L8.9941 3.20589L5.90587 6.29413" stroke="currentColor" stroke-width="1.7"
                stroke-linecap="round" stroke-linejoin="round" />
            </svg>
            <div class="upload-text">投稿</div>
          </div>
        </ul>
      </div>
    </div>
  </div>
</template>

<script>
import Headerch from "./headerCh.vue";
import { ElNotification } from "element-plus";
import { h } from "vue"
import { headerLeftMenu, headerRightMenu } from "@/constants/headerMenu";
export default {
  components: {
    Headerch,
  },
  methods: {
    async logOut() {
      try {
        const { AuthAPI } = require("@/api/index");
        await AuthAPI.logout();
      } catch (e) {
        console.error("Logout error", e);
      }
      localStorage.removeItem("isLogin");
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      localStorage.removeItem("userInfo");
      window.location.reload();
    },
    open() {
      this.$emit("openChange", true);
    },
    to(url) {
      if (!localStorage.getItem("isLogin")) {
        ElNotification({
          title: "错误",
          message: h("p", { style: "color: red" }, "您还未登录"),
        });
        return;
      }
      this.$router.push(url);
    },
    navigate(item) {
      if (item.requiresAuth) {
        this.to(item.url);
        return;
      }
      this.$router.push(item.url);
    },
    searchBtn() {
      // alert(this.inputSelect)
      if (this.searchtext) {
        if (this.inputSelect === 0) {
          let query = {
            text: this.searchtext
          };
          this.$router.push({
            path: '/sevideo',
            query
          });
        }
        if(this.inputSelect === 1) {
          let query = {
            text: this.searchtext
          };
          this.$router.push({
            path: '/seuser',
            query
          });
        }
      }
      else {
        this.$router.push('/search');
      }
    },
  },
  data() {
    return {
      rightData: [],
      leftData: [],
      searchtext: "",
      isLogin: false,
      inputSelects: [
        {
          value: 0,
          label: "视频",
        },
        {
          value: 1,
          label: "用户",
        },
      ],
      inputSelect: 0,
    };
  },
  created() {
    this.leftData = headerLeftMenu;
    this.rightData = headerRightMenu;
    this.isLogin = localStorage.getItem("isLogin");
  },
};
</script>

<style scoped lang="less">
a {
  text-decoration: none;
}

@media screen and (max-width: 1367px) {
  .download {
    display: none;
  }
}

@media screen and (max-width: 1279.9px) and (min-width: 970px) {
  .download {
    display: none;
  }

  .right-ch {
    display: none;
  }

  .upload-text {
    display: none;
  }

  .entry-right {
    min-width: 34px !important;
    width: 40px !important;
  }

  .upload {
    width: 34px !important;
    height: 34px;
    margin-left: 15px !important;
  }
}

@media screen and (max-width: 970px) {
  .download {
    display: none;
  }

  .right-ch {
    display: none;
  }

  .upload-text {
    display: none;
  }

  .entry-right {
    min-width: 34px !important;
    width: 40px !important;
  }

  .upload {
    width: 34px !important;
    height: 34px;
    margin-left: 15px !important;
  }
}

.right-sag {
  margin-right: 6px;
  height: 22px;
  width: 22px;
  transition: 0.3s;
}

.upload:hover {
  background-color: #fc8bab;
  transition: 0.3s;
}

.upload {
  margin-right: 10px;
  height: 34px;
  width: 90px;
  background-color: #fb7299;
  color: #fff;
  display: flex;
  justify-content: center;
  align-items: center;
  margin-left: 20px;
  border-radius: 6px;
  transition: 0.3s;
  cursor: pointer;
}

.entry-right {
  flex-direction: column;
  min-width: 50px;
  font-size: 13px;
}

a {
  display: block;
}

.right-ch {
  height: 20px;
}

.header-top-right {
  display: flex;
  height: 64px;
  align-items: center;
}

.header-top-right>a {
  display: flex;
  justify-content: center;
  align-items: center;
}

span {
  font-weight: 500;
  font-size: 16px;
}

.svg {
  margin-right: 3px;
  position: relative;
  top: 2px;
}

.entry-title {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 64px;
  color: #fff;
  font-size: 14px;
  cursor: pointer;
}

.av {
  border: 2px solid #ffffff;
  border-radius: 50%;
  height: 50px;
  width: 50px;
  background-color: #fff;
  margin-right: 20px;
  /* border: 2px solid #fff; */
  display: flex;
  justify-content: center;
  align-items: center;

  .av-img {
    height: 50px;
    width: 50px;
    overflow: hidden;
    border-radius: 50%;
  }

  div {
    cursor: pointer;
    font-size: 18px;
    color: #00aeec;
  }
}

.av:hover {
  cursor: pointer;
}

.av>img {
  height: 50px;
  width: 50px;
  overflow: hidden;
  border-radius: 50%;
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
  background-color: #fff;
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

.header-input:focus {
  background-color: #e3e5e7;
  transition: 0.5s;
}

.left-ch {
  margin-right: 20px;
  cursor: pointer;
  transition: 0.3s;
}

.left-ch:hover,
.right-sag:hover {
  transform: translateY(-4px);
  transition: 0.3s;
}

.header-top {
  position: absolute;
  display: flex;
  justify-content: space-around;
  align-items: center;
  color: white;
  height: 64px;
  padding: 0 24px;
  width: 97%;
}

.header-top-left {
  display: flex;
  height: 64px;
  z-index: 999;
  align-items: center;
  margin-right: 30px;
}

.header-video {
  object-fit: cover;
  transform: scale(1) translate(0px, 0px) rotate(0deg);
  width: 100%;
  aspect-ratio: auto 1728 / 162;
  // height: 162px;
  overflow: hidden;
}

.header {
  margin-bottom: 55px;
  font-size: 16px;
  /* color:white; */
  /* background-color: pink; */
  width: 100%;
  min-width: 970px;
  // z-index:999;
}
</style>
