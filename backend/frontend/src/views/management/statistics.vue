<template>
  <div class="outer">
    <h2>视频总计</h2>
    <div class="container">
      <div class="play"></div>
      <div class="right">
        <div>
          <h2>活跃用户</h2>
          <div class="right-text">今日活跃用户为{{ playOption.series[0].data[0].value[3] }}</div>
        </div>
        <div>
          <h2>通过稿件</h2>
          <div class="right-text">今日通过稿件为{{ playOption.series[0].data[0].value[2] }}</div>
        </div>
        <div>
          <h2>投稿</h2>
          <div class="right-text">今日投稿总数为{{ playOption.series[0].data[0].value[0] }}</div>
        </div>
        <div>
          <h2>播放</h2>
          <div class="right-text">今日本站播放为{{ playOption.series[0].data[0].value[1] }}</div>
        </div>
      </div>
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
    axios.get('http://172.20.10.6:8081/admin/TodayData').then(res => {
      console.log(res);
      this.dataValue = Object.values(res.data)

      this.playOption.series[0].data[0].value = this.dataValue;
      console.log(this.playOption.series[0].data[0].value);
      let chartDom = document.querySelector('.play')
      let myChart = echarts.init(chartDom);
      this.playOption && myChart.setOption(this.playOption);
    })
  },
  mounted() {
  },
  methods: {
    createPlay() {
      let chartDom = document.querySelector('.play')
      let myChart = echarts.init(chartDom);
      this.playOption && myChart.setOption(this.playOption);
    }
  },
  data() {
    return {
      dataValue: [0, 0, 0, 0],
      playOption: {
        title: {
          text: ''
        },
        legend: {
          data: ['今日数据', '历史数据']
        },
        radar: {
          shape: 'circle',
          indicator: [
            { name: '活跃', max: 30 },
            { name: '投稿', max: 30 },
            { name: '通过', max: 30 },
            { name: '播放', max: 30 },
          ]
        },
        series: [
          {
            name: 'Budget vs spending',
            type: 'radar',
            emphasis: {
              label: {
                show: true,
                color: '#fff',
                fontSize: 12,
                formatter: '{c}', // 鼠标悬浮时展示数据加上单位
                backgroundColor: '#0D1B42',
                borderRadius: 5,
                padding: 3,
                shadowBlur: 3
              }
            },
            data: [
              {
                value: [0, 0, 0, 0],
                name: '今日数据'
              },
              {
                value: [25, 25, 25, 25],
                name: '历史数据'
              }
            ]
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
    width: 70%;
    background-color: #fff;
    display: flex;
    align-items: center;

    .play {
      width: 450px;
      height: 500px;
    }

    .right {
      display: flex;
      flex-direction: column;
      height: 450px;
      margin-left: 100px;
      justify-content: space-around;

      .right-text {
        margin-top: 10px;
      }
    }
  }

}
</style>