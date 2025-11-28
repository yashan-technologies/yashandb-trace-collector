package i18n

import (
	"embed"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var localizer *i18n.Localizer
var currentLang string

func InitI18n(lang string) {
	currentLang = lang
	bundle := i18n.NewBundle(language.Chinese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// Load translations from embedded FS
	_, err := bundle.LoadMessageFileFS(localeFS, "locales/active.en.toml")
	if err != nil {
		panic("Failed to load active.en.toml: " + err.Error())
	}
	_, err = bundle.LoadMessageFileFS(localeFS, "locales/active.zh.toml")
	if err != nil {
		panic("Failed to load active.zh.toml: " + err.Error())
	}

	// Determine language
	if lang == "" {
		lang = "zh"
	}

	accepts := []string{}
	if lang != "" {
		// Handle formats like "zh_CN.UTF-8"
		base := strings.Split(lang, ".")[0]
		accepts = append(accepts, base)
		// Also try simple "zh" if it was "zh_CN"
		if strings.Contains(base, "_") {
			accepts = append(accepts, strings.Split(base, "_")[0])
		}
	}
	accepts = append(accepts, "zh")
	accepts = append(accepts, "en")

	localizer = i18n.NewLocalizer(bundle, accepts...)
}

// T translates a message ID
func T(messageID string) string {
	if localizer == nil {
		InitI18n("")
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID // Fallback to ID if not found
	}
	return msg
}

// TWithData translates a message ID with template data
func TWithData(messageID string, data map[string]interface{}) string {
	if localizer == nil {
		InitI18n("")
	}
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID // Fallback to ID
	}
	return msg
}

// GetCurrentLang returns the current language
func GetCurrentLang() string {
	if currentLang == "" {
		return "zh"
	}
	return currentLang
}
