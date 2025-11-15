#!/usr/bin/env python3
"""
Emoji优化脚本 V2：支持增量更新和版本管理
"""

import os
import json
import glob
import argparse
from PIL import Image
import math
from emoji_version_manager import EmojiVersionManager

class EmojiOptimizerV2:
    def __init__(self, emoji_dir, output_dir):
        self.emoji_dir = emoji_dir
        self.output_dir = output_dir
        self.target_size = 64  # 统一尺寸64x64
        self.sprites_per_row = 16  # 每行16个emoji
        self.emojis_per_sprite = 128  # 每个雪碧图128个emoji
        
        self.version_manager = EmojiVersionManager(
            os.path.join(output_dir, "emoji_version.json")
        )
        
        os.makedirs(output_dir, exist_ok=True)
    
    def scan_existing_emojis(self):
        """扫描现有emoji文件"""
        print("=== 扫描现有emoji文件 ===")
        
        png_files = glob.glob(os.path.join(self.emoji_dir, "s*.png"))
        existing_emojis = []
        
        for file_path in sorted(png_files):
            filename = os.path.basename(file_path)
            if filename.startswith('s') and filename.endswith('.png'):
                try:
                    old_num = int(filename[1:-4])  # s123.png -> 123
                    existing_emojis.append({
                        'old_filename': filename,
                        'old_number': old_num,
                        'file_path': file_path
                    })
                except ValueError:
                    continue
        
        existing_emojis.sort(key=lambda x: x['old_number'])
        print(f"找到 {len(existing_emojis)} 个emoji文件")
        return existing_emojis
    
    def load_existing_mapping(self):
        """加载现有的映射表"""
        mapping_file = os.path.join(self.output_dir, "emoji_frontend_mapping.json")
        if os.path.exists(mapping_file):
            with open(mapping_file, 'r') as f:
                return json.load(f)
        return {}
    
    def generate_incremental_mapping(self, new_emojis):
        """生成增量映射表"""
        print("=== 生成增量映射表 ===")
        
        # 加载现有映射
        existing_mapping = self.load_existing_mapping()
        next_index = self.version_manager.get_next_emoji_index()
        
        new_mapping = {}
        for i, emoji in enumerate(new_emojis):
            old_key = f"s{emoji['old_number']}"
            new_key = f"e{next_index + i}"
            
            if old_key not in existing_mapping:
                new_mapping[old_key] = new_key
                print(f"新增映射: {old_key} -> {new_key}")
        
        return new_mapping
    
    def resize_emoji(self, image_path, target_size=64):
        """调整emoji尺寸为统一大小"""
        with Image.open(image_path) as img:
            # 转换为RGBA模式
            if img.mode != 'RGBA':
                img = img.convert('RGBA')
            
            # 调整尺寸
            img_resized = img.resize((target_size, target_size), Image.Resampling.LANCZOS)
            return img_resized
    
    def create_incremental_sprite(self, new_emojis, new_mapping):
        """创建增量雪碧图"""
        print("=== 创建增量雪碧图 ===")
        
        if not new_emojis:
            print("没有新emoji需要处理")
            return []
        
        sprites_created = []
        sprite_id = self.version_manager.get_next_sprite_id()
        
        # 按128个emoji分组创建雪碧图
        for sprite_index in range(0, len(new_emojis), self.emojis_per_sprite):
            sprite_emojis = new_emojis[sprite_index:sprite_index + self.emojis_per_sprite]
            
            # 计算雪碧图尺寸
            rows = math.ceil(len(sprite_emojis) / self.sprites_per_row)
            sprite_width = self.sprites_per_row * self.target_size
            sprite_height = rows * self.target_size
            
            # 创建空白雪碧图
            sprite_img = Image.new('RGBA', (sprite_width, sprite_height), (0, 0, 0, 0))
            
            # 放置emoji
            for i, emoji in enumerate(sprite_emojis):
                row = i // self.sprites_per_row
                col = i % self.sprites_per_row
                
                x = col * self.target_size
                y = row * self.target_size
                
                # 调整emoji尺寸并粘贴
                emoji_img = self.resize_emoji(emoji['file_path'], self.target_size)
                sprite_img.paste(emoji_img, (x, y), emoji_img)
            
            # 保存雪碧图
            sprite_filename = f"emoji-sprite-{sprite_id}.png"
            sprite_path = os.path.join(self.output_dir, sprite_filename)
            sprite_img.save(sprite_path, 'PNG', optimize=True)
            
            sprites_created.append({
                'id': sprite_id,
                'filename': sprite_filename,
                'path': sprite_path,
                'emoji_count': len(sprite_emojis),
                'size': (sprite_width, sprite_height)
            })
            
            print(f"创建雪碧图: {sprite_filename} ({len(sprite_emojis)} 个emoji)")
            sprite_id += 1
        
        return sprites_created
    
    def generate_incremental_css(self, new_mapping, sprites_info):
        """生成增量CSS"""
        print("=== 生成增量CSS ===")
        
        css_content = f"""/* 增量Emoji CSS - {self.version_manager.get_current_version()} */
/* 生成时间: {self.version_manager.config.created_at} */

"""
        
        # 添加雪碧图背景定义
        for sprite in sprites_info:
            css_content += f""".emoji-sprite-{sprite['id']} {{
  background-image: url('/emoji/{sprite['filename']}');
  background-size: {sprite['size'][0]//2}px {sprite['size'][1]//2}px;
}}

"""
        
        # 生成位置定义
        next_index = self.version_manager.get_next_emoji_index()
        for i, (old_key, new_key) in enumerate(new_mapping.items()):
            emoji_index = next_index + i
            sprite_id = self.version_manager.get_next_sprite_id() + (emoji_index - next_index) // self.emojis_per_sprite
            pos_in_sprite = (emoji_index - next_index) % self.emojis_per_sprite
            
            row = pos_in_sprite // self.sprites_per_row
            col = pos_in_sprite % self.sprites_per_row
            
            x = col * 32  # 32px间距（显示尺寸）
            y = row * 32
            
            css_content += f""".emoji-{new_key} {{
  background-position: -{x}px -{y}px;
}}

"""
        
        # 保存增量CSS
        css_filename = f"emoji-sprites-incremental-{self.version_manager.get_current_version()}.css"
        css_path = os.path.join(self.output_dir, css_filename)
        with open(css_path, 'w', encoding='utf-8') as f:
            f.write(css_content)
        
        print(f"增量CSS已保存: {css_filename}")
        return css_path
    
    def update_mapping_file(self, new_mapping):
        """更新映射文件"""
        print("=== 更新映射文件 ===")
        
        # 加载现有映射
        existing_mapping = self.load_existing_mapping()
        
        # 合并新映射
        updated_mapping = {**existing_mapping, **new_mapping}
        
        # 保存更新后的映射
        mapping_file = os.path.join(self.output_dir, "emoji_frontend_mapping.json")
        with open(mapping_file, 'w', encoding='utf-8') as f:
            json.dump(updated_mapping, f, indent=2, ensure_ascii=False)
        
        print(f"映射文件已更新，总计 {len(updated_mapping)} 个emoji")
        return updated_mapping
    
    def generate_frontend_config(self):
        """生成前端配置文件"""
        print("=== 生成前端配置文件 ===")
        
        config = {
            "version": self.version_manager.get_current_version(),
            "total_emojis": self.version_manager.config.total_emojis,
            "sprites": self.version_manager.get_sprites_info(),
            "updated_at": self.version_manager.config.created_at
        }
        
        config_file = os.path.join(self.output_dir, "emoji_config.json")
        with open(config_file, 'w', encoding='utf-8') as f:
            json.dump(config, f, indent=2, ensure_ascii=False)
        
        print(f"前端配置已生成: emoji_config.json")
        return config_file
    
    def run_incremental_update(self, new_emoji_files=None):
        """运行增量更新"""
        print(f"=== 开始增量更新 ===")
        
        # 扫描所有emoji
        all_emojis = self.scan_existing_emojis()
        
        # 确定新增的emoji
        existing_mapping = self.load_existing_mapping()
        new_emojis = []
        
        if new_emoji_files:
            # 指定了特定的新文件
            for emoji in all_emojis:
                old_key = f"s{emoji['old_number']}"
                if old_key not in existing_mapping and emoji['old_filename'] in new_emoji_files:
                    new_emojis.append(emoji)
        else:
            # 自动检测新文件
            for emoji in all_emojis:
                old_key = f"s{emoji['old_number']}"
                if old_key not in existing_mapping:
                    new_emojis.append(emoji)
        
        if not new_emojis:
            print("没有发现新的emoji文件")
            return
        
        print(f"发现 {len(new_emojis)} 个新emoji:")
        for emoji in new_emojis:
            print(f"  - {emoji['old_filename']}")
        
        # 生成新映射
        new_mapping = self.generate_incremental_mapping(new_emojis)
        
        # 创建增量雪碧图
        sprites_created = self.create_incremental_sprite(new_emojis, new_mapping)
        
        # 更新版本管理器
        for sprite in sprites_created:
            self.version_manager.add_sprite(
                filename=sprite['filename'],
                url="",  # 上传后更新
                emoji_count=sprite['emoji_count'],
                size=sprite['size']
            )
        
        # 生成增量CSS
        css_path = self.generate_incremental_css(new_mapping, sprites_created)
        
        # 更新映射文件
        self.update_mapping_file(new_mapping)
        
        # 创建新版本
        new_version = self.version_manager.create_new_version(
            f"Added {len(new_emojis)} new emojis"
        )
        
        # 保存版本配置
        self.version_manager.save_config()
        
        # 生成前端配置
        self.generate_frontend_config()
        
        print(f"\n🎉 增量更新完成!")
        print(f"版本: {new_version}")
        print(f"新增emoji: {len(new_emojis)} 个")
        print(f"新增雪碧图: {len(sprites_created)} 个")
        print(f"下一步: 上传雪碧图到CDN并更新URL")
        
        return {
            'version': new_version,
            'new_emojis': len(new_emojis),
            'sprites_created': sprites_created,
            'css_path': css_path
        }
    
    def run_initial_setup(self):
        """运行初始设置（基于现有的emoji_optimizer.py结果）"""
        print("=== 初始设置：导入现有数据 ===")
        
        # 检查是否已有现有数据
        mapping_file = os.path.join(self.output_dir, "emoji_frontend_mapping.json")
        if not os.path.exists(mapping_file):
            print("❌ 未找到现有映射文件，请先运行原始的emoji_optimizer.py")
            return
        
        # 加载现有映射
        with open(mapping_file, 'r') as f:
            existing_mapping = json.load(f)
        
        # 检查现有雪碧图
        sprite_files = glob.glob(os.path.join(self.output_dir, "emoji-sprite-*.png"))
        
        print(f"发现 {len(existing_mapping)} 个emoji映射")
        print(f"发现 {len(sprite_files)} 个雪碧图文件")
        
        # 导入到版本管理器
        for i, sprite_file in enumerate(sorted(sprite_files)):
            filename = os.path.basename(sprite_file)
            
            # 获取雪碧图尺寸
            with Image.open(sprite_file) as img:
                size = img.size
            
            # 计算emoji数量
            emoji_count = min(128, len(existing_mapping) - i * 128)
            if emoji_count <= 0:
                break
            
            self.version_manager.add_sprite(
                filename=filename,
                url="",  # 需要手动更新
                emoji_count=emoji_count,
                size=size,
                frozen=True  # 标记为冻结
            )
        
        # 保存初始版本
        self.version_manager.save_config()
        self.generate_frontend_config()
        
        print("✅ 初始设置完成，所有现有雪碧图已标记为冻结")
        self.version_manager.print_status()

def main():
    parser = argparse.ArgumentParser(description='Emoji优化器 V2 - 支持增量更新')
    parser.add_argument('--emoji-dir', default='/media/jiang/hsk/practice/Golang/go/goCode/prac09/go_blog/web-blog/public/emoji',
                       help='Emoji文件目录')
    parser.add_argument('--output-dir', default='emoji_output',
                       help='输出目录')
    parser.add_argument('--mode', choices=['init', 'incremental'], default='incremental',
                       help='运行模式：init=初始设置，incremental=增量更新')
    parser.add_argument('--new-files', nargs='*',
                       help='指定新增的emoji文件名（可选）')
    
    args = parser.parse_args()
    
    optimizer = EmojiOptimizerV2(args.emoji_dir, args.output_dir)
    
    if args.mode == 'init':
        optimizer.run_initial_setup()
    else:
        optimizer.run_incremental_update(args.new_files)

if __name__ == "__main__":
    main()
