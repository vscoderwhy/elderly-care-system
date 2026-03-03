<template>
  <div class="file-upload">
    <!-- 上传区域 -->
    <el-upload
      ref="uploadRef"
      v-model:file-list="fileList"
      :action="action"
      :headers="headers"
      :data="data"
      :name="name"
      :accept="accept"
      :limit="limit"
      :multiple="multiple"
      :disabled="disabled"
      :drag="drag"
      :with-credentials="withCredentials"
      :before-upload="beforeUpload"
      :on-progress="handleProgress"
      :on-success="handleSuccess"
      :on-error="handleError"
      :on-exceed="handleExceed"
      :on-change="handleChange"
      :on-remove="handleRemove"
      :auto-upload="autoUpload"
      :http-request="httpRequest"
      class="upload-component"
      :class="{ 'is-drag': drag }"
    >
      <!-- 拖拽上传 -->
      <template v-if="drag">
        <div class="upload-dragger">
          <el-icon class="upload-icon"><UploadFilled /></el-icon>
          <div class="upload-text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <div v-if="tip" class="upload-tip">{{ tip }}</div>
        </div>
      </template>

      <!-- 按钮上传 -->
      <template v-else>
        <el-button type="primary" :loading="uploading">
          <el-icon><Upload /></el-icon>
          {{ uploading ? '上传中...' : '选择文件' }}
        </el-button>
        <div v-if="tip" class="upload-tip">{{ tip }}</div>
      </template>

      <!-- 插槽：自定义上传内容 -->
      <template #tip v-if="$slots.tip">
        <slot name="tip"></slot>
      </template>
    </el-upload>

    <!-- 文件列表 -->
    <div v-if="showFileList && fileList.length > 0" class="file-list">
      <transition-group name="list">
        <div
          v-for="file in fileList"
          :key="file.uid"
          class="file-item"
          :class="{ 'is-success': file.status === 'success', 'is-error': file.status === 'error' }"
        >
          <!-- 文件图标 -->
          <div class="file-icon">
            <el-icon v-if="isImage(file)">
              <Picture />
            </el-icon>
            <el-icon v-else-if="isVideo(file)">
              <VideoCamera />
            </el-icon>
            <el-icon v-else-if="isPdf(file)">
              <Document />
            </el-icon>
            <el-icon v-else>
              <Files />
            </el-icon>
          </div>

          <!-- 文件信息 -->
          <div class="file-info">
            <div class="file-name" :title="file.name">
              {{ file.name }}
            </div>
            <div class="file-size">
              {{ formatSize(file.size) }}
            </div>
            <!-- 上传进度条 -->
            <el-progress
              v-if="file.status === 'uploading'"
              :percentage="parseInt(file.percentage as string)"
              :stroke-width="2"
            />
          </div>

          <!-- 操作按钮 -->
          <div class="file-actions">
            <!-- 预览按钮（图片） -->
            <el-button
              v-if="isImage(file) && file.status === 'success'"
              size="small"
              text
              @click="handlePreview(file)"
            >
              <el-icon><ZoomIn /></el-icon>
            </el-button>

            <!-- 下载按钮 -->
            <el-button
              v-if="file.status === 'success' && file.url"
              size="small"
              text
              @click="handleDownload(file)"
            >
              <el-icon><Download /></el-icon>
            </el-button>

            <!-- 删除按钮 -->
            <el-button
              v-if="!disabled"
              size="small"
              text
              type="danger"
              @click="handleRemove(file)"
            >
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </transition-group>
    </div>

    <!-- 图片预览 -->
    <el-image-viewer
      v-if="showPreview"
      :url-list="previewUrls"
      :initial-index="previewIndex"
      @close="showPreview = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ElMessage, type UploadInstance, type UploadProps, type UploadUserFile, type UploadRawFile, type UploadRequestOptions } from 'element-plus'
import {
  Upload,
  UploadFilled,
  Picture,
  VideoCamera,
  Document,
  Files,
  ZoomIn,
  Download,
  Delete
} from '@element-plus/icons-vue'
import { ElImageViewer } from 'element-plus'
import type { AxiosProgressEvent } from 'axios'

interface Props {
  modelValue?: string[]
  action?: string
  headers?: Record<string, string>
  data?: Record<string, any>
  name?: string
  accept?: string
  limit?: number
  maxSize?: number // MB
  multiple?: boolean
  disabled?: boolean
  drag?: boolean
  withCredentials?: boolean
  autoUpload?: boolean
  showFileList?: boolean
  tip?: string
}

const props = withDefaults(defineProps<Props>(), {
  action: '/api/upload/file',
  headers: () => ({}),
  data: () => ({}),
  name: 'file',
  accept: '*',
  limit: 10,
  maxSize: 100,
  multiple: false,
  disabled: false,
  drag: false,
  withCredentials: false,
  autoUpload: true,
  showFileList: true,
  tip: ''
})

const emit = defineEmits<{
  'update:modelValue': [urls: string[]]
  'change': [files: UploadUserFile[]]
  'success': [file: UploadUserFile, response: any]
  'error': [file: UploadUserFile, error: Error]
}>()

const uploadRef = ref<UploadInstance>()
const fileList = ref<UploadUserFile[]>([])
const uploading = ref(false)
const showPreview = ref(false)
const previewIndex = ref(0)

// 预览URL列表
const previewUrls = computed(() => {
  return fileList.value
    .filter(f => isImage(f) && f.status === 'success' && f.url)
    .map(f => f.url)
})

// 初始化文件列表
watch(
  () => props.modelValue,
  (urls) => {
    if (urls && urls.length > 0) {
      fileList.value = urls.map((url, index) => ({
        uid: Date.now() + index,
        name: getFileName(url),
        url,
        status: 'success',
        size: 0
      }))
    }
  },
  { immediate: true }
)

// 监听文件列表变化
watch(
  fileList,
  (files) => {
    uploading.value = files.some(f => f.status === 'uploading')

    const urls = files
      .filter(f => f.status === 'success' && f.url)
      .map(f => f.url)

    emit('update:modelValue', urls)
    emit('change', files)
  },
  { deep: true }
)

// 上传前校验
const beforeUpload: UploadProps['beforeUpload'] = (rawFile) => {
  // 文件大小校验
  if (rawFile.size > props.maxSize * 1024 * 1024) {
    ElMessage.error(`文件大小不能超过 ${props.maxSize}MB`)
    return false
  }

  // 文件类型校验
  if (props.accept !== '*') {
    const acceptTypes = props.accept.split(',').map(t => t.trim())
    const fileType = rawFile.type
    const fileName = rawFile.name
    const ext = fileName.substring(fileName.lastIndexOf('.'))

    const valid = acceptTypes.some(type => {
      if (type.startsWith('.')) {
        return ext.toLowerCase() === type.toLowerCase()
      }
      return fileType.includes(type.replace('*', ''))
    })

    if (!valid) {
      ElMessage.error(`文件类型不符合，请上传 ${props.accept} 格式的文件`)
      return false
    }
  }

  return true
}

// 上传进度
const handleProgress: UploadProps['onProgress'] = (evt: AxiosProgressEvent, file) => {
  file.percentage = evt.percent || 0
}

// 上传成功
const handleSuccess: UploadProps['onSuccess'] = (response, file) => {
  // 假设响应格式为 { code: 0, data: { url: 'xxx' } }
  const url = response?.data?.url || response?.url || file.url
  if (url) {
    file.url = url
  }
  emit('success', file, response)
}

// 上传失败
const handleError: UploadProps['onError'] = (error, file) => {
  emit('error', file, error as Error)
}

// 超出限制
const handleExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning(`最多只能上传 ${props.limit} 个文件`)
}

// 文件状态改变
const handleChange: UploadProps['onChange'] = (file, files) => {
  fileList.value = files
}

// 移除文件
const handleRemove: UploadProps['onRemove'] = (file) => {
  const index = fileList.value.findIndex(f => f.uid === file.uid)
  if (index > -1) {
    fileList.value.splice(index, 1)
  }
}

// 自定义上传方法
const httpRequest = (options: UploadRequestOptions) => {
  const { action, file, data, headers, onProgress, onSuccess, onError } = options

  const formData = new FormData()
  if (data) {
    Object.keys(data).forEach(key => {
      formData.append(key, data[key])
    })
  }
  formData.append(options.name || 'file', file)

  // TODO: 实际上传请求
  // axios.post(action, formData, {
  //   headers,
  //   onUploadProgress: (e) => {
  //     onProgress({ percent: (e.loaded * 100) / e.total })
  //   }
  // })
  // .then(res => onSuccess(res))
  // .catch(err => onError(err))

  // 模拟上传
  return new Promise((resolve) => {
    setTimeout(() => {
      onSuccess({ data: { url: URL.createObjectURL(file) } })
      resolve({})
    }, 1000)
  })
}

// 判断是否为图片
const isImage = (file: UploadUserFile): boolean => {
  return file.name?.match(/\.(jpg|jpeg|png|gif|webp|bmp)$/i) !== null ||
         file.type?.startsWith('image/')
}

// 判断是否为视频
const isVideo = (file: UploadUserFile): boolean => {
  return file.name?.match(/\.(mp4|webm|ogg|avi|mov)$/i) !== null ||
         file.type?.startsWith('video/')
}

// 判断是否为PDF
const isPdf = (file: UploadUserFile): boolean => {
  return file.name?.match(/\.pdf$/i) !== null ||
         file.type === 'application/pdf'
}

// 格式化文件大小
const formatSize = (bytes?: number): string => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return `${(bytes / Math.pow(k, i)).toFixed(2)} ${sizes[i]}`
}

// 获取文件名
const getFileName = (url: string): string => {
  const parts = url.split('/')
  return parts[parts.length - 1] || 'file'
}

// 预览
const handlePreview = (file: UploadUserFile) => {
  const index = previewUrls.value.findIndex(url => url === file.url)
  previewIndex.value = index >= 0 ? index : 0
  showPreview.value = true
}

// 下载
const handleDownload = (file: UploadUserFile) => {
  const link = document.createElement('a')
  link.href = file.url || ''
  link.download = file.name
  link.click()
}

// 暴露方法
defineExpose({
  submitUpload: () => uploadRef.value?.submit(),
  clearFiles: () => {
    fileList.value = []
    uploadRef.value?.clearFiles()
  },
  abort: (file?: UploadUserFile) => {
    uploadRef.value?.abort(file)
  }
})
</script>

<style scoped lang="scss">
.file-upload {
  .upload-component {
    &:not(.is-drag) {
      display: inline-block;
    }
  }

  .upload-dragger {
    padding: 40px;
    text-align: center;
    border: 1px dashed var(--border-color);
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: var(--primary-color);
    }
  }

  .upload-icon {
    font-size: 48px;
    color: var(--primary-color);
    margin-bottom: 16px;
  }

  .upload-text {
    font-size: 14px;
    color: var(--text-primary);

    em {
      color: var(--primary-color);
      font-style: normal;
    }
  }

  .upload-tip {
    margin-top: 8px;
    font-size: 12px;
    color: var(--text-secondary);
  }

  .file-list {
    margin-top: 16px;
  }

  .file-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    margin-bottom: 8px;
    transition: all 0.3s;

    &:hover {
      background: var(--bg-secondary);
    }

    &.is-success {
      border-left: 3px solid var(--success-color);
    }

    &.is-error {
      border-left: 3px solid var(--danger-color);
    }
  }

  .file-icon {
    font-size: 32px;
    color: var(--primary-color);
    flex-shrink: 0;
  }

  .file-info {
    flex: 1;
    min-width: 0;
  }

  .file-name {
    font-size: 14px;
    color: var(--text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .file-size {
    font-size: 12px;
    color: var(--text-secondary);
    margin-top: 4px;
  }

  .file-actions {
    display: flex;
    gap: 4px;
    flex-shrink: 0;
  }
}

// 列表动画
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}
</style>
