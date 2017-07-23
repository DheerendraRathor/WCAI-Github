package httpClient

import "net/http"

type GithubTransport struct {
	http.RoundTripper
	ClientId     string
	ClientSecret string
}

func (gt *GithubTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	q := req.URL.Query()
	q.Add("client_id", gt.ClientId)
	q.Add("client_secret", gt.ClientSecret)
	req.URL.RawQuery = q.Encode()

	return gt.RoundTripper.RoundTrip(req)
}
