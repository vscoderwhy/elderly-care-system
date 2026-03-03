<template>
  <div class="image-upload">
    <!-- 已上传的图片列表 -->
    <draggable
      v-model="fileList"
      :disabled="!dragSort"
      item-key="uid"
      class="upload-list"
      animation="200"
      @end="handleDragEnd"
    >
      <template #item="{ element, index }">
        <div class="upload-item" :class="{ 'is-uploading': element.status === 'uploading' }">
          <!-- 图片预览 -->
          <el-image
            :src="element.url"
            :preview-src-list="previewList"
            :initial-index="index"
            fit="cover"
            class="upload-image"
          >
            <template #error>
              <div class="image-error">
                <el-icon><PictureFilled /></el-icon>
              </div>
            </template>
          </el-image>

          <!-- 遮罩层 -->
          <div class="upload-mask">
            <span class="mask-icon" @click="handlePreview(element)">
              <el-icon><ZoomIn /></el-icon>
            </span>
            <span v-if="!disabled" class="mask-icon" @click="handleRemove(element)">
              <el-icon><Delete /></el-icon>
            </span>
          </div>

          <!-- 上传进度 -->
          <div v-if="element.status === 'uploading'" class="upload-progress">
            <el-progress
              type="circle"
              :percentage="element.percent"
              :width="40"
            />
          </div>
        </div>
      </template>
    </draggable>

    <!-- 上传按钮 -->
    <div
      v-if="!disabled && fileList.length < limit"
      class="upload-trigger"
      :class="{ 'is-dragover': isDragOver }"
      @click="handleSelect"
      @drop.prevent="handleDrop"
      @dragover.prevent="isDragOver = true"
      @dragleave.prevent="isDragOver = false"
    >
      <el-icon class="upload-icon"><Plus /></el-icon>
      <div class="upload-text">点击或拖拽上传</div>
      <div v-if="tip" class="upload-tip">{{ tip }}</div>
    </div>

    <!-- 隐藏的文件输入 -->
    <input
      ref="fileInputRef"
      type="file"
      :accept="accept"
      :multiple="multiple"
      class="upload-input"
      @change="handleFileChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, ZoomIn, Delete, PictureFilled } from '@element-plus/icons-vue'
import draggable from 'vuedraggable'
import type { UploadProps, UploadUserFile } from 'element-plus'

interface UploadFile extends UploadUserFile {
  uid: number
  url: string
  status?: 'ready' | 'uploading' | 'success' | 'error'
  percent?: number
}

interface Props {
  modelValue: string[]
  accept?: string
  limit?: number
  maxSize?: number // MB
  multiple?: boolean
  disabled?: boolean
  dragSort?: boolean
  tip?: string
  action?: string
}

const props = withDefaults(defineProps<Props>(), {
  accept: 'image/*',
  limit: 9,
  maxSize: 5,
  multiple: false,
  disabled: false,
  dragSort: true,
  tip: '',
  action: '/api/upload/image'
})

const emit = defineEmits<{
  'update:modelValue': [urls: string[]]
  'change': [files: UploadFile[]]
}>()

const fileInputRef = ref<HTMLInputElement>()
const isDragOver = ref(false)

// 文件列表
const fileList = ref<UploadFile[]>([])

// 预览列表
const previewList = computed(() => {
  return fileList.value
    .filter(item => item.status === 'success' && item.url)
    .map(item => item.url)
})

// 初始化文件列表
watch(
  () => props.modelValue,
  (urls) => {
    if (urls && urls.length > 0) {
      fileList.value = urls.map((url, index) => ({
        uid: Date.now() + index,
        name: `image-${index}`,
        url,
        status: 'success',
        percent: 100
      }))
    } else {
      fileList.value = []
    }
  },
  { immediate: true }
)

// 监听文件列表变化，同步到父组件
watch(
  fileList,
  (files) => {
    const urls = files
      .filter(f => f.status === 'success' && f.url)
      .map(f => f.url)
    emit('update:modelValue', urls)
    emit('change', files)
  },
  { deep: true }
)

// 选择文件
const handleSelect = () => {
  fileInputRef.value?.click()
}

// 文件选择变化
const handleFileChange = (e: Event) => {
  const target = e.target as HTMLInputElement
  const files = Array.from(target.files || [])
  handleFiles(files)
  target.value = '' // 清空输入，允许重复选择同一文件
}

// 拖拽上传
const handleDrop = (e: DragEvent) => {
  isDragOver.value = false
  const files = Array.from(e.dataTransfer?.files || [])
  handleFiles(files)
}

// 处理文件
const handleFiles = (files: File[]) => {
  // 检查数量限制
  if (fileList.value.length + files.length > props.limit) {
    ElMessage.warning(`最多只能上传 ${props.limit} 个文件`)
    return
  }

  files.forEach(file => {
    // 检查文件类型
    if (!file.type.startsWith('image/')) {
      ElMessage.error('只能上传图片文件')
      return
    }

    // 检查文件大小
    if (file.size > props.maxSize * 1024 * 1024) {
      ElMessage.error(`图片大小不能超过 ${props.maxSize}MB`)
      return
    }

    // 创建预览
    const reader = new FileReader()
    reader.onload = (e) => {
      const uploadFile: UploadFile = {
        uid: Date.now(),
        name: file.name,
        url: e.target?.result as string,
        status: 'uploading',
        percent: 0,
        raw: file
      }
      fileList.value.push(uploadFile)
      uploadFile(uploadFile)
    }
    reader.readAsDataURL(file)
  })
}

// 上传文件
const uploadFile = async (uploadFile: UploadFile) => {
  try {
    const formData = new FormData()
    formData.append('file', uploadFile.raw as File)

    // 模拟上传进度
    const timer = setInterval(() => {
      if (uploadFile.percent && uploadFile.percent < 90) {
        uploadFile.percent += 10
      }
    }, 200)

    // TODO: 实际上传请求
    // const res = await axios.post(props.action, formData, {
    //   onUploadProgress: (e) => {
    //     uploadFile.percent = Math.round((e.loaded * 100) / e.total)
    //   }
    // })

    // 模拟上传成功
    await new Promise(resolve => setTimeout(resolve, 1500))
    clearInterval(timer)

    uploadFile.status = 'success'
    uploadFile.percent = 100
    uploadFile.url = uploadFile.url || 'https://via.placeholder.com/200'

    ElMessage.success('上传成功')
  } catch (error) {
    uploadFile.status = 'error'
    ElMessage.error('上传失败')
  }
}

// 预览
const handlePreview = (file: UploadFile) => {
  // 预览功能由 el-image 组件处理
}

// 移除
const handleRemove = (file: UploadFile) => {
  const index = fileList.value.findIndex(f => f.uid === file.uid)
  if (index > -1) {
    fileList.value.splice(index, 1)
  }
}

// 拖拽排序结束
const handleDragEnd = () => {
  // 排序已通过 v-model 自动同步
}

// 暴露方法
defineExpose({
  clearFiles: () => {
    fileList.value = []
  },
  submitUpload: () => {
    fileList.value
      .filter(f => f.status === 'ready')
      .forEach(f => uploadFile(f))
  }
})
</script>

<style scoped lang="scss">
.image-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.upload-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.upload-item {
  position: relative;
  width: 104px;
  height: 104px;
  border-radius: 4px;
  overflow: hidden;
  border: 1px solid var(--border-color-lighter);
  cursor: move;

  &.is-uploading {
    cursor: not-allowed;
  }
}

.upload-image {
  width: 100%;
  height: 100%;

  :deep(.el-image__inner) {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

.image-error {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  background: var(--bg-tertiary);
  color: var(--text-tertiary);
  font-size: 32px;
}

.upload-mask {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  background: rgba(0, 0, 0, 0.6);
  opacity: 0;
  transition: opacity 0.3s;

  .upload-item:hover & {
    opacity: 1;
  }
}

.mask-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  font-size: 18px;
  cursor: pointer;

  &:hover {
    background: rgba(255, 255, 255, 0.3);
  }
}

.upload-progress {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
}

.upload-trigger {
  position: relative;
  width: 104px;
  height: 104px;
  border-radius: 4px;
  border: 1px dashed var(--border-color);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;

  &:hover,
  &.is-dragover {
    border-color: var(--primary-color);
    background: var(--el-fill-color-light);
  }
}

.upload-icon {
  font-size: 28px;
  color: var(--text-tertiary);
  margin-bottom: 8px;
}

.upload-text {
  font-size: 12px;
  color: var(--text-secondary);
  text-align: center;
  line-height: 1.4;
}

.upload-tip {
  font-size: 10px;
  color: var(--text-tertiary);
  text-align: center;
  line-height: 1.4;
  margin-top: 4px;
}

.upload-input {
  display: none;
}
</style>
