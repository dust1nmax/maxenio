# Chain

### 一句话理解

Chain 就是把多个组件按**固定顺序**串成一条直线，并提前声明整条链的输入类型和输出类型。

### 为什么需要

不编排的时候，`main.go` 是这样写的：

```go
messages, err := template.Format(ctx, params)
response, err := model.Generate(ctx, messages)
```

两个组件还好，一旦变成「Template → ChatModel → 加工 → 再进模型」，手写中间变量就会越来越乱：顺序靠人的代码顺序保证，中间产物到处传，出错了也不容易看出来是哪一步。

Chain 把「顺序」变成显式声明，并且把类型检查交给框架。

### 具体例子

`orchestration/chain/chain.go`：

```go
// Lambda 节点：进 LLM 之前加工用户输入
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

chain := compose.NewChain[string, *schema.Message]()
chain.AppendLambda(lambda).AppendChatModel(model)

r, err := chain.Compile(ctx)
answer, err := r.Invoke(ctx, "你好，请告诉我你的名字")
fmt.Println(answer.Content)
```

实测输出：

```text
你好，我的名字是DeepSeek。今天天气很好。
```

Lambda 把消息拼成了 `"你好，请告诉我你的名字回答结尾加上 今天天气很好"`，模型把后半句当成了指令，所以回答结尾真的加了那句话。

### 工作流程

```text
Invoke(ctx, "你好，请告诉我你的名字")   ← 字符串
  ↓
Lambda 节点 (string → []*schema.Message)
  ↓
ChatModel 节点 ([]*schema.Message → *schema.Message)
  ↓
Invoke 返回 *schema.Message
  ↓
answer.Content
```

### 输入 / 输出

```go
chain := compose.NewChain[string, *schema.Message]()
//                        ↑ I  ↑ O

r, err := chain.Compile(ctx)        // r 的类型是 compose.Runnable[string, *schema.Message]
answer, err := r.Invoke(ctx, "...") // 入参类型 I，返回值类型 O
```

`Runnable[I, O]` 接口（`compose/runnable.go:32`）：

```go
type Runnable[I, O any] interface {
	Invoke(ctx context.Context, input I, opts ...Option) (output O, err error)
	Stream(ctx context.Context, input I, opts ...Option) (output *schema.StreamReader[O], err error)
	Collect(ctx context.Context, input *schema.StreamReader[I], opts ...Option) (output O, err error)
	Transform(ctx context.Context, input *schema.StreamReader[I], opts ...Option) (output *schema.StreamReader[O], err error)
}
```

四个方法的区别只是「流进流出」的组合：`Invoke` 非流进非流出，`Stream` 非流进流出，`Collect` 流进非流出，`Transform` 流进流出。

### 关键规则：`NewChain[I, O]` 的类型参数

> **`I` 必须等于第一个节点的输入类型，`O` 必须等于最后一个节点的输出类型。**

中间节点的类型由框架自己串，不用你声明，但首尾两端必须和声明一致。

这条规则是实测出来的：

```go
// 第一个节点是 ChatModel，但 I 写成了 string
chain := compose.NewChain[string, *schema.Message]()
chain.AppendChatModel(model)

r, err := chain.Compile(ctx)
```

**Go 编译能通过**（`go vet` 也不报），但运行到 `Compile` 时报：

```text
panic: graph edge[start]-[node_0]: start node's output type[string] and end node's input type[[]*schema.Message] mismatch
```

因为 ChatModel 节点的输入类型是 `[]*schema.Message`（它的 `Generate` 方法签名决定的），接不上 `string`。

两种改法都实测可用：

| 改法 | 代码 | 适用 |
|------|------|------|
| A. 让 I 等于 ChatModel 的输入 | `NewChain[[]*schema.Message, *schema.Message]()`，`Invoke(ctx, []*schema.Message{schema.UserMessage("...")})` | 不想多一层节点 |
| B. 加一个接受 string 的节点 | `NewChain[string, *schema.Message]()` + `AppendLambda(...)` | 想直接 `Invoke(ctx, "字符串")` |

`chain.go` 用的是 B。

### 容易混淆

- **`Compile` 报的错不是 Go 编译错误。** 类型参数写错在 Go 编译期**查不出来**，要等 `chain.Compile(ctx)` 返回 error。这正是"静态类型有边界"在 Chain 上的体现——见 [静态类型](./static_typing.md)。
- **`AppendXxx` 的返回值是链本身，可以直接点下去。** `chain.AppendLambda(lambda).AppendChatModel(model)` 等价于分两行写；`I`、`O` 不会因为追加节点而改变。
- **`NewChain` 不返回 error，`Compile` 才返回。** 追加节点阶段只记录顺序，所有类型检查推迟到 `Compile`。
- **Lambda 是"绕过组件接口"的通用节点。** 组件类型（ChatModel / Template / Retriever 等）都有固定的输入输出，Lambda 的类型来自你写的函数签名，所以它能做任意的类型转换和加工。
- **Chain 是"直线的 Graph"。** 顺序写死、不能分支和并行；需要分支就用 Graph（或 `AppendBranch` / `AppendParallel`）。
- **局部变量不要叫 `context`。** 会遮蔽 `context` 包名，同一函数里再想写 `context.Context` 就用不了了。

### 相关知识

- [静态类型](./static_typing.md) —— 为什么 `Compile` 的类型不匹配是运行期错误
- [ChatModel](./chat_model.md) —— 节点的输入输出类型从哪来
- [Template](./template.md) —— 编排里最常见的第一个节点
