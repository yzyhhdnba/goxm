<template>
  <Nav></Nav>
  <el-row class="tac">
    <el-col :span="5">
      <el-menu class="el-menu-vertical-demo" @open="handleOpen" @close="handleClose" :router="true"
        style="background-color:transparent;">
        <el-menu-item index="/management/videomanage1">
          <el-icon>
            <document />
          </el-icon>
          <span>未审核</span>
        </el-menu-item>
        <el-menu-item index="/management/videomanage2">
          <el-icon>
            <document />
          </el-icon>
          <span>已审核</span>
        </el-menu-item>
        <el-menu-item index="/management/statistics">
          <el-icon>
            <document />
          </el-icon>
          <span>本站统计信息</span>
        </el-menu-item>
        <el-menu-item index="/management/statistics2">
          <el-icon>
            <document />
          </el-icon>
          <span>视频分区统计</span>
        </el-menu-item>
      </el-menu>
    </el-col>
    <el-col :span="18" class="divider">
      <router-view style="width:100%"></router-view>
    </el-col>
  </el-row>
</template>

<script lang="ts">
import {
  Location,
  Document,
  Menu as IconMenu,
  Setting,
} from '@element-plus/icons-vue'
import Nav from '../../components/management/Nav.vue'
import { onMounted } from 'vue'
import { useStore } from 'vuex'
export default {
  components: {
    Nav
  },
  setup() {
    const store = useStore();
    const role = localStorage.getItem('role');
    if (localStorage.getItem('isLogin') && store.state.role === '') {
      store.commit('changeRole', localStorage.getItem('role'))
    }
    onMounted(() => {
      document.body.style.backgroundImage = "url('http://101.35.142.191:8081/image/background1.png')";
      document.body.style.backgroundSize = '100%';
      document.body.style.margin = '0 0 0 0';
      document.body.style.padding = '0 0 0 0';
    })
    const handleOpen = (key: string, keyPath: string[]) => {
      console.log(key, keyPath)
    }
    const handleClose = (key: string, keyPath: string[]) => {
      console.log(key, keyPath)
    }
    return {
      handleOpen, handleClose, role
    }
  }

}
</script>
<style scoped>
.divider {
  margin-left: 20px;
}
</style>