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
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	model, err := ark.NewChatModel(
		ctx,
		&ark.ChatModelConfig{
			APIKey: os.Getenv("ARK_API_KEY"),
			Model:  os.Getenv("MODEL"),
		})

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}"),
		&schema.Message{
			Role:    schema.User,
			Content: "回答{task}",
		},
	)

	params := map[string]any{
		"role": "编程助手",
		"task": "golang的优势",
	}

	messages, err := template.Format(ctx, params)

	// intput := []*schema.Message{
	// 	schema.SystemMessage("你是一个人"),
	// 	schema.UserMessage("你是谁？"),
	// }

	response, err := model.Generate(
		ctx,
		messages,
	)
	if err != nil {
		panic(err)
	}
	print(response.Content)

	//流式输出
	// reader,err := model.Stream(ctx, intput)
	// if err != nil{panic(err)}
	// defer reader.Close()

	// for{
	// 	chunk,err := reader.Recv()
	// 	if err != nil{panic(err)}
	// 	print(chunk.Content)
	// }
}
