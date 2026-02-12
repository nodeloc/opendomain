package middleware

import (
	"opendomain/internal/i18n"

	"github.com/gin-gonic/gin"
)

const (
	LocaleContextKey = "locale"
)

// I18nMiddleware 国际化中间件
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取语言设置
		acceptLanguage := c.GetHeader("Accept-Language")

		// 解析语言
		locale := i18n.ParseAcceptLanguage(acceptLanguage)

		// 将语言设置存入上下文
		c.Set(LocaleContextKey, locale)

		c.Next()
	}
}

// GetLocale 获取当前请求的语言
func GetLocale(c *gin.Context) string {
	if locale, exists := c.Get(LocaleContextKey); exists {
		return locale.(string)
	}
	return "zh-CN" // 默认中文
}

// T 上下文翻译函数
func T(c *gin.Context, key string, args ...interface{}) string {
	locale := GetLocale(c)
	return i18n.T(locale, key, args...)
}
