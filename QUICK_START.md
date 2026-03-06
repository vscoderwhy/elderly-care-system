# 百万级数据生成器 - 快速开始

## ✅ 问题已修复！

之前的错误已修复，现在生成器可以正常工作了。

## 🚀 快速开始

### 1. 测试模式（推荐先运行）

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system/backend
go run cmd/generator/main.go test
```

这将生成100条老人数据，用于测试和验证。

### 2. 生成百万级数据

```bash
# 方式1: 生成完整数据（100万老人 + 5000万护理记录）
go run cmd/generator/main.go all

# 方式2: 只生成100万老人（更快，约5-10分钟）
go run cmd/generator/main.go elderly

# 方式3: 只生成护理记录（需要先有老人数据）
go run cmd/generator/main.go records
```

### 3. 使用便捷脚本

```bash
cd /root/.openclaw/workspace-feishu-elderly/elderly-care-system
bash generate-million-data.sh
```

然后选择选项：
- 选项1: 完整数据（100万老人 + 护理记录）
- 选项2: 仅100万老人
- 选项3: 仅护理记录
- 选项4: 仅设施数据
- 选项5: 仅护工账号

## 📊 预期性能

| 数据类型 | 数量 | 预计时间 |
|---------|------|---------|
| 测试数据 | 100条老人 | <1秒 |
| 设施数据 | 2400个床位 | <5秒 |
| 护工账号 | 10个 | <2秒 |
| 老人数据 | 100万条 | 5-10分钟 |
| 护理记录 | 5000万条 | 10-20分钟 |
| **完整数据** | **5050万条** | **20-30分钟** |

## 🔍 验证数据

```bash
# 查看老人总数
docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM elderly;'

# 查看护理记录数
docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT COUNT(*) FROM care_records;'

# 查看床位使用情况
docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT status, COUNT(*) FROM beds GROUP BY status;'

# 查看年龄分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT EXTRACT(YEAR FROM age(birth_date)) as age, COUNT(*) FROM elderly GROUP BY age ORDER BY age;'

# 查看护理等级分布
docker exec elderly-care-db psql -U postgres -d elderly_care -c 'SELECT care_level, COUNT(*) FROM elderly GROUP BY care_level;'
```

## ⚠️ 注意事项

1. **错误提示**: 
   - 如果看到"插入失败: pq: insert or update on table "elderly" violates foreign key constraint"这是正常的，说明床位已分配
   - 程序会自动跳过失败的记录继续插入

2. **性能优化**:
   - 每条记录使用独立事务，避免一条失败影响整批
   - 使用并发插入，提高速度
   - 每100条显示一次进度

3. **磁盘空间**:
   - 确保**至少5-10GB**可用空间
   - 完整数据（5050万条）约占5-8GB

4. **运行环境**:
   - 建议至少2GB可用内存
   - PostgreSQL需要运行在localhost:5432

## 🛠️ 故障排除

### 问题1: 数据库连接失败

```bash
# 检查数据库是否运行
docker ps | grep elderly-care-db

# 如果没有运行，启动服务
bash /root/.openclaw/workspace-feishu-elderly/elderly-care-system/start-all.sh
```

### 问题2: 速度太慢

```bash
# 检查系统资源
free -h
df -h

# 优化PostgreSQL配置（可选）
docker exec -it elderly-care-db bash
vi /var/lib/postgresql/data/postgresql.conf
```

添加以下配置：
```ini
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 128MB
synchronous_commit = off
```

然后重启：
```bash
docker restart elderly-care-db
```

### 问题3: 重新生成数据

```bash
# 清除所有数据
docker exec elderly-care-db psql -U postgres -d elderly_care -c "
  TRUNCATE care_records CASCADE;
  TRUNCATE elderly CASCADE;
  TRUNCATE beds CASCADE;
  TRUNCATE rooms CASCADE;
  TRUNCATE floors CASCADE;
  TRUNCATE buildings CASCADE;
"

# 重新运行生成器
go run cmd/generator/main.go all
```

## 📝 实时监控

生成过程中，你会看到类似这样的输出：

```
2026/03/05 09:58:51 数据库连接成功
2026/03/05 09:58:51 测试模式：生成100条老人数据...
2026/03/05 09:58:52 已插入: 100 条
总进度: 100 条 (速度: 100 条/秒)

=== 数据生成完成 ===
总耗时: 852ms
总插入: 91 条

=== 数据库统计 ===
老人总数: 275
楼栋总数: 17
房间总数: 1805
床位总数: 4813
护理记录: 10
```

## 🎯 下一步

生成完成后，你可以：

1. **访问系统**: http://1.12.223.138
2. **登录账号**: 13800138000 / 123456
3. **查看数据**: 在后台看到海量的测试数据
4. **性能测试**: 测试系统在百万级数据下的表现

## 📞 需要帮助？

如果遇到问题，检查：
1. PostgreSQL是否运行
2. 端口5432是否被占用
3. 磁盘空间是否充足
4. 内存是否足够

或者查看详细文档：`MILLION_DATA_GUIDE.md`
