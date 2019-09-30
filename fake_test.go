package fakestorage_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/sinmetal/fakestorage"
	"google.golang.org/api/option"
)

func TestGetObject(t *testing.T) {
	ctx := context.Background()

	faker := fakestorage.NewFaker(t)

	stg, err := storage.NewClient(ctx, option.WithHTTPClient(faker.Client))
	if err != nil {
		t.Fatal(err)
	}

	reader, err := stg.Bucket("hoge").Object("hoge.txt").NewReader(ctx)
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
