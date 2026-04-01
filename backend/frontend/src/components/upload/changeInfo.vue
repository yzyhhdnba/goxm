<template>
  <!-- Form -->
  <el-dialog v-model="isEditOpen.isopen.isopen" title="编辑视频">
    <el-form ref="ruleFormRef" :model="ruleForm" :rules="rules" label-width="120px" class="demo-ruleForm" :size="formSize"
    status-icon>
    <el-form-item label="标题" prop="title">
      <el-input v-model="ruleForm.title" />
    </el-form-item>
    <el-form-item label="封面" prop="imgs">
      <imgsUpload></imgsUpload>
    </el-form-item>
    <el-form-item label="分区" prop="type">
      <el-select v-model="ruleForm.type" placeholder="请选择视频分区">
        <el-option label="Zone one" value="shanghai" />
        <el-option label="Zone two" value="beijing" />
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
  </el-dialog>
</template>

<script lang="ts" setup>
import { reactive, ref,defineProps } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import imgsUpload from '@/components/upload/videoUpload.vue'
import axios from 'axios'
const isEditOpen=defineProps(['isopen'])
console.log(isEditOpen.isopen);
axios.get('http://172.20.10.6:8081/area').then(res=>{
  console.log(res);
})
const dialogFormVisible = ref(false)
const formLabelWidth = '140px'
const formSize = ref('default')
    const ruleFormRef = ref<FormInstance>()
    const ruleForm = reactive({
      title: '',
      imgs:'a',
      type: '',
      desc: '',
    })

    const rules = reactive<FormRules>({
      title: [
        { required: true, message: '请输入标题', trigger: 'blur' },
        { min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur' },
      ],
      imgs:[{
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

    const submitForm = async (formEl: FormInstance | undefined) => {
      if (!formEl) return
      await formEl.validate((valid, fields) => {
        if (valid) {
          console.log('submit!')
        } else {
          console.log('error submit!', fields)
        }
      })
    }

    const resetForm = (formEl: FormInstance | undefined) => {
      if (!formEl) return
      formEl.resetFields()
    }

    const options = Array.from({ length: 10000 }).map((_, idx) => ({
      value: `${idx + 1}`,
      label: `${idx + 1}`,
    }))

</script>
<style scoped>
.el-button--text {
  margin-right: 15px;
}

.el-select {
  width: 300px;
}

.el-input {
  width: 300px;
}

.dialog-footer button:first-child {
  margin-right: 10px;
}
</style>