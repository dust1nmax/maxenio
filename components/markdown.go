package components

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/components/document"
)

// NewTransform 创建 Markdown 标题切分器，把长文档按 # / ## / ### 拆成多片。
//
// Headers 的 key 只能由 '#' 组成，value 是写进每片 MetaData 的字段名；
// TrimHeaders=false 表示标题行保留在片内容里（true 则只留正文）。
//
// 注意一个容易踩的坑：切分器默认把所有片段的 ID 都设成原文档 ID
// （header.go 里的 defaultIDGenerator 直接返回原 ID），所以 10 个片段会共用
// 同一个 ID。写入 Milvus 时 ID 就是主键，后写入的片段会覆盖先写入的，
// 最后库里只剩一片。入库前需要配置 HeaderConfig.IDGenerator 生成唯一 ID。
func NewTransform(ctx context.Context) (document.Transformer, error) {
	transform, err := markdown.NewHeaderSplitter(
		ctx,
		&markdown.HeaderConfig{
			Headers: map[string]string{
				"#":   "h1",
				"##":  "h2",
				"###": "h3",
			},
			TrimHeaders: false,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("创建 markdown 切分器失败: %w", err)
	}

	return transform, nil
}
