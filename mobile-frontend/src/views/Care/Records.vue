<template>
  <div class="care-records-page">
    <van-nav-bar title="护理记录" />

    <div class="elderly-selector" v-if="elderlyList.length > 0">
      <van-dropdown-menu>
        <van-dropdown-item v-model="selectedElderlyId" :options="elderlyOptions" @change="onElderlyChange" />
      </van-dropdown-menu>
    </div>

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list v-model:loading="loading" :finished="finished" @load="onLoad">
        <div class="care-list">
          <div class="care-item" v-for="item in list" :key="item.id">
            <div class="care-header">
              <van-tag type="primary">{{ item.care_item?.name }}</van-tag>
              <span class="time">{{ formatTime(item.recorded_at) }}</span>
            </div>
            <div class="care-body">
              <div class="staff">
                <van-icon name="user-o" />
                {{ item.staff?.nickname }}
              </div>
              <div class="notes" v-if="item.notes">{{ item.notes }}</div>
              <div class="images" v-if="item.images && item.images.length > 0">
                <van-image
                  v-for="(img, index) in item.images.slice(0, 3)"
                  :key="index"
                  width="80"
                  height="80"
                  :src="img"
                  fit="cover"
                  @click="previewImage(item.images, index)"
                />
              </div>
            </div>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <van-empty v-if="list.length === 0 && !loading" description="暂无护理记录" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { showToast, showImagePreview } from 'vant'
import api from '@/api'

const elderlyList = ref<any[]>([])
const selectedElderlyId = ref(0)
const list = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)

const elderlyOptions = computed(() => {
  return elderlyList.value.map(e => ({
    text: e.name,
    value: e.id
  }))
})

onMounted(() => {
  loadElderlyList()
})

function loadElderlyList() {
  api.elderly.familyList().then((data: any) => {
    elderlyList.value = data || []
    if (data && data.length > 0) {
      selectedElderlyId.value = data[0].id
      loadRecords()
    }
  })
}

function onElderlyChange() {
  list.value = []
  finished.value = false
  loadRecords()
}

function loadRecords() {
  if (!selectedElderlyId.value) return

  loading.value = true
  api.care.records({
    elderly_id: selectedElderlyId.value,
    page: 1,
    page_size: 20
  }).then((data: any) => {
    list.value = data.list || []
    loading.value = false
    finished.value = true
  }).catch(() => {
    loading.value = false
  })
}

function onLoad() {
  loadRecords()
}

function onRefresh() {
  finished.value = false
  loadRecords()
  refreshing.value = false
}

function formatTime(time: string) {
  const date = new Date(time)
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`
}

function previewImage(images: string[], index: number) {
  showImagePreview({
    images,
    startPosition: index
  })
}
</script>

<style scoped>
.care-records-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.elderly-selector {
  position: sticky;
  top: 46px;
  z-index: 10;
}

.care-list {
  padding: 15px;
}

.care-item {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 12px;
}

.care-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.time {
  font-size: 12px;
  color: #999;
}

.care-body {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.care-body .staff {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #666;
}

.care-body .notes {
  font-size: 13px;
  color: #999;
  line-height: 1.5;
}

.care-body .images {
  display: flex;
  gap: 8px;
}

.care-body .images :deep(.van-image) {
  border-radius: 4px;
}
</style>
