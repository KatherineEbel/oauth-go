package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	_ "github.com/mercadolibre/golang-restclient/rest"

	"github.com/KatherineEbel/oauth-go/oauth/errors"
)

var (
	restClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8080",
		Timeout: 200 * time.Millisecond,
	}
)

const (
	headerXPublic   = "X-Public"
	headerXClientId = "X-Client-Id"
	headerXUserId   = "X-User-Id"

	paramAccessToken = "access_token"
)

type IClient interface {
}

type accessToken struct {
	Id       string `json:"id"`
	UserId   int64  `json:"user_id"`
	ClientId int64  `json:"client_id"`
}

func IsPublic(r *http.Request) bool {
	if r == nil {
		return true
	}
	return r.Header.Get(headerXPublic) == "true"
}

func GetUserId(r *http.Request) int64 {
	if r == nil {
		return 0
	}
	uId, err := strconv.ParseInt(r.Header.Get(headerXUserId), 10, 64)
	if err != nil {
		return 0
	}
	return uId
}

func GetClientId(r *http.Request) int64 {
	if r == nil {
		return 0
	}
	cId, err := strconv.ParseInt(r.Header.Get(headerXClientId), 10, 64)
	if err != nil {
		return 0
	}
	return cId
}

func Authenticate(r *http.Request) *errors.RestError {
	if r == nil {
		return nil
	}
	cleanRequest(r)
	tokId := strings.TrimSpace(r.URL.Query().Get(paramAccessToken))
	if len(tokId) == 0 {
		return nil
	}
	tok, err := getAccessToken(tokId)
	if err != nil {
		return err
	}
	r.Header.Add(headerXClientId, fmt.Sprintf("%v", tok.ClientId))
	r.Header.Add(headerXUserId, fmt.Sprintf("%v", tok.UserId))
	return nil
}

func cleanRequest(r *http.Request) {
	if r == nil {
		return
	}
	r.Header.Del(headerXClientId)
	r.Header.Del(headerXUserId)
}

func getAccessToken(t string) (*accessToken, *errors.RestError) {
	res := restClient.Get(fmt.Sprintf("/oauth/access_token/%s", t))
	if res == nil || res.Response == nil {
		return nil, errors.NewInternalServerError("response timeout from rest client")
	}
	if res.StatusCode > 299 {
		var rErr errors.RestError
		if err := json.Unmarshal(res.Bytes(), &rErr); err != nil {
			return nil, errors.NewInternalServerError("unknown error type returned from rest client")
		}
		return nil, &rErr
	}
	var token *accessToken
	if err := json.Unmarshal(res.Bytes(), &token); err != nil {
		return nil, errors.NewInternalServerError("can't unmarshal response to user")
	}
	return token, nil
}
