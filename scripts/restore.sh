#!/bin/bash

# 恢复脚本

set -e

# 配置
BACKUP_DIR="./backups"
DB_CONTAINER="elderly-care-postgres"
DB_NAME="elderly_care_db"
DB_USER="elderly_care"

# 显示可用的备份
echo "=========================================="
echo "  可用的数据库备份："
echo "=========================================="
echo ""

ls -lht "$BACKUP_DIR"/db_backup_*.sql 2>/dev/null | awk '{print NR ". "$9}' | head -10

if [ $? -ne 0 ]; then
    echo "❌ 没有找到备份文件"
    exit 1
fi

echo ""
read -p "请选择要恢复的备份编号 (1-10): " choice

# 获取选中的备份文件
backup_file=$(ls -t "$BACKUP_DIR"/db_backup_*.sql | sed -n "${choice}p")

if [ -z "$backup_file" ]; then
    echo "❌ 无效的选择"
    exit 1
fi

echo ""
echo "将恢复备份: $backup_file"
echo "⚠️  此操作将覆盖当前数据库！"
echo ""
read -p "确定要继续吗？(yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    echo "❌ 操作已取消"
    exit 0
fi

# 恢复数据库
echo ""
echo "📥 恢复数据库..."
docker exec -i "$DB_CONTAINER" psql -U "$DB_USER" "$DB_NAME" < "$backup_file"

if [ $? -eq 0 ]; then
    echo "✅ 数据库恢复成功"
else
    echo "❌ 数据库恢复失败"
    exit 1
fi

# 询问是否恢复文件
echo ""
read -p "是否同时恢复上传文件？(yes/no): " restore_files

if [ "$restore_files" = "yes" ]; then
    echo ""
    echo "=========================================="
    echo "  可用的文件备份："
    echo "=========================================="
    echo ""

    ls -lht "$BACKUP_DIR"/uploads_backup_*.tar.gz 2>/dev/null | awk '{print NR ". "$9}' | head -10

    echo ""
    read -p "请选择要恢复的备份编号 (1-10): " file_choice

    file_backup=$(ls -t "$BACKUP_DIR"/uploads_backup_*.tar.gz | sed -n "${file_choice}p")

    if [ -n "$file_backup" ]; then
        echo ""
        echo "📥 恢复上传文件..."
        rm -rf uploads/*
        tar -xzf "$file_backup" -C ./
        echo "✅ 文件恢复成功"
    fi
fi

echo ""
echo "=========================================="
echo "  恢复完成"
echo "=========================================="
echo ""
