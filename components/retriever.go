package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// NewArkRetriever 创建 Milvus Retriever，负责把 Query 向量化后在向量库里检索相似文档。
//
// Retriever 和 Indexer 共用同一个 embedder：只有用同一个 embedding 模型，
// 写入时的向量和检索时的向量才在同一个语义空间里，否则检索结果没有意义。
//
// TopK 控制返回条数；OutputFields 决定 Milvus 回传哪些字段（要拿 doc.Content
// 就必须把 "content" 列出来，否则取回来是空的）。
func NewArkRetriever(ctx context.Context, embedder *ark.Embedder) (*milvus2.Retriever, error) {
	retriever, err := milvus2.NewRetriever(
		ctx,
		&milvus2.RetrieverConfig{
			ClientConfig: &milvusclient.ClientConfig{
				// 注意：这里的地址是写死的，和 NewArkIndexer 由外部传入 address 不一致，
				// 需要跟 indexer 连同一个 Milvus 实例
				Address: "localhost:19530",
				DBName:  "MaxEino",
			},
			Collection: "test",
			OutputFields: []string{
				"id",
				"content",
				"metadata",
			},
			TopK:       3,
			SearchMode: search_mode.NewApproximate(milvus2.COSINE),
			Embedding:  embedder,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("创建 retriever 失败: %w", err)
	}

	return retriever, nil
}
