<template>
  <div class="stat-card" :class="[`stat-${type}`, { 'is-loading': loading }]">
    <!-- 背景装饰 -->
    <div class="stat-bg">
      <div class="stat-bg-circle stat-bg-circle-1"></div>
      <div class="stat-bg-circle stat-bg-circle-2"></div>
    </div>

    <!-- 卡片内容 -->
    <div class="stat-content">
      <div class="stat-header">
        <div class="stat-icon" :style="{ background: iconBg }">
          <component :is="icon" />
        </div>
        <div class="stat-extra">
          <slot name="extra"></slot>
        </div>
      </div>

      <div class="stat-body">
        <div class="stat-value-wrapper">
          <count-to
            :start-val="0"
            :end-val="displayValue"
            :duration="1000"
            :decimals="decimals"
            class="stat-value"
          />
          <span v-if="unit" class="stat-unit">{{ unit }}</span>
        </div>

        <div class="stat-label">{{ label }}</div>

        <!-- 趋势指示 -->
        <div v-if="trend !== undefined" class="stat-trend" :class="trendClass">
          <el-icon>
            <component :is="trendIcon" />
          </el-icon>
          <span>{{ Math.abs(trend) }}%</span>
          <span class="trend-text">{{ trendText }}</span>
        </div>
      </div>

      <!-- 迷你图表 -->
      <div v-if="showSparkline" class="stat-sparkline">
        <slot name="sparkline">
          <ECharts :option="sparklineOption" height="60px" />
        </slot>
      </div>
    </div>

    <!-- 加载遮罩 -->
    <div v-if="loading" class="stat-loading">
      <el-icon class="is-loading"><Loading /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ArrowUp, ArrowDown, Minus, Loading } from '@element-plus/icons-vue'
import CountTo from 'vue-count-to'
import ECharts from './ECharts.vue'

interface Props {
  label: string
  value: number | string
  unit?: string
  type?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  icon?: any
  trend?: number
  decimals?: number
  loading?: boolean
  sparkline?: number[]
  showSparkline?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  unit: '',
  type: 'primary',
  decimals: 0,
  loading: false,
  showSparkline: false,
  sparkline: () => []
})

// 显示数值（处理字符串类型）
const displayValue = computed(() => {
  if (typeof props.value === 'string') {
    return parseFloat(props.value) || 0
  }
  return props.value
})

// 趋势相关
const trendClass = computed(() => {
  if (!props.trend) return ''
  return props.trend > 0 ? 'trend-up' : props.trend < 0 ? 'trend-down' : 'trend-flat'
})

const trendIcon = computed(() => {
  if (!props.trend) return Minus
  return props.trend > 0 ? ArrowUp : ArrowDown
})

const trendText = computed(() => {
  if (!props.trend) return '持平'
  return props.trend > 0 ? '较上期' : '较上期'
})

// 图标背景色
const iconBg = computed(() => {
  const gradients = {
    primary: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
    success: 'linear-gradient(135deg, #84fab0 0%, #8fd3f4 100%)',
    warning: 'linear-gradient(135deg, #fccb90 0%, #d57eeb 100%)',
    danger: 'linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%)',
    info: 'linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%)'
  }
  return gradients[props.type]
})

// 迷你图表配置
const sparklineOption = computed(() => {
  const colors = {
    primary: '#667eea',
    success: '#67c23a',
    warning: '#e6a23c',
    danger: '#f56c6c',
    info: '#909399'
  }

  return {
    grid: { top: 0, left: 0, right: 0, bottom: 0 },
    xAxis: {
      type: 'category',
      show: false,
      data: props.sparkline.map((_, i) => i)
    },
    yAxis: {
      type: 'value',
      show: false
    },
    series: [{
      type: 'line',
      data: props.sparkline,
      smooth: true,
      symbol: 'none',
      lineStyle: {
        color: colors[props.type],
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
            { offset: 0, color: colors[props.type] + '40' },
            { offset: 1, color: colors[props.type] + '05' }
          ]
        }
      }
    }]
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
  overflow: hidden;
  transition: var(--transition-base);
  cursor: pointer;

  &:hover {
    transform: translateY(-4px);
    box-shadow: var(--card-shadow-hover);

    .stat-bg-circle {
      transform: scale(1.2);
    }
  }

  &.is-loading {
    pointer-events: none;
    opacity: 0.6;
  }
}

.stat-bg {
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
  z-index: 0;
}

.stat-bg-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  transition: transform 0.6s cubic-bezier(0.4, 0, 0.2, 1);

  &-1 {
    width: 120px;
    height: 120px;
    top: -40px;
    right: -40px;
    background: var(--primary-color);
  }

  &-2 {
    width: 80px;
    height: 80px;
    bottom: -30px;
    left: -30px;
    background: var(--success-color);
  }
}

.stat-content {
  position: relative;
  z-index: 1;
}

.stat-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  color: #fff;
}

.stat-body {
  margin-top: 8px;
}

.stat-value-wrapper {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: var(--text-primary);
  line-height: 1;
}

.stat-unit {
  font-size: 14px;
  color: var(--text-secondary);
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-top: 4px;
}

.stat-trend {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  font-size: 12px;
  padding: 4px 8px;
  border-radius: 4px;
  background: var(--bg-tertiary);

  &.trend-up {
    color: var(--success-color);
    background: rgba(103, 194, 58, 0.1);
  }

  &.trend-down {
    color: var(--danger-color);
    background: rgba(245, 108, 108, 0.1);
  }

  &.trend-flat {
    color: var(--info-color);
  }

  .trend-text {
    margin-left: 4px;
    opacity: 0.8;
  }
}

.stat-sparkline {
  margin-top: 12px;
  height: 60px;
}

.stat-loading {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.8);
  z-index: 10;

  .el-icon {
    font-size: 32px;
    color: var(--primary-color);
  }
}

[data-theme='dark'] {
  .stat-loading {
    background: rgba(0, 0, 0, 0.6);
  }

  .stat-value {
    color: var(--text-primary);
  }
}
</style>
