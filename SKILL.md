---
name: maxenio-notes
description: "maxenio 自动记笔记助手。自动记录技术知识点到 docs/knowledge.md，并自动 Git 提交。触发词：学习、记录、问题、代码、讲解、解释、技术术语。"
metadata:
  author: max
  version: "1.0.0"
  license: MIT
  tags: maxenio, notes, golang, eino, learning
---

# maxenio Notes Assistant

自动化知识点记录系统，将技术知识整理到 docs/knowledge.md。

## When to use

**自动激活场景**：

1. **讨论技术内容时**：
   - Go 代码、eino 框架
   - API 调用、ARK 模型
   - Bug/报错
   - 项目架构、设计决策
   - 技术术语解释

2. **显式记录请求**：
   - 用户说："记录一下"、"保存笔记"、"写入笔记"

## How to use

### 1. 初始化检查

```bash
# 检测 Git 仓库根目录
git rev-parse --show-toplevel

# 创建笔记目录（如不存在）
mkdir -p docs

# 创建知识文件（如不存在）
# - docs/knowledge.md
```

### 2. 知识点文件结构

**docs/knowledge.md**：
```markdown
# Knowledge

## Go / Eino

### [概念名称]
[解释]

## API / ARK

### [概念名称]
[解释]
```

### 3. 内容记录

**检测到技术内容时**，提取并记录：

```python
def extract_knowledge(user_message, assistant_response):
    """
    1. 提取技术概念
       - 术语定义
       - 代码逻辑
       - API 用法

    2. 推断分类
       - Go / Eino（框架相关）
       - API / ARK（API 调用）
       - 项目架构
       - Bug 解决

    3. 更新 knowledge.md
       - 在对应分类下添加条目
       - 去重处理
    """
```

**分类关键词**：

```python
CATEGORY_KEYWORDS = {
    "Go / Eino": ["Go", "golang", "eino", "func", "interface", "struct", "context", "channel"],
    "API / ARK": ["API", "ark", "ChatModel", "Generate", "Stream", "APIKey", "Model"],
    "项目架构": ["main", "package", "import", "module", "struct", "设计"],
    "Bug 解决": ["报错", "错误", "panic", "nil", "failed", "问题"],
}
```

### 4. Git 自动化

**生成 Commit Message**：

```python
def generate_commit_message(user_message, assistant_response):
    """
    格式: "[动作] [主题]"

    动作词:
    - 添加 (新知识点)
    - 更新 (修改现有)
    - 解决 (bug 相关)
    """

    if any(kw in user_message for kw in ["报错", "错误", "panic"]):
        action = "解决"
        topic = extract_bug_topic(user_message)
    elif any(kw in user_message for kw in ["解释", "什么是", "概念"]):
        action = "添加"
        topic = extract_concept_topic(user_message)
    else:
        action = "更新"
        topic = extract_topic(user_message)

    return f"{action} {topic}"[:50]
```

**执行 Git 操作**：

```bash
cd {repo_root}
git add docs/knowledge.md
git commit -m "{generated_message}"
git push origin {current_branch}
```

**错误处理**：

```python
def safe_git_push(max_retries=3):
    for attempt in range(max_retries):
        try:
            result = run_git_push(timeout=30)
            if result.success:
                return True
        except TimeoutError:
            wait = 2 ** attempt
            sleep(wait)

    log_warning("Git push 失败，更改已提交到本地")
    return False
```

### 5. 静默运行原则

**不打扰用户**：

```python
# 完全静默，只在出错时提示
if git_push_failed:
    print("提示：更改已保存到本地，推送失败（网络问题）")
```

## Summary

**核心行为**：
1. 检测技术内容并提取知识点
2. 整理到 docs/knowledge.md
3. 自动 Git 提交和推送
4. 完全静默，不打扰用户

---

**版本**: 1.0.0
**作者**: max
**许可**: MIT