#!/usr/bin/env python3
"""
Emoji版本管理系统
支持增量更新、版本控制、回滚等功能
"""

import json
import os
from datetime import datetime
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass, asdict

@dataclass
class SpriteInfo:
    """雪碧图信息"""
    id: int
    filename: str
    url: str
    range_start: int
    range_end: int
    frozen: bool
    created_at: str
    size: Tuple[int, int]  # (width, height)

@dataclass
class EmojiVersionConfig:
    """Emoji版本配置"""
    version: str
    total_emojis: int
    sprites: List[SpriteInfo]
    created_at: str
    description: str

class EmojiVersionManager:
    """Emoji版本管理器"""
    
    def __init__(self, config_path: str = "emoji_output/emoji_version.json"):
        self.config_path = config_path
        self.config: Optional[EmojiVersionConfig] = None
        self.load_config()
    
    def load_config(self) -> None:
        """加载版本配置"""
        if os.path.exists(self.config_path):
            with open(self.config_path, 'r', encoding='utf-8') as f:
                data = json.load(f)
                sprites = [SpriteInfo(**sprite) for sprite in data['sprites']]
                self.config = EmojiVersionConfig(
                    version=data['version'],
                    total_emojis=data['total_emojis'],
                    sprites=sprites,
                    created_at=data['created_at'],
                    description=data['description']
                )
        else:
            # 初始化配置
            self.config = EmojiVersionConfig(
                version="v1.0",
                total_emojis=0,
                sprites=[],
                created_at=datetime.now().isoformat(),
                description="Initial version"
            )
    
    def save_config(self) -> None:
        """保存版本配置"""
        os.makedirs(os.path.dirname(self.config_path), exist_ok=True)
        with open(self.config_path, 'w', encoding='utf-8') as f:
            json.dump(asdict(self.config), f, indent=2, ensure_ascii=False)
    
    def get_current_version(self) -> str:
        """获取当前版本"""
        return self.config.version
    
    def get_next_sprite_id(self) -> int:
        """获取下一个雪碧图ID"""
        if not self.config.sprites:
            return 0
        return max(sprite.id for sprite in self.config.sprites) + 1
    
    def get_next_emoji_index(self) -> int:
        """获取下一个emoji索引"""
        return self.config.total_emojis
    
    def add_sprite(self, filename: str, url: str, emoji_count: int, 
                   size: Tuple[int, int], frozen: bool = False) -> SpriteInfo:
        """添加新的雪碧图"""
        sprite_id = self.get_next_sprite_id()
        range_start = self.get_next_emoji_index()
        range_end = range_start + emoji_count - 1
        
        sprite = SpriteInfo(
            id=sprite_id,
            filename=filename,
            url=url,
            range_start=range_start,
            range_end=range_end,
            frozen=frozen,
            created_at=datetime.now().isoformat(),
            size=size
        )
        
        self.config.sprites.append(sprite)
        self.config.total_emojis += emoji_count
        return sprite
    
    def freeze_sprite(self, sprite_id: int) -> bool:
        """冻结雪碧图（标记为不可修改）"""
        for sprite in self.config.sprites:
            if sprite.id == sprite_id:
                sprite.frozen = True
                return True
        return False
    
    def freeze_all_sprites(self) -> None:
        """冻结所有现有雪碧图"""
        for sprite in self.config.sprites:
            sprite.frozen = True
    
    def create_new_version(self, description: str = "") -> str:
        """创建新版本"""
        # 解析当前版本号
        current_version = self.config.version
        if current_version.startswith('v'):
            version_parts = current_version[1:].split('.')
            major, minor = int(version_parts[0]), int(version_parts[1])
            new_version = f"v{major}.{minor + 1}"
        else:
            new_version = "v1.1"
        
        self.config.version = new_version
        self.config.created_at = datetime.now().isoformat()
        self.config.description = description or f"Updated to {new_version}"
        
        return new_version
    
    def get_sprites_info(self) -> List[Dict]:
        """获取所有雪碧图信息（用于前端）"""
        return [
            {
                "id": sprite.id,
                "filename": sprite.filename,
                "url": sprite.url,
                "range": [sprite.range_start, sprite.range_end],
                "frozen": sprite.frozen,
                "size": sprite.size
            }
            for sprite in self.config.sprites
        ]
    
    def get_unfrozen_sprites(self) -> List[SpriteInfo]:
        """获取未冻结的雪碧图"""
        return [sprite for sprite in self.config.sprites if not sprite.frozen]
    
    def get_emoji_range_for_sprite(self, sprite_id: int) -> Optional[Tuple[int, int]]:
        """获取指定雪碧图的emoji范围"""
        for sprite in self.config.sprites:
            if sprite.id == sprite_id:
                return (sprite.range_start, sprite.range_end)
        return None
    
    def validate_config(self) -> List[str]:
        """验证配置的完整性"""
        errors = []
        
        # 检查emoji范围是否连续
        if self.config.sprites:
            sorted_sprites = sorted(self.config.sprites, key=lambda x: x.range_start)
            expected_start = 0
            
            for sprite in sorted_sprites:
                if sprite.range_start != expected_start:
                    errors.append(f"Sprite {sprite.id} range gap: expected {expected_start}, got {sprite.range_start}")
                expected_start = sprite.range_end + 1
            
            if expected_start != self.config.total_emojis:
                errors.append(f"Total emoji count mismatch: expected {expected_start}, got {self.config.total_emojis}")
        
        return errors
    
    def print_status(self) -> None:
        """打印当前状态"""
        print(f"=== Emoji版本状态 ===")
        print(f"版本: {self.config.version}")
        print(f"总emoji数: {self.config.total_emojis}")
        print(f"雪碧图数量: {len(self.config.sprites)}")
        print(f"创建时间: {self.config.created_at}")
        print(f"描述: {self.config.description}")
        
        print(f"\n=== 雪碧图详情 ===")
        for sprite in self.config.sprites:
            status = "🔒 冻结" if sprite.frozen else "🔓 可修改"
            print(f"ID {sprite.id}: {sprite.filename} ({sprite.range_start}-{sprite.range_end}) {status}")
        
        # 验证配置
        errors = self.validate_config()
        if errors:
            print(f"\n❌ 配置错误:")
            for error in errors:
                print(f"  - {error}")
        else:
            print(f"\n✅ 配置验证通过")

def main():
    """命令行工具"""
    import argparse
    
    parser = argparse.ArgumentParser(description='Emoji版本管理工具')
    parser.add_argument('--status', action='store_true', help='显示当前状态')
    parser.add_argument('--freeze-all', action='store_true', help='冻结所有雪碧图')
    parser.add_argument('--new-version', type=str, help='创建新版本')
    
    args = parser.parse_args()
    
    manager = EmojiVersionManager()
    
    if args.status:
        manager.print_status()
    elif args.freeze_all:
        manager.freeze_all_sprites()
        manager.save_config()
        print("✅ 所有雪碧图已冻结")
    elif args.new_version:
        new_version = manager.create_new_version(args.new_version)
        manager.save_config()
        print(f"✅ 已创建新版本: {new_version}")
    else:
        manager.print_status()

if __name__ == "__main__":
    main()
