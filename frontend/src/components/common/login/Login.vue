<template>
  <div class="login">
    <el-card class="box-card" style="height: 350px">
      <el-form :model="loginForm" class="demo-form-inline" :label-position="right" label-width="100px">
        <el-row>
          <el-col :span="5"></el-col>
          <el-col :span="12">
            <el-form-item label="账号:">
              <el-input v-model="loginForm.id" placeholder="请输入账号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="5"></el-col>
          <el-col :span="12">
            <el-form-item label="密码:">
              <el-input v-model="loginForm.userName" placeholder="请输入密码"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row>
          <el-col :span="7"></el-col>
          <el-col :span="6">
            <el-form-item>
              <el-button class="register-button" type="primary">
                <pre> 登    录</pre>
              </el-button>
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
    </el-card>
  </div>
</template>

<script>
  export default {
    props: {

    },
    components: {

    },
    data(){
      return {
        loginForm: {
          id: '',
          userName: '',
        }
      }
    },
    methods: {

    }
  }
</script>

<style scoped>
pre {
  font-family: 'Helvetica Neue', Helvetica, 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', '微软雅黑', Arial, sans-serif;
}
.register-button {
  padding: 3px 20px;
}
.el-form-item * {
  margin: 0 auto;
}
.box-card {
  background-color: rgba(255, 255, 255, 0.4);
}
.row {
  margin-top: 40px;
}
</style>