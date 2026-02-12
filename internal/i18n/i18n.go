package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Translator 翻译器接口
type Translator struct {
	locales map[string]map[string]string
	mu      sync.RWMutex
}

var (
	defaultTranslator *Translator
	once              sync.Once
)

// Init 初始化翻译器
func Init(localesPath string) error {
	var err error
	once.Do(func() {
		defaultTranslator = &Translator{
			locales: make(map[string]map[string]string),
		}
		err = defaultTranslator.LoadLocales(localesPath)
	})
	return err
}

// LoadLocales 加载翻译文件
func (t *Translator) LoadLocales(localesPath string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 支持的语言列表
	languages := []string{"zh-CN", "en-US"}

	for _, lang := range languages {
		filePath := filepath.Join(localesPath, fmt.Sprintf("%s.json", lang))
		data, err := os.ReadFile(filePath)
		if err != nil {
			// 如果文件不存在，跳过
			continue
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			return fmt.Errorf("failed to parse %s: %w", filePath, err)
		}

		t.locales[lang] = translations
	}

	return nil
}

// T 翻译函数 (Translate)
func (t *Translator) T(locale, key string, args ...interface{}) string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 获取对应语言的翻译
	if translations, ok := t.locales[locale]; ok {
		if translation, exists := translations[key]; exists {
			// 如果有参数，进行格式化
			if len(args) > 0 {
				return fmt.Sprintf(translation, args...)
			}
			return translation
		}
	}

	// 如果当前语言没有，尝试使用英文
	if locale != "en-US" {
		if translations, ok := t.locales["en-US"]; ok {
			if translation, exists := translations[key]; exists {
				if len(args) > 0 {
					return fmt.Sprintf(translation, args...)
				}
				return translation
			}
		}
	}

	// 如果都没有，返回 key
	return key
}

// T 全局翻译函数
func T(locale, key string, args ...interface{}) string {
	if defaultTranslator == nil {
		return key
	}
	return defaultTranslator.T(locale, key, args...)
}

// GetSupportedLocales 获取支持的语言列表
func GetSupportedLocales() []string {
	return []string{"zh-CN", "en-US"}
}

// ParseAcceptLanguage 解析 Accept-Language 头
func ParseAcceptLanguage(acceptLanguage string) string {
	if acceptLanguage == "" {
		return "zh-CN" // 默认中文
	}

	// 简单解析，支持 zh-CN, en-US, zh, en 等格式
	if len(acceptLanguage) >= 2 {
		lang := acceptLanguage[:2]
		switch lang {
		case "zh":
			return "zh-CN"
		case "en":
			return "en-US"
		}
	}

	// 如果是完整的语言代码
	if acceptLanguage == "zh-CN" || acceptLanguage == "en-US" {
		return acceptLanguage
	}

	return "zh-CN" // 默认中文
}
