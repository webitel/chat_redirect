package app

import (
	"chat_redirect/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/beevik/guid"
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
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var (
		id    string
		query url.Values
		uri   string
	)
	query = r.URL.Query()
	payload := query.Get("payload")
	uri = query.Get("bot")
	var general model.GeneralBody
	err := json.Unmarshal([]byte(payload), &general)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	switch general.Type {
	case "general":
		var pl model.GeneralQuestions
		bytes, err := json.Marshal(general.Payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(bytes, &pl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		id = guid.New().String()
		err = h.app.Redis.Set(context.Background(), fmt.Sprintf("%d.%s", *h.app.config.Server.DomainId, id), string(bytes), 10000)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	case "card":
		var pl model.CardTransactionQuestions
		bytes, err := json.Marshal(general.Payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(bytes, &pl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		id = guid.New().String()
		err = h.app.Redis.Set(context.Background(), fmt.Sprintf("%d.%s", *h.app.config.Server.DomainId, id), string(bytes), 10000)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	case "payment":
		var pl model.TransferDetailsQuestions
		bytes, err := json.Marshal(general.Payload)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		err = json.Unmarshal(bytes, &pl)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		id = guid.New().String()
		err = h.app.Redis.Set(context.Background(), fmt.Sprintf("%d.%s", *h.app.config.Server.DomainId, id), string(bytes), 10000)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	}
	parsedUrl, err := url.Parse(uri)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	q := parsedUrl.Query()
	for _, link := range h.app.config.Bots.ParsedLinks {
		if strings.HasPrefix(uri, link) {
			switch general.Gateway {
			case "telegram":
				q.Add("start", id)
				parsedUrl.RawQuery = q.Encode()
				http.Redirect(w, r, parsedUrl.String(), 302)
				return
			case "viber":
				q.Add("context", id)
				parsedUrl.RawQuery = q.Encode()
				url := parsedUrl.String()
				http.Redirect(w, r, url, 302)
				return
			default:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("unsupported gateway type"))
				return
			}

		}
	}
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("unsupported link type"))
	return
}

func NewApiHandler(app *Application) (*ApiHandler, error) {
	return &ApiHandler{app: app}, nil
}
