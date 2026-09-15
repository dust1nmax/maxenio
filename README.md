# Maxenio

CloudWeGo Eino + Go Agent 学习项目。边学边记，笔记即代码。

自带学习助手。SKILL.md 定义了一个 Eino Agent 学习规则：

- 陪用户学习 Eino / Go / Agent 工程
- 自动把讨论内容整理成 Markdown 笔记到 `docs/`
- 主题笔记（eino/template.md） + 每日记录（daily/2026-09-14.md）

```text
学习 → 讨论 → 自动生成 docs/笔记 → 形成可复习的知识库
```

## 项目结构

```
├── demos/              # 独立可运行的示例
│   ├── 01_basic_chat/     # Template + Generate
│   ├── 02_streaming_chat/ # Stream 流式输出
│   ├── 03_template_demo/  # Prompt Template
│   ├── 05_indexer/        # Milvus Indexer + RAG 索引
│   └── 06_retriever/      # Milvus Retriever + RAG 检索
├── docs/               # 学习笔记（自动生成）
│   ├── eino/           # Eino 主题笔记
│   │   ├── chat_model.md
│   │   ├── embedder.md
│   │   ├── milvus.md
│   │   └── template.md
│   └── daily/          # 每日学习记录
├── main.go             # 主程序入口
└── SKILL.md            # 学习助手配置
```

## 从零配置环境

### 1. 初始化项目

```bash
mkdir my-eino-project && cd my-eino-project
go mod init my-eino-project
```

### 2. 安装依赖

```bash
go get github.com/cloudwego/eino
go get github.com/cloudwego/eino-ext/components/model/ark
go get github.com/cloudwego/eino/schema
go get github.com/cloudwego/eino/components/prompt
go get github.com/joho/godotenv
```

### 3. 配置环境变量

创建 `.env` 文件：

```env
APK_API_KEY=your_api_key
MODEL=your_model_name
```

### 4. 验证安装

```bash
go mod tidy
go run .
```

## 运行 Demo

```bash
cd demos/01_basic_chat && go run main.go
cd demos/02_streaming_chat && go run main.go
cd demos/03_template_demo && go run main.go
cd demos/05_indexer && go run main.go
cd demos/06_retriever && go run main.go
```

## 技术栈

- Go
- CloudWeGo Eino
- 豆包/火山引擎 Ark 模型
- Milvus 向量数据库
- 自定义 Eino Agent Learning Skill