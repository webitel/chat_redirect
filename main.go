package main

import (
	"chat_redirect/app"
	"chat_redirect/model"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/wlog"
)

func main() {
	//flagDefine()
	log := wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  wlog.LevelDebug,
	})

	wlog.RedirectStdLog(log)
	wlog.InitGlobalLogger(log)

	config, appErr := loadConfig()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	// * Create an application layer
	application, err := app.NewApplication(config)
	if err != nil {
		panic(err.Error())
	}

	err = application.Start()
	if err != nil {
		panic(err.Error())
	}
	return

}

func loadConfig() (*model.AppConfig, error) {
	var appConfig model.AppConfig

	configurator := configuration.New(
		&appConfig,
		// order of execution will be preserved:
		configuration.NewFlagProvider(),
		configuration.NewEnvProvider(),
		configuration.NewDefaultProvider(),
	).SetOptions(configuration.OnFailFnOpt(func(err error) {
	}))

	if err := configurator.InitValues(); err != nil {
		return nil, err
	}
	err := appConfig.Validate()
	if err != nil {
		return nil, err
	}
	return &appConfig, nil
}
