<template>
  <div class="register">
    <el-card class="box-card" style="height: 490px">
      <el-form :model="registerForm" class="demo-form-inline" :label-position="right" label-width="100px">
        <!--上传用户头像-->
        <el-row class="row">
          <el-col :span="8"></el-col>
          <el-col :span="5">
            <el-form-item label="用户头像:">
              <el-upload class="avatar-uploader" action="" :show-file-list="false" :on-success="handleAvartarSuccess"
                :before-upload="beforeAvatarUpload">
                <img v-if="imageUrl" :src="imageUrl" class="avatar">
                <el-avatar src="https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png" shape="circle"
                  :size="50" />
              </el-upload>
            </el-form-item>
          </el-col>
        </el-row>
        <!--学号-->
        <el-row>
          <el-col :span="5"></el-col>
          <el-col :span="12">
            <el-form-item label="学号:">
              <el-input v-model="registerForm.id" placeholder="请输入学号"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!--用户名-->
        <el-row>
          <el-col :span="5"></el-col>

          <el-col :span="12">
            <el-form-item label="用户名:">
              <el-input v-model="registerForm.userName" placeholder="请输入用户名"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!--邮箱-->
        <el-row>
          <el-col :span="5"></el-col>

          <el-col :span="12">
            <el-form-item label="邮箱:">
              <el-input v-model="registerForm.Email" placeholder="请输入邮箱"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!--密码-->
        <el-row>
          <el-col :span="5"></el-col>

          <el-col :span="12">
            <el-form-item label="密码:">
              <el-input v-model="registerForm.password" placeholder="请输入密码" type="password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!--确认密码-->
        <el-row>
          <el-col :span="5"></el-col>

          <el-col :span="12">
            <el-form-item label="确认密码:">
              <el-input v-model="registerForm.confirmPwd" placeholder="请再次输入密码" type="password"></el-input>
            </el-form-item>
          </el-col>
        </el-row>
        <!--注册按钮-->
        <el-row>
          <el-col :span="7"></el-col>

          <el-col :span="6">
            <el-form-item>
              <el-button class="register-button" type="primary">
                <pre> 注    册 </pre>
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
        registerForm: {
          id: '',
          userName: '',
          Email: '',
          password: '',
          confirmPwd: '',
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