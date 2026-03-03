<template>
  <div class="care-assessment">
    <div class="page-header">
      <h2 class="page-title">
        <el-icon><Document /></el-icon>
        护理评估
      </h2>
      <el-button type="primary" @click="handleNewAssessment">
        <el-icon><Plus /></el-icon>
        新建评估
      </el-button>
    </div>

    <!-- 评估列表 -->
    <el-card shadow="never">
      <el-table :data="assessmentList" stripe v-loading="loading">
        <el-table-column prop="assessmentNo" label="评估编号" width="150" />
        <el-table-column prop="elderlyName" label="老人姓名" width="100" />
        <el-table-column prop="assessmentType" label="评估类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.assessmentType }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="assessmentDate" label="评估日期" width="120" />
        <el-table-column prop="nurseName" label="评估人" width="100" />
        <el-table-column prop="overallScore" label="综合评分" width="100">
          <template #default="{ row }">
            <el-tag :type="getScoreType(row.overallScore)" size="small">
              {{ row.overallScore }} 分
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="careLevel" label="建议护理等级" width="120" />
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 'completed' ? 'success' : 'warning'" size="small">
              {{ row.status === 'completed' ? '已完成' : '待审核' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button size="small" text @click="handleView(row)">查看</el-button>
            <el-button size="small" text @click="handlePrint(row)">打印</el-button>
            <el-button size="small" text type="danger" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Plus } from '@element-plus/icons-vue'

const loading = ref(false)
const assessmentList = ref([
  {
    id: 1,
    assessmentNo: 'AS202603001',
    elderlyName: '张奶奶',
    assessmentType: '入院评估',
    assessmentDate: '2026-03-01',
    nurseName: '赵护士',
    overallScore: 85,
    careLevel: '二级护理',
    status: 'completed'
  },
  {
    id: 2,
    assessmentNo: 'AS202603002',
    elderlyName: '王爷爷',
    assessmentType: '定期评估',
    assessmentDate: '2026-02-28',
    nurseName: '钱护士',
    overallScore: 72,
    careLevel: '一级护理',
    status: 'completed'
  }
])

const getScoreType = (score: number) => {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

const handleNewAssessment = () => {
  console.log('新建评估')
}

const handleView = (row: any) => {
  console.log('查看评估', row)
}

const handlePrint = (row: any) => {
  console.log('打印评估', row)
}

const handleDelete = () => {
  ElMessage.success('删除成功')
}

onMounted(() => {
  // 加载数据
})
</script>

<style scoped lang="scss">
.care-assessment {
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
}
</style>
