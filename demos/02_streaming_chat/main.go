package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件中的环境变量
	_ = godotenv.Load(".env")

	ctx := context.Background()

	// 1. 初始化 ChatModel
	model, err := ark.NewChatModel(
		ctx,
		&ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_API_KEY"),
			Model:  os.Getenv("MODEL"),
		})
	if err != nil {
		panic(err)
	}

	// 2. 准备聊天消息（直接写死，不用模板）
	messages := []*schema.Message{
		schema.SystemMessage("你是一个助手"),
		schema.UserMessage("golang的优势是什么？"),
	}

	// 3. 使用 Stream 流式调用（注意：ark 模型可能不支持 Stream，需实际测试）
	reader, err := model.Stream(ctx, messages)
	if err != nil {
		panic(err)
	}
	defer reader.Close() // 记得关闭 reader

	// 4. 循环接收流式输出
	for {
		chunk, err := reader.Recv()
		if err != nil {
			break // 流结束或出错时退出循环
		}
		// 5. 打印每个 chunk（增量输出，体验更好）
		print(chunk.Content)
	}
}