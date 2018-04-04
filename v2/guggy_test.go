// Copyright 2017 orijtech. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package guggy_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/orijtech/guggy/v2"
)

func TestSearch(t *testing.T) {
	client, err := guggy.NewClient(testAPIKey1)
	if err != nil {
		t.Fatal(err)
	}
	client.SetHTTPRoundTripper(&backend{route: searchRoute})

	tests := [...]struct {
		req     *guggy.Request
		wantErr bool
	}{
		0: {wantErr: true},
		1: {req: &guggy.Request{Query: "say what?"}},
	}

	for i, tt := range tests {
		res, err := client.Search(context.Background(), tt.req)
		if tt.wantErr {
			if err == nil {
				t.Errorf("#%d: expected non-nil error", i)
			}
			continue
		}

		if err != nil {
			t.Errorf("#%d gotErr: %v", i, err)
			continue
		}
		if res == nil {
			t.Errorf("#%d expected a non-blank response", i)
			continue
		}
	}
}

var errUnimplemented = errors.New("unimplemented")

type backend struct {
	route string
}

func (b *backend) RoundTrip(req *http.Request) (*http.Response, error) {
	switch b.route {
	case searchRoute:
		return b.searchRoundTrip(req)
	default:
		return nil, errUnimplemented
	}
}

func makeResp(status string, code int, body io.ReadCloser) *http.Response {
	return &http.Response{
		Status:     status,
		StatusCode: code,
		Body:       body,
		Header:     make(http.Header),
	}
}

func checkAuthAndMethod(req *http.Request, wantMethod string) (*http.Response, error) {
	if req.Method != wantMethod {
		return makeResp(fmt.Sprintf("got method %q want %q", req.Method, wantMethod), http.StatusMethodNotAllowed, nil), nil
	}
	apiKey := strings.TrimSpace(req.Header.Get("apiKey"))
	if apiKey == "" {
		return makeResp(`expected "apiKey" in the header`, http.StatusBadRequest, nil), nil
	}
	switch apiKey {
	case testAPIKey1, testAPIKey2:
		return nil, nil
	default:
		return makeResp("unauthorized API key", http.StatusUnauthorized, nil), nil
	}
}

var blankRequest = new(guggy.Request)

func (b *backend) searchRoundTrip(req *http.Request) (*http.Response, error) {
	if badAuthResp, err := checkAuthAndMethod(req, "POST"); badAuthResp != nil || err != nil {
		return badAuthResp, err
	}

	defer req.Body.Close()
	slurp, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return makeResp(err.Error(), http.StatusBadRequest, nil), nil
	}
	if len(slurp) < 6 {
		return makeResp("expecting a non-blank Request", http.StatusBadRequest, nil), nil
	}
	f, err := os.Open("./testdata/search-0.json")
	if err != nil {
		return makeResp(err.Error(), http.StatusBadRequest, nil), nil
	}
	return makeResp("200 OK", http.StatusOK, f), nil
}

const (
	testAPIKey1 = "test-api-key-1"
	testAPIKey2 = "test-api-key-2"

	searchRoute = "/search"
)
