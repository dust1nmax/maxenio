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

## 容易混淆

- **变量名拼错不会编译报错。** `Format` 的入参是 `map[string]any`（`components/prompt/interface.go:44`）：

  ```go
  Format(ctx context.Context, vs map[string]any, opts ...Option) ([]*schema.Message, error)
  ```

  模板里写 `{role}`，params 里却写成 `{"roles": ...}`——**编译完全通过**，运行期这个变量填不上。编译器只知道 value 是 `any`，不知道 map 里有哪些 key。这是 Eino 里静态类型检查的边界之一，详见 [静态类型](./static_typing.md)。

- **`Format` 和模板内容的对应关系没有编译期约束。** 模板字符串里写了哪些占位符，只能靠人和运行期保证，`go build` 看不出来。

## 相关知识

- [ChatModel](./chat_model.md)
- [静态类型](./static_typing.md) —— 为什么 `map[string]any` 是类型检查的断点