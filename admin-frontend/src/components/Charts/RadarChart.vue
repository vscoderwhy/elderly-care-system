<template>
  <div ref="chartRef" class="radar-chart" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import * as echarts from 'echarts'
import { useECharts, getChartColors } from '@/composables/useECharts'

interface RadarData {
  name: string
  value: number[]
}

interface Props {
  data: RadarData[]
  indicator: Array<{ name: string; max: number; min?: number }>
  height?: string
  title?: string
  subtitle?: string
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  title: '',
  subtitle: ''
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
    legend: {
      data: props.data.map((d) => d.name),
      bottom: 10,
      type: 'scroll'
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: 'transparent',
      textStyle: {
        color: '#fff'
      }
    },
    radar: {
      indicator: props.indicator,
      radius: '65%',
      center: ['50%', '50%'],
      splitNumber: 4,
      nameGap: 12,
      axisName: {
        color: 'inherit',
        fontSize: 12
      },
      splitArea: {
        areaStyle: {
          color: ['rgba(255, 255, 255, 0.05)', 'rgba(255, 255, 255, 0.1)']
        }
      },
      splitLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.2)'
        }
      },
      axisLine: {
        lineStyle: {
          color: 'rgba(255, 255, 255, 0.2)'
        }
      }
    },
    series: [
      {
        type: 'radar',
        data: props.data.map((item, index) => ({
          value: item.value,
          name: item.name,
          symbol: 'circle',
          symbolSize: 6,
          lineStyle: {
            color: colors[index % colors.length],
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
                { offset: 0, color: colors[index % colors.length] + 'cc' },
                { offset: 1, color: colors[index % colors.length] + '22' }
              ]
            }
          },
          itemStyle: {
            color: colors[index % colors.length],
            borderColor: '#fff',
            borderWidth: 2
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
  () => [props.data, props.indicator],
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
.radar-chart {
  width: 100%;
}
</style>
