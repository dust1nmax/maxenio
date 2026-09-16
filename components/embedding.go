package components

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
)

// NewEmbedder 创建 Ark Embedder，负责把文本转成向量。
//
// 依赖环境变量 ARK_API_KEY 和 EMBEDDER，所以必须在加载 .env 之后调用。
// 它是 Indexer 和 Retriever 的共同依赖（两边共用同一个 embedder 做向量化），
// 失败时返回 error 交给调用方，而不是自己终止进程——否则上层拿不到失败原因，
// 也没有机会重试或降级。
func NewEmbedder(ctx context.Context) (*ark.Embedder, error) {
	// 单次请求的超时时间
	timeout := 30 * time.Second

	// 多模态 API 类型，除文本外还能接受其他输入形式
	apiType := ark.APITypeMultiModal

	// 创建 Embedder 实例，配置 API Key、模型名称、超时时间和 API 类型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  os.Getenv("ARK_API_KEY"),
		Model:   os.Getenv("EMBEDDER"),
		Timeout: &timeout,
		APIType: &apiType,
	})
	if err != nil {
		// 用 %w 包装而不是 %v，保留原始错误链，
		// 调用方仍可用 errors.Is / errors.As 判断底层原因
		return nil, fmt.Errorf("创建 embedder 失败: %w", err)
	}

	return embedder, nil
}

// EmbedTexts 把一批文本转成向量，返回顺序与输入一致，类型是 [][]float64。
//
// 写入和检索时不需要手动调它：Indexer 会把 Document.Content、Retriever 会把
// 查询词各自交给 embedder 向量化。手动调用主要是用来验证配置，
// 比如 len(vectors[0]) 就是模型实际输出的向量维度，可以用它确认
// IndexerConfig 里的 Dimension 有没有填对（doubao-embedding 应该是 2048）。
func EmbedTexts(ctx context.Context, embedder *ark.Embedder, texts []string) ([][]float64, error) {
	vectors, err := embedder.EmbedStrings(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("文本转向量失败: %w", err)
	}

	return vectors, nil
}
