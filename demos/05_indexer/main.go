package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	milvus2 "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

var collection = "MaxEino"

func IndexerRAG() {
	ctx := context.Background()

	// 加载环境变量
	err := godotenv.Load("/home/max/maxenio/.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	// 初始化 Embedding
	timeout := 30 * time.Second

	apiType := ark.APITypeMultiModal

	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("EMBEDDER"),
		Timeout: &timeout,
		APIType: &apiType,
	})
	if err != nil {
		log.Fatalf("failed to create embedder: %v", err)
	}

	// 初始化 Milvus Indexer
	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: "localhost:19530",
			DBName: "MaxEino",
		},
		Collection: "test",

		Vector: &milvus2.VectorConfig{
			Dimension: 2048,
			MetricType: milvus2.COSINE,
		},

		Embedding: embedder,
	})
	if err != nil {
		log.Fatalf("failed to create indexer: %v", err)
	}

	// 要存储的文档
	docs := []*schema.Document{
		{
			ID:      "1",
			Content: "今天是2026年9月15日17：05",
			MetaData: map[string]any{
				"author": "max",
			},
		},
	}

	// 存储
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("failed to store documents: %v", err)
	}

	log.Printf("Stored documents with IDs: %v", ids)
}

func main() {
	IndexerRAG()
}