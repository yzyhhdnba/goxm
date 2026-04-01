<template>
  <div class="index">
    <!-- 
    index类名
    -->
    <el-row justify="center" align="middle" :gutter="20">
      <el-col :span="15">
        <!-- 
        分配空间
        -->
        <div class="information">
          <!-- 
          information类名
          -->
          <p>
            <!--
            版权声明信息
            -->
            ©PiliPili&nbsp;&nbsp;&nbsp;
            <!--
            网备信息
            -->
            <a href="https://beian.miit.gov.cn/" target="_blank">
              皖ICP备2022005784号-1&nbsp;
            </a>
            <a
              target="_blank"
              href="http://www.beian.gov.cn/portal/registerSystemInfo?recordcode=34088102000589"
              style="
                display: inline-block;
                text-decoration: none;
                height: 20px;
                line-height: 20px;
              "
            >
              <img
                alt="公安网备"
                src="https://s3.bmp.ovh/imgs/2022/04/11/ec64e41ab1c48dab.png"
                style="float: left"
              />
              <span>&nbsp;皖公网安备 34088102000589号</span>
            </a>
          </p>

          <p>
            <!--
            版权声明信息
            -->
            本站内原创内容版权属于其原创作者,除作者或特别声明外未经许可不可转载
          </p>
          <div class="text-con__line">
            <i class="item-icon-1"></i>
            <a href="http://www.12377.cn/" target="_blank"
              >网上有害信息举报专区丨中央网信办（国家网信办）违法和不良信息举报中心</a
            >&nbsp;&nbsp;&nbsp;
            <a href="http://www.cyberpolice.cn" target="_blank"
              >公安部网络违法犯罪举报网站</a
            >
          </div>
        </div>
      </el-col>
      <el-col :span="1">
        <!--
        分配1空间,存放分隔符
        -->
        <el-divider direction="vertical" />
      </el-col>
      <el-col :span="1">
        <!--
        分配1空间,存放vue-logo
        -->
        <a href="https://v3.cn.vuejs.org/" target="_blank">
          <img
            class="vue-logo"
            alt="Vue3"
            src="https://s3.bmp.ovh/imgs/2022/06/27/ddbd16f57446e396.png"
            title="本站使用Vue3框架开发"
          />
        </a>
      </el-col>
      <el-col :span="1" class="nav-plus">
        <el-icon :size="24" class="is-loading">
          <circle-plus-filled />
        </el-icon>
      </el-col>
      <el-col :span="2">
        <!--
        分配2空间,存放element-logo
        -->
        <a href="https://element-plus.gitee.io/zh-CN/" target="_blank">
          <img
            class="element-logo"
            alt="Element Plus"
            src="https://s3.bmp.ovh/imgs/2022/06/27/f70647f9f1fb5dd3.png"
            title="本站使用Element Plus组件库开发"
          />
        </a>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { CirclePlusFilled } from "@element-plus/icons-vue";

/**
 * @element-plus/icons-vue: vue3中element-plus全局使用icon图标
 * @param NULL
 * @param NULL
 * @return void
 */
export default {
  // eslint-disable-next-line vue/multi-word-component-names
  name: "Foot",
  props: {},
  components: {
    CirclePlusFilled,
  },
  setup() {
    return {};
  },
  methods: {},
};
</script>

<style scoped>
.index,
a {
  color: white;
  text-decoration: none;
}
.nav-plus {
  color: gray;
  text-align: center;
}

img.vue-logo {
  height: 100%;
  width: 100%;
}

img.element-logo {
  height: 100%;
  width: 100%;
}

.el-affix {
  height: 70px;
}

.el-row {
  height: 150px;
  background-color: #f6f7f8;
  opacity: 0.9;
}

.el-divider {
  height: 100px;
}

.information {
  text-align: right;
  color: white;
}

.index {
  margin: -18px;
  opacity: 100%;
  height: 70px;
  z-index: 114514;
  --el-header-padding: 0;
}

.item-icon-1 {
  color: #999;
  margin: 0;
  padding: 0;
  border: 0;
  font: inherit;
  font-size: 100%;
  font-style: italic;
  background: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACAAAAAQCAMAAABA3o1rAAABTVBMVEUAAADWAA8EnOfWAA/YAA8Ald7nAAHfAAgAluDYAA8Anu0AkuIAl94Alt8CoO8Am+LYAA8GmeUAo/TUBhHgAQfWAgvgAAcAieToAANbWosAnOU+dKzjAATVAA/VABAAj9kAneXTAA/lAAHsAAAAmOcDoe8Aq//QChQAluDwAAHaChAApejHEh0AmucAld/VABEAe75TW40Ahc3UAAYAl+rxAAZgWIwAn+LJAAbwAAZSYJRmVYMnis4AqPyrGyneAAoAqv3WABG1KT213evWAA/S6vMAluD6/PwAh7vg8fcFi74BdrDC4+9MqM7q9fr1+/3N6PIBgbTy+f3t9/vc7/bJ5fCNyuF4wNwxnsj////w+Puq1+cAleBpt9X2/f6Z0uYAmeM1l8QelsMakL3TABDl8PuOw93X6fj98vEAn+0AneoAkd/tcnLXABHZABBZG46EAAAAQ3RSTlMAzybp3/399+zV0dHMyL25tbCrjIVqXkFBQDIVDPvy39jY1cLAsp2dm5iTko2MiHJhWVNTUVBNSEhFNjQhIB4ZFxML8b/YSQAAAWlJREFUKM91kcVWY1EURCsQaNy1od3d3fu6Pbe4Aa3w/0PCW8kM9qzOXmdShfMQsQqCpC5wAR5JG78UscTDmIlHi4tzT/ZQ4jeimtewfhQ1/LFf/cf/LL34WoaEkSRzhel16jZNytPl5almd3V3/dYOABEQz4WDWu5CF7HUP/O3m83u2seH0yeV94BgXkubrOa1jM7rLAK2Fzjna59/PrhCj2+8gmLQskey3GipwRK8nud85v4X4Nv6JKXVLQsmddHL2qGRmiHF8vB/YQJDfqxQSisWVvYLZ48OW6G0SPH8WpfPPP0O7D+rUjr5XykhtevkxtUK01YKeHmd86m7e/sbV4d+87cI4GTYOur0w75DUAfwbr75d+nTxiV6UN0CBGnHUkszyAfG80lZ1fbcnd3HQz/75iwpIpJC6tDksSBq1OTbe9P0oPJhVDUBgiwLABJhxObN2crKzjgdEhYDHiv/L5ibpZbF5dynAIdWsI2ALm4AAAAASUVORK5CYII=")
    no-repeat 0 0;
  display: inline-block;
  width: 16px;
  height: 16px;
  margin-right: 8px;
  vertical-align: middle;
}

a,
p {
  color: #999;
  margin: 0;
  padding: 0;
  border: 0;
  font: inherit;
  font-size: 100%;
  vertical-align: baseline;
}
p {
  padding-bottom: 10px;
}
</style>
