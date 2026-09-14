---
name: maxenio-notes
description: "maxenio 自动记笔记助手。自动记录所有对话内容、技术知识、问题解决方案到 docs/ 目录，并自动 Git 提交。触发词：开始、记录、问题、代码、讲解、解释。"
metadata:
  author: max
  version: "1.0.0"
  license: MIT
  tags: maxenio, notes, golang, eino, learning
---

# maxenio Notes Assistant

自动化笔记系统，记录所有对话内容到 docs/ 目录。

## When to use

**自动激活场景**：

1. **任何对话后**：记录所有对话内容
2. **显式记录请求**：
   - 用户说："记录一下"、"保存笔记"、"写入笔记"

## How to use

### 1. 初始化检查

首次激活时，确保笔记结构存在：

```bash
# 检测 Git 仓库根目录
git rev-parse --show-toplevel

# 创建笔记目录（如不存在）
mkdir -p docs

# 创建笔记文件（如不存在）
# - docs/notes.md
# - docs/knowledge.md
```

### 2. 笔记文件结构

**docs/notes.md** - 对话记录：
```markdown
# Notes

## 2026-09-14

### 对话 1
**用户**: ...
**助手**: ...

### 对话 2
**用户**: ...
**助手**: ...
```

**docs/knowledge.md** - 知识整理：
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

**每次对话后**，提取并记录：

```python
def record_conversation(user_message, assistant_response):
    """
    1. 提取关键信息
       - 用户问题/请求
       - 助手回答/代码
       - 技术术语
       - 代码示例

    2. 更新 notes.md
       - 追加新对话记录
       - 按日期分组

    3. 更新 knowledge.md
       - 提取新知识点
       - 归类整理
    """
```

### 4. Git 自动化

**生成 Commit Message**：

```python
def generate_commit_message(user_message, assistant_response):
    """
    格式: "[动作] [主题]"

    动作词:
    - 添加 (新内容/代码)
    - 更新 (修改现有)
    - 记录 (笔记/知识)
    - 解决 (问题/bug)
    """

    # 提取主要动作
    if "报错" in user_message or "错误" in user_message:
        action = "解决"
        topic = extract_error_topic(user_message)
    elif "代码" in user_message or "实现" in user_message:
        action = "添加"
        topic = extract_code_topic(assistant_response)
    elif "解释" in user_message or "什么是" in user_message:
        action = "记录"
        topic = extract_concept_topic(user_message)
    else:
        action = "更新"
        topic = extract_topic(user_message)

    return f"{action} {topic}"
```

**执行 Git 操作**：

```bash
cd {repo_root}
git add docs/
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

## Validation

**检查笔记文件**：

```bash
# 验证笔记文件存在
ls -la docs/notes.md docs/knowledge.md

# 验证 Git 状态
git status docs/
```

## Summary

**核心行为**：
1. 自动记录所有对话内容
2. 提取并整理技术知识
3. 自动 Git 提交和推送
4. 完全静默，不打扰用户

**用户体验**：
- 对话结束 → 笔记自动保存
- 完全静默 → 无打扰
- Git 记录 → 所有修改可追溯

---

**版本**: 1.0.0
**作者**: max
**许可**: MIT