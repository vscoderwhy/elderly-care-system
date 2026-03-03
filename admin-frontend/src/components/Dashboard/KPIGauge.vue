<template>
  <div ref="chartRef" class="kpi-gauge" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from 'vue'
import { useECharts } from '@/composables/useECharts'

interface KPIRange {
  color: string
  start: number
  end: number
  label?: string
}

interface Props {
  value: number
  min?: number
  max?: number
  title?: string
  unit?: string
  height?: string
  ranges?: KPIRange[]
  splitNumber?: number
}

const props = withDefaults(defineProps<Props>(), {
  min: 0,
  max: 100,
  title: '',
  unit: '',
  height: '300px',
  splitNumber: 10
})

const chartRef = ref<HTMLElement>()
const { initChart, setOption, resize, watchTheme } = useECharts(chartRef)

const getColorByValue = (value: number): string => {
  if (props.ranges && props.ranges.length > 0) {
    for (const range of props.ranges) {
      if (value >= range.start && value < range.end) {
        return range.color
      }
    }
    return props.ranges[props.ranges.length - 1]?.color || '#409eff'
  }

  // 默认颜色逻辑
  const percentage = (value - props.min) / (props.max - props.min)
  if (percentage < 0.3) return '#67c23a'
  if (percentage < 0.7) return '#409eff'
  return '#f56c6c'
}

const getOption = () => {
  const color = getColorByValue(props.value)

  return {
    title: {
      text: props.title,
      left: 'center',
      top: '10%',
      textStyle: {
        fontSize: 14,
        color: 'inherit'
      }
    },
    series: [
      {
        type: 'gauge',
        startAngle: 180,
        endAngle: 0,
        min: props.min,
        max: props.max,
        splitNumber: props.splitNumber,
        center: ['50%', '70%'],
        radius: '90%',
        itemStyle: {
          color: color,
          shadowColor: color,
          shadowBlur: 10,
          shadowOffsetX: 2,
          shadowOffsetY: 2
        },
        progress: {
          show: true,
          roundCap: true,
          width: 18
        },
        pointer: {
          icon: 'path://M2090.36389,615.30999 L2090.36389,615.30999 C2091.48372,615.30999 2092.40383,616.194028 2092.44859,617.312956 L2096.90698,728.755929 C2097.05155,732.369577 2094.2393,735.416212 2090.62566,735.56078 C2090.53845,735.564269 2090.45117,735.566014 2090.36389,735.566014 L2090.36389,735.566014 C2086.74736,735.566014 2083.81557,732.63423 2083.81557,729.017692 C2083.81557,728.930412 2083.81732,728.84314 2083.82081,728.755929 L2088.2792,617.312956 C2088.32396,616.194028 2089.24407,615.30999 2090.36389,615.30999 Z',
          length: '75%',
          width: 16,
          offsetCenter: [0, '5%']
        },
        axisLine: {
          roundCap: true,
          lineStyle: {
            width: 18,
            color: props.ranges && props.ranges.length > 0
              ? props.ranges.map(r => [r.end / props.max, r.color])
              : [[1, 'rgba(255,255,255,0.1)']]
          }
        },
        axisTick: {
          splitNumber: 2,
          lineStyle: {
            width: 2,
            color: 'inherit'
          }
        },
        splitLine: {
          length: 12,
          lineStyle: {
            width: 3,
            color: 'inherit'
          }
        },
        axisLabel: {
          distance: 30,
          color: 'inherit',
          fontSize: 12,
          formatter: (value: number) => {
            return Math.round(value).toString()
          }
        },
        detail: {
          valueAnimation: true,
          formatter: `{value}${props.unit}`,
          color: 'inherit',
          fontSize: 32,
          offsetCenter: [0, '20%'],
          fontWeight: 'bold'
        },
        data: [
          {
            value: props.value
          }
        ]
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
  () => props.value,
  () => {
    setOption(getOption())
  }
)

defineExpose({
  resize
})
</script>

<style scoped lang="scss">
.kpi-gauge {
  width: 100%;
}
</style>
