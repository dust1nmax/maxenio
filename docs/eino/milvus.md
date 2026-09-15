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

## 相关知识

- [Embedder](./embedder.md)
- [RAG](./rag.md)