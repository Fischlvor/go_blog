/**
 * Emoji样式管理器 - 支持动态加载和版本管理
 */

import { getEmojiConfig, type EmojiConfig as ApiEmojiConfig } from '@/api/emoji'

export interface SpriteInfo {
  id: number
  filename: string
  url: string
  range: [number, number]
  frozen: boolean
  size: [number, number]
}

// 本地使用的配置类型，字段与后端返回的大致一致，但 sprites 使用内部的 SpriteInfo 结构
export interface EmojiConfig {
  version: string
  total_emojis: number
  sprites: SpriteInfo[]
  updated_at: string
  // 旧->新 key 映射，主要给兼容逻辑使用，可选
  mapping?: Record<string, string>
}

export class EmojiStyleManager {
  private static instance: EmojiStyleManager
  private config: EmojiConfig | null = null
  private loadedStyles = new Set<string>()
  private styleElement: HTMLStyleElement | null = null

  private constructor() {
    this.createStyleElement()
  }

  static getInstance(): EmojiStyleManager {
    if (!EmojiStyleManager.instance) {
      EmojiStyleManager.instance = new EmojiStyleManager()
    }
    return EmojiStyleManager.instance
  }

  /**
   * 创建样式元素
   */
  private createStyleElement(): void {
    this.styleElement = document.createElement('style')
    this.styleElement.id = 'emoji-dynamic-styles'
    document.head.appendChild(this.styleElement)
  }

  /**
   * 加载emoji配置（通过统一的 API 模块）
   */
  async loadConfig(): Promise<EmojiConfig> {
    if (this.config) return this.config

    try {
      const res = await getEmojiConfig()
      if (res.code !== 0 || !res.data) {
        throw new Error(res.msg || '获取emoji配置失败')
      }

      const apiConfig = res.data as ApiEmojiConfig

      this.config = {
        version: apiConfig.version,
        total_emojis: apiConfig.total_emojis,
        updated_at: apiConfig.updated_at,
        sprites: (apiConfig.sprites || []) as unknown as SpriteInfo[],
        mapping: apiConfig.mapping as Record<string, string> | undefined
      }

      return this.config
    } catch (error) {
      console.error('加载emoji配置失败:', error)
      throw error
    }
  }

  /**
   * 获取当前配置
   */
  getConfig(): EmojiConfig | null {
    return this.config
  }

  /**
   * 生成基础emoji样式
   */
  private generateBaseStyles(): string {
    return `
/* Emoji基础样式 */
.emoji {
  display: inline-block;
  /* 与雪碧图单元格尺寸保持一致，避免只显示 1/4 */
  width: 32px;
  height: 32px;
  background-repeat: no-repeat;
  vertical-align: -0.1em;
  line-height: 1;
}

/* Emoji选择器样式 */
.emoji-item {
  width: 32px !important;
  height: 32px !important;
  cursor: pointer;
  border-radius: 4px;
  transition: all 0.2s;
}

.emoji-item:hover {
  background-color: #f0f9ff;
  transform: scale(1.2);
}
`
  }

  /**
   * 生成雪碧图样式
   */
  private generateSpriteStyles(sprites: SpriteInfo[]): string {
    let css = ''
    
    for (const sprite of sprites) {
      css += `
.emoji-sprite-${sprite.id} {
  background-image: url('${sprite.url}');
  /* 完全按照后端给的 size 作为显示尺寸 */
  background-size: ${sprite.size[0]}px ${sprite.size[1]}px;
}
`
    }
    
    return css
  }

  /**
   * 生成位置样式
   */
  private generatePositionStyles(sprites: SpriteInfo[]): string {
    let css = ''
    
    for (const sprite of sprites) {
      const [rangeStart, rangeEnd] = sprite.range
      const totalEmojis = rangeEnd - rangeStart + 1
      const cols = 16
      const rows = Math.ceil(totalEmojis / cols)
      
      for (let i = rangeStart; i <= rangeEnd; i++) {
        const posInSprite = i - rangeStart
        const row = Math.floor(posInSprite / cols)
        const col = posInSprite % cols

        // 使用后端下发的 size 动态计算每个格子的宽高
        const cellWidth = sprite.size[0] / cols
        const cellHeight = sprite.size[1] / rows

        const x = col * cellWidth
        const y = row * cellHeight
        
        css += `
.emoji-e${i} {
  background-position: -${x}px -${y}px;
}
`
      }
    }
    
    return css
  }

  /**
   * 应用样式到页面
   */
  private applyStyles(css: string): void {
    if (this.styleElement) {
      this.styleElement.textContent = css
    }
  }

  /**
   * 加载所有样式
   */
  async loadAllStyles(): Promise<void> {
    if (!this.config) {
      await this.loadConfig()
    }

    if (!this.config) {
      throw new Error('Failed to load emoji config')
    }

    // 生成完整的CSS
    let fullCSS = this.generateBaseStyles()
    fullCSS += this.generateSpriteStyles(this.config.sprites)
    fullCSS += this.generatePositionStyles(this.config.sprites)

    // 应用样式
    this.applyStyles(fullCSS)
  }

  /**
   * 加载增量样式（用于热更新）
   */
  async loadIncrementalStyles(newSprites: SpriteInfo[]): Promise<void> {
    if (!this.config) {
      await this.loadConfig()
    }

    // 生成新增的样式
    let incrementalCSS = this.generateSpriteStyles(newSprites)
    incrementalCSS += this.generatePositionStyles(newSprites)

    // 添加到现有样式
    if (this.styleElement) {
      this.styleElement.textContent += incrementalCSS
    }

    // 更新配置
    if (this.config) {
      this.config.sprites.push(...newSprites)
      this.config.total_emojis += newSprites.reduce((sum, sprite) => 
        sum + (sprite.range[1] - sprite.range[0] + 1), 0)
    }
  }

  /**
   * 检查样式是否已加载
   */
  isStyleLoaded(spriteId: number): boolean {
    return this.loadedStyles.has(`sprite-${spriteId}`)
  }

  /**
   * 标记样式为已加载
   */
  markStyleLoaded(spriteId: number): void {
    this.loadedStyles.add(`sprite-${spriteId}`)
  }

  /**
   * 获取emoji信息
   */
  getEmojiInfo(emojiKey: string): { spriteId: number; position: [number, number] } | null {
    if (!this.config) return null

    // 解析emoji键 (e123 -> 123)
    const match = emojiKey.match(/^e(\d+)$/)
    if (!match) return null

    const emojiIndex = parseInt(match[1])

    // 查找对应的雪碧图
    for (const sprite of this.config.sprites) {
      const [rangeStart, rangeEnd] = sprite.range
      if (emojiIndex >= rangeStart && emojiIndex <= rangeEnd) {
        const posInSprite = emojiIndex - rangeStart
        const row = Math.floor(posInSprite / 16)
        const col = posInSprite % 16

        return {
          spriteId: sprite.id,
          position: [col * 32, row * 32]
        }
      }
    }

    return null
  }

  /**
   * 预加载指定范围的样式
   */
  async preloadStyles(startIndex: number, count: number): Promise<void> {
    if (!this.config) return

    const endIndex = startIndex + count - 1
    const spritesToLoad: SpriteInfo[] = []

    for (const sprite of this.config.sprites) {
      const [rangeStart, rangeEnd] = sprite.range
      
      // 检查是否与请求范围重叠
      if (rangeStart <= endIndex && rangeEnd >= startIndex) {
        if (!this.isStyleLoaded(sprite.id)) {
          spritesToLoad.push(sprite)
          this.markStyleLoaded(sprite.id)
        }
      }
    }

    if (spritesToLoad.length > 0) {
      console.log(`🔄 预加载 ${spritesToLoad.length} 个雪碧图样式`)
      // 这里可以添加实际的预加载逻辑
    }
  }

  /**
   * 清理样式
   */
  cleanup(): void {
    if (this.styleElement) {
      this.styleElement.remove()
      this.styleElement = null
    }
    this.loadedStyles.clear()
    this.config = null
  }

  /**
   * 获取版本信息
   */
  getVersion(): string {
    return this.config?.version || 'unknown'
  }

  /**
   * 检查是否有新版本
   */
  async checkForUpdates(): Promise<boolean> {
    try {
      const response = await fetch('/emoji/emoji_config.json')
      if (!response.ok) return false

      const newConfig: EmojiConfig = await response.json()
      
      if (!this.config) return true
      
      return newConfig.version !== this.config.version
    } catch (error) {
      console.error('检查更新失败:', error)
      return false
    }
  }

  /**
   * 热更新样式
   */
  async hotUpdate(): Promise<boolean> {
    try {
      const hasUpdate = await this.checkForUpdates()
      if (!hasUpdate) return false

      // 重新加载配置和样式
      await this.loadConfig()
      await this.loadAllStyles()
      return true
    } catch (error) {
      console.error('❌ 热更新失败:', error)
      return false
    }
  }
}

// 导出单例实例
export const emojiStyleManager = EmojiStyleManager.getInstance()
