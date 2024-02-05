package app

import (
	"chat_redirect/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/guid"
	"github.com/webitel/wlog"
	"net/http"
	"net/url"
	"strings"
)

type ApiHandler struct {
	app *Application
}

func (h *ApiHandler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		wlog.Debug(formatRedirectLog("accepted redirect request"))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var (
		id    string
		query url.Values
		uri   string
		//pl    any
		bytes []byte
	)
	query = r.URL.Query()
	payload := query.Get("payload")
	wlog.Debug(formatRedirectLog(fmt.Sprintf("payload - %s", payload)))
	uri = query.Get("bot")
	var general model.GeneralBody
	err := json.Unmarshal([]byte(payload), &general)
	if err != nil {
		r := err.Error()
		wlog.Debug(formatRedirectLog(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(r))
	}
	bytes, err = json.Marshal(general.Payload)
	if err != nil {
		r := err.Error()
		wlog.Debug(formatRedirectLog(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(r))
	}
	id = guid.New().String()
	err = h.app.Redis.Set(context.Background(), fmt.Sprintf("%d.%s", *h.app.config.Server.DomainId, id), string(bytes), 10000)
	if err != nil {
		r := err.Error()
		wlog.Debug(formatRedirectLog(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(r))
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		r := err.Error()
		wlog.Debug(formatRedirectLog(r))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(r))
	}
	q := parsedUrl.Query()
	for _, link := range h.app.config.Bots.ParsedLinks {
		if strings.HasPrefix(uri, link) {
			switch general.Gateway {
			case "telegram":
				q.Add("start", id)
				parsedUrl.RawQuery = q.Encode()
			case "viber":
				q.Add("context", id)
				parsedUrl.RawQuery = q.Encode()
			default:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("unsupported gateway type"))
				return
			}
			url := parsedUrl.String()
			wlog.Debug(formatRedirectLog(fmt.Sprintf("redirection to %s ...", url)))
			http.Redirect(w, r, parsedUrl.String(), 302)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("unsupported link type"))
	return
}

func formatRedirectLog(log string) string {
	return fmt.Sprintf("redirect: %s", log)
}

func NewApiHandler(app *Application) (*ApiHandler, error) {
	return &ApiHandler{app: app}, nil
}
