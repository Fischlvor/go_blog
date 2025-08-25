package utils

import (
	"math"
	"time"
)

// 计算文章分数
func CalculateArticleScore(createTime time.Time, views, comments, likes int64) float64 {
	// 1. 时间衰减因子 (0-1之间，越新内容分数越高)
	hoursSinceCreation := time.Since(createTime).Hours()
	timeDecay := math.Exp(-hoursSinceCreation / (24 * 7 * 2)) // 半衰期为2周

	// 2. 浏览数因子 (对数处理，防止过高)
	viewFactor := math.Log1p(float64(views)) / 10 // 对数处理并缩放

	// 3. 互动因子 (评论和点赞的加权组合)
	interactionFactor := float64(comments)*0.3 + float64(likes)*0.7
	interactionFactor = math.Log1p(interactionFactor) / 5 // 对数处理并缩放

	// 4. 综合评分 (加权组合)
	score := timeDecay*0.4 + viewFactor*0.3 + interactionFactor*0.3

	// 确保分数在0-1之间
	return math.Max(0, math.Min(1, score))
}
