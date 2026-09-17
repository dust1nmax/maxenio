# 静态类型（Eino 为什么强调这个）

### 为什么需要

Eino 的官方介绍里有一句"大模型应用编排框架的主流语言是 Python，这门语言以灵活性著称，灵活性给 SDK 的开发带来便利，但同时也给 SDK 的使用者带来了心智负担。基于 Golang 的 Eino 则是静态类型"。

这句话是**拿来跟 Python 系框架（LangChain / LlamaIndex 等）做对比的**，不是单纯在夸 Go。要理解它的价值，得先看清动态类型给 SDK 使用者带来了什么负担。

Python SDK 的典型手感：

```python
# 这个参数到底该传什么？看签名看不出来
chain.invoke({"question": "..."})
# 是 dict？是 str？是 HumanMessage？拼错 key 会怎样？
```

答案是：看文档、读源码、或者跑一遍试试。而框架在运行期拿到这个"什么都有可能"的对象后，只能靠 `isinstance` / `hasattr` / `getattr` 一层层猜，猜不到就抛 `TypeError` / `KeyError`——**错误发生在运行期，而不是你写代码的时候**。

### 一句话理解

静态类型 = **类型在编译期就确定，并且由编译器强制检查**。

对照：

| | 类型什么时候确定 | 传错了什么时候发现 |
|---|---|---|
| Python（动态类型） | 运行期 | 跑到那一行才报错 |
| Go（静态类型） | 编译期 | `go build` 就报错，不用运行 |

### 直觉理解

动态类型像一台**没有形状约束的插座**：说明书在文档里、在示例里、在别人的博客里。你得先读、先猜，插错了才知道。

静态类型像**插头和插座有形状**：形状不对根本插不进去，不需要谁来告诉你。

在这个比喻里，**函数签名就是那个"形状"**。

### 具体例子

#### 1. 签名本身就是说明书

Eino 的组件接口（`eino@v0.9.19/components/document/interface.go:53`）：

```go
type Transformer interface {
	Transform(ctx context.Context, src []*schema.Document, opts ...TransformerOption) ([]*schema.Document, error)
}
```

只看这一行就知道：入口是 `Document` 切片，出口也是 `Document` 切片。**不需要查文档**。

同类的还有（`components/retriever/interface.go:48`、`components/indexer/interface.go:38`）：

```go
type Retriever interface {
	Retrieve(ctx context.Context, query string, opts ...Option) ([]*schema.Document, error)
}

type Indexer interface {
	Store(ctx context.Context, docs []*schema.Document, opts ...Option) (ids []string, err error)
}
```

#### 2. 类型对上，链路才能串起来

`RAG/rag.go` 里三段代码能直接接起来，靠的就是类型吻合：

```go
doc := &schema.Document{ID: "doc1", Content: string(base)}

// Transform 出口类型 []*schema.Document
splitDocs, err := transform.Transform(ctx, []*schema.Document{doc})

// Store 入口类型 []*schema.Document —— 正好等于上一步的出口
ids, err := indexer.Store(ctx, splitDocs)

// Retrieve 出口类型也是 []*schema.Document
results, err := retriever.Retrieve(ctx, "海鸥")
```

这条链上任何一环类型不匹配，都是在 `go build` 阶段断开，而不是运行到一半才炸。

#### 3. 工具的定义方式

`tool/tool.go`：

```go
type InputParams struct {
	City string `json:"city" jsonschema:"description = name of city"`
}

func GetWeather(_ context.Context, params *InputParams) (string, error) { ... }

utils.NewTool(&schema.ToolInfo{ ... }, GetWeather)
```

`NewTool[T, D any]` 里的 `T = *InputParams`、`D = string` 是编译器从 `GetWeather` 的签名推断出来的。换一种入参形态就编译不过。（泛型机制详见 [Go 泛型](../go/generics.md)）

### 它替使用者省掉了哪些心智负担

| 心智负担 | Python 动态类型 | Go / Eino 静态类型 |
|---|---|---|
| 这个参数要传什么？ | 查文档 / 读源码 / 试 | 看签名，IDE 直接补全 |
| 我传错了怎么办？ | 运行到那行才报错 | 编译不过 |
| 这个字段叫什么？ | 靠记忆，拼错静默失败 | `schema.Message{Role: ..., Content: ...}` 字段名编译期检查 |
| 改了结构体，哪些地方受影响？ | 全局搜索 + 祈祷 | 编译器把每个出错点列出来 |
| 这条链路怎么接？ | 靠约定和文档 | 上游输出类型 == 下游输入类型 |

核心一句话：**类型信息从"文档和记忆"搬到了"代码本身"**。

### 工作流程

```text
写代码
  ↓
类型不匹配？ ──是──→ go build 报错（根本不用运行）
  ↓ 否
编译通过
  ↓
运行
```

### 重要边界：静态类型管不到的地方

这是最容易误解的一点。Go 有 `any`（= `interface{}`），**一旦落到 `any`，静态检查就断在那里**，重新变成运行期问题。

`main.go` 里就有现成的例子：

```go
params := map[string]any{
	"role": "编程助手",
	"task": "golang的优势",
}
messages, err := template.Format(ctx, params)
```

而 `Format` 的签名是（`eino@v0.9.19/components/prompt/interface.go:44`）：

```go
Format(ctx context.Context, vs map[string]any, opts ...Option) ([]*schema.Message, error)
```

模板里写的是 `{role}`。如果这里把 key 拼成 `"roles"`：

- **编译完全通过**
- 运行期模板变量填不上

因为编译器只知道 value 是 `any`，不知道这个 map 里有没有 `role` 这个 key。这就是 [Go fmt](../go/fmt.md) 里那个区分的实际后果：`any` 的**静态类型是 `any`**，里面装的什么要运行期才知道。

同样的边界出现在工具调用上。框架面向模型暴露的入口是 JSON 字符串（`components/tool/interface.go`）：

```go
InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
```

模型吐出来的 JSON 是运行期数据，静态类型在这一步必须让位。`utils.NewTool` 内部会把 JSON 反序列化成 `T`，这一步失败也是运行期错误。所以手工传进去的 `*schema.ToolInfo` 和 `T` 的字段是否一致，**编译器不检查**——这也是 `InferTool` 存在的理由。

结论：

> 静态类型不是"消灭运行期错误"，而是**把检查推到编译期能推的最远处，剩下的部分明确标出来**。

Eino 的做法正是如此：组件与组件之间靠编译期类型衔接；一旦进入 `any` / JSON 这种必须动态的边界，就明确用 `any` 标出来。

### 容易混淆

- **静态类型 vs 强类型是两根轴**。Go 是静态类型，Python 是动态类型，但两者都是强类型（`1 + "a"` 在两边都不让过）。静态/动态说的是"**什么时候**检查"，强/弱说的是"检查**多严**"。这里讨论的是前者。
- **静态类型不排斥动态**。Go 的 `any` 和接口就是可控的动态。`any` 用得越多，手感越接近动态语言，代价是回到运行期报错。
- **"Go 编译期检查"和"图编排检查"不是一回事**。`compose.NewGraph[I, O any]()` 的输入输出类型是编译期确定的（`compose/generic_graph.go:72`），但节点之间"上游输出类型必须等于下游输入类型"这个检查，是 `AddEdge` / `Compile` 阶段**返回 error** 检查的，不属于 Go 编译期。学到 Graph 时注意这个区别。
- **动态类型不等于"弱"、静态类型不等于"啰嗦"**。静态类型的"啰嗦"（要写 struct、要处理 error）换来的是上面那张表里的每一项负担被消除。

### 还没理解

- Python SDK 的"心智负担"具体在 LangChain 里长什么样（还没实际用过 Python 系框架，属于二手信息）

### 相关知识

- [Go 泛型](../go/generics.md) —— 泛型与接口的分工，就是静态类型在 Eino 里的落地方式
- [Go fmt](../go/fmt.md) —— 静态类型 vs 动态类型（`any`、`%T`）
- [Template](./template.md) —— `Format(ctx, map[string]any)` 正是静态类型的边界所在
- [ChatModel](./chat_model.md) —— `Generate(ctx, []*schema.Message)` 同上
