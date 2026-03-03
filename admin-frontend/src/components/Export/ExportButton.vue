<template>
  <div class="export-button-wrapper">
    <el-dropdown trigger="click" @command="handleCommand">
      <el-button type="primary" :loading="exporting">
        <el-icon><Download /></el-icon>
        {{ exporting ? '导出中...' : '导出' }}
        <el-icon class="el-icon--right"><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="current">
            <el-icon><Document /></el-icon>
            导出当前页
          </el-dropdown-item>
          <el-dropdown-item command="all">
            <el-icon><Files /></el-icon>
            导出全部数据
          </el-dropdown-item>
          <el-dropdown-item command="selected" :disabled="!hasSelected">
            <el-icon><Checked /></el-icon>
            导出选中项 ({{ selectedCount }})
          </el-dropdown-item>
          <el-dropdown-item divided command="custom">
            <el-icon><Setting /></el-icon>
            自定义导出
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>

    <!-- 自定义导出对话框 -->
    <el-dialog
      v-model="customDialogVisible"
      title="自定义导出"
      width="600px"
    >
      <el-form :model="customForm" label-width="100px">
        <el-form-item label="导出范围">
          <el-radio-group v-model="customForm.range">
            <el-radio label="current">当前页</el-radio>
            <el-radio label="all">全部数据</el-radio>
            <el-radio label="custom">自定义范围</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="customForm.range === 'custom'" label="数据范围">
          <el-input-number
            v-model="customForm.start"
            :min="1"
            placeholder="起始行"
            style="width: 120px"
          />
          <span style="margin: 0 8px">-</span>
          <el-input-number
            v-model="customForm.end"
            :min="1"
            placeholder="结束行"
            style="width: 120px"
          />
        </el-form-item>

        <el-form-item label="导出字段">
          <el-checkbox-group v-model="customForm.fields">
            <el-checkbox
              v-for="col in columns"
              :key="col.key"
              :label="col.key"
            >
              {{ col.title }}
            </el-checkbox>
          </el-checkbox-group>
          <div style="margin-top: 8px">
            <el-link type="primary" @click="selectAllFields">全选</el-link>
            <el-divider direction="vertical" />
            <el-link type="primary" @click="clearAllFields">清空</el-link>
          </div>
        </el-form-item>

        <el-form-item label="文件名">
          <el-input v-model="customForm.filename" placeholder="请输入文件名" />
        </el-form-item>

        <el-form-item label="导出格式">
          <el-radio-group v-model="customForm.format">
            <el-radio label="xlsx">Excel (.xlsx)</el-radio>
            <el-radio label="csv">CSV (.csv)</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="customDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCustomExport" :loading="exporting">
          导出
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Download,
  ArrowDown,
  Document,
  Files,
  Checked,
  Setting
} from '@element-plus/icons-vue'
import { exportExcel } from '@/utils/export'

interface Column {
  key: string
  title: string
  width?: number
}

interface Props {
  // 数据总数
  total?: number
  // 当前页数据
  data?: any[]
  // 选中项数据
  selectedData?: any[]
  // 列配置
  columns?: Column[]
  // 默认文件名
  filename?: string
  // 导出函数（自定义导出逻辑）
  exportFunction?: (params: {
    type: string
    data: any[]
    columns?: Column[]
    filename: string
  }) => Promise<void> | void
}

const props = withDefaults(defineProps<Props>(), {
  total: 0,
  data: () => [],
  selectedData: () => [],
  columns: () => [],
  filename: '导出数据'
})

const exporting = ref(false)
const customDialogVisible = ref(false)
const selectedCount = ref(0)

// 自定义导出表单
const customForm = ref({
  range: 'current',
  start: 1,
  end: 100,
  fields: [] as string[],
  filename: props.filename,
  format: 'xlsx'
})

// 是否有选中项
const hasSelected = computed(() => props.selectedData.length > 0)

// 监听选中项变化
const updateSelectedCount = () => {
  selectedCount.value = props.selectedData.length
}

// 导出命令处理
const handleCommand = async (command: string) => {
  switch (command) {
    case 'current':
      await exportCurrentPage()
      break
    case 'all':
      await exportAll()
      break
    case 'selected':
      await exportSelected()
      break
    case 'custom':
      showCustomDialog()
      break
  }
}

// 导出当前页
const exportCurrentPage = async () => {
  if (props.data.length === 0) {
    ElMessage.warning('当前页暂无数据')
    return
  }

  exporting.value = true
  try {
    if (props.exportFunction) {
      await props.exportFunction({
        type: 'current',
        data: props.data,
        columns: props.columns,
        filename: props.filename
      })
    } else {
      exportExcel({
        data: props.data,
        columns: props.columns,
        filename: `${props.filename}_当前页`
      })
    }
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 导出全部
const exportAll = async () => {
  if (props.total === 0) {
    ElMessage.warning('暂无数据可导出')
    return
  }

  exporting.value = true
  try {
    // 全部导出需要从服务器获取所有数据
    if (props.exportFunction) {
      await props.exportFunction({
        type: 'all',
        data: [],
        columns: props.columns,
        filename: props.filename
      })
    } else {
      ElMessage.warning('请实现 exportFunction 方法以导出全部数据')
    }
  } catch (error) {
    console.error('导出失败', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 导出选中项
const exportSelected = async () => {
  if (props.selectedData.length === 0) {
    ElMessage.warning('请先选择要导出的数据')
    return
  }

  exporting.value = true
  try {
    if (props.exportFunction) {
      await props.exportFunction({
        type: 'selected',
        data: props.selectedData,
        columns: props.columns,
        filename: props.filename
      })
    } else {
      exportExcel({
        data: props.selectedData,
        columns: props.columns,
        filename: `${props.filename}_选中项`
      })
    }
    ElMessage.success('导出成功')
  } catch (error) {
    console.error('导出失败', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 显示自定义导出对话框
const showCustomDialog = () => {
  // 初始化字段选择（默认全选）
  customForm.value.fields = props.columns.map(col => col.key)
  customForm.value.filename = props.filename
  customDialogVisible.value = true
}

// 全选字段
const selectAllFields = () => {
  customForm.value.fields = props.columns.map(col => col.key)
}

// 清空字段
const clearAllFields = () => {
  customForm.value.fields = []
}

// 自定义导出
const handleCustomExport = async () => {
  if (customForm.value.fields.length === 0) {
    ElMessage.warning('请至少选择一个导出字段')
    return
  }

  // 筛选列
  const selectedColumns = props.columns.filter(col =>
    customForm.value.fields.includes(col.key)
  )

  exporting.value = true
  try {
    let data: any[] = []

    // 根据范围获取数据
    switch (customForm.value.range) {
      case 'current':
        data = props.data
        break
      case 'all':
      case 'custom':
        // 需要从服务器获取
        if (props.exportFunction) {
          await props.exportFunction({
            type: customForm.value.range,
            data: [],
            columns: selectedColumns,
            filename: customForm.value.filename,
            range: customForm.value.range === 'custom'
              ? { start: customForm.value.start, end: customForm.value.end }
              : undefined
          })
        } else {
          ElMessage.warning('请实现 exportFunction 方法')
          exporting.value = false
          return
        }
        customDialogVisible.value = false
        return
    }

    exportExcel({
      data,
      columns: selectedColumns,
      filename: customForm.value.filename
    })

    ElMessage.success('导出成功')
    customDialogVisible.value = false
  } catch (error) {
    console.error('导出失败', error)
    ElMessage.error('导出失败')
  } finally {
    exporting.value = false
  }
}

// 暴露方法供外部调用
defineExpose({
  exportCurrentPage,
  exportAll,
  exportSelected,
  updateSelectedCount
})
</script>

<style scoped lang="scss">
.export-button-wrapper {
  display: inline-block;
}

:deep(.el-checkbox-group) {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}
</style>
