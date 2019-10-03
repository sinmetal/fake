package storage

import (
	"fmt"
	"io/ioutil"
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
		fakeResponses: &fakeResponses{
			m: make(map[string]*http.Response),
		},
	}
	return &Faker{
		transport: transport,
		Client: &http.Client{
			Transport: transport,
		},
	}
}

var _ http.RoundTripper = &Transport{}

type Transport struct {
	t             *testing.T
	Transport     http.RoundTripper
	fakeResponses *fakeResponses
}

func (tran *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return tran.fakeResponses.Get(req.URL.String(), req.Method), nil
}

type fakeResponses struct {
	m map[string]*http.Response
}

func (f *fakeResponses) key(url string, method string) string {
	return fmt.Sprintf("%s-%s", url, strings.ToUpper(method))
}

func (f *fakeResponses) Add(url string, method string, response *http.Response) {
	f.m[f.key(url, method)] = response
}

func (f *fakeResponses) Get(url string, method string) *http.Response {
	v, ok := f.m[f.key(url, method)]
	if !ok {
		// TODO 適当なやつを返す
		return GetObjectOKResponse()
	}
	return v
}

func GetObjectOKResponse() *http.Response {
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
