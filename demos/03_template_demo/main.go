package main

import (
	"context"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func main() {
	ctx := context.Background()

	// 1. 定义 Prompt 模板
	// prompt.FromMessages: 从多个消息构建模板
	// schema.FString: 表示使用 {变量名} 的语法
	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"), // 系统消息：你是[什么角色]
		&schema.Message{
			Role:    schema.User,
			Content: "回答{task}", // 用户消息：回答[什么问题]
		},
	)

	// 2. 准备变量值（key 要和模板里的 {role}, {task} 对应）
	params := map[string]any{
		"role": "编程助手",
		"task": "golang的优势",
	}

	// 3. 调用 Format 填充模板，得到 []*schema.Message
	// 这是 Template 最重要的方法：固定模板 + 变量 = 最终消息
	messages, err := template.Format(ctx, params)
	if err != nil {
		panic(err)
	}

	// 4. 查看生成的消息（可以打印出来理解 Template 的作用）
	for i, msg := range messages {
		println(i, msg.Role, msg.Content)
	}
}