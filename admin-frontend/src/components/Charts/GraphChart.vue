<template>
  <div ref="chartRef" class="graph-chart" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useECharts, getChartColors } from '@/composables/useECharts'

interface GraphNode {
  id: string
  name: string
  symbolSize?: number
  category?: number
  itemStyle?: { color?: string }
  x?: number
  y?: number
  value?: number
}

interface GraphLink {
  source: string | number
  target: string | number
  value?: number
  lineStyle?: { color?: string; width?: number; curveness?: number }
}

interface GraphCategory {
  name: string
}

interface Props {
  nodes: GraphNode[]
  links: GraphLink[]
  categories?: GraphCategory[]
  height?: string
  title?: string
  layout?: 'none' | 'circular' | 'force'
  edgeSymbol?: [string, string]
  roam?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  height: '500px',
  title: '',
  layout: 'force',
  edgeSymbol: () => ['none', 'arrow'],
  roam: true
})

const chartRef = ref<HTMLElement>()
const { initChart, setOption, resize, watchTheme } = useECharts(chartRef)

const getOption = () => {
  const colors = getChartColors()

  // 为节点分配颜色
  const nodesWithColor = props.nodes.map((node) => ({
    ...node,
    itemStyle: {
      color: node.itemStyle?.color || colors[(node.category || 0) % colors.length]
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
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: 'transparent',
      textStyle: {
        color: '#fff'
      },
      formatter: (params: any) => {
        if (params.dataType === 'node') {
          return `${params.data.name}<br/>ID: ${params.data.id}`
        } else if (params.dataType === 'edge') {
          return `${params.data.source} → ${params.data.target}`
        }
        return ''
      }
    },
    legend: props.categories
      ? [{
          data: props.categories.map((c) => c.name),
          bottom: 10
        }]
      : [],
    series: [
      {
        type: 'graph',
        layout: props.layout,
        data: nodesWithColor,
        links: props.links,
        categories: props.categories || [],
        title: {
          text: props.title
        },
        roam: props.roam,
        label: {
          show: true,
          position: 'right',
          formatter: '{b}',
          fontSize: 12
        },
        labelLayout: {
          hideOverlap: true
        },
        scaleLimit: {
          min: 0.4,
          max: 2
        },
        lineStyle: {
          color: 'source',
          curveness: 0.3,
          width: 1,
          opacity: 0.5
        },
        edgeSymbol: props.edgeSymbol,
        edgeSymbolSize: [4, 8],
        edgeLabel: {
          fontSize: 12
        },
        emphasis: {
          focus: 'adjacency',
          lineStyle: {
            width: 3
          }
        },
        force: props.layout === 'force'
          ? {
              repulsion: 200,
              gravity: 0.1,
              edgeLength: 100,
              layoutAnimation: true
            }
          : undefined
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
  () => [props.nodes, props.links, props.categories],
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
.graph-chart {
  width: 100%;
}
</style>
