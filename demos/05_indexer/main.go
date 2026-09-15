package main

// 该示例演示了如何使用 Milvus2 Indexer 和 Ark Embedder 构建 RAG 索引系统
// 主要功能：创建嵌入向量并存储文档到 Milvus 向量数据库

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

// collection 定义 Milvus 集合名称
var collection = "MaxEino"

func IndexerRAG() {
	// 创建上下文，用于控制超时和取消操作
	ctx := context.Background()

	// 加载 .env 环境变量文件，获取 ARK_API_KEY 和 EMBEDDER 配置
	err := godotenv.Load("/home/max/maxenio/.env")
	if err != nil {
		log.Fatalf("failed to load .env: %v", err)
	}

	// 初始化 Ark Embedder，用于将文本内容转换为向量嵌入
	// 超时设置为 30 秒，APIType 使用多模态模式支持多种输入类型
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

	// 初始化 Milvus Indexer，用于将文档向量存储到 Milvus 向量数据库
	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: "localhost:19530", // Milvus 服务地址和端口
			DBName: "MaxEino",          // 目标数据库名称
		},
		Collection: "test", // Milvus 集合名称

		// 向量配置：Dimension 指定向量维度为 2048，MetricType 使用余弦相似度
		Vector: &milvus2.VectorConfig{
			Dimension: 2048,
			MetricType: milvus2.COSINE,
		},

		// 关联 Embedder，Indexer 会自动使用它对文档内容进行向量化
		Embedding: embedder,
	})
	if err != nil {
		log.Fatalf("failed to create indexer: %v", err)
	}

	// 准备要存储的文档，包含文档 ID、内容和元数据信息
	docs := []*schema.Document{
		{
			ID:      "1", // 文档唯一标识符
			Content: "今天是2026年9月15日17：05", // 文档内容
			MetaData: map[string]any{
				"author": "max", // 文档作者元数据
			},
		},
	}

	// 调用 Indexer.Store 将文档向量化并存储到 Milvus，返回存储后的文档 ID
	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatalf("failed to store documents: %v", err)
	}

	log.Printf("Stored documents with IDs: %v", ids)
}

// main 程序入口函数
func main() {
	IndexerRAG()
}