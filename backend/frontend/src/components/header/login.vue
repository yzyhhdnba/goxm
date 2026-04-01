<template>
  <!-- Form -->
  <el-dialog v-model="props.isopen.dialogFormVisible" custom-class="dialog" :modal="true">
    <div class="login">
      <div class="login-bar">
        <div :class="{ active: isLoginPage }" @click="toRegister">密码登录</div>
        <div class="fenge"></div>
        <div :class="{ active: !isLoginPage }" @click="toLogin">立即注册</div>
      </div>
      <div class="login-input" v-if="isLoginPage">
        <div>
          <span class="text">账号</span>
          <input v-model="loginForm.username" placeholder="请输入用户名" />
        </div>
        <div>
          <span class="text">密码</span>
          <input v-model="loginForm.password" placeholder="请输入密码" @blur="bluring" @focus="focusing" type="password" />
        </div>
      </div>
      <div class="register-form" v-if="!isLoginPage">
        <div>
          <span class="text">账号</span>
          <input v-model="registerForm.username" placeholder="请输入用户名" @blur='checkUserName' />
          <span class="icnos" v-if="nameCheck">
            <el-icon :size="16" color="red">
              <CircleClose />
            </el-icon>
            <span>用户名已注册</span>
          </span>
        </div>
        <div>
          <span class="text">邮箱</span>
          <input v-model="registerForm.email" placeholder="请输入邮箱" @blur="checkEmail" />
          <span class="icnos" v-if="emailCheck">
            <el-icon :size="16" color="red">
              <CircleClose />
            </el-icon>
            <span>邮箱已被注册</span>
          </span>
        </div>

        <div>
          <span class="text">密码</span>
          <input v-model="registerForm.password1" placeholder="请输入密码" @blur="bluring" @focus="focusing"
            type="password" />
        </div>
        <div>
          <span class="text">确认密码</span>
          <input v-model="registerForm.password2" placeholder="请输入确认密码" @blur="bluring" @focus="focusing"
            type="password" />
        </div>
      </div>
      <div class="boxs">
        <el-checkbox label="记住密码" name="type" v-model="checkboxs.rememberMe" size="large" v-if="isLoginPage"
          @change="boxChange(1)" />
        <div @click="to('/managementLogin')" v-if="isLoginPage">管理员登录</div>
      </div>
      <div style="width: 100%; display: flex; justify-content: center">
        <el-button type="primary" style="width: 350px; height: 40px; border-radius: 10px" v-if="isLoginPage"
          @click="login" :loading="isloading">登录</el-button>
        <el-button type="primary" style="width: 350px; height: 40px; border-radius: 10px" v-if="!isLoginPage"
          @click="register" :loading="isloading" :disabled="emailCheck || nameCheck || !hasNetwork1 || !hasNetwork2">注册
        </el-button>
      </div>
      <div style="margin-top: 60px; display: flex; justify-content: center">
        <p>登录或完成注册即代表你同意用户协议和隐私政策</p>
      </div>
    </div>
    <div class="buttom">
      <div class="buttom-left" ref="b1"></div>
      <div class="buttom-right" ref="b2"></div>
    </div>
  </el-dialog>
</template>

<script>
import { useRouter } from "vue-router";
import axios from "axios";
import { reactive, ref, h } from "vue";
import { ElNotification } from "element-plus";
import qs from 'qs'
export default {
  props: {
    isopen: Object,
  },
  emits: ['loginSuccess'],
  setup(props) {
    const router = useRouter();
    const formLabelWidth = "140px";
    const loginForm = reactive({
      username: "",
      password: "",
    });
    let hasNetwork1 = ref(false);
    let hasNetwork2 = ref(false);
    let emailCheck = ref(false);
    let nameCheck = ref(false);

    const checkUserName = () => {

      axios.get('http://172.20.10.6:8081/auth/check/name', {
        params: {
          userName: registerForm.username
        },
      }).then(res => {
        if (res.data.status === 0) {
          nameCheck.value = true;
          // console.log('重复');
          hasNetwork1.value = true
        }
        else {
          hasNetwork2.value = true;
        }
        console.log(res);
      })
    }
    const checkEmail = () => {
      axios.get('http://172.20.10.6:8081/auth/check/email', {
        params: {
          email: registerForm.email
        }
      }).then(res => {
        console.log(res);
        if (res.data.status === 0) {
          emailCheck.value = true;
          hasNetwork1.value = true;
        }
        else {
          hasNetwork1.value = true;
        }
      })
    }
    const registerForm = reactive({
      username: "",
      password1: "",
      password2: "",
      email: "",
    });
    let isLoginPage = ref(true);
    const checkboxs = reactive({
      rememberMe: false,
      autoLogin: false,
    });
    const toRegister = () => {
      isLoginPage.value = true;
      console.log((document.querySelector(".dialog").style.height = "460px"));
    };
    const toLogin = () => {
      isLoginPage.value = false;
      console.log((document.querySelector(".dialog").style.height = "524px"));
    };
    const boxChange = (index) => {

    };
    let isloading = ref(false);
    const login = () => {
      if (/[0-9]+/g.test(loginForm.username)) {
        ElNotification({
          title: "登录失败",
          message: h("p", { style: "color: red" }, "用户名不能是纯数字"),
        });
        return;
      }
      isloading.value = true;
      if (checkboxs.rememberMe) {
        localStorage.setItem("rememberMe", loginForm.password)
      }
      if (checkboxs.autoLogin) {

      }
      if (loginForm.username && loginForm.password) {
        axios({
          method: 'post',
          url: "http://172.20.10.6:8081/login",
          data: qs.stringify({
            password: loginForm.password,
            username: loginForm.username,
          }),
          headers: { 'content-type': 'application/x-www-form-urlencoded' },
        }).then((res) => {
          console.log(res);
          if (res.data.hasOwnProperty('status')) {
            if (res.data.status === -1) {
              isloading.value = false;
              ElNotification({
                title: "登录失败",
                message: h("p", { style: "color: red" }, "用户不存在"),
              });
            }
            else if (res.data.status === -2) {
              isloading.value = false;
              ElNotification({
                title: "登录失败",
                message: h("p", { style: "color: red" }, "密码输入错误"),
              });
            }
            else if (res.data.status === 0) {
              isloading.value = false;
              ElNotification({
                title: "登录失败",
                message: h("p", { style: "color: red" }, "用户未激活"),
              });
            }
          }
          else {
            ElNotification({
              title: "登录成功",
              message: h("p", { style: "color: green" }, "登录成功ヽ(✿ﾟ▽ﾟ)ノ"),
            });
            isloading.value = false;
            props.isopen.dialogFormVisible = false;
            localStorage.setItem('isLogin', true);
            localStorage.setItem('userInfo', JSON.stringify(res.data))
           location.reload();
            emit('loginSuccess', true)
            
          }
        }).catch((err) => {
          isloading.value = false;
          console.log(err);
        });
      } else if (!loginForm.username && loginForm.password) {
        ElNotification({
          title: "登录失败",
          message: h("p", { style: "color: red" }, "请输入用户名"),
        });
      } else if (loginForm.username && !loginForm.password) {
        ElNotification({
          title: "登录失败",
          message: h("p", { style: "color: red" }, "请输入密码"),
        });
      } else {
        ElNotification({
          title: "登录失败",
          message: h("p", { style: "color: red" }, "请输入用户名和密码"),
        });
      }
    };
    const register = () => {
      if (/[0-9]+/g.test(registerForm.username)) {
        ElNotification({
          title: "注册失败",
          message: h("p", { style: "color: red" }, "用户名不能是纯数字"),
        });
      }
      else if (!registerForm.username || !registerForm.password1) {
        ElNotification({
          title: "注册失败",
          message: h("p", { style: "color: red" }, "请检查您的用户名或密码"),
        });
      } else if (registerForm.password1 !== registerForm.password2) {
        ElNotification({
          title: "注册失败",
          message: h("p", { style: "color: red" }, "两次密码不匹配"),
        });
      } else if (!registerForm.email) {
        ElNotification({
          title: "注册失败",
          message: h("p", { style: "color: red" }, "请输入邮箱"),
        });
      } else if (
        registerForm.username &&
        registerForm.password1 &&
        registerForm.email
      ) {
        isloading.value = true;
        axios
          .post("http://172.20.10.6:8081/auth/register", {
            userName: registerForm.username,
            password: registerForm.password1,
            email: registerForm.email,
          })
          .then((res) => {
            console.log(res);
            if (res.data.status == 0) {
              ElNotification({
                title: "注册失败",
                message: h("p", { style: "color: red" }, "邮箱或者用户名重复"),
              });
              return;
            }
            else {
              isloading.value = false;
              ElNotification({
                title: "注册成功",
                message: h("p", { style: "color: green" }, "一封带有激活链接的邮件已发送，请注意查收"),
              });
            }
            // window.location.reload();
          })
          .catch((res) => {
            console.log(res);
            isloading.value = false;
            ElNotification({
              title: "注册失败",
              message: h("p", { style: "color: red" }, "网络错误"),
            });
          });
      }
    };
    const bluring = () => {
      // console.log(document.querySelector('.buttom-left'));
      document.querySelector(".buttom-left").style.backgroundImage =
        `url(` + require("../../assets/image/22_open.png") + `)`;
      document.querySelector(".buttom-right").style.backgroundImage =
        `url(` + require("../../assets/image/33_open.png") + `)`;
    };
    const focusing = () => {
      // console.log(document.querySelector('.buttom-left'));
      document.querySelector(".buttom-left").style.backgroundImage =
        `url(` + require("../../assets/image/22_close.png") + `)`;
      document.querySelector(".buttom-right").style.backgroundImage =
        `url(` + require("../../assets/image/33_close.png") + `)`;
    };
    const to = (url) => {
      router.push(url)
    }
    return {
      to,
      formLabelWidth,
      loginForm,
      props,
      checkboxs,
      bluring,
      focusing,
      toRegister,
      toLogin,
      isLoginPage,
      registerForm,
      login,
      register,
      isloading,
      boxChange,
      checkUserName,
      checkEmail,
      emailCheck,
      nameCheck,
      hasNetwork1,
      hasNetwork2
    };
  },
};
</script>
<style>
.dialog {
  border-radius: 10px !important;
  height: 460px;
  display: flex !important;
  justify-content: center !important;
  width: 820px !important;
  border-radius: 8px !important;
  -webkit-box-shadow: 0 0 6px rgb(0 0 0 / 10%) !important;
  box-shadow: 0 0 6px rgb(0 0 0 / 10%) !important;
  padding: 52px 65px 29px 92px !important;
  -webkit-box-sizing: border-box !important;
  box-sizing: border-box !important;
}
</style>
<style scoped lang="less">
* {
  font-family: -apple-system, BlinkMacSystemFont, Helvetica Neue, Helvetica,
    Arial, PingFang SC, Hiragino Sans GB, Microsoft YaHei, sans-serif;
  padding: 0;
  margin: 0;
}

.icnos {
  position: relative;
  left: 40px;
  width: 95px;
  display: inline-flex;

  span {
    font-size: 12px;
  }
}

.buttom {
  border-radius: 8px;
  display: flex;
  width: 820px;
  justify-content: space-between;
  position: relative;
  left: -42px;
  top: -94px;

  .buttom-left,
  .buttom-right {
    border-radius: 8px;
    position: relative;
    height: 115px;
    width: 115px;
    background-color: #fff;
    background-repeat: no-repeat, no-repeat !important;
    background-size: 100% !important;
    -webkit-box-orient: vertical !important;
    -webkit-box-direction: normal !important;
    -ms-flex-direction: column !important;
    flex-direction: column !important;
    background-position: 0 100%, 100% 100%;
  }

  .buttom-right {
    background-image: url("../../assets/image/33_open.png");
  }

  .buttom-left {
    background-image: url("../../assets/image/22_open.png");
  }
}

.login-bar {
  display: flex;
  color: #505050;
  display: -webkit-box;
  display: -ms-flexbox;
  -webkit-box-pack: center;
  -ms-flex-pack: center;
  justify-content: center;
  height: 25px;
  // margin-top: 30px;
  margin-bottom: 35px;

  div {
    color: #505050;
    font-size: 23px;
    cursor: pointer;
  }

  .active {
    color: #4fa5d9;
  }

  .fenge {
    margin: 0 20px;
    height: 20px;
    width: 0px;
    border-left: 1px solid #e7e7e7;
  }
}

.login {
  display: flex;
  width: 100%;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

input {
  border: 1px solid #000;
}

.boxs {
  display: flex;
  justify-content: space-around;
  width: 60%;
  margin-bottom: 20px;
  align-items: center;

  div {
    cursor: pointer;
    font-size: 16px;
    text-decoration: underline;
  }

  div:hover {
    color: #4fa5d9;
  }

}

.register-form {
  div:nth-child(1) {
    border-radius: 8px 8px 0 0;
    border-bottom: none;
  }

  div:nth-child(4) {
    border-radius: 0 0 8px 8px;
    // border-bottom: none;
  }
}

.login-input {
  div:nth-child(1) {
    border-radius: 8px 8px 0 0;
    border-bottom: none;
  }

  div:nth-child(2) {
    border-radius: 0 0 8px 8px;
    // border-bottom: none;
  }
}

.login-input,
.register-form {
  display: flex;
  flex-direction: column;
  margin-bottom: 15px;

  div {
    height: 50px;
    width: 400px;
    display: -webkit-box;
    display: -ms-flexbox;
    display: flex;
    padding: 0 20px;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    border: 1px solid #e7e7e7;
  }

  .text {
    margin-right: 20px;
    font-size: 17px;
    color: #212121;
  }

  input {
    width: 230px;
    outline: none;
    border: none;
    font-size: 17px;
    color: #212121;
    -webkit-box-shadow: 0 0 0 20px #fff inset;
    box-shadow: inset 0 0 0 20px #fff;
  }
}
</style>
