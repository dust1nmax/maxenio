package components

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
)

// StreamChat 以流式方式调用模型，立刻返回一个流读取器。
//
// 和 Generate 的区别：Generate 会阻塞到模型把整段回复生成完才返回；
// Stream 马上返回，由调用方不断从 reader 里取增量片段，边生成边展示。
//
// 这里只负责创建和错误包装，"怎么消费"交给 ConsumeStream。
// 注意一个 StreamReader 只能读一次，多处要用需要先 Copy。
func StreamChat(ctx context.Context, model *ark.ChatModel, messages []*schema.Message) (*schema.StreamReader[*schema.Message], error) {
	reader, err := model.Stream(ctx, messages)
	if err != nil {
		return nil, fmt.Errorf("流式调用模型失败: %w", err)
	}

	return reader, nil
}

// ConsumeStream 把流里的增量片段逐个交给 onChunk，并负责关流和结束判断。
//
// 封装的是 eino 官方推荐的那套固定循环：
//
//	defer reader.Close()                     // 读到 io.EOF 之后也要关
//	for {
//	    chunk, err := reader.Recv()
//	    if errors.Is(err, io.EOF) { break }  // 流正常结束
//	    if err != nil { ... }                // 真出错
//	    // 处理 chunk
//	}
//
// 关键点：Recv 返回的 io.EOF 表示"流正常结束"，不是错误。
// 如果只写 `if err != nil { break }`，正常结束和真出错就分不开了，
// 网络中断、鉴权失效这类失败会被静默当成"生成完了"。
//
// onChunk 返回 error 时会中断消费并把它抛出去，
// 让调用方能在下游写失败时（比如 HTTP 连接断了）及时止损。
//
// 关流由 ConsumeStream 负责，调用方不要再自己 defer reader.Close()，避免重复关闭。
func ConsumeStream(reader *schema.StreamReader[*schema.Message], onChunk func(chunk string) error) error {
	defer reader.Close()

	for {
		chunk, err := reader.Recv()
		if errors.Is(err, io.EOF) {
			// 流正常读完
			return nil
		}
		if err != nil {
			return fmt.Errorf("读取流失败: %w", err)
		}

		if err := onChunk(chunk.Content); err != nil {
			return err
		}
	}
}
