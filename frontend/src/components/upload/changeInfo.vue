<template>
  <el-dialog v-model="dialogVisible" title="编辑视频" width="560px">
    <el-form label-width="80px">
      <el-form-item label="标题">
        <el-input v-model="form.title" maxlength="128" />
      </el-form-item>
      <el-form-item label="分区">
        <el-select v-model="form.area_id" placeholder="请选择分区" style="width: 100%;">
          <el-option v-for="item in areas" :key="item.id" :label="item.name" :value="item.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="简介">
        <el-input v-model="form.description" type="textarea" :rows="4" />
      </el-form-item>
    </el-form>
    <div class="notice">
      保存后稿件将重新进入待审核队列，原审核结果会被清空。
    </div>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submit">保存修改</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { computed, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { CreatorAPI, VideoAPI } from '@/api/index'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false,
  },
  video: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:visible', 'updated'])

const dialogVisible = computed({
  get: () => props.visible,
  set: (value: boolean) => emit('update:visible', value),
})

const areas = ref<any[]>([])
const form = reactive({
  title: '',
  area_id: undefined as number | undefined,
  description: '',
})

VideoAPI.getAreas().then((res: any) => {
  areas.value = res || []
}).catch((error: any) => {
  console.error('load areas failed', error)
})

watch(
  () => props.video,
  (video: any) => {
    form.title = video?.title || ''
    form.area_id = video?.area_id
    form.description = video?.description || ''
  },
  { immediate: true }
)

const submit = async () => {
  if (!props.video?.id) {
    return
  }
  if (!form.title.trim() || !form.area_id) {
    ElMessage.warning('请先补全标题和分区')
    return
  }

  try {
    await CreatorAPI.updateVideo(props.video.id, {
      title: form.title.trim(),
      area_id: Number(form.area_id),
      description: form.description.trim(),
    })
    ElMessage.success('稿件已更新，已重新进入审核队列')
    emit('updated')
    dialogVisible.value = false
  } catch (error) {
    console.error('update creator video failed', error)
  }
}
</script>

<style scoped>
.notice {
  color: #909399;
  line-height: 22px;
}
</style>
