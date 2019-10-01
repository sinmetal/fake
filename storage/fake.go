package storage

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"testing"
)

type Faker struct {
	transport *Transport
	Client    *http.Client
}

func NewFaker(t *testing.T) *Faker {
	t.Helper()

	transport := &Transport{
		DropRequest:  true,
		FakeResponse: defaultFakeResponse(),
	}
	return &Faker{
		transport: transport,
		Client: &http.Client{
			Transport: transport,
		},
	}
}

func defaultFakeResponse() *http.Response {
	header := make(map[string][]string)
	header["Accept-Ranges"] = []string{"bytes"}
	header["Age"] = []string{"268"}
	header["Alt-Svc"] = []string{`quic=":443"; ma=2592000; v="46,43",h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000`}
	header["Cache-Control"] = []string{"public", "max-age=3600"}
	header["Content-Length"] = []string{"25"}
	header["Content-Type"] = []string{"text/plain"}
	header["Date"] = []string{"Mon, 30 Sep 2019 10:23:16 GMT"}
	header["Etag"] = []string{"c4d22707e0d79bd01e33fe19a5e21487"}
	header["Expires"] = []string{"Mon, 30 Sep 2019 11:23:16 GMT"}
	header["Last-Modified"] = []string{"Mon, 30 Sep 2019 10:01:47 GMT"}
	header["X-Goog-Generation"] = []string{"1569837707444808"}
	header["X-Goog-Hash"] = []string{"crc32c=CrEDEg== md5=xNInB+DXm9AeM/4ZpeIUhw=="}
	header["X-Goog-Metageneration"] = []string{"2"}
	header["X-Goog-Storage-Class"] = []string{"REGIONAL"}
	header["X-Goog-Stored-Content-Encoding"] = []string{"identity"}
	header["X-Goog-Stored-Content-Length"] = []string{"25"}
	header["X-Guploader-Uploadid"] = []string{"AEnB2UoygSa1dB8aXstLosALQoifLpXnQ5kIx_lyzTyIvk5bFuIcG7nqk-sR5GdihmWdTtHDuiKCtSgxyRJ9iLJmHnQ7RHmvoQ"}

	r := ioutil.NopCloser(strings.NewReader(`{"message":"Hello Hoge"}`))

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Header:        header,
		Body:          r,
		ContentLength: 25,
	}
}

func NewFakeHTTPClient(t *testing.T) *http.Client {
	t.Helper()

	return &http.Client{
		Transport: &Transport{},
	}
}

// Transport implements http.RoundTripper. When set as Transport of http.Client, it executes HTTP requests with logging.
// No field is mandatory.
type Transport struct {
	Transport    http.RoundTripper
	LogRequest   func(req *http.Request)
	LogResponse  func(resp *http.Response)
	DropRequest  bool
	FakeResponse *http.Response
}

var DefaultTransport = &Transport{
	Transport: http.DefaultTransport,
}

var DefaultLogRequest = func(req *http.Request) {
	log.Printf("--> %s %s", req.Method, req.URL)
}

var DefaultLogResponse = func(resp *http.Response) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	log.Printf("fakestorage : response status = %v, header = %v, body %s", resp.Status, resp.Header, string(body))
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.DropRequest {
		return t.FakeResponse, nil
	}
	t.logRequest(req)

	resp, err := t.transport().RoundTrip(req)
	if err != nil {
		return resp, err
	}

	t.logResponse(resp)

	if t.FakeResponse != nil {
		fmt.Println("Fake Response !!!")
		return t.FakeResponse, nil
	}
	return resp, err
}

func (t *Transport) logRequest(req *http.Request) {
	if t.LogRequest != nil {
		t.LogRequest(req)
	} else {
		DefaultLogRequest(req)
	}
}

func (t *Transport) logResponse(resp *http.Response) {
	if t.LogResponse != nil {
		t.LogResponse(resp)
	} else {
		DefaultLogResponse(resp)
	}
}

func (t *Transport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}

	return http.DefaultTransport
}
