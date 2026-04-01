<template>
  <el-form ref="ruleFormRef" :model="ruleForm" :rules="rules" label-width="120px" class="demo-ruleForm" :size="formSize"
    status-icon>
    <el-form-item label="标题" prop="title">
      <el-input v-model="ruleForm.title" />
    </el-form-item>
    <el-form-item label="分区" prop="type">
      <el-select v-model="ruleForm.type" placeholder="请选择视频分区">
        <el-option :label="item.label" :value='item.value' v-for='item in options' :key='item.value'/>
      </el-select>
    </el-form-item>
    <el-form-item label="简介" prop="desc">
      <el-input v-model="ruleForm.desc" type="textarea" />
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="submitForm(ruleFormRef)">上传</el-button>
      <el-button @click="resetForm(ruleFormRef)">清空</el-button>
    </el-form-item>
  </el-form>
</template>

<script lang="ts">
import { reactive, ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { useStore } from 'vuex'
import { UploadAPI, VideoAPI } from '@/api/index'
// import imgsUpload from '../../components/upload/imgsUpload.vue'
export default {
  components: {

  },
  setup(props: any, ctx: any) {
    let options=ref<any[]>([]);
    VideoAPI.getAreas().then(res => {
      const list = (res as any) || [];
      options.value = list.map((item: any) => ({
        label: item.name,
        value: item.id,
      }));
    })
    const store = useStore();
    const formSize = ref('default')
    const ruleFormRef = ref<FormInstance>()
    const ruleForm = reactive({
      title: '',
      type: '',
      desc: '',
    })

    const rules = reactive<FormRules>({
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 3, max: 20, message: '标题长度应在3-20个字符之间', trigger: 'blur' },
      ],
      imgs: [{
        required: true,
      }],
      type: [
        {
          required: true,
          message: '请选择视频分区',
          trigger: 'change',
        },
      ],
      desc: [
        { required: true, message: '请填写视频简介', trigger: 'blur' },
      ],
    })
    // submitForm 是投稿两段式流程的第一步。
    // 对应文档“投稿链路：先建稿件，再传文件”，这里只创建元数据并把 videoId 写进 Store。
    const submitForm = async (formEl: FormInstance | undefined) => {
      if (!formEl) return
      await formEl.validate((valid, fields) => {
        if (valid) {
          UploadAPI.createVideoMetadata({
            title: ruleForm.title,
            area_id: Number(ruleForm.type),
            description: ruleForm.desc,
          }).then((res: any) => {
            store.commit('changeVideoId', {
              id: res.id,
              title: ruleForm.title,
              area_id: Number(ruleForm.type),
              description: ruleForm.desc,
            });
            store.commit('addVideoId', res.id);
            ctx.emit("finish", true);
          })
        } else {
          console.log('error submit!', fields)
        }
      })
    }

    const resetForm = (formEl: FormInstance | undefined) => {
      if (!formEl) return
      formEl.resetFields()
    }

    return {
      formSize, ruleFormRef, ruleForm, rules, submitForm, resetForm, options,
    }
  }
}

</script>
