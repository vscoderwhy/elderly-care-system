<template>
  <div class="time-slider">
    <div class="time-slider__header">
      <div class="time-slider__current">
        {{ formattedCurrentTime }}
      </div>
      <div class="time-slider__controls">
        <el-button
          :icon="isPlaying ? VideoPause : VideoPlay"
          circle
          @click="togglePlay"
        />
        <el-button
          :icon="RefreshLeft"
          circle
          @click="reset"
        />
      </div>
      <div class="time-slider__speed">
        <el-select v-model="speed" size="small" style="width: 100px">
          <el-option label="0.5x" :value="0.5" />
          <el-option label="1x" :value="1" />
          <el-option label="2x" :value="2" />
          <el-option label="5x" :value="5" />
        </el-select>
      </div>
    </div>
    <div class="time-slider__body">
      <el-slider
        v-model="currentIndex"
        :min="0"
        :max="maxIndex"
        :step="1"
        :show-tooltip="false"
        @change="handleSliderChange"
      />
      <div class="time-slider__timeline">
        <div
          v-for="(point, index) in timePoints"
          :key="index"
          class="time-slider__point"
          :class="{ 'is-active': index === currentIndex }"
          :style="{ left: `${(index / maxIndex) * 100}%` }"
          @click="setCurrentIndex(index)"
        >
          <el-tooltip :content="point.label" placement="top">
            <div class="time-slider__dot"></div>
          </el-tooltip>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onUnmounted } from 'vue'
import { VideoPlay, VideoPause, RefreshLeft } from '@element-plus/icons-vue'

interface TimePoint {
  time: Date
  label: string
  data?: any
}

interface Props {
  timePoints: TimePoint[]
  modelValue: number
  autoplay?: boolean
  interval?: number
}

const props = withDefaults(defineProps<Props>(), {
  autoplay: false,
  interval: 2000
})

const emit = defineEmits<{
  'update:modelValue': [value: number]
  'change': [index: number, point: TimePoint]
}>()

const isPlaying = ref(props.autoplay)
const speed = ref(1)
const currentIndex = ref(props.modelValue)
const timer = ref<number>()

const maxIndex = computed(() => Math.max(0, props.timePoints.length - 1))

const formattedCurrentTime = computed(() => {
  const point = props.timePoints[currentIndex.value]
  if (!point) return '--'
  return point.label
})

const togglePlay = () => {
  isPlaying.value = !isPlaying.value
  if (isPlaying.value) {
    startPlay()
  } else {
    stopPlay()
  }
}

const startPlay = () => {
  stopPlay()
  const interval = props.interval / speed.value
  timer.value = window.setInterval(() => {
    if (currentIndex.value < maxIndex.value) {
      currentIndex.value++
      emitChange()
    } else {
      // 播放结束，停止或循环
      isPlaying.value = false
      stopPlay()
    }
  }, interval)
}

const stopPlay = () => {
  if (timer.value) {
    clearInterval(timer.value)
    timer.value = undefined
  }
}

const reset = () => {
  stopPlay()
  isPlaying.value = false
  currentIndex.value = 0
  emitChange()
}

const handleSliderChange = () => {
  emitChange()
}

const setCurrentIndex = (index: number) => {
  currentIndex.value = index
  emitChange()
}

const emitChange = () => {
  emit('update:modelValue', currentIndex.value)
  const point = props.timePoints[currentIndex.value]
  if (point) {
    emit('change', currentIndex.value, point)
  }
}

// 监听速度变化
watch(speed, () => {
  if (isPlaying.value) {
    startPlay()
  }
})

// 监听播放状态变化
watch(isPlaying, (newVal) => {
  if (newVal) {
    startPlay()
  } else {
    stopPlay()
  }
})

// 组件卸载时清理定时器
onUnmounted(() => {
  stopPlay()
})

// 监听外部值变化
watch(
  () => props.modelValue,
  (newVal) => {
    currentIndex.value = newVal
  }
)
</script>

<style scoped lang="scss">
.time-slider {
  padding: 16px;
  background: var(--card-bg);
  border-radius: var(--card-border-radius);
  box-shadow: var(--card-shadow);

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  &__current {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
  }

  &__controls {
    display: flex;
    gap: 8px;
  }

  &__speed {
    display: flex;
    align-items: center;
  }

  &__body {
    position: relative;
    padding-top: 20px;
  }

  &__timeline {
    position: absolute;
    top: 24px;
    left: 0;
    right: 0;
    height: 20px;
    pointer-events: none;
  }

  &__point {
    position: absolute;
    transform: translateX(-50%);
    pointer-events: auto;
    cursor: pointer;

    &.is-active {
      .time-slider__dot {
        background-color: var(--primary-color);
        transform: scale(1.5);
      }
    }
  }

  &__dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    background-color: var(--border-color);
    transition: var(--transition-base);
  }
}
</style>
