<template>
  <div ref="chartRef" :style="{ height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import * as echarts from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { RadarChart } from 'echarts/charts'
import { TitleComponent, TooltipComponent, LegendComponent } from 'echarts/components'

echarts.use([CanvasRenderer, RadarChart, TitleComponent, TooltipComponent, LegendComponent])

interface Props {
  data: any[]
  indicator: any[]
  height?: string
  title?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  title: ''
})

const chartRef = ref<HTMLElement>()
let chartInstance: any

const initChart = () => {
  if (!chartRef.value) return
  chartInstance = echarts.init(chartRef.value)
  updateChart()
}

const updateChart = () => {
  if (!chartInstance) return
  chartInstance.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c}' },
    legend: { data: props.data.map(d => d.name), bottom: 0 },
    radar: { indicator: props.indicator, center: ['50%', '50%'], radius: '65%' },
    series: [{
      type: 'radar',
      data: props.data.map(d => ({ name: d.name, value: d.value }))
    }]
  })
}

onMounted(initChart)
watch(() => props.data, updateChart, { deep: true })
</script>
