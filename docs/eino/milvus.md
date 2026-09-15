# Milvus

## 一句话理解

开源向量数据库，用于存储和检索高维向量，支持相似度搜索。

## 架构组件

| 组件 | 作用 |
|------|------|
| etcd | 保存 Milvus 元数据 |
| MinIO | 对象存储，保存持久化数据 |
| Milvus Standalone | 提供向量数据库服务 |
| Attu | Web 可视化管理界面 |

### 关键理解

- Attu 只是管理界面，不是 Milvus 数据库本身
- Milvus 服务端口：`localhost:19530`
- Attu Web 界面：`localhost:8000`

## 数据结构

```
Database → Collection → Field
```

示例：
```
MaxEino → test → id / vector / content / metadata
```

- Database：逻辑分类
- Collection：类似关系数据库的 Table
- Field：数据字段

## 向量类型

### FloatVector

- 由浮点数构成的向量
- 适合普通 dense Embedding
- 当前使用这个

### BinaryVector

- 二进制向量
- 与 FloatVector 是不同的数据类型

### MetricType

用于计算向量相似度：

- `COSINE`：余弦相似度
- `L2`：欧氏距离
- `IP`：内积

当前使用 `COSINE`。

## Docker 部署

```yaml
# docker-compose.yml
services:
  milvus:
    volumes:
      - ./volumes/milvus:/var/lib/milvus
```

### 数据持久化

`docker compose down` 删除容器时，宿主机 `volumes/` 中的数据仍然保留。

### 端口

- Milvus：`19530`
- Attu：`8000`

## Eino milvus2

使用：
```go
github.com/cloudwego/eino-ext/components/indexer/milvus2
```

### 核心配置

- Address：`localhost:19530`
- Database：`MaxEino`
- Collection：`test`
- Dimension：`2048`
- Metric：`COSINE`

### Indexer 工作流程

```
Document → Embedding → Vector → Indexer → Milvus
```

### Store 用法

```go
err = indexer.Store(ctx, docs)
```

### Document 存储格式

`schema.Document` 是 Eino 定义的文档结构：

```go
docs := []*schema.Document{
    {
        ID:      "1",                              // 文档唯一标识
        Content: "今天是2026年9月15日17：05",      // 文本内容（会被转成向量）
        MetaData: map[string]any{
            "author": "max",                       // 自定义元数据
        },
    },
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `ID` | string | 文档唯一标识符 |
| `Content` | string | 文档内容（Indexer 自动调用 Embedder 转成向量） |
| `MetaData` | `map[string]any` | 自定义元数据（如 author、timestamp） |

### 存储流程

```
输入 schema.Document（包含文本内容）
        ↓
Indexer 内部调用关联的 Embedder
        ↓
Content 被转换为向量（2048 维）
        ↓
向量 + ID + MetaData 存入 Milvus
```

### milvus2 vs milvus（v1）

| | milvus | milvus2 |
|---|---|---|
| 客户端路径 | `milvus-io/milvus/client/v1` | `milvus-io/milvus/client/v2` |
| API 风格 | 旧版 | 现代化 |
| Eino 封装 | `milvus` | `milvus2` |

Eino 使用 `milvus2` 包作为 v2 客户端的封装，API 更简洁。

## 相关知识

- [Embedder](./embedder.md)
- [RAG](./rag.md)