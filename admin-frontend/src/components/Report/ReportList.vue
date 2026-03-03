<template>
  <div class="report-list">
    <el-table :data="reports" stripe style="width: 100%">
      <el-table-column prop="name" label="报表名称" min-width="200">
        <template #default="{ row }">
          <div class="report-name">
            <el-icon class="report-icon"><Document /></el-icon>
            <span>{{ row.name }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="300" show-overflow-tooltip />
      <el-table-column prop="updateAt" label="更新时间" width="120" />
      <el-table-column prop="schedule" label="生成周期" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.schedule" size="small">{{ row.schedule }}</el-tag>
          <span v-else class="text-muted">手动</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button size="small" text @click="$emit('view', row)">
            查看
          </el-button>
          <el-button size="small" text @click="handleGenerate(row)">
            生成
          </el-button>
          <el-dropdown trigger="click">
            <el-button size="small" text>
              更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="$emit('export', row)">
                  <el-icon><Download /></el-icon>
                  导出Excel
                </el-dropdown-item>
                <el-dropdown-item @click="handleEdit(row)">
                  <el-icon><Edit /></el-icon>
                  编辑
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleDelete(row)">
                  <el-icon><Delete /></el-icon>
                  删除
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="reports.length === 0" description="暂无报表" />
  </div>
</template>

<script setup lang="ts">
import { Document, ArrowDown, Download, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessageBox, ElMessage } from 'element-plus'

interface Report {
  id: string | number
  name: string
  category: string
  description: string
  updateAt: string
  schedule?: string
}

defineProps<{
  reports: Report[]
}>()

const emit = defineEmits<{
  view: [report: Report]
  export: [report: Report]
}>()

const handleGenerate = (report: Report) => {
  ElMessage.success(`正在生成报表: ${report.name}`)
}

const handleEdit = (report: Report) => {
  console.log('编辑报表', report)
}

const handleDelete = (report: Report) => {
  ElMessageBox.confirm(
    `确定要删除报表"${report.name}"吗？`,
    '删除确认',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(() => {
    ElMessage.success('删除成功')
  })
}
</script>

<style scoped lang="scss">
.report-list {
  .report-name {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .report-icon {
    color: var(--primary-color);
  }

  .text-muted {
    color: var(--text-tertiary);
  }
}
</style>
