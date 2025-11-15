#!/usr/bin/env python3
"""
数据库emoji格式迁移脚本
将旧的emoji格式批量转换为新格式
"""

import json
import re
import os

class EmojiDBMigrator:
    def __init__(self, mapping_file):
        """初始化迁移器"""
        with open(mapping_file, 'r', encoding='utf-8') as f:
            self.emoji_mapping = json.load(f)
        
        print(f"加载emoji映射: {len(self.emoji_mapping)} 个映射关系")
    
    def migrate_text(self, text):
        """迁移单个文本内容"""
        if not text:
            return text
        
        result = text
        changes = []
        
        # 1. 转换 :emoji:s123: -> :emoji:e456:
        def replace_colon_format(match):
            old_key = f"s{match.group(1)}"
            new_key = self.emoji_mapping.get(old_key)
            if new_key:
                changes.append(f":emoji:{old_key}: -> :emoji:{new_key}:")
                return f":emoji:{new_key}:"
            return match.group(0)
        
        result = re.sub(r':emoji:s(\d+):', replace_colon_format, result)
        
        # 2. 转换 ![](/emoji/s123.png) -> :emoji:e456:
        def replace_markdown_absolute(match):
            old_key = f"s{match.group(1)}"
            new_key = self.emoji_mapping.get(old_key)
            if new_key:
                changes.append(f"![](/emoji/{old_key}.png) -> :emoji:{new_key}:")
                return f":emoji:{new_key}:"
            return match.group(0)
        
        result = re.sub(r'!\[\]\(/emoji/s(\d+)\.png\)', replace_markdown_absolute, result)
        
        # 3. 转换 ![](emoji/s123.png) -> :emoji:e456:
        def replace_markdown_relative(match):
            old_key = f"s{match.group(1)}"
            new_key = self.emoji_mapping.get(old_key)
            if new_key:
                changes.append(f"![](emoji/{old_key}.png) -> :emoji:{new_key}:")
                return f":emoji:{new_key}:"
            return match.group(0)
        
        result = re.sub(r'!\[\]\(emoji/s(\d+)\.png\)', replace_markdown_relative, result)
        
        return result, changes
    
    def generate_sql_script(self, output_file):
        """生成SQL迁移脚本"""
        print("=== 生成SQL迁移脚本 ===")
        
        sql_statements = []
        
        # 添加脚本头部
        sql_statements.append("-- Emoji格式迁移脚本")
        sql_statements.append("-- 将旧的emoji格式转换为新的雪碧图格式")
        sql_statements.append("-- 执行前请备份数据库！")
        sql_statements.append("")
        sql_statements.append("START TRANSACTION;")
        sql_statements.append("")
        
        # 为每个映射生成UPDATE语句（只针对comments表）
        for old_key, new_key in self.emoji_mapping.items():
            # 更新冒号格式
            sql_statements.append(f"-- 更新 :emoji:{old_key}: -> :emoji:{new_key}:")
            sql_statements.append(f"UPDATE comments SET content = REPLACE(content, ':emoji:{old_key}:', ':emoji:{new_key}:');")
            sql_statements.append("")
            
            # 更新Markdown格式
            sql_statements.append(f"-- 更新 ![](/emoji/{old_key}.png) -> :emoji:{new_key}:")
            sql_statements.append(f"UPDATE comments SET content = REPLACE(content, '![](/emoji/{old_key}.png)', ':emoji:{new_key}:');")
            sql_statements.append("")
            
            # 更新相对路径格式
            sql_statements.append(f"-- 更新 ![](emoji/{old_key}.png) -> :emoji:{new_key}:")
            sql_statements.append(f"UPDATE comments SET content = REPLACE(content, '![](emoji/{old_key}.png)', ':emoji:{new_key}:');")
            sql_statements.append("")
        
        # 添加脚本尾部
        sql_statements.append("-- 提交事务")
        sql_statements.append("COMMIT;")
        sql_statements.append("")
        sql_statements.append("-- 验证迁移结果")
        sql_statements.append("SELECT COUNT(*) as old_format_count FROM comments WHERE content LIKE '%:emoji:s%' OR content LIKE '%![](emoji/s%' OR content LIKE '%![](/emoji/s%';")
        sql_statements.append("SELECT COUNT(*) as new_format_count FROM comments WHERE content LIKE '%:emoji:e%';")
        
        # 写入文件
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write('\n'.join(sql_statements))
        
        print(f"✅ SQL脚本已生成: {output_file}")
        print(f"   包含 {len(self.emoji_mapping)} 个emoji的迁移语句")
        return output_file
    
    def generate_rollback_script(self, output_file):
        """生成回滚脚本"""
        print("=== 生成回滚脚本 ===")
        
        sql_statements = []
        
        sql_statements.append("-- Emoji格式回滚脚本")
        sql_statements.append("-- 将新格式回滚为旧格式")
        sql_statements.append("-- 仅在迁移出现问题时使用！")
        sql_statements.append("")
        sql_statements.append("START TRANSACTION;")
        sql_statements.append("")
        
        # 生成回滚语句（只针对comments表）
        for old_key, new_key in self.emoji_mapping.items():
            sql_statements.append(f"-- 回滚 :emoji:{new_key}: -> :emoji:{old_key}:")
            sql_statements.append(f"UPDATE comments SET content = REPLACE(content, ':emoji:{new_key}:', ':emoji:{old_key}:');")
            sql_statements.append("")
        
        sql_statements.append("COMMIT;")
        
        with open(output_file, 'w', encoding='utf-8') as f:
            f.write('\n'.join(sql_statements))
        
        print(f"✅ 回滚脚本已生成: {output_file}")
        return output_file
    
    def test_migration(self, test_texts):
        """测试迁移效果"""
        print("=== 测试迁移效果 ===")
        
        for i, text in enumerate(test_texts, 1):
            print(f"\n测试 {i}:")
            print(f"原文: {text}")
            
            migrated, changes = self.migrate_text(text)
            print(f"迁移后: {migrated}")
            
            if changes:
                print("变更:")
                for change in changes:
                    print(f"  - {change}")
            else:
                print("无变更")

def main():
    # 配置文件路径
    mapping_file = "/media/jiang/hsk/practice/Golang/go/goCode/prac09/go_blog/scripts/emoji_output/emoji_frontend_mapping.json"
    output_dir = "/media/jiang/hsk/practice/Golang/go/goCode/prac09/go_blog/scripts/emoji_output"
    
    # 创建迁移器
    migrator = EmojiDBMigrator(mapping_file)
    
    # 生成迁移脚本
    migrate_sql = os.path.join(output_dir, "migrate_emoji.sql")
    rollback_sql = os.path.join(output_dir, "rollback_emoji.sql")
    
    migrator.generate_sql_script(migrate_sql)
    migrator.generate_rollback_script(rollback_sql)
    
    # 测试迁移
    test_texts = [
        "这是一个测试 :emoji:s1: 文本",
        "包含图片 ![](/emoji/s123.png) 的内容",
        "相对路径 ![](emoji/s45.png) 测试",
        "混合格式 :emoji:s1: 和 ![](emoji/s2.png) 一起",
        "已经是新格式 :emoji:e1: 不应该改变"
    ]
    
    migrator.test_migration(test_texts)
    
    print("\n🎉 迁移脚本生成完成！")
    print("📋 执行步骤:")
    print("  1. 备份数据库")
    print(f"  2. 执行: {migrate_sql}")
    print("  3. 验证迁移结果")
    print(f"  4. 如有问题，执行回滚: {rollback_sql}")

if __name__ == "__main__":
    main()
