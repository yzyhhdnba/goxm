const fs = require('fs');
const path = './frontend/src/components/body/partition.vue';
let content = fs.readFileSync(path, 'utf8');

const newScript = `<script>
import Card from "@/components/common/card/Card";
import { VideoAPI } from "@/api/index";

export default {
  async created() {
    try {
      const res1 = await VideoAPI.getRecommend({ limit: 6 });
      if (res1 && res1.items) {
          this.animationList1 = res1.items;
      }
      this.animationListFlag1 = true;
      
      const res2 = await VideoAPI.getRecommend({ limit: 6, cursor: Math.random().toString() }); // Or similar to offset
      if (res2 && res2.items) {
          this.animationList2 = res2.items;
      }
      this.animationListFlag2 = true;

      const res3 = await VideoAPI.getRecommend({ limit: 6, cursor: Math.random().toString() });
      if (res3 && res3.items) {
          this.animationList3 = res3.items;
      }
      this.animationListFlag3 = true;
    } catch (err) {
      console.error(err);
      this.animationListFlag1 = true;
      this.animationListFlag2 = true;
      this.animationListFlag3 = true;
    }
  },
  name: "Partition",
  components: {
    Card,
  },
  data() {
    return {
      swiperList: [
        { title: '精选强档', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917513797492.jpg' },
        { title: '番剧出击', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917531190280.jpg' },
        { title: '前方高能', pic: 'https://uploadstatic.mihoyo.com/contentweb/20210719/2021071917534678902.jpg' }
      ],
      animationList1: [],
      animationListFlag1: false,
      animationList2: [],
      animationListFlag2: false,
      animationList3: [],
      animationListFlag3: false,
    }
  },
  methods: {
    jumpPath(title) {
      console.log('Jumping to', title);
    }
  }
}
</script>`;

content = content.replace(/<script>[\s\S]*?<\/script>/, newScript);
fs.writeFileSync(path, content, 'utf8');
console.log('Partition patched');
