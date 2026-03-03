<template>
  <div class="equipment-management">
    <div class="page-header">
      <div class="header-left">
        <h2 class="page-title">
          <el-icon><Connection /></el-icon>
          设备管理
        </h2>
      </div>
      <div class="header-actions">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          添加设备
        </el-button>
        <el-button @click="handleImport">
          <el-icon><Upload /></el-icon>
          批量导入
        </el-button>
      </div>
    </div>

    <!-- 设备统计 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="12" :sm="6" v-for="stat in deviceStats" :key="stat.key">
        <div class="device-stat" :class="`stat-${stat.status}`">
          <div class="stat-icon">
            <component :is="stat.icon" />
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </el-col>
    </el-row>

    <!-- 设备分类 -->
    <el-card shadow="never" class="category-card">
      <template #header>
        <div class="card-header">
          <span>设备分类</span>
        </div>
      </template>

      <div class="category-grid">
        <div
          v-for="category in categories"
          :key="category.key"
          class="category-item"
          @click="handleCategoryClick(category)"
        >
          <div class="category-icon" :style="{ background: category.color }">
            <component :is="category.icon" />
          </div>
          <div class="category-info">
            <div class="category-name">{{ category.name }}</div>
            <div class="category-count">{{ category.count }}台</div>
          </div>
          <div class="category-status">
            <span class="status-dot" :class="`dot-${category.status}`"></span>
            <span class="status-text">{{ category.statusText }}</span>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 设备列表 -->
    <el-card shadow="never" class="table-card">
      <template #header>
        <div class="card-header">
          <span>设备列表</span>
          <div class="header-actions">
            <el-select v-model="filterForm.status" placeholder="全部状态" clearable style="width: 120px" @change="loadDevices">
              <el-option label="全部" value="" />
              <el-option label="在线" value="online" />
              <el-option label="离线" value="offline" />
              <el-option label="维修中" value="maintenance" />
              <el-option label="已报废" value="scrapped" />
            </el-select>
            <el-select v-model="filterForm.type" placeholder="设备类型" clearable style="width: 150px" @change="loadDevices">
              <el-option label="全部类型" value="" />
              <el-option label="医疗设备" value="medical" />
              <el-option label="护理设备" value="care" />
              <el-option label="康复设备" value="rehab" />
              <el-option label="智能设备" value="smart" />
              <el-option label="其他" value="other" />
            </el-select>
          </div>
        </div>
      </template>

      <el-table :data="equipmentList" stripe v-loading="loading">
        <el-table-column prop="deviceNo" label="设备编号" width="150" />
        <el-table-column prop="name" label="设备名称" width="150" />
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="location" label="位置" width="150" />
        <el-table-column prop="purchaseDate" label="购入日期" width="120" />
        <el-table-column prop="lastMaintenance" label="上次维护" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handleEdit(row)">编辑</el-button>
            <el-dropdown trigger="click">
              <el-button size="small" text>
                更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="handleMaintenance(row)">
                    <el-icon><Tools /></el-icon>
                    维护记录
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleQRCode(row)">
                    <el-icon><QrCode /></el-icon>
                    二维码
                  </el-dropdown-item>
                  <el-dropdown-item @click="handleCalibrate(row)">
                    <el-icon><Setting /></el-icon>
                    校准
                  </el-dropdown-item>
                  <el-dropdown-item divided @click="handleScrap(row)" v-if="row.status !== 'scrapped'">
                    <el-icon><Delete /></el-icon>
                    报废
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Connection,
  Plus,
  Upload,
  ArrowDown,
  Tools,
  QrCode,
  Setting,
  Delete
} from '@element-plus/icons-vue'

// 设备统计
const deviceStats = ref([
  { key: 'total', label: '设备总数', value: 156, status: 'all', icon: 'Box' },
  { key: 'online', label: '在线', value: 142, status: 'online', icon: 'Connection' },
  { key: 'offline', label: '离线', value: 8, status: 'offline', icon: 'Connection' },
  { key: 'maintenance', label: '维修中', value: 6, status: 'warning', icon: 'Tools' }
])

// 设备分类
const categories = ref([
  {
    key: 'medical',
    name: '医疗设备',
    count: 42,
    status: 'online',
    statusText: '运行正常',
    icon: 'FirstAidKit',
    color: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
  },
  {
    key: 'care',
    name: '护理设备',
    count: 35,
    status: 'online',
    statusText: '运行正常',
    icon: 'Monitor',
    color: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)'
  },
  {
    key: 'rehab',
    name: '康复设备',
    count: 28,
    status: 'online',
    statusText: '运行正常',
    icon: 'Trophy',
    color: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)'
  },
  {
    key: 'smart',
    name: '智能设备',
    count: 51,
    status: 'offline',
    statusText: '部分离线',
    icon: 'MagicStick',
    color: 'linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)'
  }
])

const filterForm = reactive({
  status: '',
  type: ''
})

const loading = ref(false)
const equipmentList = ref([
  {
    id: 1,
    deviceNo: 'EQ001',
    name: '心电监护仪',
    category: '医疗设备',
    location: '医务室',
    purchaseDate: '2025-01-15',
    lastMaintenance: '2026-02-15',
    status: 'online'
  },
  {
    id: 2,
    deviceNo: 'EQ002',
    name: '智能床垫',
    category: '智能设备',
    location: '3号楼201',
    purchaseDate: '2025-03-20',
    lastMaintenance: '2026-01-20',
    status: 'offline'
  },
  {
    id: 3,
    deviceNo: 'EQ003',
    name: '康复训练仪',
    category: '康复设备',
    location: '康复室',
    purchaseDate: '2025-02-10',
    lastMaintenance: '2026-02-28',
    status: 'maintenance'
  }
])

const pagination = reactive({
  page: 1,
  pageSize: 20
})
const total = ref(156)

const getStatusType = (status: string) => {
  const map: Record<string, any> = {
    online: 'success',
    offline: 'info',
    maintenance: 'warning',
    scrapped: 'danger'
  }
  return map[status] || ''
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    online: '在线',
    offline: '离线',
    maintenance: '维修中',
    scrapped: '已报废'
  }
  return map[status] || status
}

const handleAdd = () => {
  console.log('添加设备')
}

const handleImport = () => {
  console.log('批量导入')
}

const handleCategoryClick = (category: any) => {
  filterForm.type = category.key
  loadDevices()
}

const handleView = (row: any) => {
  console.log('查看设备', row)
}

const handleEdit = (row: any) => {
  console.log('编辑设备', row)
}

const handleMaintenance = (row: any) => {
  console.log('维护记录', row)
}

const handleQRCode = (row: any) => {
  console.log('二维码', row)
}

const handleCalibrate = (row: any) => {
  console.log('校准设备', row)
}

const handleScrap = async (row: any) => {
  try {
    await ElMessageBox.confirm(
      `确定要将设备"${row.name}"报废吗？此操作不可恢复。`,
      '报废确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    ElMessage.success('已标记为报废')
  } catch {
    // 取消
  }
}

const loadDevices = () => {
  loading.value = true
  // TODO: 加载设备列表
  setTimeout(() => {
    loading.value = false
  }, 500)
}

onMounted(() => {
  loadDevices()
})
</script>

<style scoped lang="scss">
.equipment-management {
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
    margin: 0;
  }

  .stats-row {
    margin-bottom: 20px;

    :deep(.el-col) {
      margin-bottom: 12px;
    }
  }

  .device-stat {
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

      .stat-all & { background: #f0f2f5; color: #909399; }
      .stat-online & { background: #e1f3e8; color: #67c23a; }
      .stat-offline & { background: #fef0f0; color: #909399; }
      .stat-warning & { background: #fdf6ec; color: #e6a23c; }
    }

    .stat-info {
      flex: 1;
    }

    .stat-value {
      font-size: 28px;
      font-weight: 600;
      color: var(--text-primary);
      margin-bottom: 4px;
    }

    .stat-label {
      font-size: 12px;
      color: var(--text-secondary);
    }
  }

  .category-card,
  .table-card {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .category-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
    gap: 16px;
  }

  .category-item {
    position: relative;
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    background: var(--card-bg);
    border: 1px solid var(--border-color-lighter);
    border-radius: var(--card-border-radius);
    cursor: pointer;
    transition: var(--transition-base);

    &:hover {
      border-color: var(--primary-color);
      box-shadow: var(--card-shadow-hover);
      transform: translateY(-2px);
    }

    .category-icon {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      font-size: 20px;
    }

    .category-info {
      flex: 1;
    }

    .category-name {
      font-size: 14px;
      font-weight: 500;
      color: var(--text-primary);
      margin-bottom: 4px;
    }

    .category-count {
      font-size: 12px;
      color: var(--text-secondary);
    }

    .category-status {
      display: flex;
      align-items: center;
      gap: 6px;
      font-size: 12px;
      color: var(--text-secondary);
    }

    .status-dot {
      width: 8px;
      height: 8px;
      border-radius: 50%;

      &.dot-online { background: var(--success-color); }
      &.dot-offline { background: var(--info-color); }
      &.dot-warning { background: var(--warning-color); }
    }
  }

  .header-actions {
    display: flex;
    gap: 12px;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }
}

@media (max-width: 768px) {
  .equipment-management {
    padding: 10px;

    .page-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 12px;
    }

    .category-grid {
      grid-template-columns: 1fr;
    }

    .header-actions {
      flex-wrap: wrap;
    }
  }
}
</style>
