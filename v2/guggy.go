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

package guggy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"

	"github.com/orijtech/otils"
)

const (
	baseURL = "https://text2gif.guggy.com/v2"
)

type Client struct {
	sync.RWMutex
	apiKey string

	rt http.RoundTripper
}

func (c *Client) httpClient() *http.Client {
	c.RLock()
	rt := c.rt
	c.RUnlock()

	return &http.Client{
		Transport: rt,
	}
}

func (c *Client) SetHTTPRoundTripper(rt http.RoundTripper) {
	c.Lock()
	c.rt = rt
	c.Unlock()
}

func (c *Client) SetAPIKey(apiKey string) {
	c.Lock()
	c.apiKey = apiKey
	c.Unlock()
}

var (
	envAPIKeyKey      = "GUGGY_API_KEY"
	errBlankEnvAPIKey = fmt.Errorf("expected %q to have been set in your environment", envAPIKeyKey)
	errBlankAPIKey    = errors.New("expecting a non-blank apiKey")
)

func NewClientFromEnv() (*Client, error) {
	apiKey := strings.TrimSpace(os.Getenv(envAPIKeyKey))
	if apiKey == "" {
		return nil, errBlankEnvAPIKey
	}
	return &Client{apiKey: apiKey}, nil
}

func NewClient(apiKey string) (*Client, error) {
	if apiKey == "" {
		return nil, errBlankAPIKey
	}
	return &Client{apiKey: apiKey}, nil
}

type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Image struct {
	URL string `json:"secureUrl,omitempty"`

	Dimensions *Dimensions `json:"dimensions,omitempty"`
}

type Size struct {
	Original *Image `json:"original"`
	Preview  *Image `json:"preview"`

	LowQuality *Image `json:"lowQuality"`

	HighResolution *Image `json:"hires"`
}

type Collection struct {
	GIF  *Size `json:"gif,omitempty"`
	MP4  *Size `json:"mp4,omitempty"`
	PNG  *Size `json:"png,omitempty"`
	WEBP *Size `json:"webp,omitempty"`

	Original *Size `json:"original"`

	Thumbnail *Size `json:"thumbnail"`
}

type Response struct {
	RequestID string        `json:"reqId"`
	Stickers  []*Collection `json:"stickers"`
	Gifs      []*Collection `json:"animated"`
}

type Request struct {
	Query string `json:"query"`

	Language Language `json:"lang"`
}

func (c *Client) _apiKey() string {
	c.RLock()
	defer c.RUnlock()

	return c.apiKey
}

type request struct {
	Terms    string   `json:"sentence"`
	Language Language `json:"lang"`
}

var blankResponse = new(Response)

var (
	errBlankRequest  = errors.New("expecting a non-blank request")
	errBlankResponse = errors.New("server sent back a blank response")
)

func (c *Client) Search(req *Request) (*Response, error) {
	if req == nil {
		return nil, errBlankRequest
	}
	rreq := &request{Terms: req.Query, Language: req.Language}
	blob, err := json.Marshal(rreq)
	if err != nil {
		return nil, err
	}
	theURL := fmt.Sprintf("%s/guggify", baseURL)
	httpReq, err := http.NewRequest("POST", theURL, bytes.NewReader(blob))
	if err != nil {
		return nil, err
	}
	slurp, _, err := c.doAuthAndHTTPReq(httpReq)
	if err != nil {
		return nil, err
	}
	resp := new(Response)
	if err := json.Unmarshal(slurp, resp); err != nil {
		return nil, err
	}
	if reflect.DeepEqual(resp, blankResponse) {
		return nil, errBlankResponse
	}
	return resp, nil
}

func (c *Client) doAuthAndHTTPReq(req *http.Request) ([]byte, http.Header, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apiKey", c._apiKey())
	res, err := c.httpClient().Do(req)
	if err != nil {
		return nil, nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	if !otils.StatusOK(res.StatusCode) {
		return nil, res.Header, errors.New(res.Status)
	}
	slurp, err := ioutil.ReadAll(res.Body)
	return slurp, res.Header, err
}

// Languages
type Language string

const (
	LangSpanish            Language = "es"
	LangPortuguese         Language = "pt"
	LangIndonesian         Language = "id"
	LangFrench             Language = "fr"
	LangArabic             Language = "ar"
	LangTurkish            Language = "tr"
	LangThai               Language = "th"
	LangVietnamese         Language = "vi"
	LangGerman             Language = "de"
	LangItalian            Language = "it"
	LangJapanese           Language = "ja"
	LangChineseSimplified  Language = "zh-CN"
	LangChineseTraditional Language = "zh-TW"
	LangRussian            Language = "ru"
	LangKorean             Language = "ko"
	LangPolish             Language = "pl"
	LangDutch              Language = "nl"
	LangRomanian           Language = "ro"
	LangHungarian          Language = "hu"
	LangSwedish            Language = "sv"
	LangCzech              Language = "cs"
	LangHindi              Language = "hi"
	LangBengali            Language = "bn"
	LangDanish             Language = "da"
	LangFarsi              Language = "fa"
	LangFilipino           Language = "tl"
	LangFinnish            Language = "fi"
	LangHebrew             Language = "iw"
	LangMalay              Language = "ms"
	LangNorwegian          Language = "no"
	LangUkrainian          Language = "uk"
)
