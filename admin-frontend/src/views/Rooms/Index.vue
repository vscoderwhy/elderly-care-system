<template>
  <div class="rooms-page">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-value">{{ stats.total }}</div>
          <div class="stat-label">总床位</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card occupied">
          <div class="stat-value">{{ stats.occupied }}</div>
          <div class="stat-label">已入住</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card empty">
          <div class="stat-value">{{ stats.empty }}</div>
          <div class="stat-label">空床位</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card rate">
          <div class="stat-value">{{ occupancyRate }}%</div>
          <div class="stat-label">入住率</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 楼栋选择 -->
    <el-card class="building-selector">
      <el-radio-group v-model="selectedBuilding" @change="loadBuildingData">
        <el-radio-button v-for="b in buildings" :key="b.id" :label="b.id">
          {{ b.name }}
        </el-radio-button>
      </el-radio-group>
    </el-card>

    <!-- 房间可视化 -->
    <el-card v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>房间布局</span>
          <el-select v-model="selectedFloor" placeholder="选择楼层" @change="filterFloor" style="width: 150px;">
            <el-option label="全部楼层" :value="0" />
            <el-option v-for="f in floors" :key="f.id" :label="f.name" :value="f.id" />
          </el-select>
        </div>
      </template>

      <!-- 楼层视图 -->
      <div class="floors-container">
        <div v-for="floor in displayFloors" :key="floor.id" class="floor-section">
          <div class="floor-title">{{ floor.name }}</div>
          <div class="rooms-grid">
            <div v-for="room in floor.rooms" :key="room.id" class="room-card">
              <div class="room-header">
                <span class="room-name">{{ room.name }}</span>
                <span class="room-capacity">{{ getOccupiedCount(room) }}/{{ room.bed_capacity }}床</span>
              </div>
              <div class="beds-grid">
                <div
                  v-for="bed in room.beds"
                  :key="bed.id"
                  class="bed-item"
                  :class="bed.status"
                  @click="showBedDetail(bed)"
                >
                  <div class="bed-name">{{ bed.name }}</div>
                  <div class="bed-elderly" v-if="bed.elderly">
                    {{ bed.elderly.name }}
                  </div>
                  <div class="bed-status">
                    {{ bed.status === 'occupied' ? '已住' : '空' }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 床位详情对话框 -->
    <el-dialog v-model="showBedDialog" title="床位详情" width="400px">
      <div class="bed-detail" v-if="selectedBed">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="床位">{{ selectedBed.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="selectedBed.status === 'occupied' ? 'success' : 'info'">
              {{ selectedBed.status === 'occupied' ? '已入住' : '空床位' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="入住老人" v-if="selectedBed.elderly">
            {{ selectedBed.elderly.name }}
          </el-descriptions-item>
          <el-descriptions-item label="护理等级" v-if="selectedBed.elderly">
            {{ selectedBed.elderly.care_level }}级护理
          </el-descriptions-item>
        </el-descriptions>
      </div>
      <template #footer>
        <el-button @click="showBedDialog = false">关闭</el-button>
        <el-button type="danger" v-if="selectedBed?.status === 'occupied'" @click="releaseBed">
          释放床位
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import instance from '@/api'

interface Bed {
  id: number
  name: string
  status: string
  elderly?: {
    id: number
    name: string
    care_level: number
  }
}

interface Room {
  id: number
  name: string
  bed_capacity: number
  beds: Bed[]
}

interface Floor {
  id: number
  name: string
  rooms: Room[]
}

interface Building {
  id: number
  name: string
  floors: Floor[]
}

const loading = ref(false)
const buildings = ref<Building[]>([])
const selectedBuilding = ref<number>(0)
const selectedFloor = ref<number>(0)
const floors = ref<Floor[]>([])
const buildingData = ref<Building | null>(null)
const showBedDialog = ref(false)
const selectedBed = ref<Bed | null>(null)

const stats = ref({
  total: 0,
  occupied: 0,
  empty: 0
})

const occupancyRate = computed(() => {
  if (stats.value.total === 0) return 0
  return Math.round((stats.value.occupied / stats.value.total) * 100)
})

const displayFloors = computed(() => {
  if (!buildingData.value) return []
  if (selectedFloor.value === 0) return buildingData.value.floors
  return buildingData.value.floors.filter(f => f.id === selectedFloor.value)
})

const getOccupiedCount = (room: Room) => {
  return room.beds.filter(b => b.status === 'occupied').length
}

const loadBuildings = async () => {
  try {
    const result = await instance.get('/rooms/buildings')
    buildings.value = result || []
    if (buildings.value.length > 0) {
      selectedBuilding.value = buildings.value[0].id
      loadBuildingData()
    }
  } catch (error) {
    console.error(error)
  }
}

const loadBuildingData = async () => {
  if (!selectedBuilding.value) return
  loading.value = true
  try {
    const result = await instance.get(`/rooms/buildings/${selectedBuilding.value}`)
    buildingData.value = result
    floors.value = result?.floors || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const result = await instance.get('/rooms/stats')
    stats.value = result || { total: 0, occupied: 0, empty: 0 }
  } catch (error) {
    console.error(error)
  }
}

const filterFloor = () => {
  // Floor filtering is handled by computed property
}

const showBedDetail = (bed: Bed) => {
  selectedBed.value = bed
  showBedDialog.value = true
}

const releaseBed = async () => {
  if (!selectedBed.value) return
  try {
    await instance.post(`/rooms/beds/${selectedBed.value.id}/release`)
    ElMessage.success('床位已释放')
    showBedDialog.value = false
    loadBuildingData()
    loadStats()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

onMounted(() => {
  loadBuildings()
  loadStats()
})
</script>

<style scoped>
.rooms-page {
  padding: 20px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
  padding: 20px;
}

.stat-value {
  font-size: 36px;
  font-weight: bold;
  color: #409EFF;
}

.stat-card.occupied .stat-value {
  color: #67C23A;
}

.stat-card.empty .stat-value {
  color: #909399;
}

.stat-card.rate .stat-value {
  color: #E6A23C;
}

.stat-label {
  color: #666;
  margin-top: 8px;
}

.building-selector {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.floors-container {
  padding: 10px;
}

.floor-section {
  margin-bottom: 30px;
}

.floor-title {
  font-size: 18px;
  font-weight: bold;
  color: #333;
  margin-bottom: 15px;
  padding-bottom: 10px;
  border-bottom: 2px solid #409EFF;
}

.rooms-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}

.room-card {
  width: 280px;
  background: #f8f9fa;
  border-radius: 8px;
  padding: 15px;
  border: 1px solid #e4e7ed;
}

.room-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #dcdfe6;
}

.room-name {
  font-weight: bold;
  color: #303133;
}

.room-capacity {
  font-size: 12px;
  color: #909399;
}

.beds-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.bed-item {
  width: calc(50% - 4px);
  padding: 10px;
  border-radius: 6px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.bed-item.empty {
  background: #f0f9eb;
  border: 1px solid #c2e7b0;
}

.bed-item.occupied {
  background: #ecf5ff;
  border: 1px solid #b3d8ff;
}

.bed-item:hover {
  transform: scale(1.02);
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.bed-name {
  font-size: 12px;
  color: #606266;
  margin-bottom: 4px;
}

.bed-elderly {
  font-size: 14px;
  font-weight: 500;
  color: #303133;
  margin-bottom: 2px;
}

.bed-status {
  font-size: 11px;
  color: #909399;
}

.bed-detail {
  padding: 10px 0;
}
</style>
