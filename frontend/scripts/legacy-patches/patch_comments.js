const fs = require('fs');
const path = require('path');

function replaceInFile(filePath, isRootFloor) {
    const fullPath = path.join(__dirname, 'frontend/src/components/detail', filePath);
    let content = fs.readFileSync(fullPath, 'utf8');

    // Remove axios
    content = content.replace(/import axios from "axios";\n/g, '');

    // Refactor data rendering variables in templates to match backend struct
    // backend comment model: `id`, `user` { username, avatar_url }, `content`, `created_at`, `like_count`, `reply_count`, `viewer_state`
    content = content.replace(/item\.userName/g, '(item.user && item.user.username) || item.userName');
    
    // We already do item.commentContent -> item.content (so let's support both)
    content = content.replace(/item\.commentContent/g, 'item.content || item.commentContent');
    content = content.replace(/item\.replyContent/g, 'item.content || item.replyContent'); // for within

    content = content.replace(/item\.commentTime\.slice\(0, 10\)/g, '(item.created_at || item.commentTime || "").slice(0, 10)');
    content = content.replace(/item\.replyTime\.slice\(0, 10\)/g, '(item.created_at || item.replyTime || "").slice(0, 10)');

    content = content.replace(/item\.commentLikes/g, 'item.like_count !== undefined ? item.like_count : item.commentLikes');
    content = content.replace(/item\.replyLikes/g, 'item.like_count !== undefined ? item.like_count : item.replyLikes');

    content = content.replace(/item\.commentReplies/g, 'item.reply_count !== undefined ? item.reply_count : item.commentReplies');

    if (isRootFloor) {
        // refactor `created()` for floor.vue
        const createdRegex = /created\(\) \{[\s\S]*?\}\)/;
        const newCreated = `created() {
    import('@/api/index').then(({ InteractAPI }) => {
      const parentId = this.item.id || this.item.commentId;
      InteractAPI.getReplies(parentId, { page: 1, limit: 10 }).then(res => {
        if (res && res.data) {
          this.within = res.data.items || res.data || [];
          this.within.forEach(r => {
             if (r.viewer_state) {
                 r.isLike = r.viewer_state.liked;
             }
          });
          this.flag = true;
        }
      }).catch(() => { this.flag = true; });
    });
  }`;
        content = content.replace(createdRegex, newCreated);
    } else {
        // for floorWithin, it also gets deeper replies?
        // Wait, floorWithin might just be displaying a reply, and fetching more? "floorWithinMore" ?
        // if it has similar created
        const createdRegex = /created\(\) \{[\s\S]*?\}\)/;
        const newCreated = `created() {
    import('@/api/index').then(({ InteractAPI }) => {
      const parentId = this.item.id || this.item.replyId;
      InteractAPI.getReplies(parentId, { page: 1, limit: 10 }).then(res => {
        if (res && res.data) {
          this.withinMore = res.data.items || res.data || [];
          this.withinMore.forEach(r => {
             r.isLike = r.viewer_state && r.viewer_state.liked;
          });
          this.flag = true;
        }
      }).catch(() => { this.flag = true; });
    });
  }`;
        if (content.match(createdRegex)) {
            content = content.replace(createdRegex, newCreated);
        }
    }

    // Refactor `submitComment` or `submitReply`
    // It might be submitComment(commentId) or submitReply(replyId)
    const submitRegex = /(submitComment|submitReply)\([^)]*\) \{[\s\S]*?\}\)/;
    const newSubmit = `$1(id) {
    if (!localStorage.getItem('userInfo')) {
      ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
      return;
    }
    const targetId = this.item.id || this.item.commentId || this.item.replyId;
    import('@/api/index').then(({ InteractAPI }) => {
      InteractAPI.postReply(targetId, { content: this.replyText }).then(res => {
        if (res) {
          ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "回复成功"), type: "success" });
          location.reload();
        }
      });
    });
  }`;
    if (content.match(submitRegex)) {
        content = content.replace(submitRegex, newSubmit);
    }

    // Refactor dianzan()
    const dianzanRegex = /dianzan\(\) \{[\s\S]*?\}\s*\},/m;
    const newDianzan = `dianzan() {
    if (!localStorage.getItem('userInfo')) {
      ElNotification({ title: "错误", message: h("p", { style: "color: red" }, "请先登录！"), type: "error" });
      return;
    }
    const targetId = this.item.id || this.item.commentId || this.item.replyId;
    import('@/api/index').then(({ InteractAPI }) => {
      if (!this.item.isLike) {
        InteractAPI.likeComment(targetId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "点赞成功！"), type: "success" });
            if (this.item.like_count !== undefined) this.item.like_count++;
            else if (this.item.commentLikes !== undefined) this.item.commentLikes++;
            else if (this.item.replyLikes !== undefined) this.item.replyLikes++;
            this.item.isLike = true;
          }
        });
      } else {
        InteractAPI.unlikeComment(targetId).then(res => {
          if (res) {
            ElNotification({ title: "成功", message: h("p", { style: "color: green" }, "取消点赞成功！"), type: "success" });
            if (this.item.like_count !== undefined) this.item.like_count--;
            else if (this.item.commentLikes !== undefined) this.item.commentLikes--;
            else if (this.item.replyLikes !== undefined) this.item.replyLikes--;
            this.item.isLike = false;
          }
        });
      }
    });
  },`;
    if (content.match(dianzanRegex)) {
        content = content.replace(dianzanRegex, newDianzan);
    }

    // Sometimes they don't have created() (e.g. floorWithinMore)
    // We just patch what's there
    fs.writeFileSync(fullPath, content, 'utf8');
    console.log(`Patched ${filePath}`);
}

replaceInFile('floor.vue', true);
replaceInFile('floorWithin.vue', false);
replaceInFile('floorWithinMore.vue', false);
