package components

import (
	"context"
	"fmt"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
)

// NewarkModel 创建 Ark ChatModel，用于对话和最终的答案生成。
//
// 依赖环境变量 ARK_API_KEY 和 MODEL，必须在加载 .env 之后调用。
func NewarkModel(ctx context.Context) (*ark.ChatModel, error) {
	model, err := ark.NewChatModel(
		ctx,
		&ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_API_KEY"),
			Model:  os.Getenv("MODEL"),
		})
	if err != nil {
		return nil, fmt.Errorf("创建 chat model 失败: %w", err)
	}

	return model, nil
}
