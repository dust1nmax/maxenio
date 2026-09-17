package main

import (
	"context"
	"fmt"
	"maxenio/components"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("/home/max/maxenio/.env")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	model, err := components.NewarkModel(ctx)
	if err != nil {
		panic(err)
	}

	// 编写lambda 节点 在送往llm 前 对用户输入进行加工
	lambda := compose.InvokableLambda(
		func(ctx context.Context, input string) (output []*schema.Message, err error) {
			content := input + "回答结尾加上 今天天气很好"
			output = []*schema.Message{
				{
					Role:    schema.User,
					Content: content,
				},
			}
			return output, nil
		})

	// 注册一个链条
	// 提前定义好 chain的 输入 和 输出 类型
	chain := compose.NewChain[string, *schema.Message]()
	chain.AppendLambda(lambda).AppendChatModel(model)
	r, err := chain.Compile(ctx)
	if err != nil {
		panic(err)
	}
	answer, err := r.Invoke(ctx, "你好，请告诉我你的名字")
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Content)
}
