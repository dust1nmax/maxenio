package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// NewArkIndexer 创建 Milvus Indexer，负责把 Document 写入向量库。
//
// Indexer 内部持有 embedder，Store 时会自动把 Document.Content 转成向量，
// 所以调用方不需要自己算向量，直接交 Document 就行。
//
// address 是 Milvus 服务地址（例如 localhost:19530），由外部传入而不是写死，
// 方便在本地直连和容器内访问之间切换。
func NewArkIndexer(ctx context.Context, embedder *ark.Embedder, address string) (*milvus2.Indexer, error) {
	// 初始化 Milvus Indexer，用于将文档向量存储到 Milvus 向量数据库
	indexer, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: address,   // Milvus 服务地址和端口
			DBName:  "MaxEino", // 目标数据库名称
		},
		Collection: "test", // Milvus 集合名称

		// 向量配置：Dimension 必须与 embedder 的实际输出维度一致
		// （doubao-embedding 输出 2048 维，填错会在写入时报维度不匹配）
		Vector: &milvus2.VectorConfig{
			Dimension:  2048,
			MetricType: milvus2.COSINE, // 相似度度量：余弦
		},
		// 关联 Embedder，Indexer 会自动使用它对文档内容进行向量化
		Embedding: embedder,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 indexer 失败: %w", err)
	}

	return indexer, nil
}
