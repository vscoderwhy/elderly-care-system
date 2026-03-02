<template>
  <div class="elderly-create">
    <el-card>
      <template #header>
        <span>新增老人档案</span>
      </template>

      <el-form :model="form" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="form.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="性别" prop="gender">
          <el-radio-group v-model="form.gender">
            <el-radio label="男">男</el-radio>
            <el-radio label="女">女</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="身份证号" prop="id_card">
          <el-input v-model="form.id_card" placeholder="请输入身份证号" />
        </el-form-item>
        <el-form-item label="联系电话" prop="phone">
          <el-input v-model="form.phone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="紧急联系人" prop="emergency_contact">
          <el-input v-model="form.emergency_contact" placeholder="请输入紧急联系人" />
        </el-form-item>
        <el-form-item label="紧急联系电话" prop="emergency_phone">
          <el-input v-model="form.emergency_phone" placeholder="请输入紧急联系电话" />
        </el-form-item>
        <el-form-item label="护理等级" prop="care_level">
          <el-select v-model="form.care_level" placeholder="请选择护理等级">
            <el-option label="一级护理" :value="1" />
            <el-option label="二级护理" :value="2" />
            <el-option label="三级护理" :value="3" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
          <el-button @click="handleCancel">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance } from 'element-plus'
import { elderlyApi } from '@/api'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  name: '',
  gender: '',
  id_card: '',
  phone: '',
  emergency_contact: '',
  emergency_phone: '',
  care_level: 1
})

const rules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  gender: [{ required: true, message: '请选择性别', trigger: 'change' }],
  phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await elderlyApi.create(form)
        ElMessage.success('创建成功')
        router.push('/elderly')
      } catch (error) {
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleCancel = () => {
  router.push('/elderly')
}
</script>
