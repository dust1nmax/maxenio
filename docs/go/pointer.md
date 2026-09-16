# 指针与自动解引用

### 为什么需要

C 里访问指针指向的结构体字段要写 `p->field`，Go 把它简化成了 `p.field`——编译器知道 `p` 是指针，会自动替你解引用。这是为了让指针用起来不那么"硌手"。

### 一句话理解

`params.City` 就是 `(*params).City` 的简写。规则只有一条：**如果 x 是指向结构体的指针，`x.f` 等价于 `(*x).f`。**

### 具体例子

```go
type InputParams struct {
	City string `json:"city"`
}

func GetWeather(ctx context.Context, params *InputParams) (string, error) {
	// 下面两行完全等价，写第一种就行
	_ = params.City
	_ = (*params).City
	return "", nil
}
```

### 方法调用也自动，而且方向是双向的

```go
type Counter struct{ N int }

func (c Counter) Get() int  { return c.N }  // 值接收者
func (c *Counter) Inc()     { c.N++ }       // 指针接收者

p := &Counter{N: 1}
p.Get()   // 指针调值接收者 ⇒ (*p).Get()

v := Counter{N: 1}
v.Inc()   // 值调指针接收者 ⇒ (&v).Inc()，前提是 v 可寻址
```

两条规则合起来就是 Go 的"方法集自动补齐"：`*T` 的方法集包含 `T` 的所有方法，反过来不行——所以 `p.Get()` 合法，但如果换成不可寻址的值（比如函数返回的临时值）调指针接收者方法就会编译失败。

### 只自动一层

```go
var p **InputParams
_ = p.City        // 编译错误：type **InputParams has no field or method City
_ = (*p).City     // 对：先解一层拿到 *InputParams，再自动解第二层
```

自动解引用只对"指向结构体的一层指针"生效，`**T` 不行。

### 指针没有 nil 检查

```go
var p *InputParams
_ = p.City // panic: invalid memory address or nil pointer dereference
```

自动解引用只是语法糖，**不会**帮你挡 nil。指针类型的参数进函数后，如果来源不可靠，第一件事往往是判空。

### 为什么这里用指针而不是值

对一个只有一个字段的结构体，两者都能用。指针的实际意义是：

- **避免复制**：结构体大时省一次内存拷贝
- **允许修改调用方的数据**：函数内改字段，外面看得到
- **和 nil 区分"没传"和"传了零值"**：值类型无法表达"缺失"

Eino 里这个指针是安全的：`invokableTool.InvokableRun` 用
`generic.NewInstance[T]()` 构造实例，而 `NewInstance` 对指针类型做了
`reflect.New(elem)`（`internal/generic/generic.go:35`），所以 `T = *InputParams`
会拿到一个已分配的指针，再把 JSON 反序列化进去。

### 容易混淆

- **`p.City` 不是"指针上的字段"**，它就是被解引用后的字段，和 `(*p).City` 是同一个东西。
- **字段必须导出（大写）才能被 `encoding/json` / 反射填充**。小写字段在同一个包里能访问，编译器不报错，但反序列化时会被**静默忽略**——`p.City` 永远是零值。这个坑和指针无关，但经常一起出现。
- **取地址 `&v` 和自动解引用方向相反**：前者是"值 → 指针"（需要 v 可寻址），后者是"指针 → 值"（自动发生）。
- **切片、map、channel 本身是引用语义**，传参时不需要再加 `*`；只有 struct 和基本类型加指针才有"避免复制/允许修改"的效果。

### 相关知识

- `docs/go/generics.md`：`NewTool[T, D any]` 里 `T` 就是这个输入类型
- `internal/generic/generic.go`：`NewInstance[T]()` 对各类型的零值构造
