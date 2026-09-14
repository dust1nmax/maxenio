# Template

## 一句话理解

Template 就是"填空题"——固定内容 + 动态变量 = 最终 Prompt。

## 具体例子

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
```

## 工作流程

```text
输入: params = {role: "编程助手", task: "golang的优势"}
  ↓
Template.Format
  ↓
输出: []*schema.Message
  ↓
ChatModel.Generate
```

## Eino 中怎么用

### 1. 定义模板

```go
prompt.FromMessages(
    schema.FString,  // 模板语法：{variable}
    schema.SystemMessage("你是一个{role}"),
    &schema.Message{
        Role:    schema.User,
        Content: "回答{task}",
    },
)
```

### 2. 填充变量

```go
messages, err := template.Format(ctx, map[string]any{
    "role": "编程助手",
    "task": "golang的优势",
})
```

## 相关知识

- [ChatModel](./chat_model.md)