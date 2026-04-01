const fs = require('fs');

const loginFile = 'frontend/src/components/management/login.vue';
let content = fs.readFileSync(loginFile, 'utf8');

content = content.replace(
  /AuthAPI\.login\(ruleForm\.username,\s*ruleForm\.password\)/, 
  "AuthAPI.login({ username: ruleForm.username, password: ruleForm.password })"
);

// Also add ElMessage or ElNotification on catch to show errors better
content = content.replace(
  /\.catch\(err => \{\s*console\.error\(err\);\s*loadingInstance\.close\(\);\s*\}\);/,
  `.catch(err => {
        console.error(err);
        ElNotification({ title: "登录失败", message: h("p", { style: "color: red" }, err.response?.data?.error || "账号或密码错误"), type: "error" });
        loadingInstance.close();
      });`
);

fs.writeFileSync(loginFile, content);
console.log('Admin login params fixed.');
