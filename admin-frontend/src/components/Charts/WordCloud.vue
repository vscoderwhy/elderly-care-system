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
    tooltip: {},
    series: [{
      type: 'wordCloud',
      shape: 'circle',
      left: 'center', top: 'center',
      width: '70%', height: '80%',
      sizeRange: [12, 50],
      rotationRange: [-90, 90],
      data: props.data
    }]
  })
})
</script>
