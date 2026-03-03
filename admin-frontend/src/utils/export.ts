import * as XLSX from 'xlsx'

export interface ExcelColumn {
  key: string
  title: string
  width?: number
  formatter?: (row: any) => any
}

/**
 * 导出Excel
 */
export const exportToExcel = (
  data: any[],
  columns: ExcelColumn[],
  filename: string
) => {
  if (!data || data.length === 0) {
    console.warn('No data to export')
    return
  }

  // 格式化数据
  const formattedData = data.map(row => {
    const obj: any = {}
    columns.forEach(col => {
      obj[col.title] = col.formatter ? col.formatter(row) : row[col.key]
    })
    return obj
  })

  // 创建工作表
  const worksheet = XLSX.utils.json_to_sheet(formattedData)

  // 设置列宽
  const colWidths = columns.map(col => ({ wch: col.width || 15 }))
  worksheet['!cols'] = colWidths

  // 创建工作簿
  const workbook = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(workbook, worksheet, 'Sheet1')

  // 生成文件名（带时间戳）
  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  const fullFilename = `${filename}_${timestamp}.xlsx`

  // 下载
  XLSX.writeFile(workbook, fullFilename)
}

/**
 * 导出老人列表
 */
export const exportElderlyList = (data: any[]) => {
  const columns: ExcelColumn[] = [
    { key: 'id', title: 'ID', width: 8 },
    { key: 'name', title: '姓名', width: 12 },
    { key: 'gender', title: '性别', width: 8 },
    { key: 'age', title: '年龄', width: 8 },
    { key: 'bedNumber', title: '床位号', width: 15 },
    { key: 'careLevel', title: '护理等级', width: 12 },
    { key: 'healthScore', title: '健康评分', width: 12 },
    { key: 'status', title: '状态', width: 10, formatter: (row) => {
      const map: Record<string, string> = {
        'in_hospital': '在院',
        'hospitalized': '住院',
        'leave': '请假',
        'discharged': '出院'
      }
      return map[row.status] || row.status
    }},
    { key: 'admissionDate', title: '入院日期', width: 15 },
    { key: 'emergencyContact', title: '紧急联系人', width: 15 },
    { key: 'emergencyPhone', title: '紧急电话', width: 15 },
    { key: 'medicalHistory', title: '既往病史', width: 30 }
  ]
  exportToExcel(data, columns, '老人名单')
}

/**
 * 导出护理记录
 */
export const exportCareRecords = (data: any[]) => {
  const columns: ExcelColumn[] = [
    { key: 'id', title: 'ID', width: 8 },
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'careType', title: '护理类型', width: 15 },
    { key: 'content', title: '护理内容', width: 40 },
    { key: 'careTime', title: '护理时间', width: 18 },
    { key: 'nurseName', title: '护理人员', width: 12 },
    { key: 'evaluation', title: '效果评价', width: 15 },
    { key: 'tags', title: '标签', width: 20, formatter: (row) => {
      return Array.isArray(row.tags) ? row.tags.join(', ') : ''
    }}
  ]
  exportToExcel(data, columns, '护理记录')
}

/**
 * 导出财务报表
 */
export const exportFinancialReport = (data: any[]) => {
  const columns: ExcelColumn[] = [
    { key: 'billNo', title: '账单号', width: 20 },
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'bedNumber', title: '床位号', width: 12 },
    { key: 'billType', title: '费用类型', width: 12 },
    { key: 'amount', title: '金额(元)', width: 12 },
    { key: 'billDate', title: '账单日期', width: 15 },
    { key: 'dueDate', title: '应付日期', width: 15 },
    { key: 'paidDate', title: '支付日期', width: 15 },
    { key: 'paymentMethod', title: '支付方式', width: 12, formatter: (row) => {
      const map: Record<string, string> = {
        'wechat': '微信支付',
        'alipay': '支付宝',
        'cash': '现金',
        'pos': 'POS机',
        'bank': '银行转账'
      }
      return map[row.paymentMethod] || row.paymentMethod || '-'
    }},
    { key: 'status', title: '状态', width: 10, formatter: (row) => {
      const map: Record<string, string> = {
        'paid': '已支付',
        'unpaid': '未支付',
        'overdue': '已逾期'
      }
      return map[row.status] || row.status
    }},
    { key: 'remark', title: '备注', width: 30 }
  ]
  exportToExcel(data, columns, '财务报表')
}

/**
 * 导出健康数据
 */
export const exportHealthData = (data: any[]) => {
  const columns: ExcelColumn[] = [
    { key: 'id', title: 'ID', width: 8 },
    { key: 'elderlyName', title: '老人姓名', width: 12 },
    { key: 'type', title: '数据类型', width: 12 },
    { key: 'value', title: '数值', width: 12 },
    { key: 'unit', title: '单位', width: 8 },
    { key: 'recordTime', title: '记录时间', width: 18 },
    { key: 'recordBy', title: '记录人', width: 12 },
    { key: 'isAbnormal', title: '是否异常', width: 10, formatter: (row) => row.isAbnormal ? '是' : '否' },
    { key: 'remark', title: '备注', width: 30 }
  ]
  exportToExcel(data, columns, '健康数据')
}

/**
 * 导出员工列表
 */
export const exportStaffList = (data: any[]) => {
  const columns: ExcelColumn[] = [
    { key: 'id', title: 'ID', width: 8 },
    { key: 'name', title: '姓名', width: 12 },
    { key: 'gender', title: '性别', width: 8 },
    { key: 'department', title: '部门', width: 15 },
    { key: 'position', title: '职位', width: 12 },
    { key: 'phone', title: '联系电话', width: 15 },
    { key: 'status', title: '状态', width: 10, formatter: (row) => {
      const map: Record<string, string> = {
        'active': '在职',
        'leave': '请假',
        'resigned': '离职'
      }
      return map[row.status] || row.status
    }},
    { key: 'hireDate', title: '入职日期', width: 15 }
  ]
  exportToExcel(data, columns, '员工名单')
}

/**
 * 导出为CSV
 */
export const exportToCSV = (
  data: any[],
  columns: ExcelColumn[],
  filename: string
) => {
  if (!data || data.length === 0) {
    console.warn('No data to export')
    return
  }

  // CSV头部
  let csv = columns.map(col => col.title).join(',') + '\n'

  // CSV数据
  data.forEach(row => {
    const values = columns.map(col => {
      let value = col.formatter ? col.formatter(row) : row[col.key]
      // 处理包含逗号的字段
      if (typeof value === 'string' && (value.includes(',') || value.includes('"'))) {
        value = `"${value.replace(/"/g, '""')}"`
      }
      return value ?? ''
    })
    csv += values.join(',') + '\n'
  })

  // 添加BOM以支持中文
  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csv], { type: 'text/csv;charset=utf-8;' })

  // 下载
  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  const fullFilename = `${filename}_${timestamp}.csv`
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = fullFilename
  link.click()
  URL.revokeObjectURL(link.href)
}

/**
 * 导出为JSON
 */
export const exportToJSON = (
  data: any[],
  filename: string
) => {
  if (!data || data.length === 0) {
    console.warn('No data to export')
    return
  }

  const json = JSON.stringify(data, null, 2)
  const blob = new Blob([json], { type: 'application/json;charset=utf-8;' })

  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, '')
  const fullFilename = `${filename}_${timestamp}.json`
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = fullFilename
  link.click()
  URL.revokeObjectURL(link.href)
}

/**
 * 导出Excel（别名，供ExportButton组件使用）
 */
export const exportExcel = (params: {
  data: any[]
  columns: ExcelColumn[]
  filename: string
}) => {
  exportToExcel(params.data, params.columns, params.filename)
}
