# 养老院管理系统

一个功能完整的养老院管理系统，包含管理后台（BI 数据看板、暗黑模式）和 uniapp 移动端（H5 + 小程序）。

## 项目结构

```
elderly-care-system/
├── admin-frontend/          # 管理后台（Vue 3 + Element Plus + ECharts）
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   │   ├── Dashboard/  # 仪表盘（概览、指挥中心）
│   │   │   ├── Statistics/ # 数据分析（高级看板、报表）
│   │   │   ├── Monitor/    # 实时监控
│   │   │   └── Layout.vue  # 布局组件
│   │   ├── components/     # 公共组件
│   │   │   ├── Charts/     # 图表组件（雷达图、桑基图、词云等）
│   │   │   └── Dashboard/  # 看板组件（卡片、仪表盘、时间轴）
│   │   ├── composables/    # 组合式函数
│   │   │   ├── useDarkMode.ts      # 暗黑模式
│   │   │   ├── useECharts.ts       # ECharts 封装
│   │   │   └── useWebSocket.ts     # WebSocket 封装
│   │   ├── styles/         # 样式文件
│   │   │   └── theme.css   # 主题变量（浅色/深色）
│   │   ├── router/         # 路由配置
│   │   ├── api/            # API 接口
│   │   └── main.ts         # 入口文件
│   ├── package.json
│   └── vite.config.ts
│
├── uniapp-elderly-care/    # 移动端（uniapp + Vue 3）
│   ├── pages/              # 页面
│   │   ├── index/          # 首页
│   │   ├── elderly/        # 老人信息
│   │   ├── care/           # 护理记录
│   │   ├── health/         # 健康数据
│   │   ├── bills/          # 费用账单
│   │   ├── visits/         # 探视预约
│   │   └── profile/        # 个人中心
│   ├── store/              # Pinia 状态管理
│   ├── api/                # API 封装
│   ├── utils/              # 工具函数
│   ├── components/         # 公共组件
│   ├── manifest.json       # 应用配置
│   ├── pages.json          # 页面配置
│   └── App.vue
│
└── backend/                # 后端（Go + PostgreSQL）
    └── (待实现)
```

## 功能特性

### 管理后台

#### 1. 企业级 BI 数据看板
- **概览仪表盘**：关键指标统计卡片、趋势图表、KPI 仪表盘
- **高级数据分析**：
  - 雷达图：护理质量多维度评估
  - 桑基图：老人流转路径分析
  - 关系图：老人-家属-护工关系网络
  - 词云图：护理记录关键词分析
  - 漏斗图：费用支付转化率
- **实时监控中心**：WebSocket 实时数据推送、护理任务进度、健康数据异常告警、设备状态监控

#### 2. 暗黑模式
- CSS 变量实现主题切换
- 支持自动跟随系统主题
- 手动切换（浅色/深色/自动）
- 持久化用户偏好

#### 3. 响应式设计
- 移动端/平板/桌面自适应
- 侧边栏自动折叠
- 图表自适应大小

### 移动端（uniapp）

#### 家属端功能
- 老人信息查看（基本信息、照片、护理等级）
- 护理记录浏览（图文记录、服务评价）
- 健康数据查看（血压、血糖、体温趋势图）
- 费用账单（账单列表、在线支付、发票下载）
- 探视预约（在线预约、预约历史、取消预约）
- 消息通知（护理提醒、费用提醒、公告通知）

#### 护工端功能
- 今日任务（待办护理、已完成、逾期）
- 快速记录（拍照上传、语音输入）
- 服务请求（家属请求处理、进度更新）
- 考勤打卡（上班打卡、下班打卡、请假申请）
- 用药提醒（用药列表、服药记录）
- 排班查询（本周排班、排班日历）

## 技术栈

### 管理后台
- **框架**：Vue 3 + TypeScript + Vite
- **UI 库**：Element Plus
- **图表**：ECharts 5.6 + ECharts-GL + Vue-ECharts
- **状态管理**：Pinia
- **路由**：Vue Router 4
- **HTTP 客户端**：Axios
- **实时通信**：WebSocket

### 移动端
- **框架**：uniapp（Vue 3 + TypeScript）
- **状态管理**：Pinia
- **UI 组件**：uni-ui
- **平台**：H5 + 微信小程序

### 后端（待实现）
- **语言**：Go
- **数据库**：PostgreSQL
- **ORM**：GORM
- **API**：RESTful API

## 快速开始

### 管理后台

```bash
cd admin-frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

### 移动端

```bash
cd uniapp-elderly-care

# 安装依赖
npm install

# H5 开发
npm run dev:h5

# 微信小程序开发
npm run dev:mp-weixin

# 构建
npm run build:h5
npm run build:mp-weixin
```

## 依赖安装

### 管理后台主要依赖

```bash
npm install element-plus @element-plus/icons-vue
npm install echarts echarts-gl vue-echarts
npm install pinia vue-router
npm install axios
```

### uniapp 主要依赖

```bash
npm install pinia
npm install @dcloudio/uni-app
npm install @dcloudio/uni-ui
```

## 开发规范

### 代码规范
- 使用 TypeScript 进行类型检查
- 遵循 ESLint 代码规范
- 使用 Vue 3 Composition API

### 命名规范
- 组件：PascalCase（如 `StatCard.vue`）
- 文件夹：PascalCase（如 `Dashboard/`）
- 函数/变量：camelCase（如 `getUserInfo`）
- 常量：UPPER_SNAKE_CASE（如 `API_BASE_URL`）

### 提交规范
- feat: 新功能
- fix: 修复 bug
- docs: 文档更新
- style: 代码格式调整
- refactor: 重构
- test: 测试相关
- chore: 构建/工具变动

## 浏览器支持

管理后台支持以下浏览器：
- Chrome >= 90
- Firefox >= 88
- Safari >= 14
- Edge >= 90

## 移动端平台

- H5：现代浏览器
- 微信小程序：基础库 >= 2.0

## 待完成功能

- [ ] 后端 API 开发
- [ ] 文件上传组件
- [ ] 数据导出功能（Excel）
- [ ] 高级搜索和批量操作
- [ ] 审核流程配置
- [ ] 操作日志
- [ ] 数据备份恢复

## 贡献指南

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

## 联系方式

如有问题，请联系开发团队。
