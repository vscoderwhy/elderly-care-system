<template>
  <div class="care-records">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Document /></el-icon>
          护理记录
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新增记录
        </el-button>
        <ExportButton
          :data="recordList"
          :columns="exportColumns"
          :total="total"
          filename="护理记录"
          :export-function="handleExport"
        />
      </div>
    </div>

    <!-- 筛选表单 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="老人">
          <el-select
            v-model="filterForm.elderlyId"
            placeholder="请选择"
            clearable
            filterable
          >
            <el-option
              v-for="elderly in elderlyList"
              :key="elderly.id"
              :label="`${elderly.name} (${elderly.bedNumber})`"
              :value="elderly.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="护理类型">
          <el-select v-model="filterForm.careType" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="日常护理" value="日常护理" />
            <el-option label="康复训练" value="康复训练" />
            <el-option label="健康监测" value="健康监测" />
            <el-option label="医疗护理" value="医疗护理" />
          </el-select>
        </el-form-item>
        <el-form-item label="记录人">
          <el-select v-model="filterForm.nurseId" placeholder="请选择" clearable filterable>
            <el-option
              v-for="nurse in nurses"
              :key="nurse.id"
              :label="nurse.name"
              :value="nurse.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleFilter">查询</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 记录列表 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="recordList"
        stripe
        style="width: 100%"
      >
        <el-table-column prop="recordTime" label="记录时间" width="180" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="bedNumber" label="床位" width="100" />
        <el-table-column prop="careType" label="护理类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeColor(row.careType)" size="small">
              {{ row.careType }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="careContent" label="护理内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="nurseName" label="记录人" width="100" />
        <el-table-column prop="result" label="护理结果" width="120" show-overflow-tooltip />
        <el-table-column prop="evaluation" label="评价" width="80">
          <template #default="{ row }">
            <el-rate v-model="row.evaluation" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleEdit(row)">编辑</el-button>
            <el-button size="small" text type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="护理记录详情" width="700px">
      <el-descriptions v-if="currentRecord" :column="2" border>
        <el-descriptions-item label="记录时间" :span="2">
          {{ currentRecord.recordTime }}
        </el-descriptions-item>
        <el-descriptions-item label="老人姓名">
          {{ currentRecord.elderlyName }}
        </el-descriptions-item>
        <el-descriptions-item label="床位号">
          {{ currentRecord.bedNumber }}
        </el-descriptions-item>
        <el-descriptions-item label="护理类型">
          <el-tag :type="getTypeColor(currentRecord.careType)" size="small">
            {{ currentRecord.careType }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="记录人">
          {{ currentRecord.nurseName }}
        </el-descriptions-item>
        <el-descriptions-item label="护理内容" :span="2">
          {{ currentRecord.careContent }}
        </el-descriptions-item>
        <el-descriptions-item label="护理结果" :span="2">
          {{ currentRecord.result || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="服务评价" :span="2">
          <el-rate v-model="currentRecord.evaluation" disabled />
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ currentRecord.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <div v-if="currentRecord?.images?.length" class="record-images">
        <div class="section-title">相关照片</div>
        <div class="image-grid">
          <el-image
            v-for="(img, index) in currentRecord.images"
            :key="index"
            :src="img"
            :preview-src-list="currentRecord.images"
            fit="cover"
            class="record-image"
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Plus } from '@element-plus/icons-vue'
import ExportButton from '@/components/Export/ExportButton.vue'
import { exportCareRecords } from '@/utils/export'
import { getCareRecords, getElderlyList } from '@/utils/seedData'

const loading = ref(false)
const recordList = ref([])
const total = ref(0)
const detailVisible = ref(false)
const currentRecord = ref<any>(null)

const filterForm = reactive({
  elderlyId: '',
  careType: '',
  nurseId: '',
  dateRange: []
})

const pagination = reactive({
  page: 1,
  pageSize: 20
})

const elderlyList = ref(getElderlyList().map(e => ({
  id: e.id,
  name: e.name,
  bedNumber: e.bedNumber
})))

const nurses = ref([
  { id: 1, name: '赵护士' },
  { id: 2, name: '李护士' },
  { id: 3, name: '周护士' },
  { id: 4, name: '吴护士' },
  { id: 5, name: '郑护士' }
])

const exportColumns = [
  { key: 'recordTime', title: '记录时间', width: 18 },
  { key: 'elderlyName', title: '老人姓名', width: 12 },
  { key: 'careType', title: '护理类型', width: 15 },
  { key: 'careContent', title: '护理内容', width: 30 },
  { key: 'nurseName', title: '护理员', width: 12 },
  { key: 'result', title: '护理结果', width: 20 },
  { key: 'evaluation', title: '评价', width: 15 }
]

const getTypeColor = (type: string) => {
  const map: Record<string, any> = {
    '日常护理': 'info',
    '康复训练': 'success',
    '健康监测': 'warning',
    '医疗护理': 'danger'
  }
  return map[type] || ''
}

const loadData = async () => {
  loading.value = true
  try {
    // 使用种子数据
    await new Promise(resolve => setTimeout(resolve, 300))
    let data = getCareRecords()

    // 应用筛选
    if (filterForm.elderlyId) {
      data = data.filter(r => r.elderlyId === filterForm.elderlyId)
    }
    if (filterForm.careType) {
      data = data.filter(r => r.careType === filterForm.careType)
    }

    // 分页
    const start = (pagination.page - 1) * pagination.pageSize
    const end = start + pagination.pageSize
    recordList.value = data.slice(start, end).map(r => ({
      ...r,
      recordTime: r.careTime,
      careContent: r.content,
      result: r.remarks || r.content
    }))
    total.value = data.length
  } finally {
    loading.value = false
  }
}

const handleFilter = () => {
  pagination.page = 1
  loadData()
}

const handleReset = () => {
  Object.assign(filterForm, {
    elderlyId: '',
    careType: '',
    nurseId: '',
    dateRange: []
  })
  handleFilter()
}

const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  loadData()
}

const handleAdd = () => {
  console.log('新增记录')
}

const handleView = (row: any) => {
  currentRecord.value = row
  detailVisible.value = true
}

const handleEdit = (row: any) => {
  console.log('编辑记录', row)
}

const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除这条记录吗？', '提示', {
      type: 'warning'
    })
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // 取消
  }
}

const handleExport = async (params: any) => {
  if (params.type === 'all') {
    // 获取全部数据
    exportCareRecords([])
  } else {
    exportCareRecords(params.data)
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.care-records {
  padding: 20px;

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .page-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 24px;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .filter-card {
    margin-bottom: 20px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .record-images {
    margin-top: 24px;
  }

  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 12px;
  }

  .image-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 12px;
  }

  .record-image {
    width: 100%;
    height: 120px;
    border-radius: 4px;
  }
}
</style>
