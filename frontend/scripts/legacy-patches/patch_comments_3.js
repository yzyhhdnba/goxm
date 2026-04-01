const fs = require('fs');
const path = require('path');

let floorPath = path.join(__dirname, 'frontend/src/components/detail/floor.vue');
let floorContent = fs.readFileSync(floorPath, 'utf8');

const regexFloor = /created\(\)\s*\{[\s\S]*?this\.within = res\.data\s*\}\)/;
const newCreatedFloor = `created() {
    import('@/api/index').then(({ InteractAPI }) => {
      const parentId = this.item.id || this.item.commentId;
      if (!parentId) return;
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

floorContent = floorContent.replace(regexFloor, newCreatedFloor);
fs.writeFileSync(floorPath, floorContent, 'utf8');


let floorPathMore = path.join(__dirname, 'frontend/src/components/detail/floorWithinMore.vue');
let floorContentMore = fs.readFileSync(floorPathMore, 'utf8');
floorContentMore = floorContentMore.replace(/axios/g, 'removed_axios'); // just to kill remaining references
fs.writeFileSync(floorPathMore, floorContentMore, 'utf8');

console.log('Fixed floor files');