#!/usr/bin/env python3
"""
上传雪碧图到七牛云脚本
"""

import os
import requests
import json

class SpriteUploader:
    def __init__(self, base_url="http://localhost:8081"):
        self.base_url = base_url
        self.upload_url = f"{base_url}/api/image/upload"
        self.upload_token = "37395c61-a2ec-464e-9567-ce6fa92630f7"
        
    def upload_file(self, file_path):
        """上传单个文件到七牛云"""
        filename = os.path.basename(file_path)
        print(f"上传文件: {filename}")
        
        try:
            with open(file_path, 'rb') as f:
                files = {'image': (filename, f, 'image/png')}
                headers = {'Authorization': f'Bearer {self.upload_token}'}
                
                response = requests.post(self.upload_url, files=files, headers=headers)
                
                if response.status_code == 200:
                    result = response.json()
                    if result.get('code') == 0:
                        data = result.get('data', {})
                        print(f"  ✅ 上传成功: {data.get('url')}")
                        return {
                            'success': True,
                            'url': data.get('url'),
                            'oss_type': data.get('ossType'),
                            'filename': filename
                        }
                    else:
                        print(f"  ❌ 上传失败: {result.get('msg')}")
                        return {'success': False, 'error': result.get('msg')}
                else:
                    print(f"  ❌ HTTP错误: {response.status_code}")
                    return {'success': False, 'error': f'HTTP {response.status_code}'}
                    
        except Exception as e:
            print(f"  ❌ 上传异常: {e}")
            return {'success': False, 'error': str(e)}
    
    def upload_sprites(self, output_dir):
        """上传所有雪碧图文件"""
        print("=== 开始上传雪碧图到七牛云 ===")
        
        # 查找所有雪碧图文件
        sprite_files = []
        for filename in os.listdir(output_dir):
            if filename.startswith('emoji-sprite-') and filename.endswith('.png'):
                file_path = os.path.join(output_dir, filename)
                sprite_files.append(file_path)
        
        sprite_files.sort()  # 按文件名排序
        
        if not sprite_files:
            print("❌ 未找到雪碧图文件")
            return []
        
        print(f"找到 {len(sprite_files)} 个雪碧图文件")
        
        # 上传每个文件
        upload_results = []
        for file_path in sprite_files:
            result = self.upload_file(file_path)
            upload_results.append(result)
        
        # 统计结果
        success_count = sum(1 for r in upload_results if r['success'])
        print(f"\n=== 上传完成 ===")
        print(f"✅ 成功: {success_count}/{len(upload_results)}")
        
        if success_count < len(upload_results):
            print("❌ 部分文件上传失败")
        
        return upload_results
    
    def update_css_urls(self, output_dir, upload_results):
        """更新CSS文件中的图片URL"""
        print("=== 更新CSS文件URL ===")
        
        css_file = os.path.join(output_dir, "emoji-sprites.css")
        if not os.path.exists(css_file):
            print("❌ CSS文件不存在")
            return
        
        # 读取CSS内容
        with open(css_file, 'r', encoding='utf-8') as f:
            css_content = f.read()
        
        # 替换URL
        for result in upload_results:
            if result['success']:
                filename = result['filename']
                url = result['url']
                
                # 替换相对路径为CDN URL
                old_url = f"url('{filename}')"
                new_url = f"url('{url}')"
                css_content = css_content.replace(old_url, new_url)
                print(f"  替换: {filename} -> {url}")
        
        # 保存更新后的CSS
        updated_css_file = os.path.join(output_dir, "emoji-sprites-cdn.css")
        with open(updated_css_file, 'w', encoding='utf-8') as f:
            f.write(css_content)
        
        print(f"✅ CDN版本CSS已保存: {updated_css_file}")
        return updated_css_file
    
    def generate_config(self, output_dir, upload_results):
        """生成前端配置文件"""
        print("=== 生成前端配置 ===")
        
        config = {
            'sprite_urls': {},
            'sprite_info': {
                'target_size': 64,
                'sprites_per_row': 16,
                'emojis_per_sprite': 128
            },
            'upload_time': __import__('datetime').datetime.now().isoformat()
        }
        
        for result in upload_results:
            if result['success']:
                filename = result['filename']
                # 提取sprite组号 emoji-sprite-0.png -> 0
                group_id = int(filename.split('-')[2].split('.')[0])
                config['sprite_urls'][group_id] = result['url']
        
        config_file = os.path.join(output_dir, "emoji-config.json")
        with open(config_file, 'w', encoding='utf-8') as f:
            json.dump(config, f, indent=2, ensure_ascii=False)
        
        print(f"✅ 前端配置已生成: {config_file}")
        return config_file

def main():
    output_dir = "/media/jiang/hsk/practice/Golang/go/goCode/prac09/go_blog/scripts/emoji_output"
    
    # 创建上传器
    uploader = SpriteUploader()
    
    # 上传雪碧图
    upload_results = uploader.upload_sprites(output_dir)
    
    if any(r['success'] for r in upload_results):
        # 更新CSS文件
        uploader.update_css_urls(output_dir, upload_results)
        
        # 生成前端配置
        uploader.generate_config(output_dir, upload_results)
        
        print("\n🎉 雪碧图优化完成！")
        print("📋 下一步：")
        print("  1. 将 emoji-sprites-cdn.css 集成到前端")
        print("  2. 更新前端emoji解析逻辑")
        print("  3. 执行数据库迁移")
    else:
        print("\n❌ 上传失败，请检查网络连接和服务器状态")

if __name__ == "__main__":
    main()
