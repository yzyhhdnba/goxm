<template>
  <div class="outer">
    <div class="container">
      <div class="play"></div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { AdminAPI } from '@/api/index'
export default {
  created() {
    AdminAPI.getAreaStats().then(res => {
      const data = res || [];
      this.playOption.legend.data = data.map((item) => item.area_name);
      this.playOption.series = data.map((item) => ({
        name: item.area_name,
        type: 'line',
        stack: 'Total',
        areaStyle: {},
        emphasis: {
          focus: 'series'
        },
        data: [
          item.approved_count,
          item.pending_count,
          item.rejected_count,
          item.total_count,
        ]
      }));
      let chartDom = document.querySelector('.play')
      let myChart = echarts.init(chartDom);
      this.playOption && myChart.setOption(this.playOption);
    })
  },

  data() {
    return {
      playOption: {
        title: {
          text: '视频分区稿件统计'
        },
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            type: 'cross',
            label: {
              backgroundColor: '#6a7985'
            }
          }
        },
        legend: {
          data: []
        },
        toolbox: {
          feature: {
            saveAsImage: {}
          }
        },
        grid: {
          left: '3%',
          right: '4%',
          bottom: '3%',
          containLabel: true
        },
        xAxis: [
          {
            type: 'category',
            boundaryGap: false,
            data: ['通过', '待审', '驳回', '总稿件']
          }
        ],
        yAxis: [
          {
            type: 'value'
          }
        ],
        series: []
      }
    }
  }
};
</script>

<style lang="less" scoped>
.outer {
  display: flex;
  justify-content: center;
  flex-direction: column;

  .container {
    width: 80%;
    background-color: #fff;
    display: flex;
    align-items: center;

    .play {
      width: 850px;
      height: 600px;
    }
  }
}
</style>
