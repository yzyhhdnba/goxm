<template>
  <el-form ref="ruleFormRef" :model="ruleForm" status-icon :rules="rules" label-width="120px" class="demo-ruleForm">
    <el-form-item prop="companyId" size="large">
      <el-input v-model="ruleForm.username" clearable class="pointer" placeholder="PiliPili ID"></el-input>
    </el-form-item>
    <el-form-item prop="passWord" size="large">
      <el-input v-model="ruleForm.password" type="password" autocomplete="off" @keyup.enter="submitForm(ruleFormRef)"
        class="pointer" placeholder="密码"></el-input>
    </el-form-item>
    <el-form-item>
      <div class="btns">
        <el-button type="primary" @click="submitForm(ruleFormRef)">登录</el-button>
        <el-button @click="resetForm(ruleFormRef)">重置</el-button>
      </div>
    </el-form-item>
  </el-form>
  <div style="position: relative;left:50px">
    <el-checkbox v-model="isStill" label="保持我的登录状态" size="large" />
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, h } from "vue";
import type { ElForm } from "element-plus";
import { ElNotification } from 'element-plus'
import axios from "axios"
import qs from "qs"
import { useRouter } from "vue-router";
import { useStore } from 'vuex'
import { CircleCheck } from '@element-plus/icons-vue'
import { ElLoading } from 'element-plus'

type FormInstance = InstanceType<typeof ElForm>;
const ruleFormRef = ref<FormInstance>();
let isLoading: any = ref(false)
const router = useRouter();
let isStill = ref(false);
const store = useStore();
const validatePass = (_rule: any, value: any, callback: any) => {
  if (value === "") {
    callback(new Error("请输入密码"));
  }
  else if (value.length < 8 || value.length > 16) {
    callback(new Error("密码要在8-16位之间"));
  }
  else callback();
};
localStorage.clear();
const ruleForm = reactive({
  username: "",
  password: "",

});

const rules = reactive({
  username: [
    {
      required: true,
      message: '请输入用户名',
      trigger: 'blur',
    },
    {
      min: 3,
      max: 10,
      message: '用户名要在3-10位之间!',
      trigger: 'blur',
    },
  ],

  password: [{ required: true, validator: validatePass, trigger: "blur" }],
});

const submitForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.validate((valid) => {
    if (valid) {
      const loadingInstance = ElLoading.service({ fullscreen: true });
      axios({
        method: 'post',
        url: "http://172.20.10.6:8081/login",
        data: qs.stringify({
          password: ruleForm.password,
          username: ruleForm.username,
        }),
        headers: { 'content-type': 'application/x-www-form-urlencoded' },
      }).then((res) => {
        console.log(res);
        if (res.data.hasOwnProperty('status')) {
          if (res.data.status === -1) {

            ElNotification({
              title: "登录失败",
              message: h("p", { style: "color: red" }, "用户不存在"),
            });
          }
          else if (res.data.status === -2) {
            ElNotification({
              title: "登录失败",
              message: h("p", { style: "color: red" }, "密码输入错误"),
            });
          }
          else if (res.data.status === 0) {
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
          //@ts-ignore
          localStorage.setItem('isLogin', true);
          localStorage.setItem('userInfo', JSON.stringify(res.data))
          loadingInstance.close()
          router.replace('/management')
        }
      }).catch((err) => {
        console.log(err);
      });

    } else {

      ElNotification({
        title: '格式错误',
        message: h('text', { style: 'color: red' }, '信息填写格式错误'),
        type: 'error'
      })
      return false;
    }
  });
};

const resetForm = (formEl: FormInstance | undefined) => {
  if (!formEl) return;
  formEl.resetFields();
};
</script>
<style scoped>
.demo-ruleForm {
  width: 100%;
}

.pointer {
  height: 30px;
}

.btns {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 250px;
  position: relative;
  left: 70px;
}
</style>