// This is the main package for ytcctl.
// Ytcctl is used to manage the yashan trace collector.
package main

import (
	"fmt"
	"strings"

	"ytc/commons/flags"
	"ytc/commons/std"
	"ytc/defs/compiledef"
	"ytc/defs/confdef"
	"ytc/defs/runtimedef"
	"ytc/i18n"
	"ytc/log"

	"git.yasdb.com/go/yaserr"
	"github.com/alecthomas/kong"
)

const (
	_APP_NAME        = "ytcctl"
	_APP_DESCRIPTION = "Ytcctl is used to manage the yashan trace collector."
)

func main() {
	var app App
	options := flags.NewAppOptions(_APP_NAME, _APP_DESCRIPTION, compiledef.GetAPPVersion())
	ctx := kong.Parse(&app, options...)
	if err := initApp(app); err != nil {
		ctx.FatalIfErrorf(err)
	}
	finalize := std.GetRedirecter().RedirectStd()
	defer finalize()
	std.WriteToFile(fmt.Sprintf("execute: %s %s\n", _APP_NAME, strings.Join(ctx.Args, " ")))
	if err := ctx.Run(); err != nil {
		fmt.Println(yaserr.Unwrap(err))
	}
}

func initLogger(logPath, level string) error {
	optFuncs := []log.OptFunc{
		log.SetLogPath(logPath),
		log.SetLevel(level),
	}
	return log.InitLogger(_APP_NAME, log.NewLogOption(optFuncs...))
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
	if lang == "" {
		lang = confdef.GetYTCConf().Language
	}
	if lang == "" {
		lang = "zh"
	}
	i18n.InitI18n(lang)
	if err := initLogger(runtimedef.GetLogPath(), confdef.GetYTCConf().LogLevel); err != nil {
		return err
	}
	if err := std.InitRedirecter(); err != nil {
		return err
	}
	return nil
}
