# Go 泛型（类型参数）

### 为什么需要

有些逻辑和"具体是什么类型"无关，但又不能放弃类型检查。

Eino 的 `utils.NewTool` 就是典型例子：它要接收用户写的业务函数，再包装成框架统一的工具接口。每个用户的参数和返回值类型都不一样，如果不用泛型，就只能收 `any`，然后在内部做类型断言——类型错了要到运行期才知道。

### 一句话理解

泛型就是"**类型也可以当参数传**"。函数名后面方括号里的 `[T, D any]` 是在声明"这个函数有两个类型参数，具体是什么类型由调用方决定"。

### 语法拆解

拿 Eino 的这行签名举例（`eino@v0.9.19/components/tool/utils/invokable_func.go:143`）：

```go
func NewTool[T, D any](desc *schema.ToolInfo, i InvokeFunc[T, D], opts ...Option) tool.InvokableTool
```

| 部分 | 名称 | 含义 |
|------|------|------|
| `func NewTool` | 函数名 | 和普通函数一样 |
| `[T, D any]` | **类型参数列表** | `T`、`D` 是类型参数，`any` 是它们的**约束** |
| `desc *schema.ToolInfo` | 普通参数 | 工具描述（名字、参数 schema） |
| `i InvokeFunc[T, D]` | 普通参数 | 你的业务函数，已被 `T`、`D` 参数化 |
| `opts ...Option` | 可变参数 | 选项 |
| `tool.InvokableTool` | 返回值 | 普通接口 |

`[T, D any]` 是简写，完整写法是 `[T any, D any]`——两个类型参数共用一个约束。位置在**函数名之后、圆括号之前**。

再看被参数化的那个函数类型（同文件 `:33`）：

```go
type InvokeFunc[T, D any] func(ctx context.Context, input T) (output D, err error)
```

所以在这套约定里，`T` 是工具的**输入类型**，`D` 是**输出类型**。

### 调用时类型是推断出来的

平时不用手写方括号，编译器会根据你传的函数推断：

```go
type WeatherArgs struct {
	City string `json:"city" jsonschema:"description=城市名，例如 大连"`
}

type WeatherResult struct {
	Temp int `json:"temp"`
}

func queryWeather(ctx context.Context, args WeatherArgs) (WeatherResult, error) {
	return WeatherResult{Temp: 26}, nil
}

// queryWeather 的类型是 func(context.Context, WeatherArgs) (WeatherResult, error)
// 它正好匹配 InvokeFunc[WeatherArgs, WeatherResult]，所以 T 和 D 被自动推断出来
t := utils.NewTool(
	&schema.ToolInfo{Name: "get_weather", Desc: "查询指定城市的天气"},
	queryWeather,
)
```

也可以显式指定，两种写法等价：

```go
t := utils.NewTool[WeatherArgs, WeatherResult](desc, queryWeather)
```

### 泛型和接口的分工（关键）

这一点是理解 Eino 为什么要用泛型的核心。

框架对外暴露的是**接口**，接口里参数是原始的 JSON 字符串（`components/tool/interface.go:42`）：

```go
type InvokableTool interface {
	BaseTool // 提供 Info(ctx) (*schema.ToolInfo, error)

	InvokableRun(ctx context.Context, argumentsInJSON string, opts ...Option) (string, error)
}
```

模型（LLM）产生的是 JSON 字符串，框架调度工具时也只能传字符串。但业务代码想操作的是**结构体**。

`NewTool` 用泛型把这两端接起来：

```text
模型 / 框架侧                    业务代码侧
字符串 JSON                     类型化的 struct
InvokableRun(ctx, jsonStr)  →   InvokeFunc[T, D](ctx, input T)
                                 ↑ 内部自动 JSON 反序列化 + 序列化
```

分工总结：

- **泛型**负责编译期：把你的 `T` / `D` 接上，类型错了编译不过
- **接口**负责运行期：给框架一个统一的 `InvokableRun` 调用入口，框架不需要知道你的类型

注意返回的 `tool.InvokableTool` 是个接口，泛型参数在这里被"擦掉"了——调用方拿到的就是一个能接收 JSON 字符串的工具。

### 泛型类型（不只是函数）

结构体也可以带类型参数，Eino 内部就是这么存的（`invokable_func.go:170`）：

```go
type invokableTool[T, D any] struct {
	info *schema.ToolInfo
	Fn   OptionableInvokeFunc[T, D]
}
```

所以 `NewTool` 返回它的时候，`T`、`D` 作为接口被擦除，这也解释了为什么返回类型写的是接口而不是 `*invokableTool[T, D]`。

### 约束（constraint）

`any` 是最宽松的约束，表示"什么类型都行"。约束还可以是：

| 约束写法 | 含义 |
|----------|------|
| `any` | 任意类型（等价于 `interface{}`） |
| `comparable` | 可以用 `==` 比较的类型 |
| 接口名 | 必须实现该接口 |
| `~int \| ~string` | 底层类型是 int 或 string 之一（`~` 表示含自定义类型） |

例如想约束"能比较大小"就得自己定义接口约束，因为 Go 没有内置的算术/排序约束。

### 容易混淆

- **`[T, D any]` 不是数组或切片**。方括号在**类型参数位置**（函数名后）时是类型参数列表；在**类型位置**（如 `[]int`、`[3]int`）时才是切片或数组。
- **`any` 就是 `interface{}` 的别名**，Go 1.18 起可用，两者完全等价。
- **同一个包的 `NewTool` 和 `InferTool` 分工不同**：`NewTool` 要你手工提供 `*schema.ToolInfo`，**没有编译期检查**保证 schema 和 `T` 的字段一致（源码注释明确提示了这点）；`InferTool` 从 `T` 的 struct tag 自动推断参数 schema（`json` 决定字段名，`jsonschema` 决定描述），所以它返回 `(tool.InvokableTool, error)`——推断可能失败。
- **泛型不是"运行时多态"**。Go 的泛型在编译期展开成具体类型，跟 Java 的类型擦除、C++ 的模板实例化都不是一回事。

### 相关知识

- `components/tool/utils/invokable_func.go`：`NewTool`、`InferTool`、`InvokeFunc` 的完整实现
- `components/tool/interface.go`：`InvokableTool`、`StreamableTool` 接口
