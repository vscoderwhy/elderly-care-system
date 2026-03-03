<template>
  <div ref="chartRef" class="echarts-wrapper" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption } from 'echarts'

interface Props {
  option: EChartsOption
  height?: string
  width?: string
  theme?: string | object
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  width: '100%',
  theme: 'default'
})

const chartRef = ref<HTMLElement>()
let chartInstance: echarts.ECharts | null = null

const initChart = () => {
  if (!chartRef.value) return

  // 销毁已存在的实例
  if (chartInstance) {
    chartInstance.dispose()
  }

  // 创建新实例
  chartInstance = echarts.init(chartRef.value, props.theme, {
    renderer: 'canvas'
  })

  // 设置选项
  chartInstance.setOption(props.option)

  // 监听窗口大小变化
  window.addEventListener('resize', handleResize)
}

const handleResize = () => {
  chartInstance?.resize()
}

const updateChart = () => {
  if (chartInstance) {
    chartInstance.setOption(props.option, true)
  }
}

const watchTheme = () => {
  const observer = new MutationObserver(() => {
    if (chartInstance) {
      const option = chartInstance.getOption()
      chartInstance.dispose()
      initChart()
      chartInstance?.setOption(option)
    }
  })

  observer.observe(document.documentElement, {
    attributes: true,
    attributeFilter: ['data-theme']
  })

  return observer
}

onMounted(() => {
  nextTick(() => {
    initChart()
    watchTheme()
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chartInstance?.dispose()
  chartInstance = null
})

watch(
  () => props.option,
  () => {
    updateChart()
  },
  { deep: true }
)

defineExpose({
  resize: handleResize,
  getInstance: () => chartInstance
})
</script>

<style scoped lang="scss">
.echarts-wrapper {
  width: 100%;
  min-height: 200px;
}
</style>
