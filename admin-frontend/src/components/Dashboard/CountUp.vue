<template>
  <span>{{ displayValue }}</span>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'

interface CountUpOptions {
  startVal?: number
  decimals?: number
  duration?: number
  useEasing?: boolean
  separator?: string
}

interface Props {
  endVal: number
  options?: CountUpOptions
}

const props = withDefaults(defineProps<Props>(), {
  options: () => ({
    startVal: 0,
    decimals: 0,
    duration: 2000,
    useEasing: true,
    separator: ','
  })
})

const displayValue = ref('0')

const easeOutExpo = (t: number): number => {
  return t === 1 ? 1 : 1 - Math.pow(2, -10 * t)
}

const formatNumber = (num: number, decimals: number, separator: string): string => {
  const fixed = num.toFixed(decimals)
  const parts = fixed.split('.')
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, separator)
  return parts.join('.')
}

const animate = () => {
  const {
    startVal = 0,
    decimals = 0,
    duration = 2000,
    useEasing = true,
    separator = ','
  } = props.options

  const range = props.endVal - startVal
  const startTime = performance.now()

  const update = (currentTime: number) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    const easedValue = useEasing ? easeOutExpo(progress) : progress
    const currentValue = startVal + range * easedValue

    displayValue.value = formatNumber(currentValue, decimals, separator)

    if (progress < 1) {
      requestAnimationFrame(update)
    }
  }

  requestAnimationFrame(update)
}

onMounted(() => {
  animate()
})

watch(
  () => props.endVal,
  () => {
    animate()
  }
)
</script>
