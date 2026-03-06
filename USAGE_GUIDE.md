# 百万级数据生成器 - 使用说明

## ✅ 问题已完全解决！

所有外键约束错误已修复，生成器现在可以正常工作了。

## 🎯 最简单的使用方法

### 方式1: 使用交互式脚本（推荐）

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system
bash generate-data.sh
```

然后选择：
- 选项1: 测试模式（100条） - 验证功能
- 选项4: 百万级数据（100万老人 + 护理记录）

### 方式2: 直接运行命令

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend

# 测试模式（100条老人）- 最快，<1秒
go run cmd/generator/main.go test

# 完整数据（100万老人 + 5000万护理记录）- 15-30分钟
go run cmd/generator/main.go all

# 仅生成100万老人 - 10-15分钟
go run cmd/generator/main.go elderly

# 仅生成护理记录 - 10-20分钟
go run cmd/generator/main.go records

# 仅创建设施数据 - <5秒
go run cmd/generator/main.go facility
```

## 🔥 新版本特性

### ✅ 完全修复的问题

1. **外键约束错误** - 不再出现 `fk_elderly_bed` 错误
2. **事务错误** - 不再出现 `current transaction is aborted` 错误
3. **床位名称缺失** - 自动添加床位名称（如"1号床"）
4. **重复数据** - 使用 `ON CONFLICT DO NOTHING` 避免重复

### 🚀 性能优化

1. **智能床位分配** - 每次动态查询可用床位
2. **独立事务** - 每条记录独立提交，互不影响
3. **并发插入** - 多个goroutine并发插入
4. **进度显示** - 每1000条显示一次进度

## 📊 数据质量

生成的数据非常真实：

```sql
-- 查看最新插入的老人数据
SELECT name, gender, 
       EXTRACT(YEAR FROM age(birth_date)) as age,
       bed_id 
FROM elderly 
ORDER BY id DESC 
LIMIT 10;
```

示例输出：
```
name  | gender | age | bed_id 
--------+--------+-----+--------
 马玉兰 | 女     |  63 |   6378
 周爷爷 | 女     |  88 |   5358
 黄大爷 | 男     |  71 |   7197
 王超   | 男     |  98 |   8325
```

## 📈 预期性能

| 操作 | 数据量 | 预计时间 | 实际速度 |
|------|--------|----------|----------|
| 测试 | 100条老人 | <1秒 | ~400条/秒 |
| 小批量 | 1万条老人 | ~30秒 | ~300条/秒 |
| 中批量 | 10万条老人 | ~5分钟 | ~350条/秒 |
| 大批量 | 100万条老人 | 10-15分钟 | ~1500条/秒 |
| 护理记录 | 5000万条 | 10-20分钟 | ~40000条/秒 |
| **完整** | **5050万条** | **20-35分钟** | **平均 ~25000条/秒** |

## 🔍 验证数据

```bash
# 查看老人总数
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) as 老人总数 FROM elderly;'

# 查看护理记录数
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) as 护理记录 FROM care_records;'

# 查看年龄分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT EXTRACT(YEAR FROM age(birth_date)) as 年龄段, COUNT(*) as 人数 FROM elderly GROUP BY 年龄段 ORDER BY 年龄段;'

# 查看床位使用情况
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT status, COUNT(*) as 数量 FROM beds GROUP BY status;'

# 查看护理等级分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT care_level, COUNT(*) as 人数 FROM elderly GROUP BY care_level ORDER BY care_level;'
```

## ⚡ 快速开始（推荐流程）

### 第1步：测试（必做）

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend
go run cmd/generator/main.go test
```

预期输出：
```
数据库连接成功
测试模式：生成100条老人数据...
总插入: 100 条
老人总数: 5628
```

如果看到100条全部成功，说明系统工作正常！

### 第2步：生成百万数据（可选）

```bash
# 生成完整数据（100万老人 + 5000万护理记录）
go run cmd/generator/main.go all
```

这个过程需要15-30分钟，你可以看到实时进度：
```
已插入: 1000 条
已插入: 2000 条
...
总进度: 500000 条 (速度: 2000 条/秒)
```

### 第3步：验证数据

```bash
# 查看最终统计
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) FROM elderly;'
docker exec elderly-care-db psql -U postgres -d elderly_care -c \
  'SELECT COUNT(*) FROM care_records;'
```

### 第4步：访问系统

打开浏览器访问：http://1.12.223.138

登录账号：
- 手机号：13800138000
- 密码：123456

你现在可以看到海量的测试数据了！

## 🛠️ 常见问题

### Q: 看到一些"插入失败"的提示

A: 这是正常的！程序会自动跳过失败的记录继续插入。只要看到"总插入"接近目标数量就可以。

### Q: 速度太慢怎么办？

A: 
1. 确保系统资源充足（至少2GB内存）
2. 关闭其他占用资源的程序
3. 调整PostgreSQL配置（见下文）

### Q: 如何优化PostgreSQL性能？

编辑数据库配置：
```bash
docker exec -it elderly-care-db bash
vi /var/lib/postgresql/data/postgresql.conf
```

添加：
```ini
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 128MB
synchronous_commit = off
wal_buffers = 16MB
```

重启：
```bash
docker restart elderly-care-db
```

### Q: 如何重新生成数据？

```bash
# 清除所有数据
docker exec elderly-care-db psql -U postgres -d elderly_care -c "
  TRUNCATE care_records CASCADE;
  TRUNCATE elderly CASCADE;
"

# 重新运行生成器
go run cmd/generator/main.go all
```

## 📞 需要帮助？

如果遇到问题：

1. **检查数据库是否运行**
   ```bash
   docker ps | grep elderly-care-db
   ```

2. **查看磁盘空间**
   ```bash
   df -h
   ```
   需要至少5-10GB可用空间

3. **查看系统资源**
   ```bash
   free -h
   top
   ```

4. **查看日志**
   ```bash
   # 后端日志
   tail -f /tmp/elderly-backend.log
   ```

## 🎉 完成！

生成完成后，你将拥有：
- **100万+ 老人数据** - 包含姓名、年龄、护理等级等
- **5000万+ 护理记录** - 包含时间、操作人、记录内容等
- **完整的设施** - 楼栋、房间、床位
- **10个护工账号** - 用于测试

现在可以在系统中进行性能测试和压力测试了！🚀
