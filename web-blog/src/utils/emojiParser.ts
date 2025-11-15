/**
 * Emoji解析工具 - 支持动态加载和版本管理
 */

import { emojiStyleManager, type EmojiConfig } from './emojiStyleManager'

export interface EmojiInfo {
  oldKey: string    // 原始键名 如 's1'
  newKey: string    // 新键名 如 'e1'
  spriteGroup: number  // 雪碧图组号
  index: number     // 在组内的索引
}

class EmojiParser {
  private emojiMapping: Record<string, string> = {}
  private reverseMapping: Record<string, string> = {}
  private isInitialized = false

  /**
   * 初始化解析器
   */
  async initialize(): Promise<void> {
    if (this.isInitialized) return

    try {
      // 加载emoji配置
      await emojiStyleManager.loadConfig()
      await emojiStyleManager.loadAllStyles()

      // 加载映射表
      await this.loadMapping()
      
      this.isInitialized = true
    } catch (error) {
      console.error('EmojiParser 初始化失败:', error)
      throw error
    }
  }

  /**
   * 加载映射表
   * 直接复用 EmojiStyleManager 中的配置，不再单独拉静态 JSON
   */
  private async loadMapping(): Promise<void> {
    try {
      const config = await emojiStyleManager.loadConfig()

      this.emojiMapping = config.mapping || {}

      // 创建反向映射
      this.reverseMapping = {}
      for (const [oldKey, newKey] of Object.entries(this.emojiMapping)) {
        this.reverseMapping[newKey] = oldKey
      }
    } catch (error) {
      console.error('加载emoji映射失败:', error)
      throw error
    }
  }

  /**
   * 确保已初始化
   */
  private async ensureInitialized(): Promise<void> {
    if (!this.isInitialized) {
      await this.initialize()
    }
  }

  /**
   * 解析emoji字符串，支持多种格式
   */
  async parseEmojis(text: string): Promise<string> {
    await this.ensureInitialized()

    // 支持的格式：
    // 1. :emoji:s123: -> :emoji:e456:
    // 2. ![](/emoji/s123.png) -> :emoji:e456:
    // 3. ![](emoji/s123.png) -> :emoji:e456:
    
    let result = text

    // 格式1: :emoji:s123:
    result = result.replace(/:emoji:s(\d+):/g, (match, num) => {
      const oldKey = `s${num}`
      const newKey = this.emojiMapping[oldKey]
      return newKey ? `:emoji:${newKey}:` : match
    })

    // 格式2: ![](/emoji/s123.png)
    result = result.replace(/!\[\]\(\/emoji\/s(\d+)\.png\)/g, (match, num) => {
      const oldKey = `s${num}`
      const newKey = this.emojiMapping[oldKey]
      return newKey ? `:emoji:${newKey}:` : match
    })

    // 格式3: ![](emoji/s123.png)
    result = result.replace(/!\[\]\(emoji\/s(\d+)\.png\)/g, (match, num) => {
      const oldKey = `s${num}`
      const newKey = this.emojiMapping[oldKey]
      return newKey ? `:emoji:${newKey}:` : match
    })

    return result
  }

  /**
   * 获取所有emoji信息
   */
  async getAllEmojis(): Promise<EmojiInfo[]> {
    await this.ensureInitialized()

    const emojis: EmojiInfo[] = []
    
    for (const [oldKey, newKey] of Object.entries(this.emojiMapping)) {
      // 解析新键获取索引 (e123 -> 123)
      const match = newKey.match(/^e(\d+)$/)
      if (!match) continue

      const index = parseInt(match[1])
      const spriteGroup = Math.floor(index / 128)

      emojis.push({
        oldKey,
        newKey,
        spriteGroup,
        index
      })
    }

    // 按新键排序
    emojis.sort((a, b) => {
      const aIndex = parseInt(a.newKey.substring(1))
      const bIndex = parseInt(b.newKey.substring(1))
      return aIndex - bIndex
    })

    return emojis
  }

  /**
   * 验证emoji键是否存在
   */
  async isValidEmojiKey(key: string): Promise<boolean> {
    await this.ensureInitialized()
    
    // 检查是否是有效的新键
    if (key.startsWith('e') && /^e\d+$/.test(key)) {
      return Object.values(this.emojiMapping).includes(key)
    }
    
    // 检查是否是有效的旧键
    if (key.startsWith('s') && /^s\d+$/.test(key)) {
      return key in this.emojiMapping
    }
    
    return false
  }

  /**
   * 获取emoji信息
   */
  async getEmojiInfo(key: string): Promise<EmojiInfo | null> {
    await this.ensureInitialized()

    let oldKey: string
    let newKey: string

    if (key.startsWith('s')) {
      oldKey = key
      newKey = this.emojiMapping[key]
      if (!newKey) return null
    } else if (key.startsWith('e')) {
      newKey = key
      oldKey = this.reverseMapping[key]
      if (!oldKey) return null
    } else {
      return null
    }

    // 解析索引
    const match = newKey.match(/^e(\d+)$/)
    if (!match) return null

    const index = parseInt(match[1])
    const spriteGroup = Math.floor(index / 128)

    return {
      oldKey,
      newKey,
      spriteGroup,
      index
    }
  }

  /**
   * 转换旧格式到新格式
   */
  async convertOldToNew(oldKey: string): Promise<string | null> {
    await this.ensureInitialized()
    return this.emojiMapping[oldKey] || null
  }

  /**
   * 转换新格式到旧格式
   */
  async convertNewToOld(newKey: string): Promise<string | null> {
    await this.ensureInitialized()
    return this.reverseMapping[newKey] || null
  }

  /**
   * 渲染emoji为HTML
   */
  async renderEmojiToHTML(emojiKey: string): Promise<string> {
    const info = await this.getEmojiInfo(emojiKey)
    if (!info) return emojiKey

    return `<span class="emoji emoji-sprite-${info.spriteGroup} emoji-${info.newKey}" title="${info.oldKey} -> ${info.newKey}"></span>`
  }

  /**
   * 解析并渲染文本中的emoji
   */
  async renderTextWithEmojis(text: string): Promise<string> {
    await this.ensureInitialized()

    // 解析 :emoji:e123: 格式
    return text.replace(/:emoji:(e\d+):/g, (match, emojiKey) => {
      const info = emojiStyleManager.getEmojiInfo(emojiKey)
      if (!info) return match

      return `<span class="emoji emoji-sprite-${info.spriteId} emoji-${emojiKey}" title="${emojiKey}"></span>`
    })
  }

  /**
   * 获取映射统计信息
   */
  async getStats(): Promise<{
    totalEmojis: number
    version: string
    spriteCount: number
  }> {
    await this.ensureInitialized()
    
    const config = emojiStyleManager.getConfig()
    
    return {
      totalEmojis: Object.keys(this.emojiMapping).length,
      version: config?.version || 'unknown',
      spriteCount: config?.sprites.length || 0
    }
  }

  /**
   * 热更新映射表
   */
  async hotUpdate(): Promise<boolean> {
    try {
      // 重新加载映射表
      await this.loadMapping()
      
      // 触发样式管理器热更新
      const updated = await emojiStyleManager.hotUpdate()
      
      return updated
    } catch (error) {
      console.error('❌ Emoji映射热更新失败:', error)
      return false
    }
  }

  /**
   * 清理资源
   */
  cleanup(): void {
    this.emojiMapping = {}
    this.reverseMapping = {}
    this.isInitialized = false
    emojiStyleManager.cleanup()
  }
}

// 导出单例实例
export const emojiParser = new EmojiParser()

// 兼容性导出（保持与旧版本的接口兼容）
export async function getAllEmojis(): Promise<EmojiInfo[]> {
  return emojiParser.getAllEmojis()
}

export async function parseEmojis(text: string): Promise<string> {
  return emojiParser.parseEmojis(text)
}

export async function isValidEmojiKey(key: string): Promise<boolean> {
  return emojiParser.isValidEmojiKey(key)
}

export async function getEmojiInfo(key: string): Promise<EmojiInfo | null> {
  return emojiParser.getEmojiInfo(key)
}

// 渲染文本中的 :emoji:e123: 为带样式的 span
export async function renderTextWithEmojisForText(text: string): Promise<string> {
  return emojiParser.renderTextWithEmojis(text)
}
