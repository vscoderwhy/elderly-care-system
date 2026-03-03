<template>
  <div class="stat-card" :class="{ 'is-dark': isDark }">
    <div class="stat-card__header">
      <div class="stat-card__icon" :style="{ background: iconBg }">
        <component :is="icon" class="stat-card__icon-svg" />
      </div>
      <el-tag v-if="trend" :type="trendType" size="small">
        <template v-if="trend > 0">
          <el-icon><ArrowUp /></el-icon>
          {{ trend }}%
        </template>
        <template v-else>
          <el-icon><ArrowDown /></el-icon>
          {{ Math.abs(trend) }}%
        </template>
      </el-tag>
    </div>
    <div class="stat-card__content">
      <div class="stat-card__value">
        <CountUp :end-val="value" :duration="1500" :options="{ separator: ',' }" />
        <span v-if="unit" class="stat-card__unit">{{ unit }}</span>
      </div>
      <div class="stat-card__label">{{ label }}</div>
    </div>
    <div v-if="showSparkline" class="stat-card__sparkline">
      <ECharts :option="sparklineOption" :height="40" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { ArrowUp, ArrowDown } from '@element-plus/icons-vue'
import { useDarkMode } from '@/composables/useDarkMode'
import CountUp from './CountUp.vue'
import ECharts from './ECharts.vue'

interface Props {
  label: string
  value: number
  unit?: string
  icon?: any
  color?: string
  trend?: number
  showSparkline?: boolean
  sparklineData?: number[]
}

const props = withDefaults(defineProps<Props>(), {
  unit: '',
  color: '#409eff',
  trend: 0,
  showSparkline: false,
  sparklineData: () => []
})

const { isDark } = useDarkMode()

const iconBg = computed(() => {
  return `linear-gradient(135deg, ${props.color}22 0%, ${props.color}44 100%)`
})

const trendType = computed(() => {
  if (props.trend > 0) return 'success'
  if (props.trend < 0) return 'danger'
  return 'info'
})

const sparklineOption = computed(() => {
  const isDarkMode = isDark.value
  return {
    grid: {
      top: 0,
      left: 0,
      right: 0,
      bottom: 0
    },
    xAxis: {
      type: 'category',
      show: false,
      data: props.sparklineData.map((_, i) => i)
    },
    yAxis: {
      type: 'value',
      show: false
    },
    series: [
      {
        data: props.sparklineData,
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: {
          color: props.color,
          width: 2
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              { offset: 0, color: `${props.color}44` },
              { offset: 1, color: `${props.color}00` }
            ]
          }
        }
      }
    ]
  }
})
</script>

<style scoped lang="scss">
.stat-card {
  position: relative;
  padding: 20px;
  background: var(--card-bg);
  border-radius: var(--card-border-radius);
  box-shadow: var(--card-shadow);
  transition: var(--transition-base);
  overflow: hidden;

  &:hover {
    box-shadow: var(--card-shadow-hover);
    transform: translateY(-2px);
  }

  &__header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  &__icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  &__icon-svg {
    width: 24px;
    height: 24px;
    color: v-bind(color);
  }

  &__content {
    position: relative;
    z-index: 1;
  }

  &__value {
    display: flex;
    align-items: baseline;
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 4px;
  }

  &__unit {
    font-size: 14px;
    font-weight: 400;
    color: var(--text-secondary);
    margin-left: 4px;
  }

  &__label {
    font-size: 14px;
    color: var(--text-secondary);
  }

  &__sparkline {
    margin-top: 12px;
  }
}

// 背景装饰
.stat-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 200%;
  height: 200%;
  background: radial-gradient(
    circle,
    v-bind(color)11 0%,
    transparent 70%
  );
  pointer-events: none;
}
</style>
