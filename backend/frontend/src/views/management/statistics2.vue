<template>
  <div class="outer">
    <div class="container">
      <div class="play"></div>
      <div>
        <div class="play-buttom"></div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import axios from 'axios'
export default {
  created() {
    axios.get('http://172.20.10.6:8081/admin/AreaData').then(res=>{
      console.log(res);
      this.playOption.xAxis[0].data=Object.keys(res.data[0].map);
      let legends=[]
      for(let i=0; i<res.data.length; i++){
        legends.push(res.data[i].areaName)
        this.playOption.series[i].data=Object.values(res.data[i].map)
        this.playOption.series[i].name=res.data[i].areaName;
      }
      this.playOption.legend.data=legends
      let chartDom = document.querySelector('.play')
      let myChart = echarts.init(chartDom);
      this.playOption && myChart.setOption(this.playOption);
    })
  },

  methods: {

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
          data: ['测试一','测试二','测试三']
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
            data: []
          }
        ],
        yAxis: [
          {
            type: 'value'
          }
        ],
        series: [
          {
            name: '测试一',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
              focus: 'series'
            },
            data: []
          },
          {
            name: '测试二',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
              focus: 'series'
            },
            data: []
          },
          {
            name: '测试三',
            type: 'line',
            stack: 'Total',
            areaStyle: {},
            emphasis: {
              focus: 'series'
            },
            data: []
          }
        ]
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