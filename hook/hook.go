package hook

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type Hooker struct {
	Transport *Transport
}

func (h *Hooker) GetRequest() *HookRequest {
	return h.Transport.Request
}

func (h *Hooker) GetResponse() *HookResponse {
	return h.Transport.Response
}

func NewHooker(t *testing.T) *Hooker {
	t.Helper()

	return &Hooker{
		Transport: &Transport{},
	}
}

type HookRequest struct {
	URL    *url.URL
	Method string
	Header map[string][]string
	Body   []byte
}

type HookResponse struct {
	Header map[string][]string
	Body   []byte
}

var _ http.RoundTripper = &Transport{}

type Transport struct {
	Transport http.RoundTripper
	Request   *HookRequest
	Response  *HookResponse
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.copyRequest(req)

	baseRoundTripper := t.Transport
	if baseRoundTripper == nil {
		baseRoundTripper = http.DefaultTransport
	}

	resp, err := baseRoundTripper.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	t.copyResponse(resp)

	return resp, nil
}

func (t *Transport) copyRequest(req *http.Request) {
	hr := &HookRequest{
		URL:    req.URL,
		Method: req.Method,
		Header: req.Header,
	}

	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		hr.Body = b

		req.Body = ioutil.NopCloser(bytes.NewReader(b))
	}

	t.Request = hr
}

func (t *Transport) copyResponse(resp *http.Response) {
	hr := &HookResponse{
		Header: resp.Header,
		Body:   nil,
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	hr.Body = b

	resp.Body = ioutil.NopCloser(bytes.NewReader(b))

	t.Response = hr
}
