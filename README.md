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
│   └── 03_template_demo/  # Prompt Template
├── docs/               # 学习笔记（自动生成）
│   ├── eino/           # Eino 主题笔记
│   │   ├── chat_model.md
│   │   └── template.md
│   └── daily/          # 每日学习记录
├── .vscode/            # VSCode 配置
├── main.go             # 主程序入口
└── SKILL.md            # 学习助手配置
```

## 运行 Demo

```bash
cd demos/01_basic_chat && go run main.go
cd demos/02_streaming_chat && go run main.go
cd demos/03_template_demo && go run main.go
```

## VSCode 环境设置

### 1. 安装扩展

- Go (by Go Team at Google)
- Delve Debugger (by Go Team)

### 2. 环境变量

确保 shell 中有 `APK_API_KEY` 和 `MODEL` 环境变量，或在 `.env` 文件中配置：

```env
APK_API_KEY=your_api_key
MODEL=your_model_name
```

### 3. Debug

按 `F5` 选择对应的 Launch 配置即可调试：

- `Launch basic_chat`
- `Launch streaming_chat`
- `Launch template_demo`

## 技术栈

- Go
- CloudWeGo Eino
- 豆包/火山引擎 Ark 模型
- 自定义 Eino Agent Learning Skill