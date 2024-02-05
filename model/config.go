package model

import (
	"errors"
	"strings"
)

const (
	SERVICE_NAME = "webitel.chat.redirect"
)

type AppConfig struct {
	Redis  *RedisSettings
	Server *HttpServerSettings
	Bots   *BotsSettings
}
type RedisSettings struct {
	Host     *string `flag:"redis_host|127.0.0.1|Redis server host"`
	Port     *int    `flag:"redis_port|6379|Redis server port"`
	Password *string `flag:"redis_password||Redis password"`
	Database *int    `flag:"redis_db|0|Redis database"`
}

type BotsSettings struct {
	Links       string `flag:"allow_link||Viber bot deeplink to identify the right links for redirection"`
	ParsedLinks []string
}

type HttpServerSettings struct {
	Port     *int ` flag:"http_port|10040|Http server port"`
	DomainId *int `flag:"domain|1|Domain id "`
}

func (a *AppConfig) Validate() error {
	if a.Redis == nil || a.Server == nil || a.Bots == nil {
		return errors.New("missing setup parameters")
	}
	if a.Redis.Host == nil || a.Redis.Port == nil {
		return errors.New("missing required redis setup parameters")
	}
	if a.Redis.Port == nil {
		return errors.New("missing port parameter")
	}
	if a.Redis.Password == nil {
		pass := ""
		a.Redis.Password = &pass
	}

	splittedLinks := strings.Split(a.Bots.Links, ",")
	if len(splittedLinks) <= 0 {
		return errors.New("at least one bot link should be set")
	}
	a.Bots.ParsedLinks = splittedLinks

	return nil
}
