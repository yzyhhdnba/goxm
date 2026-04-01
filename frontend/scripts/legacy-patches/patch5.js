const fs = require('fs');
const path = './frontend/src/components/common/card/cardList.vue';
let content = fs.readFileSync(path, 'utf8');

const regex = /setup\(props\)\s*\{[\s\S]*?return\s*\{/;
const newSetup = `setup(props) {
    let videoList=ref(props.maindataItem.slice(0, 5))
    const change = async () => {
      try {
        const { VideoAPI } = require('@/api/index');
        const res = await VideoAPI.getRecommend({ limit: 5, cursor: Math.random().toString() });
        if (res && res.items) {
          videoList.value = res.items;
        }
      } catch (err) {
        console.error(err);
      }
    }
    return {`;

content = content.replace(regex, newSetup);

fs.writeFileSync(path, content, 'utf8');
console.log('cardList patched');
