# Maxenio

CloudWeGo Eino + Go Agent 学习项目。边学边记，笔记即代码。

自带学习助手，定义在 `.zcode/skills/eino-agent-learning/SKILL.md`（ZCode 的 skill 目录）：

- 陪用户学习 Eino / Go / Agent 工程
- 自动把讨论内容整理成 Markdown 笔记到 `docs/`
- 主题笔记（`docs/eino/*.md`）+ 每日记录（`docs/daily/*.md`）

```text
学习 → 讨论 → 自动生成 docs/笔记 → 形成可复习的知识库
```

## 项目结构

```
├── components/         # 自己封装的 Eino 组件，每个组件一个构造函数
│   ├── chat_model.go      # NewChatModel：对话模型
│   ├── chat.go            # StreamChat / ConsumeStream：流式调用与消费
│   ├── embedding.go       # NewEmbedder / EmbedTexts：文本转向量
│   ├── indexer.go         # NewArkIndexer：写入 Milvus
│   ├── retriever.go       # NewArkRetriever：从 Milvus 检索
│   └── markdown.go        # NewTransform：Markdown 标题切分
├── RAG/                # 完整 RAG 流程：切分 → 向量化 → 入库 → 检索
│   ├── Rag.go
│   └── document.md        # 测试用的 Markdown 文档
├── docs/               # 学习笔记
│   ├── eino/              # Eino 主题笔记
│   │   ├── chat_model.md
│   │   ├── document.md
│   │   ├── embedder.md
│   │   ├── milvus.md
│   │   └── template.md
│   ├── go/                # Go 语言基础
│   │   └── fmt.md
│   └── daily/             # 每日学习记录
├── main.go             # 入口示例：Template + Generate
└── docker-compose.yml  # Milvus 服务（etcd + MinIO + standalone + Attu）
```

`components/` 里的函数只负责创建组件并把 error 往上报，不自己终止进程；
"失败怎么办"由调用方决定——`RAG/Rag.go` 里用 `log.Fatalf`，换成 HTTP 服务时改成返回错误即可。

## 从零配置环境

### 1. 初始化项目

```bash
mkdir my-eino-project && cd my-eino-project
go mod init my-eino-project
```

### 2. 安装依赖

```bash
go get github.com/cloudwego/eino
go get github.com/cloudwego/eino/schema
go get github.com/cloudwego/eino/components/prompt

go get github.com/cloudwego/eino-ext/components/model/ark
go get github.com/cloudwego/eino-ext/components/embedding/ark
go get github.com/cloudwego/eino-ext/components/indexer/milvus2
go get github.com/cloudwego/eino-ext/components/retriever/milvus2
go get github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown

go get github.com/milvus-io/milvus/client/v2
go get github.com/joho/godotenv
```

### 3. 配置环境变量

创建 `.env` 文件：

```env
ARK_API_KEY=your_api_key
MODEL=your_chat_model_name
EMBEDDER=your_embedding_model_name
ADDRESS=localhost:19530
```

四个键分别对应：对话模型用的 API Key、对话模型名、嵌入模型名、Milvus 地址。

### 4. 启动 Milvus

```bash
docker compose up -d
```

四个容器：etcd（元数据）、MinIO（对象存储）、standalone（Milvus 主服务）、Attu（Web 管理界面，`localhost:8000`）。
数据通过 Bind Mount 持久化在 `volumes/`，`docker compose down` 不会丢。

### 5. 验证安装

```bash
go mod tidy
go build ./...
```

## 运行

```bash
# 入口示例：Prompt Template + Generate
go run .

# 完整 RAG：切分 → 向量化 → 写入 Milvus → 检索
cd RAG && go run .
```

`RAG/Rag.go` 读的是相对当前目录的 `./document.md`（也就是 `RAG/document.md`），
所以要 `cd RAG` 之后再运行。

## 技术栈

- Go
- CloudWeGo Eino
- 豆包/火山引擎 Ark 模型（对话 + 嵌入）
- Milvus 向量数据库
- 自定义 Eino Agent Learning Skill
