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
	"fmt"
	"time"

	"github.com/orijtech/guggy/v2"
)

func Example_client_Search() {
	client, err := guggy.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.Search(ctx, &guggy.Request{
		Query: "What is the weather up there?",
	})
	if err != nil {
		panic(err)
	}

	for i, gif := range res.Gifs {
		fmt.Printf("GIF: %d MP4: %#v GIF: %#v\n", i, gif.MP4, gif.GIF.HighResolution)
	}

	for i, sticker := range res.Stickers {
		fmt.Printf("Sticker: %d MP4: %#v GIF: %#v\n", i, sticker.MP4, sticker.GIF)
	}

	// Output:
	// WW
}
