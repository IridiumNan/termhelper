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

---

## 测试逻辑

将所有小于复习时间小于当前时间戳的都取出， 顺序进行测试

现有的疑问: 是否维护一个 提供随机中文意思的生成器？如果没有这个生成器， 那么剩下的选项从哪里来？或者说我们可以直接取出的那些单词作为词库来进行中文意思的生成， 并在现实答案的时候使用蓝色显示. 但是当单词的个数不够的时候， 我们应该再次从词库中取出单词， 但是取出多少？ 3 个似乎已经足够了
