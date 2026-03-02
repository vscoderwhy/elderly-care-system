<template>
  <div class="layout">
    <router-view v-slot="{ Component }">
      <component :is="Component" />
    </router-view>

    <van-tabbar v-model="active" route v-if="showTabbar">
      <van-tabbar-item to="/home" icon="wap-home-o">首页</van-tabbar-item>
      <van-tabbar-item to="/tasks" icon="todo-list-o" badge="">任务</van-tabbar-item>
      <van-tabbar-item to="/care" icon="records">护理</van-tabbar-item>
      <van-tabbar-item to="/bills" icon="balance-pay-o">账单</van-tabbar-item>
      <van-tabbar-item to="/profile" icon="user-o">我的</van-tabbar-item>
    </van-tabbar>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const active = ref(0)

const showTabbar = computed(() => {
  return !!route.meta?.showTabbar
})

watch(
  () => route.path,
  (path) => {
    const index = ['/home', '/tasks', '/care', '/bills', '/profile'].indexOf(path)
    if (index !== -1) {
      active.value = index
    }
  }
)
</script>

<style scoped>
.layout {
  min-height: 100vh;
  padding-bottom: 50px;
}
</style>
