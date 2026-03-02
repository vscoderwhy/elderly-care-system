<template>
  <div class="elderly-list-page">
    <van-nav-bar title="我的家人" />

    <van-pull-refresh v-model="refreshing" @refresh="onRefresh">
      <van-list v-model:loading="loading" :finished="finished" @load="onLoad">
        <div class="elderly-cards">
          <div
            class="elderly-card"
            v-for="item in list"
            :key="item.id"
            @click="goToDetail(item.id)"
          >
            <div class="card-header">
              <van-image
                round
                width="50"
                height="50"
                :src="item.avatar || 'https://fastly.jsdelivr.net/npm/@vant/assets/user-active.png'"
              />
              <div class="info">
                <div class="name">{{ item.name }}</div>
                <div class="tags">
                  <van-tag size="small">{{ item.gender }}</van-tag>
                  <van-tag size="small" type="primary">{{ item.care_level }}级护理</van-tag>
                </div>
              </div>
            </div>
            <div class="card-body">
              <div class="info-row">
                <van-icon name="location-o" />
                <span>{{ item.bed?.name || '未分配床位' }}</span>
              </div>
              <div class="info-row" v-if="item.phone">
                <van-icon name="phone-o" />
                <span>{{ item.phone }}</span>
              </div>
            </div>
          </div>
        </div>
      </van-list>
    </van-pull-refresh>

    <van-empty v-if="list.length === 0 && !loading" description="暂无关联老人" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { showToast } from 'vant'
import api from '@/api'

const router = useRouter()
const list = ref<any[]>([])
const loading = ref(false)
const finished = ref(false)
const refreshing = ref(false)

function onLoad() {
  api.elderly.familyList().then((data: any) => {
    list.value = data || []
    loading.value = false
    finished.value = true
  }).catch(() => {
    loading.value = false
  })
}

function onRefresh() {
  finished.value = false
  loading.value = true
  onLoad()
  refreshing.value = false
}

function goToDetail(id: number) {
  router.push(`/elderly/${id}`)
}
</script>

<style scoped>
.elderly-list-page {
  min-height: 100vh;
  background-color: #f7f8fa;
}

.elderly-cards {
  padding: 15px;
}

.elderly-card {
  background: #fff;
  border-radius: 12px;
  padding: 15px;
  margin-bottom: 12px;
}

.card-header {
  display: flex;
  gap: 12px;
  margin-bottom: 12px;
}

.info {
  flex: 1;
}

.name {
  font-size: 18px;
  font-weight: 600;
  color: #333;
  margin-bottom: 8px;
}

.tags {
  display: flex;
  gap: 8px;
}

.card-body {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}
</style>
