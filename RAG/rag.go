package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"

	"maxenio/components"
)

func main() {
	// 所有组件都依赖 .env 里的 ARK_API_KEY / MODEL / EMBEDDER / ADDRESS，
	// 所以必须在装配组件之前加载
	if err := godotenv.Load("/home/max/maxenio/.env"); err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 装配阶段：组件只负责把 error 往上报，由这里（程序顶层）统一决定失败就退出。
	// 这也是把 log.Fatalf 从组件里移出来的意义——换成 HTTP 服务时，
	// 这里可以改成返回 500 而不是杀掉进程。
	embedder, err := components.NewEmbedder(ctx)
	if err != nil {
		log.Fatalf("初始化 embedder 失败: %v", err)
	}

	indexer, err := components.NewArkIndexer(ctx, embedder, os.Getenv("ADDRESS"))
	if err != nil {
		log.Fatalf("初始化 indexer 失败: %v", err)
	}

	retriever, err := components.NewArkRetriever(ctx, embedder)
	if err != nil {
		log.Fatalf("初始化 retriever 失败: %v", err)
	}

	transform, err := components.NewTransform(ctx)
	if err != nil {
		log.Fatalf("初始化切分器失败: %v", err)
	}

	/////////////////////////////////////////////////////////////
	//    把文档传入 milvus
	/////////////////////////////////////////////////////////////
	base, err := os.ReadFile("./document.md")
	if err != nil {
		log.Fatalf("读取 document.md 失败: %v", err)
	}

	doc := &schema.Document{
		ID:      "doc1",
		Content: string(base),
	}

	splitDocs, err := transform.Transform(
		ctx,
		[]*schema.Document{doc},
	)
	if err != nil {
		log.Fatalf("切分失败: %v", err)
	}

	ids, err := indexer.Store(ctx, splitDocs)
	if err != nil {
		log.Fatalf("写入 Milvus 失败: %v", err)
	}
	fmt.Println("已写入:", ids)

	/////////////////////////////////////////////////////////////
	//    把文档从 milvus 读出
	/////////////////////////////////////////////////////////////
	results, err := retriever.Retrieve(ctx, "海鸥")
	if err != nil {
		log.Fatalf("检索失败: %v", err)
	}

	for _, doc := range results {
		fmt.Println(doc.ID)
		fmt.Println(doc.Content)
	}
}
