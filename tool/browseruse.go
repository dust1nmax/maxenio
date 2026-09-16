// browseruse 让浏览器跳转到 curl


package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/tool/browseruse"
)

func newbrowseruse(){
	ctx := context.Background()
	but, err := browseruse.NewBrowserUseTool(ctx, &browseruse.Config{})
	if err != nil {
		log.Fatal(err)
	}

	url := "https://www.google.com"
	result, err := but.Execute(&browseruse.Param{
		Action: browseruse.ActionGoToURL,
		URL:    &url,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	time.Sleep(10 * time.Second)
	but.Cleanup()
}