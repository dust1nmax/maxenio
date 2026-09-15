# Embedder

## 一句话理解

把文字转换成向量（Embedding）的组件，用于相似度搜索和 RAG。

## 具体例子

```go
// 输入文本
input := "你好"

// Embedder 输出
vector := [0.123, -0.456, 0.789, ...]
```

## 工作流程

```text
文本
  ↓
Embedder
  ↓
向量 [0.123, -0.456, ...]
  ↓
存储 / 相似度计算
```

## 应用场景

- RAG：Retrieval 阶段用 embedding 找相似文档
- 语义搜索：对比向量相似度
- 聚类分析

## Eino 中的使用

通过 Embedding 把文本转成向量后，配合 Indexer 存入 Milvus。

## 相关知识

- [Milvus](./milvus.md)
- [RAG](./rag.md)