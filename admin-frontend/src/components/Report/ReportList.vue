<template>
  <div class="report-list">
    <el-empty v-if="reports.length === 0" description="暂无报表" />

    <div v-else class="report-grid">
      <el-card
        v-for="report in reports"
        :key="report.id"
        shadow="hover"
        class="report-card"
        :class="{ 'is-scheduled': report.schedule }"
      >
        <template #header>
          <div class="card-header">
            <div class="header-left">
              <h3 class="report-name">{{ report.name }}</h3>
              <el-tag v-if="report.schedule" size="small" type="info">
                <el-icon><Clock /></el-icon>
                {{ report.schedule }}
              </el-tag>
            </div>
            <div class="header-actions">
              <el-dropdown trigger="click">
                <el-button size="small" text>
                  <el-icon><MoreFilled /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="$emit('view', report)">
                      <el-icon><View /></el-icon>
                      查看
                    </el-dropdown-item>
                    <el-dropdown-item @click="$emit('export', report)">
                      <el-icon><Download /></el-icon>
                      导出
                    </el-dropdown-item>
                    <el-dropdown-item divided @click="handleSchedule(report)">
                      <el-icon><Clock /></el-icon>
                      定时任务
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
        </template>

        <div class="card-body">
          <p class="report-description">{{ report.description }}</p>

          <div class="report-meta">
            <div class="meta-item">
              <el-icon><Calendar /></el-icon>
              <span>更新于 {{ report.updateAt }}</span>
            </div>
            <div class="meta-item" v-if="report.category">
              <el-icon><Folder /></el-icon>
              <span>{{ getCategoryText(report.category) }}</span>
            </div>
          </div>
        </div>

        <template #footer>
          <div class="card-footer">
            <el-button size="small" @click="$emit('view', report)">
              <el-icon><View /></el-icon>
              查看报表
            </el-button>
            <el-button size="small" type="primary" @click="$emit('export', report)">
              <el-icon><Download /></el-icon>
              导出数据
            </el-button>
          </div>
        </template>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Clock, MoreFilled, View, Download, Calendar, Folder } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

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

defineEmits<{
  view: [report: Report]
  export: [report: Report]
}>()

const getCategoryText = (category: string) => {
  const map: Record<string, string> = {
    operation: '运营报表',
    nursing: '护理报表',
    finance: '财务报表',
    health: '健康报表'
  }
  return map[category] || category
}

const handleSchedule = (report: Report) => {
  ElMessage.info(`设置"${report.name}"的定时任务`)
}
</script>

<style scoped lang="scss">
.report-list {
  .report-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 20px;
  }

  .report-card {
    transition: var(--transition-base);

    &.is-scheduled {
      border-top: 3px solid var(--primary-color);
    }

    &:hover {
      transform: translateY(-4px);
      box-shadow: var(--card-shadow-hover);
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;

    .header-left {
      flex: 1;

      .report-name {
        font-size: 16px;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 8px 0;
      }

      :deep(.el-tag) {
        display: inline-flex;
        align-items: center;
        gap: 4px;
      }
    }

    .header-actions {
      margin-left: 12px;
    }
  }

  .card-body {
    .report-description {
      font-size: 14px;
      color: var(--text-secondary);
      line-height: 1.6;
      margin: 0 0 16px 0;
      min-height: 48px;
    }

    .report-meta {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;

      .meta-item {
        display: flex;
        align-items: center;
        gap: 4px;
        font-size: 12px;
        color: var(--text-tertiary);

        .el-icon {
          font-size: 14px;
        }
      }
    }
  }

  .card-footer {
    display: flex;
    gap: 8px;

    .el-button {
      flex: 1;
    }
  }
}

@media (max-width: 768px) {
  .report-list {
    .report-grid {
      grid-template-columns: 1fr;
    }

    .card-footer {
      flex-direction: column;

      .el-button {
        width: 100%;
      }
    }
  }
}
</style>
