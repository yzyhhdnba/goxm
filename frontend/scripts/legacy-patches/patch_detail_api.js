const fs = require('fs');
const path = require('path');

function replaceInFile(filePath, replacements) {
    const fullPath = path.join(__dirname, 'frontend/src/components', filePath);
    let content = fs.readFileSync(fullPath, 'utf8');
    for (const [search, replace] of replacements) {
        content = content.replace(search, replace);
    }
    fs.writeFileSync(fullPath, content, 'utf8');
    console.log(`Patched ${filePath}`);
}

// 1. Patch recommend.vue
const recommendReplacements = [
    [
        /<img :src="'http:\/\/172\.20\.10\.6:8081\/cover\/cover-' \+ item\.videoId \+ '\.jpeg'" class="pic"/g,
        `<img :src="item.cover_url || item.pic || ('http://172.20.10.6:8081/cover/cover-' + (item.videoId || item.id) + '.jpeg')" class="pic"`
    ],
    [
        /\{\{\s*item\.videoTitle\s*\}\}/g,
        `{{ item.title || item.videoTitle }}`
    ],
    [
        /\{\{\s*item\.userName\s*\}\}/g,
        `{{ (item.author && item.author.username) || item.userName || item.up }}`
    ],
    [
        /播放量:\{\{ item\.videoHits \}\}\s*评论:\{\{ item\.videoComments \}\}/g,
        `播放量:{{ item.view_count !== undefined ? item.view_count : item.videoHits }} 评论:{{ item.comment_count !== undefined ? item.comment_count : item.videoComments }}`
    ]
];
replaceInFile('detail/recommend.vue', recommendReplacements);

// 2. Patch videoMore.vue
const videoMoreScriptReplace = [
    [
        /this\.videoUrl = 'http:\/\/172\.20\.10\.6:8081\/video\/processed\/video-' \+ this\.\$route\.query\.videoId \+ '\/ts\/index\.m3u8';\s*console\.log\(this\.videoUrl\);\s*[\s\S]*?(?=},\s*components: \{)/m,
        `this.videoId = this.$route.query.videoId;
    import('@/api/index').then(({ VideoAPI }) => {
      VideoAPI.getVideoDetail(this.videoId).then(res => {
        if (res && res.data) {
          const data = res.data;
          this.videoUrl = data.play_url || ('http://172.20.10.6:8081/video/processed/video-' + this.videoId + '/ts/index.m3u8');
          this.childList[0].num = data.like_count !== undefined ? data.like_count : (data.videoLikes || 0);
          this.childList[1].num = data.favorite_count !== undefined ? data.favorite_count : (data.videoCollects || 0);
          if (data.viewer_state) {
            this.childList[0].status = data.viewer_state.liked;
            this.childList[1].status = data.viewer_state.favorited;
          }
        }
      });
    });
  `
    ],
    [
        /if \(!localStorage\.getItem\('userInfo'\)\) \{\s*ElNotification[\s\S]*?let \{ userId \} = JSON\.parse\(localStorage\.getItem\('userInfo'\)\)[\s\S]*?if \(e === 0\) \{/m,
        `if (!localStorage.getItem('userInfo')) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录再进行操作哦！"), type: "error" });
        return;
      }
      import('@/api/index').then(({ InteractAPI }) => {
      if (e === 0) {`
    ],
    [
        /axios\.get\('http:\/\/172\.20\.10\.6:8081\/video\/dislike',\s*\{\s*params: \{\s*videoId: this\.videoId,\s*userId\s*\}\s*\}\)/g,
        `InteractAPI.unlikeVideo(this.videoId)`
    ],
    [
        /axios\.get\('http:\/\/172\.20\.10\.6:8081\/video\/like',\s*\{\s*params: \{\s*videoId: this\.videoId,\s*userId\s*\}\s*\}\)/g,
        `InteractAPI.likeVideo(this.videoId)`
    ],
    [
        /axios\.get\('http:\/\/172\.20\.10\.6:8081\/video\/cancel',\s*\{\s*params: \{\s*videoId: this\.videoId,\s*userId\s*\}\s*\}\)/g,
        `InteractAPI.uncollectVideo(this.videoId)`
    ],
    [
        /axios\.get\('http:\/\/172\.20\.10\.6:8081\/video\/collect',\s*\{\s*params: \{\s*videoId: this\.videoId,\s*userId\s*\}\s*\}\)/g,
        `InteractAPI.collectVideo(this.videoId)`
    ],
    [
        /if \(res\.data\.status === 1\)/g,
        `if (res)` // new APIs usually return direct success or throw error, so res is truthy
    ]
];

const videoMoreContent = fs.readFileSync(path.join(__dirname, 'frontend/src/components/video/videoMore.vue'), 'utf8');

// Also handle sanlian function APIs
let modifiedVideoMore = videoMoreContent.replace(/axios\.get\('http:\/\/172\.20\.10\.6:8081\/video\/(like|collect)'[^\)]+\)/g, (match, type) => {
    return type === 'like' ? `InteractAPI.likeVideo(that.videoId)` : `InteractAPI.collectVideo(that.videoId)`;
});

for (const [search, replace] of videoMoreScriptReplace) {
    modifiedVideoMore = modifiedVideoMore.replace(search, replace);
}
// Inject InteractAPI import effectively in sanlian()
modifiedVideoMore = modifiedVideoMore.replace(
    /if \(that\.sanlianNum >= 100 && that\.timeEnd - that\.timeStart >= 3000\) \{[\s\S]*?let \{ userId \} = JSON\.parse\(localStorage\.getItem\('userInfo'\)\)/,
    `if (that.sanlianNum >= 100 && that.timeEnd - that.timeStart >= 3000) {
          import('@/api/index').then(({ InteractAPI }) => {`
);
modifiedVideoMore = modifiedVideoMore.replace(
    /that\.sanlianNum = 0;\s*clearInterval\(that\.timer\);\s*\/\/ alert\('三连成功'\)\s*\}/,
    `that.sanlianNum = 0;
          clearInterval(that.timer);
          });
        }`
);

modifiedVideoMore = modifiedVideoMore.replace(/import axios from 'axios';/g, '');

// Finally close the promise from change(e)
modifiedVideoMore = modifiedVideoMore.replace(
    /this\.childList\[e\]\.status = !this\.childList\[e\]\.status;\s*this\.sanlianNum = 0;\s*\/\/ console\.log\(this\.childList\);\s*\}/,
    `this.childList[e].status = !this.childList[e].status;
      this.sanlianNum = 0;
      });
    }`
);

fs.writeFileSync(path.join(__dirname, 'frontend/src/components/video/videoMore.vue'), modifiedVideoMore, 'utf8');
console.log('Patched videoMore.vue');
