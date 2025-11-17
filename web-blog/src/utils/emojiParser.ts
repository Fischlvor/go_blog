/**
 * Emoji解析工具 - 支持动态加载和版本管理
 */

import { emojiStyleManager, type EmojiConfig } from './emojiStyleManager'

export interface EmojiInfo {
  key: string          // emoji键名 如 'e1'
  spriteGroup: number  // 雪碧图组号
  index: number        // 全局索引
}

class EmojiParser {
  private emojiList: EmojiInfo[] = []
  private isInitialized = false

  /**
   * 初始化解析器
   */
  async initialize(): Promise<void> {
    if (this.isInitialized) return

    try {
      // 加载emoji配置（只调用一次）
      const config = await emojiStyleManager.loadConfig()
      await emojiStyleManager.loadAllStyles()

      // 从 sprites 的 range 生成 emoji 列表
      this.emojiList = []
      for (const sprite of config.sprites) {
        const [start, end] = sprite.range
        for (let i = start; i <= end; i++) {
          this.emojiList.push({
            key: `e${i}`,
            spriteGroup: sprite.id,
            index: i
          })
        }
      }
      
      this.isInitialized = true
    } catch (error) {
      console.error('EmojiParser 初始化失败:', error)
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
   * 解析emoji字符串（现在只是直接返回，因为已经是 e123 格式）
   */
  async parseEmojis(text: string): Promise<string> {
    await this.ensureInitialized()
    // 数据库中已经是 :emoji:e123: 格式，不需要转换
    return text
  }

  /**
   * 获取所有emoji信息
   */
  async getAllEmojis(): Promise<EmojiInfo[]> {
    await this.ensureInitialized()
    return this.emojiList
  }

  /**
   * 验证emoji键是否存在
   */
  async isValidEmojiKey(key: string): Promise<boolean> {
    await this.ensureInitialized()
    
    // 检查是否是有效的键
    if (key.startsWith('e') && /^e\d+$/.test(key)) {
      return this.emojiList.some(emoji => emoji.key === key)
    }
    
    return false
  }

  /**
   * 获取emoji信息
   */
  async getEmojiInfo(key: string): Promise<EmojiInfo | null> {
    await this.ensureInitialized()

    if (!key.startsWith('e')) {
      return null
    }

    return this.emojiList.find(emoji => emoji.key === key) || null
  }


  /**
   * 渲染emoji为HTML
   */
  async renderEmojiToHTML(emojiKey: string): Promise<string> {
    const info = await this.getEmojiInfo(emojiKey)
    if (!info) return emojiKey

    return `<span class="emoji emoji-sprite-${info.spriteGroup} emoji-${info.key}" title="${info.key}"></span>`
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
   * 获取统计信息
   */
  async getStats(): Promise<{
    totalEmojis: number
    version: string
    spriteCount: number
  }> {
    await this.ensureInitialized()
    
    const config = emojiStyleManager.getConfig()
    
    return {
      totalEmojis: this.emojiList.length,
      version: config?.version || 'unknown',
      spriteCount: config?.sprites.length || 0
    }
  }

  /**
   * 热更新
   */
  async hotUpdate(): Promise<boolean> {
    try {
      // 重置初始化状态
      this.isInitialized = false
      
      // 重新初始化
      await this.initialize()
      
      // 触发样式管理器热更新
      const updated = await emojiStyleManager.hotUpdate()
      
      return updated
    } catch (error) {
      console.error('❌ Emoji热更新失败:', error)
      return false
    }
  }

  /**
   * 清理资源
   */
  cleanup(): void {
    this.emojiList = []
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
