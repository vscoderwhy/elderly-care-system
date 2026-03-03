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
  layout?: string
  categories?: any[]
}

const props = withDefaults(defineProps<Props>(), { height: '450px' })
const chartRef = ref<HTMLElement>()

onMounted(() => {
  if (!chartRef.value) return
  const chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: {},
    series: [{
      type: 'graph',
      layout: props.layout || 'force',
      data: props.nodes.map(n => ({ ...n, itemStyle: { color: ['#409eff', '#67c23a', '#e6a23c'][n.category] } })),
      links: props.links,
      categories: props.categories?.map(c => ({ name: c.name })),
      roam: true,
      label: { show: true, position: 'right' },
      lineStyle: { color: '#ccc', curveness: 0.3 }
    }]
  })
})
</script>
