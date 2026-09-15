package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

func main() {

	ctx := context.Background()

	// 加载 .env 环境变量文件，获取 ARK_API_KEY 和 EMBEDDER 配置
	err := godotenv.Load("/home/max/maxenio/.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}


	timeout := 30 * time.Second

	apiType := ark.APITypeMultiModal

	// 创建 Embedder 实例，配置 API Key、模型名称、超时时间和 API 类型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("EMBEDDER"),
		Timeout: &timeout,
		APIType: &apiType,
	})
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	retriever, err := milvus2.NewRetriever(
		ctx,
		&milvus2.RetrieverConfig{
			ClientConfig: &milvusclient.ClientConfig{
				Address: "localhost:19530",
				DBName:  "MaxEino",
			},
			Collection: "test",
			OutputFields: []string{
				"id",
				"content",
				"metadata",
			},
			TopK: 3,
			SearchMode: search_mode.NewApproximate(milvus2.COSINE),
			Embedding:  embedder,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create retriever: %v", err)
		return
	}
	log.Printf("Retriever created successfully")

	// document, err := retriever.Retrieve(ctx, "search query")
	// if err != nil {
	// 	panic(err)
	// }

	// for i, doc := range document {
	// 	fmt.Printf("Document %d:\n", i)
	// 	fmt.Printf("  ID: %s\n", doc.ID)
	// 	fmt.Printf("  Content: %s\n", doc.Content)
	// 	fmt.Printf("  Score: %v\n", doc.Score())
	// }

	resules, err := retriever.Retrieve(ctx, "可乐")
	if err != nil {
		panic(err)
	}

	fmt.Println(resules[0].Content)
	fmt.Println(resules[1].Content)
}
