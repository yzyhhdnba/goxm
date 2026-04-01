<template>
  <el-container class="se-all">
    <el-main class="search">
      <el-row gutter="10">
        <img
          src="@/assets/image/logo.png"
          class="img-logo"
        />
        <div><p class="vertical"></p></div>
        <div><p class="font-search">搜索</p></div>
      </el-row>
      <el-col>
        <div class="searchi">
          <el-select
            v-model="inputSelect"
            style="margin-left: 5px; width: 100px"
          >
            <el-option
              v-for="item in inputSelects"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            ></el-option>
          </el-select>
          <div class="search-input">
            <input
              class="header-input"
              placeholder="输入关键字搜索"
              v-model="searchtext"
            />
          </div>

          <div class="search-btn" @click="searchBtn">
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
      </el-col>
    </el-main>
  </el-container>
</template>

<script>
export default {
  props: {
    islogin: Object,
  },
  components: {},
  methods: {
    open() {
      this.$emit("openChange", true);
    },
    to(url) {
      this.$router.push(url);
    },
    searchBtn() {
      // alert(this.inputSelect)
      if (this.searchtext) {
        if (this.inputSelect === 0) {
          let query = {
            text: this.searchtext,
          };
          this.$router.push({
            path: "/sevideo",
            query,
          });
        }
        if (this.inputSelect === 1) {
          let query = {
            text: this.searchtext,
          };
          this.$router.push({
            path: "/seuser",
            query,
          });
        }
      }
    },
   
  },
  data() {
    return {
      searchtext: "",
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
  },
};
</script>

<style scoped>
.se-all {
  margin: auto;
  height: 700px;
  padding-top: 50px;
  padding-bottom: 50px;
}

.search {
  margin-top: 200px;
  margin-left: 450px;
}

.img-logo {
  width: 180px;
  height: 75px;
  padding-bottom: 30px;
  padding-left: 110px;
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
  /* background-color: #e3e5e7; */
  transition: 0.5s;
  cursor: pointer;
}

.search-input {
  width: 100%;
  align-items: center;
  justify-content: center;
  align-items: center;
  border-radius: 6px;
  /* background-color: #f1f2f3; */
  transition: 0.5s;
}

.searchi {
  display: flex;
  width: 500px;
  height: 40px;
  border: 1px;
  /* background-color: #f1f2f3; */
  /* padding: 0 48px; */
  border-radius: 8px;
  justify-content: space-between;
  align-items: center;
  margin-right: 30px;
  padding-right: 5px;
  flex: 1;
  transition: 0.5s;
  opacity: 0.8;
  border-color: #ff2250;
  border-style: solid;
}

.font-search {
  font-size: 45px;
  color: #ff2250;
  padding-left: 20px;
  padding-top: 1px;
  font-weight: 600;
}

.search:hover,
.search-input:hover {
  transition: 0.5s;
  background-color: #fff;
}

.vertical {
  height: 60px;
  width: 0px;
  border-color: #ff2250;
  border-style: solid;
}
</style>
