const fs = require('fs');

const path = 'frontend/src/components/detail/videoDetail.vue';
let content = fs.readFileSync(path, 'utf8');

// replace invalid src binding in avatar
content = content.replace(/src="\{\{ \(videoDetail\.author && videoDetail\.author\.avatar_url\) \|\| "https:\/\/cube\.elemecdn\.com\/0\/88\/03b0d39583f48206768a7534e55bcpng\.png" \}\}"/g, 
    `:src="(videoDetail.author && videoDetail.author.avatar_url) || 'https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png'"`);

fs.writeFileSync(path, content, 'utf8');
console.log('Fixed src binding in videoDetail.vue');
