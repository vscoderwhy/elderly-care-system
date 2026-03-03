#!/bin/bash

# 备份脚本

set -e

# 配置
BACKUP_DIR="./backups"
DATE=$(date +%Y%m%d_%H%M%S)
DB_CONTAINER="elderly-care-postgres"
DB_NAME="elderly_care_db"
DB_USER="elderly_care"

# 创建备份目录
mkdir -p "$BACKUP_DIR"

echo "=========================================="
echo "  开始备份"
echo "=========================================="
echo ""

# 备份数据库
echo "📦 备份数据库..."
docker exec "$DB_CONTAINER" pg_dump -U "$DB_USER" "$DB_NAME" > "$BACKUP_DIR/db_backup_$DATE.sql"

if [ $? -eq 0 ]; then
    echo "✅ 数据库备份成功: db_backup_$DATE.sql"
else
    echo "❌ 数据库备份失败"
    exit 1
fi

# 备份上传文件
echo "📦 备份上传文件..."
tar -czf "$BACKUP_DIR/uploads_backup_$DATE.tar.gz" uploads/

if [ $? -eq 0 ]; then
    echo "✅ 文件备份成功: uploads_backup_$DATE.tar.gz"
else
    echo "❌ 文件备份失败"
    exit 1
fi

# 清理 30 天前的备份
echo ""
echo "🧹 清理旧备份..."
find "$BACKUP_DIR" -name "db_backup_*.sql" -mtime +30 -delete
find "$BACKUP_DIR" -name "uploads_backup_*.tar.gz" -mtime +30 -delete

echo ""
echo "=========================================="
echo "  备份完成"
echo "=========================================="
echo ""
echo "备份文件位置: $BACKUP_DIR"
echo ""
