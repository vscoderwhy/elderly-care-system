<template>
  <div ref="chartRef" class="funnel-chart" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useECharts, getChartColors } from '@/composables/useECharts'

interface FunnelData {
  name: string
  value: number
}

interface Props {
  data: FunnelData[]
  height?: string
  title?: string
  subtitle?: string
  sort?: 'ascending' | 'descending' | 'none'
  gap?: number
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  title: '',
  subtitle: '',
  sort: 'descending',
  gap: 0
})

const chartRef = ref<HTMLElement>()
const { initChart, setOption, resize, watchTheme } = useECharts(chartRef)

const getOption = () => {
  const colors = getChartColors()

  return {
    title: {
      text: props.title,
      subtext: props.subtitle,
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
      formatter: '{a} <br/>{b} : {c}%'
    },
    legend: {
      orient: 'vertical',
      left: 'left',
      data: props.data.map((d) => d.name)
    },
    series: [
      {
        name: props.title || '漏斗图',
        type: 'funnel',
        left: '20%',
        top: '15%',
        bottom: '15%',
        width: '60%',
        min: 0,
        max: Math.max(...props.data.map((d) => d.value)),
        minSize: '0%',
        maxSize: '100%',
        sort: props.sort,
        gap: props.gap,
        label: {
          show: true,
          fontSize: 12,
          formatter: '{b}: {c}'
        },
        labelLine: {
          length: 10,
          lineStyle: {
            width: 1,
            type: 'solid'
          }
        },
        itemStyle: {
          borderColor: '#fff',
          borderWidth: 1
        },
        emphasis: {
          label: {
            fontSize: 14,
            fontWeight: 'bold'
          },
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        },
        data: props.data.map((item, index) => ({
          ...item,
          itemStyle: {
            color: colors[index % colors.length]
          }
        }))
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
  () => props.data,
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
.funnel-chart {
  width: 100%;
}
</style>
