<template>
  <div class="care-tasks">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><ListDone /></el-icon>
          护理任务
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          新建任务
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in taskStats" :key="stat.key">
        <div class="task-stat" :class="`stat-${stat.type}`">
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 筛选和搜索 -->
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="任务状态">
          <el-radio-group v-model="filterForm.status">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="pending">待处理</el-radio-button>
            <el-radio-button label="in-progress">进行中</el-radio-button>
            <el-radio-button label="completed">已完成</el-radio-button>
            <el-radio-button label="overdue">已逾期</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select v-model="filterForm.taskType" placeholder="请选择" clearable>
            <el-option label="全部" value="" />
            <el-option label="日常护理" value="daily" />
            <el-option label="康复训练" value="rehab" />
            <el-option label="健康监测" value="health" />
            <el-option label="医疗护理" value="medical" />
            <el-option label="用药管理" value="medication" />
            <el-option label="营养配餐" value="meal" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行人">
          <el-select v-model="filterForm.nurseId" placeholder="请选择" clearable filterable>
            <el-option label="全部" value="" />
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
          <el-button type="primary" @click="handleFilter">筛选</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 任务列表 -->
    <div class="task-list">
      <div
        v-for="task in taskList"
        :key="task.id"
        class="task-card"
        :class="`task-${task.status}`"
      >
        <div class="task-header">
          <div class="task-title-row">
            <el-tag :type="getTypeColor(task.taskType)" size="small">
              {{ task.taskType }}
            </el-tag>
            <h3 class="task-title">{{ task.title }}</h3>
            <el-tag :type="getStatusColor(task.status)" size="small">
              {{ getStatusText(task.status) }}
            </el-tag>
          </div>
          <div class="task-time">
            <el-icon><Clock /></el-icon>
            {{ task.planTime }}
          </div>
        </div>

        <div class="task-body">
          <div class="task-info">
            <div class="info-item">
              <span class="info-label">老人：</span>
              <span class="info-value">{{ task.elderlyName }}</span>
              <el-tag size="small" class="ml-2">{{ task.bedNumber }}</el-tag>
            </div>
            <div class="info-item">
              <span class="info-label">执行人：</span>
              <span class="info-value">{{ task.nurseName }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">内容：</span>
              <span class="info-value">{{ task.content }}</span>
            </div>
          </div>

          <div v-if="task.status !== 'completed'" class="task-progress">
            <el-progress
              :percentage="task.progress || 0"
              :status="task.status === 'overdue' ? 'exception' : undefined"
            />
          </div>
        </div>

        <div class="task-footer">
          <div class="task-meta">
            <span v-if="task.priority !== 'normal'" class="priority-tag" :class="`priority-${task.priority}`">
              {{ task.priority === 'urgent' ? '紧急' : '重要' }}
            </span>
            <span class="duration-text">
              <el-icon><Timer /></el-icon>
              预计 {{ task.duration }} 分钟
            </span>
          </div>
          <div class="task-actions">
            <el-button size="small" @click="handleView(task)">查看</el-button>
            <el-button
              v-if="task.status === 'pending'"
              size="small"
              type="primary"
              @click="handleStart(task)"
            >
              开始
            </el-button>
            <el-button
              v-if="task.status === 'in-progress'"
              size="small"
              type="success"
              @click="handleComplete(task)"
            >
              完成
            </el-button>
            <el-dropdown trigger="click">
              <el-button size="small">
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleEdit(task)">
                    <el-icon><Edit /></el-icon>
                    编辑
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleAssign(task)">
                    <el-icon><User /></el-icon>
                    分配
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleRecord(task)">
                    <el-icon><Document /></el-icon>
                    记录
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleDelete(task)">
                    <el-icon><Delete /></el-icon>
                    删除
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页 -->
    <div class="pagination-wrapper">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="currentTask?.title"
      width="700px"
    >
      <div v-if="currentTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="任务编号">
            {{ currentTask.taskNo }}
          </el-descriptions-item>
          <el-descriptions-item label="任务类型">
            <el-tag :type="getTypeColor(currentTask.taskType)" size="small">
              {{ currentTask.taskType }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="老人姓名">
            {{ currentTask.elderlyName }}
          </el-descriptions-item>
          <el-descriptions-item label="床位号">
            {{ currentTask.bedNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="执行人">
            {{ currentTask.nurseName }}
          </el-descriptions-item>
          <el-descriptions-item label="计划时间">
            {{ currentTask.planTime }}
          </el-descriptions-item>
          <el-descriptions-item label="预计时长" :span="2">
            {{ currentTask.duration }} 分钟
          </el-descriptions-item>
          <el-descriptions-item label="任务内容" :span="2">
            {{ currentTask.content }}
          </el-descriptions-item>
          <el-descriptions-item label="执行结果" :span="2">
            {{ currentTask.result || '未执行' }}
          </el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">
            {{ currentTask.remark || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentTask.images && currentTask.images.length > 0" class="task-images">
          <div class="section-title">执行照片</div>
          <div class="image-list">
            <el-image
              v-for="(img, index) in currentTask.images"
              :key="index"
              :src="img"
              :preview-src-list="currentTask.images"
              fit="cover"
              class="task-image"
            />
          </div>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
        <el-button v-if="currentTask?.status === 'pending'" type="primary" @click="handleStartFromDetail">
          开始任务
        </el-button>
        <el-button v-if="currentTask?.status === 'in-progress'" type="success" @click="handleCompleteFromDetail">
          完成任务
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  ListDone,
  Plus,
  Clock,
  Timer,
  ArrowDown,
  Edit,
  User,
  Document,
  Delete
} from '@element-plus/icons-vue'

// 统计数据
const taskStats = ref([
  { key: 'today', label: '今日任务', value: 156, type: 'primary', icon: 'List' },
  { key: 'pending', label: '待处理', value: 42, type: 'warning', icon: 'Clock' },
  { key: 'in-progress', label: '进行中', value: 18, type: 'primary', icon: 'Loading' },
  { key: 'completed', label: '已完成', value: 96, type: 'success', icon: 'Check' }
])

// 筛选表单
const filterForm = reactive({
  status: '',
  taskType: '',
  nurseId: '',
  dateRange: []
})

// 护理员列表
const nurses = ref([
  { id: 1, name: '赵护士' },
  { id: 2, name: '钱护士' },
  { id: 3, name: '孙护士' },
  { id: 4, name: '李护士' },
  { id: 5, name: '周护士' }
])

// 任务列表
const taskList = ref([
  {
    id: 1,
    taskNo: 'TASK20260303001',
    title: '测量血压血糖',
    taskType: '日常护理',
    elderlyName: '张奶奶',
    bedNumber: '3号楼201',
    nurseName: '赵护士',
    nurseId: 1,
    planTime: '2026-03-03 14:30',
    duration: 15,
    priority: 'normal',
    status: 'pending',
    progress: 0,
    content: '测量血压、血糖并记录，注意老人是否有不适症状',
    result: '',
    remark: ''
  },
  {
    id: 2,
    taskNo: 'TASK20260303002',
    title: '上肢康复训练',
    taskType: '康复训练',
    elderlyName: '王爷爷',
    bedNumber: '3号楼202',
    nurseName: '陈康复师',
    nurseId: 3,
    planTime: '2026-03-03 15:00',
    duration: 30,
    priority: 'normal',
    status: 'in-progress',
    progress: 60,
    content: '协助老人进行上肢关节活动训练，每组动作15次，共3组',
    result: '',
    remark: '老人配合度良好'
  },
  {
    id: 3,
    taskNo: 'TASK20260303003',
    title: '伤口换药',
    taskType: '医疗护理',
    elderlyName: '李奶奶',
    bedNumber: '2号楼105',
    nurseName: '周护士',
    nurseId: 5,
    planTime: '2026-03-03 16:00',
    duration: 20,
    priority: 'urgent',
    status: 'pending',
    progress: 0,
    content: '清洁伤口，更换敷料，观察伤口愈合情况',
    result: '',
    remark: '注意无菌操作'
  },
  {
    id: 4,
    taskNo: 'TASK20260303004',
    title: '协助用餐',
    taskType: '日常护理',
    elderlyName: '刘爷爷',
    bedNumber: '2号楼108',
    nurseName: '钱护士',
    nurseId: 2,
    planTime: '2026-03-03 11:30',
    duration: 30,
    priority: 'normal',
    status: 'completed',
    progress: 100,
    content: '协助老人用餐，注意饮食禁忌',
    result: '老人用餐正常，食欲良好',
    remark: ''
  },
  {
    id: 5,
    taskNo: 'TASK20260303005',
    title: '康复训练指导',
    taskType: '康复训练',
    elderlyName: '孙奶奶',
    bedNumber: '1号楼301',
    nurseName: '陈康复师',
    nurseId: 3,
    planTime: '2026-03-03 10:00',
    duration: 45,
    priority: 'important',
    status: 'overdue',
    progress: 0,
    content: '指导老人进行下肢康复训练，包括站立平衡练习',
    result: '',
    remark: '需要家属配合'
  }
])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(156)

// 对话框
const detailDialogVisible = ref(false)
const currentTask = ref<any>(null)

// 获取类型颜色
const getTypeColor = (type: string) => {
  const map: Record<string, any> = {
    '日常护理': 'info',
    '康复训练': 'success',
    '健康监测': 'warning',
    '医疗护理': 'danger',
    '用药管理': 'primary',
    '营养配餐': ''
  }
  return map[type] || ''
}

// 获取状态颜色
const getStatusColor = (status: string) => {
  const map: Record<string, any> = {
    pending: 'info',
    'in-progress': 'primary',
    completed: 'success',
    overdue: 'danger'
  }
  return map[status] || ''
}

// 获取状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    'in-progress': '进行中',
    completed: '已完成',
    overdue: '已逾期'
  }
  return map[status] || status
}

// 筛选
const handleFilter = () => {
  pagination.page = 1
  loadData()
}

// 重置
const handleReset = () => {
  Object.assign(filterForm, {
    status: '',
    taskType: '',
    nurseId: '',
    dateRange: []
  })
  handleFilter()
}

// 分页变化
const handlePageChange = (page: number) => {
  pagination.page = page
  loadData()
}

const handleSizeChange = (size: number) => {
  pagination.pageSize = size
  pagination.page = 1
  loadData()
}

// 加载数据
const loadData = async () => {
  // TODO: 实际 API 请求
  console.log('加载任务列表', filterForm, pagination)
}

// 新建任务
const handleAdd = () => {
  console.log('新建任务')
}

// 查看详情
const handleView = (task: any) => {
  currentTask.value = task
  detailDialogVisible.value = true
}

// 开始任务
const handleStart = async (task: any) => {
  try {
    await ElMessageBox.confirm(`确定要开始执行"${task.title}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })

    // TODO: API 请求
    task.status = 'in-progress'
    task.progress = 0
    ElMessage.success('任务已开始')
  } catch {
    // 用户取消
  }
}

// 完成任务
const handleComplete = async (task: any) => {
  try {
    await ElMessageBox.prompt('请输入执行结果', '完成任务', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /.+/,
      inputErrorMessage: '请输入执行结果'
    })

    // TODO: API 请求
    task.status = 'completed'
    task.progress = 100
    ElMessage.success('任务已完成')
  } catch {
    // 用户取消
  }
}

const handleStartFromDetail = () => {
  if (currentTask.value) {
    handleStart(currentTask.value)
    detailDialogVisible.value = false
  }
}

const handleCompleteFromDetail = () => {
  if (currentTask.value) {
    handleComplete(currentTask.value)
  }
}

// 编辑
const handleEdit = (task: any) => {
  console.log('编辑任务', task)
}

// 分配
const handleAssign = (task: any) => {
  console.log('分配任务', task)
}

// 记录
const handleRecord = (task: any) => {
  console.log('记录任务', task)
}

// 删除
const handleDelete = async (task: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除任务"${task.title}"吗？`,
      '删除确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    // TODO: API 请求
    ElMessage.success('删除成功')
    loadData()
  } catch {
    // 用户取消
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.care-tasks {
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

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .task-stat {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);

    .stat-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 24px;
    }

    &.stat-primary .stat-icon {
      background: var(--gradient-blue);
      color: #fff;
    }

    &.stat-success .stat-icon {
      background: var(--gradient-green);
      color: #fff;
    }

    &.stat-warning .stat-icon {
      background: var(--gradient-orange);
      color: #fff;
    }

    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: var(--text-primary);
    }

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
    }
  }

  .filter-card {
    margin-bottom: 20px;
  }

  .task-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
    margin-bottom: 20px;
  }

  .task-card {
    background: var(--card-bg);
    border-radius: var(--card-border-radius);
    box-shadow: var(--card-shadow);
    overflow: hidden;
    border-left: 4px solid transparent;

    &.task-pending {
      border-left-color: var(--info-color);
    }

    &.task-in-progress {
      border-left-color: var(--primary-color);
    }

    &.task-completed {
      border-left-color: var(--success-color);
    }

    &.task-overdue {
      border-left-color: var(--danger-color);
    }
  }

  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px;
    border-bottom: 1px solid var(--border-color-lighter);
  }

  .task-title-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .task-title {
    font-size: 16px;
    font-weight: 600;
    margin: 0;
  }

  .task-time {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 14px;
    color: var(--text-secondary);
  }

  .task-body {
    padding: 16px;
  }

  .task-info {
    margin-bottom: 12px;
  }

  .info-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    font-size: 14px;

    &:last-child {
      margin-bottom: 0;
    }
  }

  .info-label {
    color: var(--text-secondary);
    margin-right: 8px;
  }

  .info-value {
    color: var(--text-primary);
  }

  .task-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 16px;
    background: var(--bg-tertiary);
  }

  .task-meta {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .priority-tag {
    padding: 2px 8px;
    border-radius: 4px;
    font-size: 12px;

    &.priority-urgent {
      background: var(--danger-color);
      color: #fff;
    }

    &.priority-important {
      background: var(--warning-color);
      color: #fff;
    }
  }

  .duration-text {
    display: flex;
    align-items: center;
    gap: 4px;
    font-size: 12px;
    color: var(--text-secondary);
  }

  .task-actions {
    display: flex;
    gap: 8px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
  }

  .task-detail {
    .task-images {
      margin-top: 24px;
    }

    .section-title {
      font-size: 14px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 12px;
    }

    .image-list {
      display: flex;
      gap: 8px;
    }

    .task-image {
      width: 100px;
      height: 100px;
      border-radius: 4px;
    }
  }
}

@media (max-width: 768px) {
  .care-tasks {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .task-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .task-footer {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }
  }
}
</style>
