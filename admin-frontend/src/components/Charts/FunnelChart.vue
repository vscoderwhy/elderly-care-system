<template>
  <div ref="chartRef" :style="{ height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts/core'

interface Props {
  data: any[]
  height?: string
}

const props = withDefaults(defineProps<Props>(), { height: '350px' })
const chartRef = ref<HTMLElement>()

onMounted(() => {
  if (!chartRef.value) return
  const chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c}' },
    series: [{
      type: 'funnel',
      left: '10%', top: 60, bottom: 60, width: '80%',
      data: props.data,
      label: { show: true, position: 'inside' }
    }]
  })
})
</script>
