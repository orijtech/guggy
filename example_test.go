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
	"fmt"
	"log"

	"github.com/orijtech/guggy/v2"
)

func Example_client_Search() {
	client, err := guggy.NewClientFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	res, err := client.Search(&guggy.Request{
		Query: "What is the weather up there?",
	})
	if err != nil {
		log.Fatal(err)
	}

	for i, gif := range res.Gifs {
		fmt.Printf("GIF: %d %#v\n", i, gif)
	}

	for i, sticker := range res.Stickers {
		fmt.Printf("Sticker: %d %#v\n", i, sticker)
	}
}