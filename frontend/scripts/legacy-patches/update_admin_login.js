const fs = require('fs');

const loginFile = 'frontend/src/components/management/login.vue';
let content = fs.readFileSync(loginFile, 'utf8');

// Replace old axios import and logic with new AuthAPI login
content = content.replace(/import axios from "axios"/g, "import { AuthAPI } from '@/api/index'");
content = content.replace(/import qs from "qs"/g, "");

content = content.replace(/axios\(\{\s*method:\s*'post',\s*url:\s*"http:\/\/172\.20\.10\.6:8081\/login",[\s\S]*?headers:\s*\{[^}]*\},?\s*\}\)\.then\(\(res\)\s*=>\s*\{[\s\S]*?\.catch\(\(err\)\s*=>\s*\{[\s\S]*?\}\);/m, 
`AuthAPI.login(ruleForm.username, ruleForm.password).then(res => {
        if (res && res.data) {
          ElNotification({
            title: "登录成功",
            message: h("p", { style: "color: green" }, "按新后端 API 登录成功ヽ(✿ﾟ▽ﾟ)ノ"),
          });
          // 管理员状态判断通常由后续接口报错或我们可以在前端简单拦截，但鉴于阶段，先让他进去再说
          localStorage.setItem('isLogin', 'true');
          localStorage.setItem('userInfo', JSON.stringify(res.data.user || res.data));
          if (res.data.access_token) {
            localStorage.setItem('token', res.data.access_token);
            localStorage.setItem('refresh_token', res.data.refresh_token);
          }
          loadingInstance.close();
          router.replace('/management');
        }
      }).catch(err => {
        console.error(err);
        loadingInstance.close();
      });`);


fs.writeFileSync(loginFile, content);
console.log('Admin login updated.');
