const fs = require('fs');

const file1 = 'frontend/src/components/detail/floor.vue';
let content = fs.readFileSync(file1, 'utf8');

const regex = /\}\.then\(res => \{\s*console\.log\(res\);\s*if \(res\.data\.status === 1\) \{\s*ElNotification\(\{\s*title: "成功",\s*message: h\("p", \{ style: "color: green" \}, "发布评论成功"\),\s*type: "success",\s*\}\);\s*location\.reload\(\)\s*\}\s*\}\)/g;

content = content.replace(regex, '');

fs.writeFileSync(file1, content);
console.log("Fixed floor.vue!");
