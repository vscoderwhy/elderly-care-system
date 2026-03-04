<template>
  <transition name="loading-fade">
    <div v-if="isLoading" class="global-loading" :style="{ background: background }">
      <div class="loading-content">
        <el-icon class="is-loading loading-icon" :size="40">
          <Loading />
        </el-icon>
        <p v-if="text" class="loading-text">{{ text }}</p>
      </div>
    </div>
  </transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Loading } from '@element-plus/icons-vue'
import { useLoading } from '@/composables/useLoading'

const { isLoading, text } = useLoading()
const background = computed(() => 'rgba(0, 0, 0, 0.7)')
</script>

<style scoped lang="scss">
.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  backdrop-filter: blur(4px);

  .loading-content {
    text-align: center;

    .loading-icon {
      color: var(--primary-color);
      font-size: 40px;
      margin-bottom: 16px;
    }

    .loading-text {
      color: #fff;
      font-size: 16px;
      margin: 0;
    }
  }
}

.loading-fade-enter-active,
.loading-fade-leave-active {
  transition: opacity 0.3s;
}

.loading-fade-enter-from,
.loading-fade-leave-to {
  opacity: 0;
}
</style>
