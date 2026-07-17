# storage

## 处理的核心结构体

```go
type WordEntry struct {
    Word string
    SimpleDefinition string
    DetailedExplanation string

    Proficienty float32
}
```

---

## 操作类型

- 初始化 Bucket initBucket -> 初始化数据库文件

- 查找单词， 返回 WordEntry -> DBProvider 查找单词

- 查找单词并更新 Proficiency

- 插入新的单词

---

## Test 模块 调用storage中的接口

- 遍历单词， 找到下一个达到测试时间戳的单词
使用迭代器的方式，每次预先缓存一些单词和释义
维护一个用于生成随机选项的 切片以及一个 达到时间戳的切片， 两个切片协同出选项
