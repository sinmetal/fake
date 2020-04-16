package storage_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/sinmetal/fake/hook/hars"
	"github.com/vvakame/go-harlog"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"

	. "github.com/sinmetal/fake/storage"
)

func TestGetObject(t *testing.T) {
	ctx := context.Background()

	faker := NewFaker(t)

	stg, err := storage.NewClient(ctx, option.WithHTTPClient(faker.Client))
	if err != nil {
		t.Fatal(err)
	}

	reader, err := stg.Bucket("sinmetal-ci-fake").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := reader.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()
	body, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v", string(body))
}

func TestRealGetObjectHar(t *testing.T) {
	ctx := context.Background()

	hc, err := google.DefaultClient(ctx, storage.ScopeReadWrite)
	if err != nil {
		panic(err)
	}

	// inject HAR logger!
	har := &harlog.Transport{
		Transport: hc.Transport,
	}
	hc.Transport = har
	stg, err := storage.NewClient(ctx, option.WithHTTPClient(hc))
	if err != nil {
		t.Fatal(err)
	}

	_, err = stg.Bucket("sinmetal-ci-fake").Object("hoge.txt").NewReader(ctx)
	if err != nil {
		t.Fatal(err)
	}

	hars.Compare(t, "object.get.har.golden", har.HAR())
}

func TestPostObjectHar(t *testing.T) {
	ctx := context.Background()

	hc, err := google.DefaultClient(ctx, storage.ScopeReadWrite)
	if err != nil {
		panic(err)
	}

	// inject HAR logger!
	har := &harlog.Transport{
		Transport: hc.Transport,
	}
	hc.Transport = har
	stg, err := storage.NewClient(ctx, option.WithHTTPClient(hc))
	if err != nil {
		t.Fatal(err)
	}

	w := stg.Bucket("sinmetal-ci-fake").Object("post.txt").NewWriter(ctx)
	_, err = w.Write([]byte(`{"message":"hello fake"}`))
	if err != nil {
		t.Fatal("unexpected: ", err)
	}
	w.ContentType = "application/json"
	if err := w.Close(); err != nil {
		t.Fatal("unexpected: ", err)
	}

	hars.Compare(t, "object.post.har.golden", har.HAR())
}

func TestPostObjectHarToCode(t *testing.T) {
	hars.LogFakeResponseCode(t, "object.post.har.golden")
}
