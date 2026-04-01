const fs = require('fs');
const path = require('path');

function replaceInFile(filePath) {
    const fullPath = path.join(__dirname, 'frontend/src/components/detail', filePath);
    if (!fs.existsSync(fullPath)) return;
    let content = fs.readFileSync(fullPath, 'utf8');

    // Remove any leftover axios imports
    content = content.replace(/import axios from "axios";\n/g, '');

    // Patch submit(replyFor, commentId)
    const submitRegex = /submit\([^)]*\)\s*\{[\s\S]*?\}\s*\)(?:\.catch[^}]*\})?\n\s*\}/g;
    const newSubmit = `submit(replyFor, commentId) {
      if (!localStorage.getItem('userInfo')) {
        ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
        return;
      }
      const targetId = commentId || this.item.commentId || this.item.id;
      import('@/api/index').then(({ InteractAPI }) => {
        InteractAPI.postReply(targetId, { content: this.replyText, reply_to: replyFor }).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "发布评论成功"), type: "success" });
            location.reload();
          }
        });
      });
    }`;
    content = content.replace(submitRegex, newSubmit);

    // Patch dianzan
    const dianzanRegex = /dianzan:?.*?function\s*\(\)\s*\{[\s\S]*?\}\s*\},/g;
    const newDianzan = `dianzan: function () {
    if (!localStorage.getItem('userInfo')) {
      ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录再点赞哦！"), type: "error" });
      return;
    }
    const targetId = this.item.id || this.item.replyId || this.item.commentId;
    import('@/api/index').then(({ InteractAPI }) => {
      if (!this.item.isLike) {
        InteractAPI.likeComment(targetId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "点赞成功！"), type: "success" });
            if (this.item.like_count !== undefined) this.item.like_count++;
            else if (this.item.replyLikes !== undefined) this.item.replyLikes++;
            this.item.isLike = true;
          }
        });
      } else {
        InteractAPI.unlikeComment(targetId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "取消点赞成功！"), type: "success" });
            if (this.item.like_count !== undefined) this.item.like_count--;
            else if (this.item.replyLikes !== undefined) this.item.replyLikes--;
            this.item.isLike = false;
          }
        });
      }
    });
  },`;
    content = content.replace(dianzanRegex, newDianzan);

    fs.writeFileSync(fullPath, content, 'utf8');
    console.log(`Patched ${filePath} v2`);
}

replaceInFile('floorWithin.vue');
replaceInFile('floorWithinMore.vue');
