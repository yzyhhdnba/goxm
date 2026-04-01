const fs = require('fs');

const path = 'frontend/src/components/detail/videoDetail.vue';
let content = fs.readFileSync(path, 'utf8');

content = content.replace(/\)\.catch\(err => \{\}\);\n\s*\}\);\n\s*\}\)\n\s*\}/, 
                          `).catch(err => {});\n      });\n    }`);

fs.writeFileSync(path, content, 'utf8');
console.log('Fixed submit dangling parenthesis');
