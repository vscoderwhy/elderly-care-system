<template>
  <div ref="chartRef" class="sankey-chart" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useECharts, getChartColors } from '@/composables/useECharts'

interface SankeyNode {
  name: string
  itemStyle?: { color?: string }
}

interface SankeyLink {
  source: string | number
  target: string | number
  value: number
}

interface Props {
  nodes: SankeyNode[]
  links: SankeyLink[]
  height?: string
  title?: string
  nodeWidth?: number
  nodeGap?: number
}

const props = withDefaults(defineProps<Props>(), {
  height: '500px',
  title: '',
  nodeWidth: 20,
  nodeGap: 8
})

const chartRef = ref<HTMLElement>()
const { initChart, setOption, resize, watchTheme } = useECharts(chartRef)

const getOption = () => {
  const colors = getChartColors()

  // 为节点分配颜色
  const nodesWithColor = props.nodes.map((node, index) => ({
    ...node,
    itemStyle: {
      color: node.itemStyle?.color || colors[index % colors.length]
    }
  }))

  return {
    title: {
      text: props.title,
      left: 'center',
      textStyle: {
        fontSize: 16,
        fontWeight: 600
      }
    },
    tooltip: {
      trigger: 'item',
      triggerOn: 'mousemove',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: 'transparent',
      textStyle: {
        color: '#fff'
      },
      formatter: (params: any) => {
        if (params.dataType === 'edge') {
          return `${params.data.source} > ${params.data.target}<br/>流量: ${params.data.value}`
        }
        return `${params.data.name}<br/>流量: ${params.data.value}`
      }
    },
    series: [
      {
        type: 'sankey',
        layout: 'none',
        emphasis: {
          focus: 'adjacency'
        },
        data: nodesWithColor,
        links: props.links,
        top: '10%',
        bottom: '10%',
        left: '10%',
        right: '10%',
        nodeWidth: props.nodeWidth,
        nodeGap: props.nodeGap,
        label: {
          fontSize: 12,
          color: 'inherit'
        },
        lineStyle: {
          color: 'gradient',
          curveness: 0.5,
          opacity: 0.3
        },
        itemStyle: {
          borderWidth: 0,
          borderColor: '#aaa'
        }
      }
    ]
  }
}

onMounted(() => {
  setOption(getOption())
  watchTheme(() => {
    setOption(getOption())
  })
})

watch(
  () => [props.nodes, props.links],
  () => {
    setOption(getOption())
  },
  { deep: true }
)

defineExpose({
  resize
})
</script>

<style scoped lang="scss">
.sankey-chart {
  width: 100%;
}
</style>
