/**
 * OAuth State 参数工具函数
 */

/**
 * 生成UUID v4
 * @returns {string} UUID字符串
 */
export function generateNonce() {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

/**
 * 编码state参数
 * @param {Object} data - state数据对象
 * @param {string} data.nonce - 防重放攻击的随机值
 * @param {string} data.app_id - 应用ID
 * @param {string} data.device_id - 设备ID
 * @param {string} data.redirect_uri - 回调地址
 * @param {string} data.return_url - 用户目标页面（可选）
 * @param {number} data.iat - 签发时间戳（秒，可选）
 * @param {number} data.exp - 过期时间戳（秒，可选）
 * @returns {string} Base64编码的state字符串
 */
export function encodeState(data) {
  // 验证必需字段
  if (!data.nonce) {
    throw new Error('nonce不能为空')
  }
  if (!data.app_id) {
    throw new Error('app_id不能为空')
  }
  if (!data.device_id) {
    throw new Error('device_id不能为空')
  }
  if (!data.redirect_uri) {
    throw new Error('redirect_uri不能为空')
  }

  // 设置时间戳（如果未设置）
  const now = Math.floor(Date.now() / 1000)
  if (!data.iat) {
    data.iat = now
  }
  if (!data.exp) {
    data.exp = now + 300 // 默认5分钟过期
  }

  // JSON序列化并Base64编码
  const jsonString = JSON.stringify(data)
  const base64String = btoa(jsonString)
  
  return base64String
}

/**
 * 解码state参数
 * @param {string} state - Base64编码的state字符串
 * @returns {Object} 解码后的state数据对象
 */
export function decodeState(state) {
  if (!state) {
    throw new Error('state参数为空')
  }

  try {
    // Base64解码
    const jsonString = atob(state)
    
    // JSON解析
    const data = JSON.parse(jsonString)
    
    return data
  } catch (error) {
    throw new Error('state参数解析失败: ' + error.message)
  }
}

/**
 * 验证state参数（仅前端基础验证，完整验证在后端）
 * @param {string} state - Base64编码的state字符串
 * @returns {Object} 解码后的state数据对象
 */
export function validateState(state) {
  const data = decodeState(state)

  // 验证必需字段
  if (!data.nonce) {
    throw new Error('state中缺少nonce')
  }
  if (!data.app_id) {
    throw new Error('state中缺少app_id')
  }
  if (!data.redirect_uri) {
    throw new Error('state中缺少redirect_uri')
  }

  // 验证过期时间
  const now = Math.floor(Date.now() / 1000)
  if (data.exp && data.exp < now) {
    throw new Error('state已过期')
  }

  return data
}
