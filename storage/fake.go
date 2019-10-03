package storage

import (
	"errors"
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
	fake, err := tran.fakeResponses.Get(req.URL.String(), req.Method)
	if err != nil {
		tran.t.Fatal("unexpected: ", err)
	}
	return fake, nil
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

func (f *fakeResponses) Get(url string, method string) (*http.Response, error) {
	v, ok := f.m[f.key(url, method)]
	if !ok {
		// TODO 適当なやつを返す
		switch method {
		case http.MethodGet:
			return GetObjectOKResponse(), nil
		case http.MethodPost:
			return PostObjectOKResponse(), nil
		default:
			return nil, errors.New("unsupported method")
		}
	}
	return v, nil
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

func PostObjectOKResponse() *http.Response {
	header := make(map[string][]string)
	header["Server"] = []string{"UploadServer"}
	header["Alt-Svc"] = []string{`quic=":443"; ma=2592000; v="46,43",h3-Q048=":443"; ma=2592000,h3-Q046=":443"; ma=2592000,h3-Q043=":443"; ma=2592000`}
	header["Vary"] = []string{"Origin"}
	header["Vary"] = []string{"X-Origin"}
	header["Content-Type"] = []string{"application/json; charset=UTF-8"}
	header["Cache-Control"] = []string{"no-cache, no-store, max-age=0, must-revalidate"}
	header["Date"] = []string{"Thu, 03 Oct 2019 07:31:44 GMT"}
	header["Content-Length"] = []string{"2324"}
	header["X-Guploader-Uploadid"] = []string{"AEnB2UpF0rRDJSlY8seVYqjxCchiX2GwvYwiGqkFfaduXRlzuNpGEDdlCsKtpvVe5gn0WMsW3HSeqFw4nqyNZ0v3apu9Il_VMw"}
	header["Etag"] = []string{"CMXdo57J/+QCEAE="}
	header["Pragma"] = []string{"no-cache"}
	header["Expires"] = []string{"Mon, 01 Jan 1990 00:00:00 GMT"}
	r := ioutil.NopCloser(strings.NewReader(`eyJraW5kIjoic3RvcmFnZSNvYmplY3QiLCJpZCI6ImhvZ2UvcG9zdC50eHQvMTU3MDA4NzkwNDAxNDAyMSIsInNlbGZMaW5rIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vc3RvcmFnZS92MS9iL2hvZ2Uvby9wb3N0LnR4dCIsIm5hbWUiOiJwb3N0LnR4dCIsImJ1Y2tldCI6ImhvZ2UiLCJnZW5lcmF0aW9uIjoiMTU3MDA4NzkwNDAxNDAyMSIsIm1ldGFnZW5lcmF0aW9uIjoiMSIsImNvbnRlbnRUeXBlIjoidGV4dC9wbGFpbjsgY2hhcnNldD11dGYtOCIsInRpbWVDcmVhdGVkIjoiMjAxOS0xMC0wM1QwNzozMTo0NC4wMTNaIiwidXBkYXRlZCI6IjIwMTktMTAtMDNUMDc6MzE6NDQuMDEzWiIsInN0b3JhZ2VDbGFzcyI6IlJFR0lPTkFMIiwidGltZVN0b3JhZ2VDbGFzc1VwZGF0ZWQiOiIyMDE5LTEwLTAzVDA3OjMxOjQ0LjAxM1oiLCJzaXplIjoiMjQiLCJtZDVIYXNoIjoiM2Z2MFZYSGprM25DYzN6blZOcmNSdz09IiwibWVkaWFMaW5rIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vZG93bmxvYWQvc3RvcmFnZS92MS9iL2hvZ2Uvby9wb3N0LnR4dD9nZW5lcmF0aW9uPTE1NzAwODc5MDQwMTQwMjEmYWx0PW1lZGlhIiwiYWNsIjpbeyJraW5kIjoic3RvcmFnZSNvYmplY3RBY2Nlc3NDb250cm9sIiwiaWQiOiJob2dlL3Bvc3QudHh0LzE1NzAwODc5MDQwMTQwMjEvcHJvamVjdC1vd25lcnMtMTY4NjEwOTE2ODAxIiwic2VsZkxpbmsiOiJodHRwczovL3d3dy5nb29nbGVhcGlzLmNvbS9zdG9yYWdlL3YxL2IvaG9nZS9vL3Bvc3QudHh0L2FjbC9wcm9qZWN0LW93bmVycy0xNjg2MTA5MTY4MDEiLCJidWNrZXQiOiJob2dlIiwib2JqZWN0IjoicG9zdC50eHQiLCJnZW5lcmF0aW9uIjoiMTU3MDA4NzkwNDAxNDAyMSIsImVudGl0eSI6InByb2plY3Qtb3duZXJzLTE2ODYxMDkxNjgwMSIsInJvbGUiOiJPV05FUiIsInByb2plY3RUZWFtIjp7InByb2plY3ROdW1iZXIiOiIxNjg2MTA5MTY4MDEiLCJ0ZWFtIjoib3duZXJzIn0sImV0YWciOiJDTVhkbzU3Si8rUUNFQUU9In0seyJraW5kIjoic3RvcmFnZSNvYmplY3RBY2Nlc3NDb250cm9sIiwiaWQiOiJob2dlL3Bvc3QudHh0LzE1NzAwODc5MDQwMTQwMjEvcHJvamVjdC1lZGl0b3JzLTE2ODYxMDkxNjgwMSIsInNlbGZMaW5rIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vc3RvcmFnZS92MS9iL2hvZ2Uvby9wb3N0LnR4dC9hY2wvcHJvamVjdC1lZGl0b3JzLTE2ODYxMDkxNjgwMSIsImJ1Y2tldCI6ImhvZ2UiLCJvYmplY3QiOiJwb3N0LnR4dCIsImdlbmVyYXRpb24iOiIxNTcwMDg3OTA0MDE0MDIxIiwiZW50aXR5IjoicHJvamVjdC1lZGl0b3JzLTE2ODYxMDkxNjgwMSIsInJvbGUiOiJPV05FUiIsInByb2plY3RUZWFtIjp7InByb2plY3ROdW1iZXIiOiIxNjg2MTA5MTY4MDEiLCJ0ZWFtIjoiZWRpdG9ycyJ9LCJldGFnIjoiQ01YZG81N0ovK1FDRUFFPSJ9LHsia2luZCI6InN0b3JhZ2Ujb2JqZWN0QWNjZXNzQ29udHJvbCIsImlkIjoiaG9nZS9wb3N0LnR4dC8xNTcwMDg3OTA0MDE0MDIxL3Byb2plY3Qtdmlld2Vycy0xNjg2MTA5MTY4MDEiLCJzZWxmTGluayI6Imh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3N0b3JhZ2UvdjEvYi9ob2dlL28vcG9zdC50eHQvYWNsL3Byb2plY3Qtdmlld2Vycy0xNjg2MTA5MTY4MDEiLCJidWNrZXQiOiJob2dlIiwib2JqZWN0IjoicG9zdC50eHQiLCJnZW5lcmF0aW9uIjoiMTU3MDA4NzkwNDAxNDAyMSIsImVudGl0eSI6InByb2plY3Qtdmlld2Vycy0xNjg2MTA5MTY4MDEiLCJyb2xlIjoiUkVBREVSIiwicHJvamVjdFRlYW0iOnsicHJvamVjdE51bWJlciI6IjE2ODYxMDkxNjgwMSIsInRlYW0iOiJ2aWV3ZXJzIn0sImV0YWciOiJDTVhkbzU3Si8rUUNFQUU9In0seyJraW5kIjoic3RvcmFnZSNvYmplY3RBY2Nlc3NDb250cm9sIiwiaWQiOiJob2dlL3Bvc3QudHh0LzE1NzAwODc5MDQwMTQwMjEvdXNlci1zaW5tZXRhbEBtZXJjYXJpLmNvbSIsInNlbGZMaW5rIjoiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vc3RvcmFnZS92MS9iL2hvZ2Uvby9wb3N0LnR4dC9hY2wvdXNlci1zaW5tZXRhbEBtZXJjYXJpLmNvbSIsImJ1Y2tldCI6ImhvZ2UiLCJvYmplY3QiOiJwb3N0LnR4dCIsImdlbmVyYXRpb24iOiIxNTcwMDg3OTA0MDE0MDIxIiwiZW50aXR5IjoidXNlci1zaW5tZXRhbEBtZXJjYXJpLmNvbSIsInJvbGUiOiJPV05FUiIsImVtYWlsIjoic2lubWV0YWxAbWVyY2FyaS5jb20iLCJldGFnIjoiQ01YZG81N0ovK1FDRUFFPSJ9XSwib3duZXIiOnsiZW50aXR5IjoidXNlci1zaW5tZXRhbEBtZXJjYXJpLmNvbSJ9LCJjcmMzMmMiOiJ2T011NVE9PSIsImV0YWciOiJDTVhkbzU3Si8rUUNFQUU9In0=`))

	return &http.Response{
		Status:        "200 OK",
		StatusCode:    http.StatusOK,
		Header:        header,
		Body:          r,
		ContentLength: 2324,
	}
}
