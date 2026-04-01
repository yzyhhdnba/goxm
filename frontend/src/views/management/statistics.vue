<template>
  <div class="outer">
    <h2>视频总计</h2>
    <div class="container">
      <div class="play"></div>
      <div class="right">
        <div>
          <h2>活跃用户</h2>
          <div class="right-text">今日活跃用户为{{ stats.active_user_count }}</div>
        </div>
        <div>
          <h2>通过稿件</h2>
          <div class="right-text">今日通过稿件为{{ stats.approved_video_count }}</div>
        </div>
        <div>
          <h2>投稿</h2>
          <div class="right-text">今日投稿总数为{{ stats.submitted_video_count }}</div>
        </div>
        <div>
          <h2>播放</h2>
          <div class="right-text">今日本站播放为{{ stats.play_count }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import * as echarts from 'echarts';
import { AdminAPI } from '@/api/index'
export default {
  created() {
    AdminAPI.getTodayStats().then(res => {
      this.stats = res || this.stats;
      this.playOption.series[0].data[0].value = [
        this.stats.active_user_count,
        this.stats.submitted_video_count,
        this.stats.approved_video_count,
        this.stats.play_count,
      ];
      let chartDom = document.querySelector('.play')
      let myChart = echarts.init(chartDom);
      this.playOption && myChart.setOption(this.playOption);
    })
  },
  data() {
    return {
      stats: {
        active_user_count: 0,
        submitted_video_count: 0,
        approved_video_count: 0,
        play_count: 0,
        comment_count: 0,
      },
      playOption: {
        title: {
          text: ''
        },
        legend: {
          data: ['今日数据', '基准线']
        },
        radar: {
          shape: 'circle',
          indicator: [
            { name: '活跃', max: 50 },
            { name: '投稿', max: 50 },
            { name: '通过', max: 50 },
            { name: '播放', max: 50 },
          ]
        },
        series: [
          {
            name: 'today-stats',
            type: 'radar',
            emphasis: {
              label: {
                show: true,
                color: '#fff',
                fontSize: 12,
                formatter: '{c}',
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
                value: [10, 10, 10, 10],
                name: '基准线'
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
