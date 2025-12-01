// This is the main package for ytcd.
// Ytcd means ytc daemon, which is used to do some scheduled collector works.
package main

import (
	"ytc/commons/flags"
	"ytc/defs/compiledef"
	"ytc/defs/confdef"
	"ytc/defs/runtimedef"
	"ytc/i18n"

	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "ytcd"
	_APP_DESCRIPTION = "Ytcd means ytc daemon, which is used to do some scheduled collector works."
)

func main() {
	var app App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := initApp(app); err != nil {
		ctx.FatalIfErrorf(err)
	}
	if err := ctx.Run(); err != nil {
		ctx.FatalIfErrorf(err)
	}
}

func initApp(app App) error {
	if err := runtimedef.InitRuntime(); err != nil {
		return err
	}
	if err := confdef.InitConf(app.Config); err != nil {
		return err
	}
	// 优先使用命令行参数指定的语言，如果命令行参数为默认值，则使用配置文件中的语言
	lang := app.Lang
	if lang == "zh" && confdef.GetYTCConf().Language != "" {
		lang = confdef.GetYTCConf().Language
	}
	i18n.InitI18n(lang)
	return nil
}
