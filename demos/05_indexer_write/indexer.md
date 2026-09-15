# 05_indexer 使用说明

## 1. 项目说明

05_indexer 基于 Eino 框架，将 Document 通过 Embedding 转成向量，存入 Milvus 向量数据库。

```
Document → Embedding → Vector → Indexer → Milvus
```

这是 RAG 的**索引（写入）**部分。

## 2. 项目结构

```
/home/max/maxenio/
├── docker-compose.yml    # Milvus 服务配置（在根目录）
├── .env                  # 环境变量（在根目录）
├── go.mod                # Go 依赖
├── demos/
│   └── 05_indexer/
│       └── main.go       # 主程序
└── volumes/              # Milvus 数据持久化
```

## 3. 详细流程

### 第一步：启动 Milvus

Milvus 由 4 个组件构成：

- **etcd**：保存元数据
- **MinIO**：对象存储
- **standalone**：Milvus 主服务
- **Attu**：Web 管理界面

在 `/home/max/maxenio` 目录执行：

```bash
docker compose up -d
```

等待约 1 分钟让服务全部就绪。

### 第二步：确认 Milvus 就绪

通过 Docker Desktop 查看：

1. 打开 Docker Desktop
2. 点击 **Containers**
3. 确认 4 个容器状态都是绿色 Running：
   - milvus-etcd
   - milvus-minio
   - milvus-standalone
   - milvus-attu

然后通过 Attu 打开 Milvus：

1. 点击 milvus-attu 的 **Port(s)** 或直接在浏览器打开 http://localhost:8000
2. Attu 界面中选择 Database: **MaxEino**
3. 此时 Collection 列表应该是空的，等代码运行后会自动创建 Collection "test"

### 第三步：运行 Go 程序

```bash
go run demos/05_indexer/main.go
```

程序运行后会输出类似：

```
Stored documents with IDs: [1]
```

### 第四步：查看存储结果

刷新 Attu 页面 http://localhost:8000：
- Database: MaxEino
- 查看 Collection: test
- 可以看到刚存入的文档

## 4. 配置说明

| 配置项 | 值 |
|--------|-----|
| Milvus 地址 | localhost:19530 |
| Database | MaxEino |
| Collection | test |
| 向量维度 | 2048 |
| 相似度度量 | COSINE |
| Embedder 模型 | doubao-embedding-vision-251215 |

## 5. 停止 Milvus

```bash
docker compose down
```

数据保存在 `volumes/` 目录，不会丢失。下次 `docker compose up -d` 可继续使用。

## 6. 后续学习

当前项目完成的是 RAG 的**写入**部分：

```
Document → Embedding → Indexer → Milvus
```

下一步学习**读取**部分，即 Retriever：
- 把用户查询转成向量
- 在 Milvus 中搜索相似文档
- 返回相关文档给 LLM