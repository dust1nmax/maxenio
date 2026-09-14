---
name: maxenio-notes
description: "maxenio 手动记笔记助手。用户说「记录」时将知识点写入 docs/knowledge.md，并自动 Git 提交。"
metadata:
  author: max
  version: "1.0.0"
  license: MIT
  tags: maxenio, notes, golang, eino, learning
---

# maxenio Notes Assistant

知识点记录系统，当用户说「记录」时保存到 docs/knowledge.md。

## When to use

**显式触发**：
- 用户说："记录一下"、"保存笔记"、"记一下"、"记录这个"

## How to use

### 1. 检测触发词

```python
TRIGGER_WORDS = ["记录", "保存", "记一下", "记录这个", "写入笔记"]

if any(word in user_message for word in TRIGGER_WORDS):
    extract_and_save_knowledge()
```

### 2. 提取知识点

从对话中提取：
- 技术概念/术语
- 代码逻辑
- API 用法
- 问题解决方案

### 3. 更新 knowledge.md

```markdown
## [分类]

### [概念名称]
[解释]
[代码示例（如有）]
```

### 4. Git 自动化

```bash
git add docs/knowledge.md
git commit -m "添加 [概念名称] 知识点"
git push
```

### 5. 静默原则

成功时不提示，失败时简短提示。

## Summary

- 仅在用户说「记录」时触发
- 保存知识点到 docs/knowledge.md
- 自动 Git 提交推送

---

**版本**: 1.0.0
**作者**: max
**许可**: MIT