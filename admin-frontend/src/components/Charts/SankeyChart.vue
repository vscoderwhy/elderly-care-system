<template>
  <div ref="chartRef" :style="{ height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import * as echarts from 'echarts/core'

interface Props {
  nodes: any[]
  links: any[]
  height?: string
}

const props = withDefaults(defineProps<Props>(), { height: '400px' })
const chartRef = ref<HTMLElement>()

onMounted(() => {
  if (!chartRef.value) return
  const chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'item', triggerOn: 'mousemove' },
    series: [{
      type: 'sankey',
      layout: 'none',
      data: props.nodes,
      links: props.links,
      itemStyle: { borderWidth: 1, borderColor: '#aaa' },
      lineStyle: { color: 'gradient', curveness: 0.5 }
    }]
  })
})
</script>
