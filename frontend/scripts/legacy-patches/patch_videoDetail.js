const fs = require('fs');
const path = require('path');

const filePath = path.join(__dirname, 'frontend/src/components/detail/videoDetail.vue');
let content = fs.readFileSync(filePath, 'utf8');

// 1. replace template rendering variables
content = content.replace(/videoDetail\.videoTitle/g, 'videoDetail.title || videoDetail.videoTitle');
content = content.replace(/videoDetail\.videoHits/g, 'videoDetail.view_count !== undefined ? videoDetail.view_count : videoDetail.videoHits');
content = content.replace(/videoDetail\.videoComments/g, 'videoDetail.comment_count !== undefined ? videoDetail.comment_count : videoDetail.videoComments');
content = content.replace(/videoDetail\.videoContribute/g, 'videoDetail.published_at || videoDetail.videoContribute');
content = content.replace(/videoDetail\.videoContent/g, 'videoDetail.description || videoDetail.videoContent');
content = content.replace(/\{\{\s*videoUser\s*\}\}/g, '{{ (videoDetail.author && videoDetail.author.username) || videoUser }}');
content = content.replace(/https:\/\/cube\.elemecdn\.com\/0\/88\/03b0d39583f48206768a7534e55bcpng\.png/g, '{{ (videoDetail.author && videoDetail.author.avatar_url) || "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" }}');


// 2. replace script setup logic
const setupRegex = /setup\(\)\s*\{[\s\S]*?let videoDetail = ref\(\{\}\)/;

const newSetup = `setup() {
    let isFollow = ref(false)
    let signal = ref('')
    let comments = ref([])
    let flag = ref(false)
    let router = useRouter();
    let route = useRoute();
    let { videoId } = route.query;
    let commentContent = ref('')
    // console.log(videoId);
    let userId = localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo')).userId : 'null';
    let followId = ''
    let videoUser = ref('')
    let videoDetail = ref({})

    import('@/api/index').then(({ VideoAPI, InteractAPI }) => {
      // 获取视频详情及部分用户状态
      VideoAPI.getVideoDetail(videoId).then(res => {
        if (res && res.data) {
          const data = res.data;
          videoUser.value = (data.author && data.author.username) || data.userName;
          if (data.published_at) {
             data.published_at = data.published_at.slice(0, 10);
          } else if (data.videoContribute) {
             data.videoContribute = data.videoContribute.slice(0, 10);
          }
          videoDetail.value = data;
          followId = data.userId || (data.author && data.author.id);
          if (data.viewer_state) {
            isFollow.value = data.viewer_state.followed;
          } else {
             InteractAPI.getFollowStatus(followId).then(fRes => {
               if(fRes && fRes.data) {
                 isFollow.value = fRes.data.is_following;
               }
             });
          }
        }
      });

      // 获取评论
      InteractAPI.getComments(videoId, { page: 1, limit: 20 }).then(res => {
        if (res && res.data) {
           comments.value = res.data.items || res.data || [];
           // 后端自带 viewer_state.liked，所以不需要再循环获取
           comments.value.forEach(c => {
              if (c.viewer_state) {
                 c.isLike = c.viewer_state.liked;
              }
           });
           flag.value = true;
        }
      });
    });
`;

content = content.replace(setupRegex, newSetup);

// 3. submit comment and follow/unfollow
content = content.replace(/axios\.post\('http:\/\/172\.20\.10\.6:8081\/comment\/make'.+userId,\s*commentContent:\s*commentContent\.value\s*\}\)/s, 
    `import('@/api/index').then(({ InteractAPI }) => { InteractAPI.postComment(videoId, { content: commentContent.value })`);
    
content = content.replace(/if \(res\.data\.status === 1\)/g, `if (res)`);

// ensure we close the wrapper import
content = content.replace(/ElNotification\(\{\n\s*title: "成功",\n\s*message: h\("p", \{ style: "color: green" \}, "发布评论成功"\),\n\s*type: "success",\n\s*\}\);\n\s*location\.reload\(\)\n\s*\}\n\s*\}\)\.catch\(err => \{\n\n\s*\}\)/, 
`ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "发布评论成功"), type: "success" });
          location.reload();
        }
      }).catch(err => {})
      })`); // this logic is brittle, let's substitute more robustly.

// Better way to replace the whole submit function
const submitRegex = /const submit = \(\) => \{[\s\S]*?\}\)\.catch\(err => \{\s*\}\)\s*\}/;
const newSubmit = `const submit = () => {
      if (!localStorage.getItem('userInfo')) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录再发评论哦！"), type: "error" });
        return;
      }
      import('@/api/index').then(({ InteractAPI }) => {
        InteractAPI.postComment(videoId, { content: commentContent.value }).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "发布评论成功"), type: "success" });
            location.reload();
          }
        }).catch(err => {});
      });
    }`;
content = content.replace(submitRegex, newSubmit);

const unfollowRegex = /const unfollow = \(\) => \{[\s\S]*?是Follow\.value = false;\s*\}\n\s*\}\)/;
const newUnfollow = `const unfollow = () => {
      import('@/api/index').then(({ InteractAPI }) => {
        InteractAPI.unfollowUser(followId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "取消关注成功"), type: "success" });
            isFollow.value = false;
          }
        });
      });
    }`;
content = content.replace(/const unfollow[\s\S]*?\}\)\n    \}/, newUnfollow);

const followRegex = /const follow[\s\S]*?\}\)\n    \}/;
const newFollow = `const follow = () => {
      import('@/api/index').then(({ InteractAPI }) => {
        InteractAPI.followUser(followId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "关注成功"), type: "success" });
            isFollow.value = true;
          }
        });
      });
    }`;
content = content.replace(followRegex, newFollow);

// right side videoList recommendation
const videoListRegex = /axios\.post\('http:\/\/172\.20\.10\.6:8081\/videoList'[\s\S]*?\}\)\.catch\(err => \{\s*console\.log\(err\);\s*\}\)/;
const newVideoList = `import('@/api/index').then(({ VideoAPI }) => {
      VideoAPI.getRecommend({ limit: 6 }).then(res => {
        if (res && res.items) posts.value = res.items;
      });
    })`;
content = content.replace(videoListRegex, newVideoList);


fs.writeFileSync(filePath, content, 'utf8');
console.log('Patched videoDetail.vue');
