/**
 * Emoji 全局初始化
 * 在应用启动时预加载 emoji 配置，避免首次使用时的延迟
 */

import { emojiStyleManager } from './emojiStyleManager'

/**
 * 预加载 emoji 配置
 * 可以在应用启动时调用，提前加载配置
 */
export async function preloadEmojiConfig(): Promise<void> {
  try {
    console.log('🎯 预加载 emoji 配置...')
    await emojiStyleManager.loadConfig()
    await emojiStyleManager.loadAllStyles()
    console.log('✅ emoji 配置预加载完成')
  } catch (error) {
    console.error('❌ emoji 配置预加载失败:', error)
    // 不抛出错误，允许后续按需加载
  }
}

/**
 * 检查 emoji 配置是否已加载
 */
export function isEmojiConfigLoaded(): boolean {
  return emojiStyleManager.getConfig() !== null
}

/**
 * 获取 emoji 配置状态
 */
export function getEmojiConfigStatus(): {
  loaded: boolean
  totalEmojis: number
  version: string
  spriteCount: number
} {
  const config = emojiStyleManager.getConfig()
  
  if (!config) {
    return {
      loaded: false,
      totalEmojis: 0,
      version: 'unknown',
      spriteCount: 0
    }
  }

  return {
    loaded: true,
    totalEmojis: config.total_emojis,
    version: config.version,
    spriteCount: config.sprites.length
  }
}
