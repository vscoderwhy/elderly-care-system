# 百万级数据生成器使用指南

## 功能特点

✅ **高性能生成** - 使用并发插入，快速生成百万级测试数据
✅ **真实数据** - 随机生成符合养老院业务场景的数据
✅ **灵活配置** - 支持单独生成各类数据
✅ **进度显示** - 实时显示插入进度和速度

## 数据结构

### 设施数据
- **楼栋**: 5栋 (A-E栋)
- **楼层**: 每栋6层
- **房间**: 每层20个房间
- **床位**: 每房间4个床位
- **总床位**: 5 × 6 × 20 × 4 = **2,400个床位**

### 老人数据
- **数量**: 100万条
- **字段**: 姓名、性别、出生日期、身份证、电话、地址、紧急联系人、入院日期、床位、护理等级、状态
- **年龄分布**: 60-100岁
- **性别比例**: 约50:50

### 护理记录
- **数量**: 每个老人50条记录 = **5000万条**
- **时间范围**: 过去1年内
- **记录类型**: 健康监测、用药、日常护理、康复训练、饮食

### 护工账号
- **数量**: 10个
- **职位**: 护士长、主管护工、护工、护理员
- **密码**: 默认123456

## 使用方法

### 方式1: 交互式脚本 (推荐)

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system
bash generate-million-data.sh
```

然后选择要生成的数据类型。

### 方式2: 直接运行

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

# 生成完整数据 (100万老人 + 护理记录)
go run cmd/generator/main.go all

# 仅生成100万老人
go run cmd/generator/main.go elderly

# 仅生成护理记录
go run cmd/generator/main.go records

# 仅生成设施数据
go run cmd/generator/main.go facility

# 仅生成护工账号
go run cmd/generator/main.go staff
```

## 预期性能

### 完整数据生成 (all)
- **老人数据**: 100万条
- **护理记录**: 5000万条
- **预计时间**: 5-15分钟
- **插入速度**: 约10,000-50,000条/秒

### 仅老人数据 (elderly)
- **数据量**: 100万条
- **预计时间**: 3-8分钟

## 性能优化

### 数据库配置优化

在生成大量数据前，建议优化PostgreSQL配置：

```bash
# 编辑 postgresql.conf
docker exec -it elderly-care-db bash
vi /var/lib/postgresql/data/postgresql.conf
```

添加以下配置：

```ini
# 优化批量插入
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 128MB
checkpoint_segments = 16
wal_buffers = 16MB

# 禁用同步日志 (仅测试环境)
synchronous_commit = off
fsync = off
```

重启数据库：

```bash
docker restart elderly-care-db
```

### 批量大小调整

在 `generator/main.go` 中调整：

```go
batchSize := 10000  // 默认每批1万条
```

- 增加 batch size 可以提高速度，但会占用更多内存
- 减少 batch size 可以降低内存占用，但速度会变慢

## 数据验证

### 查看数据统计

```bash
# 老人总数
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) as elderly_count FROM elderly;'

# 护理记录数
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) as record_count FROM care_records;'

# 床位使用情况
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT status, COUNT(*) FROM beds GROUP BY status;'

# 年龄分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT EXTRACT(YEAR FROM age(birth_date)) as age, COUNT(*) FROM elderly GROUP BY EXTRACT(YEAR FROM age(birth_date)) ORDER BY age;'

# 护理等级分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT care_level, COUNT(*) FROM elderly GROUP BY care_level;'
```

### 测试查询性能

```bash
# 分页查询
time docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  "SELECT * FROM elderly LIMIT 20 OFFSET 100000;"

# 条件查询
time docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  "SELECT * FROM elderly WHERE care_level = 5 LIMIT 100;"

# 关联查询
time docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  "SELECT e.name, r.room_number, b.bed_number FROM elderly e JOIN beds b ON e.bed_id = b.id JOIN rooms r ON b.room_id = r.id LIMIT 100;"
```

## 常见问题

### 1. 内存不足

**症状**: `Out of memory` 错误

**解决方法**:
- 减少 batch size
- 降低并发数
- 关闭其他程序

### 2. 速度太慢

**解决方法**:
- 增加数据库 shared_buffers
- 关闭 fsync 和 synchronous_commit
- 增加 batch size
- 使用 SSD 硬盘

### 3. 数据库连接失败

**症状**: `connection refused`

**解决方法**:
```bash
# 检查数据库是否运行
docker ps | grep elderly-care-db

# 启动数据库
docker start elderly-care-db
```

## 清除数据

如果需要重新生成：

```bash
# 删除所有数据
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  "TRUNCATE care_records CASCADE; \
   TRUNCATE elderly CASCADE; \
   TRUNCATE beds CASCADE; \
   TRUNCATE rooms CASCADE; \
   TRUNCATE floors CASCADE; \
   TRUNCATE buildings CASCADE; \
   TRUNCATE users WHERE phone LIKE '138001380%';"

# 或者重置序列
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  "TRUNCATE care_records, elderly, beds, rooms, floors, buildings RESTART IDENTITY CASCADE;"
```

## 生产环境使用

⚠️ **警告**: 此工具仅用于测试和开发环境，不要在生产环境使用！

生产环境数据生成应考虑：
- 使用脱敏的真实数据
- 遵守数据隐私法规
- 添加必要的索引和约束
- 考虑数据分片和归档策略

## 技术细节

### 并发控制
- 使用 goroutine 并发插入
- 使用 sync.WaitGroup 等待完成
- 使用 atomic.Int64 统计进度

### 事务策略
- 每个批次一个事务
- 失败自动回滚
- 批量提交减少 I/O

### 随机数据生成
- 姓名库: 20个常见姓氏 + 40个名字
- 地址库: 8个省份 + 8个区县
- 年龄: 60-100岁均匀分布
- 日期: 过去10年随机

## 扩展

如需自定义数据生成规则，编辑 `backend/cmd/generator/main.go`:

```go
// 修改老人数量
elderlyCount := 1000000  // 改为需要的数量

// 修改护理记录数量
recordPerElderly := 50  // 改为每个老人的记录数

// 添加新字段
// 在 INSERT 语句中添加新字段
```

## 联系方式

如有问题，请联系开发团队。
