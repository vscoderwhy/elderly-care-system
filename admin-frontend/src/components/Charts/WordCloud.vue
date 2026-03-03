<template>
  <div ref="chartRef" class="word-cloud" :style="{ height: height }"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useECharts, getChartColors } from '@/composables/useECharts'

interface WordCloudData {
  name: string
  value: number
  category?: number
}

interface Props {
  data: WordCloudData[]
  height?: string
  title?: string
  sizeRange?: [number, number]
  rotationRange?: [number, number]
  shape?: 'circle' | 'cardioid' | 'diamond' | 'triangle-forward' | 'triangle' | 'pentagon' | 'star'
}

const props = withDefaults(defineProps<Props>(), {
  height: '400px',
  title: '',
  sizeRange: [12, 60],
  rotationRange: [-90, 90],
  shape: 'circle'
})

const chartRef = ref<HTMLElement>()
const { initChart, setOption, resize, watchTheme } = useECharts(chartRef)

const getOption = () => {
  const colors = getChartColors()

  // 为每个词分配颜色
  const dataWithColor = props.data.map((item, index) => ({
    ...item,
    textStyle: {
      color: colors[index % colors.length]
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
      show: true,
      backgroundColor: 'rgba(0, 0, 0, 0.8)',
      borderColor: 'transparent',
      textStyle: {
        color: '#fff'
      },
      formatter: (params: any) => {
        return `${params.name}<br/>频率: ${params.value}`
      }
    },
    series: [
      {
        type: 'wordCloud',
        shape: props.shape,
        left: 'center',
        top: 'center',
        width: '80%',
        height: '80%',
        right: null,
        bottom: null,
        sizeRange: props.sizeRange,
        rotationRange: props.rotationRange,
        rotationStep: 45,
        gridSize: 8,
        drawOutOfBound: false,
        textStyle: {
          fontFamily: 'sans-serif',
          fontWeight: 'bold',
          color: () => {
            return colors[Math.floor(Math.random() * colors.length)]
          }
        },
        emphasis: {
          textStyle: {
            textShadowBlur: 10,
            textShadowColor: '#333'
          }
        },
        data: dataWithColor
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
.word-cloud {
  width: 100%;
}
</style>
