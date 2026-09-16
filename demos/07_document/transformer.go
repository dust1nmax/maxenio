package main

import (
	"context"
	"log"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main(){
	ctx := context.Background()

	// 加载 .env 环境变量文件，获取 ARK_API_KEY 和 EMBEDDER 配置
	err := godotenv.Load("/home/max/maxenio/.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	transformer, err:= markdown.NewHeaderSplitter(
		ctx, 
		&markdown.HeaderConfig{
			Headers: map[string]string{
				"#": "h1",
				"##": "h2",
				"###": "h3",	
			},
			TrimHeaders: true,
		},
	)
	if err != nil {
        log.Fatalf("创建拆分器失败: %v", err)
    }

	doc := &schema.Document{
		Content: "# 标题\n引言内容\n## 第一节\n章节内容",
	}
	
	splitDocs, err:= transformer.Transform(
		ctx, 
		[]*schema.Document{doc},
	)

    if err != nil {
        log.Fatalf("转换失败: %v", err)
    }
    
    for _, doc := range splitDocs {
        log.Printf("内容: %s, 元数据: %v\n", doc.Content, doc.MetaData)
    }
}