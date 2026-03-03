import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption, ECharts } from 'echarts'

export interface ChartTheme {
  backgroundColor?: string
  textStyle?: any
  title?: any
  legend?: any
  categoryAxis?: any
  valueAxis?: any
}

export const useECharts = (
  domRef: Ref<HTMLElement | undefined>,
  theme: string | ChartTheme = 'default'
) => {
  const chartInstance = ref<ECharts>()
  const loading = ref(false)

  // 根据当前主题获取 ECharts 主题配置
  const getChartTheme = (): ChartTheme => {
    const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
    return {
      backgroundColor: 'transparent',
      textStyle: {
        fontFamily: 'Inter, system-ui, -apple-system, sans-serif'
      },
      title: {
        textStyle: {
          color: isDark ? '#e5eaf3' : '#303133'
        },
        subtextStyle: {
          color: isDark ? '#cfd3dc' : '#606266'
        }
      },
      legend: {
        textStyle: {
          color: isDark ? '#cfd3dc' : '#606266'
        }
      },
      categoryAxis: {
        axisLine: {
          lineStyle: {
            color: isDark ? '#363636' : '#dcdfe6'
          }
        },
        axisLabel: {
          color: isDark ? '#cfd3dc' : '#606266'
        },
        splitLine: {
          lineStyle: {
            color: isDark ? '#262727' : '#ebeef5'
          }
        }
      },
      valueAxis: {
        axisLine: {
          lineStyle: {
            color: isDark ? '#363636' : '#dcdfe6'
          }
        },
        axisLabel: {
          color: isDark ? '#cfd3dc' : '#606266'
        },
        splitLine: {
          lineStyle: {
            color: isDark ? '#262727' : '#ebeef5'
          }
        }
      }
    }
  }

  const initChart = () => {
    if (!domRef.value) return

    // 销毁已存在的实例
    if (chartInstance.value) {
      chartInstance.value.dispose()
    }

    // 创建新实例
    chartInstance.value = echarts.init(domRef.value, theme === 'dark' ? 'dark' : undefined, {
      renderer: 'canvas',
      devicePixelRatio: window.devicePixelRatio
    })

    // 监听窗口大小变化
    window.addEventListener('resize', handleResize)
  }

  const setOption = (option: EChartsOption, notMerge?: boolean, lazyUpdate?: boolean) => {
    if (!chartInstance.value) {
      console.warn('[ECharts] Chart instance not initialized')
      return
    }

    const themeOption = getChartTheme()
    const mergedOption = {
      ...option,
      backgroundColor: option.backgroundColor || themeOption.backgroundColor
    } as EChartsOption

    chartInstance.value.setOption(mergedOption, notMerge, lazyUpdate)
  }

  const resize = (opts?: { width?: number | string; height?: number | string }) => {
    chartInstance.value?.resize(opts)
  }

  const handleResize = () => {
    chartInstance.value?.resize()
  }

  const showLoading = (type?: string, opts?: any) => {
    loading.value = true
    chartInstance.value?.showLoading(type, opts)
  }

  const hideLoading = () => {
    loading.value = false
    chartInstance.value?.hideLoading()
  }

  const clear = () => {
    chartInstance.value?.clear()
  }

  const dispose = () => {
    window.removeEventListener('resize', handleResize)
    chartInstance.value?.dispose()
    chartInstance.value = undefined
  }

  // 获取图表实例
  const getInstance = () => chartInstance.value

  // 导出图片
  const getDataURL = (opts?: {
    type?: string
    pixelRatio?: number
    backgroundColor?: string
  }) => {
    return chartInstance.value?.getDataURL(opts)
  }

  // 监听主题变化，重新初始化图表
  const watchTheme = (callback: () => void) => {
    const observer = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === 'data-theme') {
          // 重新设置选项以应用新主题
          if (chartInstance.value) {
            const option = chartInstance.value.getOption()
            chartInstance.value.dispose()
            initChart()
            setOption(option)
          }
          callback?.()
        }
      })
    })

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ['data-theme']
    })

    return observer
  }

  onMounted(() => {
    initChart()
  })

  onUnmounted(() => {
    dispose()
  })

  return {
    chartInstance,
    loading,
    initChart,
    setOption,
    resize,
    showLoading,
    hideLoading,
    clear,
    dispose,
    getInstance,
    getDataURL,
    watchTheme
  }
}

// 图表颜色配置
export const chartColors = {
  light: [
    '#5470c6', '#91cc75', '#fac858', '#ee6666', '#73c0de',
    '#3ba272', '#fc8452', '#9a60b4', '#ea7ccc'
  ],
  dark: [
    '#5a8fff', '#75c968', '#ffcc00', '#ff6b6b', '#4ecdc4',
    '#45b7d1', '#96ceb4', '#d4a5a5', '#9b59b6'
  ]
}

export const getChartColors = () => {
  const isDark = document.documentElement.getAttribute('data-theme') === 'dark'
  return isDark ? chartColors.dark : chartColors.light
}
