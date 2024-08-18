package cloudflare

import (
	"net/http"

	"github.com/kamontat/cloudflare-exporter/configs"
	"github.com/spf13/viper"
)

type transport struct {
	token    string
	apiEmail string
	apiKey   string
	wrapped  http.RoundTripper
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := cloneRequest(req)

	if t.token != "" {
		req2.Header.Set("Authorization", "Bearer "+t.token)
	} else {
		req2.Header.Set("X-AUTH-EMAIL", t.apiEmail)
		req2.Header.Set("X-AUTH-KEY", t.apiKey)
	}

	return t.wrapped.RoundTrip(req)
}

func NewHttpClient(config *viper.Viper) *http.Client {
	token := config.GetString(configs.CONF_CF_API_TOKEN)
	apiEmail := config.GetString(configs.CONF_CF_API_EMAIL)
	apiKey := config.GetString(configs.CONF_CF_API_KEY)

	return &http.Client{
		Timeout: config.GetDuration(configs.CONF_CF_TIMEOUT),
		Transport: &transport{
			token:    token,
			apiEmail: apiEmail,
			apiKey:   apiKey,
			wrapped:  http.DefaultTransport,
		},
	}
}

// cloneRequest returns a clone of the provided *http.Request.
// The clone is a shallow copy of the struct and its Header map.
func cloneRequest(r *http.Request) *http.Request {
	// shallow copy of the struct
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		r2.Header[k] = append([]string(nil), s...)
	}
	return r2
}
