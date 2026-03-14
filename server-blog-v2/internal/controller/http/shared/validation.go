package shared

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// 字段名称映射
var fieldNames = map[string]string{
	"Title":         "标题",
	"Content":       "内容",
	"Slug":          "标识",
	"Excerpt":       "摘要",
	"FeaturedImage": "封面图",
	"CategoryID":    "分类",
	"TagIDs":        "标签",
	"Status":        "状态",
	"Visibility":    "可见性",
	"Name":          "名称",
	"Email":         "邮箱",
	"Password":      "密码",
	"URL":           "链接",
	"UserUUID":      "用户",
}

// 验证规则映射
var tagMessages = map[string]string{
	"required": "不能为空",
	"max":      "长度超出限制",
	"min":      "长度不足",
	"email":    "格式不正确",
	"url":      "格式不正确",
	"oneof":    "值不在允许范围内",
}

// TranslateValidationErrors 将验证错误翻译为友好的中文提示
func TranslateValidationErrors(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var messages []string
	for _, e := range validationErrors {
		fieldName := e.Field()
		if name, exists := fieldNames[fieldName]; exists {
			fieldName = name
		}

		tagMsg := tagMessages[e.Tag()]
		if tagMsg == "" {
			tagMsg = fmt.Sprintf("验证失败(%s)", e.Tag())
		}

		messages = append(messages, fmt.Sprintf("%s%s", fieldName, tagMsg))
	}

	return strings.Join(messages, "; ")
}
