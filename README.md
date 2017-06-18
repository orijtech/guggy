# guggy

Sample usage:
```go
package main

import (
	"fmt"
	"log"

	"github.com/orijtech/guggy/v2"
)

func main() {
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
```
