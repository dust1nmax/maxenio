# Document 与 Markdown 文档切分

### 为什么需要

RAG 入库前要把长文档切成小片，原因有两个：

1. Embedding 模型有输入长度上限，整篇文档塞不进去。
2. 一整篇文档转成一个向量，语义会被"平均"掉，检索粒度也太粗——用户问"星海广场"，命中的却是整篇随笔。

切分的目的是让**每一片是一个独立的语义单元**，向量才有区分度。

### 一句话理解

HeaderSplitter 顺着 Markdown 的标题层级走，遇到标题就开新片，同时把"标题路径"写进这一片的 MetaData。

### 数据流

```text
document.md
 ↓ os.ReadFile
schema.Document{ID, Content}
 ↓ markdown.NewHeaderSplitter → Transform
[]*schema.Document  (每片带 h1/h2/h3 元数据)
 ↓ 下一步：Embedder
向量 → Indexer → Milvus
```

### schema.Document

| 字段 | 说明 |
|------|------|
| ID | 文档唯一标识 |
| Content | 文本内容 |
| MetaData | 自定义元数据（切分后用来存标题层级） |

注意 `Document.String()` **只返回 Content**（`eino@v0.9.19/schema/document.go:50`）。

### Eino 中怎么用

```go
transformer, err := markdown.NewHeaderSplitter(
	ctx,
	&markdown.HeaderConfig{
		Headers: map[string]string{
			"#":   "h1",
			"##":  "h2",
			"###": "h3",
		},
		TrimHeaders: false,
	},
)

base, err := os.ReadFile("./document.md")

doc := &schema.Document{ID: "doc1", Content: string(base)}

splitDocs, err := transformer.Transform(ctx, []*schema.Document{doc})
```

`Headers` 的 key 是标题标记（**只能由 `#` 组成**），value 是写进 MetaData 的字段名。

### 关键行为（读 header.go 源码确认）

1. **遇到标题行就开新片**：把标题之前的累积内容作为一片输出，然后开始新的一片。
2. **标题紧挨标题会产生"只有标题"的片**：`# 大连秋日随笔` 下面直接是 `## 一、清晨`，于是 h1 自己成一片，没有正文。这不是 bug，是规则的自然结果。
3. **空行被丢弃**：`if len(line) == 0 { continue }`，而且每行都过 `strings.TrimSpace`。所以片内不保留空行、也不保留缩进。
4. **MetaData 是"标题路径"，逐层累积**：遇到同级或更高级的标题时，更深层的记录会被删掉。所以"### 夕阳与灯塔"片的元数据是 `h1=大连秋日随笔 h2=二、傍晚 h3=夕阳与灯塔`，而进到"## 三、结语"后 h3 就没了。
5. **代码块内不切分**：```` ``` ```` / `~~~` 围栏内的行按普通行处理，里面的 `#` 不会触发切分。
6. **TrimHeaders 控制标题行是否留在正文里**：`false` 保留（片内容以标题行开头），`true` 丢弃。两种情况下 MetaData 都会带标题。
7. **不会误匹配**：判断条件是"以 header 开头，且后面紧跟空格或就是行尾"。所以 `#### 四级标题` 不会被 `###` 匹配到。

### 输入 / 输出

输入：`RAG/document.md`（869 字节，3 个 h2、6 个 h3）

输出：10 片，例如

```text
片段 0 | id=doc1 | h1=大连秋日随笔 h2=<nil> h3=<nil>
  内容: "# 大连秋日随笔"
片段 2 | id=doc1 | h1=大连秋日随笔 h2=一、清晨 h3=星海广场
  内容: "### 星海广场\n九月中旬的大连，暑气渐消……"
```

### 容易混淆 / 踩过的坑

- **所有切片共用原文档 ID。** 默认的 `defaultIDGenerator` 直接返回原 ID（`header.go:32`），所以 10 片全是 `doc1`。下一步入库 Milvus 时会主键冲突、互相覆盖，必须配 `IDGenerator` 让每片拿到唯一 ID。这是 07 → 05 衔接时最容易踩的坑。
- **`fmt.Println(splitDocs)` 看不出切分结果。** 因为 `Document.String()` 返回 Content，切片元素又是用空格连接的，打印出来像是"整篇没切开"。要逐片打印，并用 `%q` 才能看见里面的换行。
- **`os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0775)` 会真的创建文件。** 路径是相对 cwd 的，在别的目录跑一次就会凭空建出一个空文件（本项目根目录的 `document.md` 就是这么来的）。只读场景直接用 `os.ReadFile` 就够了。

### 还没理解 / 下一步

- `IDGenerator` 具体怎么写，才能保证入库不冲突
- 切好的片如何接进 05_indexer
- 切分粒度（h1/h2/h3、是否保留标题）对检索效果的实际影响

### 相关知识

- `docs/eino/milvus.md`（向量库与数据层级）
- `components/markdown.go`（切分器的封装）
- `RAG/Rag.go`（切分 → 入库 → 检索的完整调用）
