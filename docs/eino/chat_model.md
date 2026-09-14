# ChatModel

## 一句话理解

ChatModel 是 Eino 中与 LLM 交互的核心组件，负责接收消息并返回回复。

## 具体例子

### 基础用法：直接传 Messages

```go
model, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
    APIKey: os.Getenv("APK_API_KEY"),
    Model:  os.Getenv("MODEL"),
})

messages := []*schema.Message{
    schema.SystemMessage("你是一个助手"),
    schema.UserMessage("你是谁？"),
}

response, err := model.Generate(ctx, messages)
```

### 结合 Prompt Template

```go
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
response, err := model.Generate(ctx, messages)
```

### 流式输出

```go
reader, err := model.Stream(ctx, messages)
defer reader.Close()

for {
    chunk, err := reader.Recv()
    if err != nil {
        break
    }
    print(chunk.Content)
}
```

## 工作流程

```text
输入 Messages
  ↓
ChatModel
  ↓
LLM (ark/other)
  ↓
Generate → Response
或
Stream → StreamReader → chunk
```

## 输入 / 输出

**Generate 输入：**
- ctx: context.Context
- messages: []*schema.Message

**Generate 输出：**
- response: *schema.Message
- error

**Stream 输入：**
- ctx: context.Context
- messages: []*schema.Message

**Stream 输出：**
- reader: *StreamReader
- error

## Eino 中怎么用

1. **创建 ChatModel 实例**（不同 provider 有不同实现）
   - ark: `ark.NewChatModel`
   - openai: `openai.NewChatModel`

2. **调用 Generate 或 Stream**
   - Generate: 一次性返回完整响应
   - Stream: 流式返回 chunks

## 相关知识

- [Template](./template.md)
- schema.Message