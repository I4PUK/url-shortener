package tests

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"path"
	"testing"
	"url-shortener/internal/http-server/logger/handlers/url/save"
	"url-shortener/internal/lib/random"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.GenerateRandomString(10),
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")

}

func TestURLShortener_SaveRedirectRemove(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url",
			alias: gofakeit.Word(),
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			req := e.POST("/url").
				WithJSON(save.Request{
					URL:   testCase.url,
					Alias: testCase.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().
				Status(http.StatusOK).
				JSON().
				Object()

			if testCase.error != "" {
				req.NotContainsKey("alias")

				req.Value("error").String().IsEqual(testCase.error)
			}

			alias := testCase.alias

			if testCase.alias != "" {
				req.Value("alias").String().IsEqual(testCase.alias)
			} else {
				req.Value("alias").String().NotEmpty()

				alias = req.Value("alias").String().Raw()
			}

			//testRedirect(t, req, alias)

			reqDelete := e.DELETE("/"+path.Join("url", alias)).
				WithBasicAuth("myuser", "mypass").
				Expect().
				Status(http.StatusOK).
				JSON().
				Object()

			reqDelete.Value("status").String().IsEqual("OK")

			//	testRedirectNotFound(t, alias)
		})
	}

}
