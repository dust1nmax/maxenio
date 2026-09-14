package main

import (
	"context"
	"os"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {
	// 加载 .env 文件中的环境变量（APK_API_KEY, MODEL）
	_ = godotenv.Load(".env")

	ctx := context.Background()

	// 1. 初始化 ChatModel（使用豆包/火山引擎的 ark 模型）
	model, err := ark.NewChatModel(
		ctx,
		&ark.ChatModelConfig{
			APIKey: os.Getenv("APK_API_KEY"),
			Model:  os.Getenv("MODEL"),
		})
	if err != nil {
		panic(err)
	}

	// 2. 定义 Prompt 模板（填空题形式）
	// schema.FString 表示使用 {variable} 语法进行变量替换
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"), // 系统消息模板，{role} 会替换
		&schema.Message{
			Role:    schema.User,
			Content: "回答{task}", // 用户消息模板，{task} 会替换
		},
	)

	// 3. 准备变量值
	params := map[string]any{
		"role": "编程助手",
		"task": "golang的优势",
	}

	// 4. 填充模板，得到最终 messages
	messages, err := template.Format(ctx, params)
	if err != nil {
		panic(err)
	}

	// 5. 调用模型生成回复（非流式，一次性返回完整结果）
	response, err := model.Generate(ctx, messages)
	if err != nil {
		panic(err)
	}

	// 6. 打印回复内容
	print(response.Content)
}