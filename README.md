# bug-carrot

## 项目进度

- [x] dice 占卜
- [x] food 吃什么
- [x] goodmorning 早安 & goodnight 晚安
- [x] homework 作业
- [x] repeat 复读
- [x] weather 天气

## Q&A

以下说明均以 `/src` 作为根目录

**Q0: 完全不懂要做什么诶...**
- 省流：回忆一下 cqhttp 是啥，bug-carrot 又是啥；然后看看 `/plugin/default.go`，里面写了一些注释，或许能够帮你找到灵感！
- 不省流：介绍一下项目结构好啦！
  - `/config` 用来读取 `../../config` 目录里的配置文件，插件不需要配置文件的话**不需要看**
  - `/constant` 放了一些静态变量，比如 bot 的语言部分，写好插件之前**不需要看**
  - `/controller` 是项目核心逻辑的目录，绝大多数情况下**不需要看**
  - `/model` 是数据库，插件没有用数据库的需求的话**不需要看**
  - `/param` 是结构体仓库，写好插件之前**不需要看**
  - `/plugin` 是插件主体目录，里面的 `default.go` 好好看，它写了函数注释，~~你也可以通过这个目录下多个插件文件的对比来速通 golang!~~
  - `/router` 是路由注册部分，绝对**不需要看**
  - `/util` 是工具部分，其中 `qq.go` 里放了一些常用的可以给 cqhttp 通信的函数，写了注释可以看看，`word.go` 不需要看

**Q0.5: 还是云里雾里？**
- 戳戳 GC，她可以手把手教你写 plugin!

**Q1: 如果我想写一个新的 plugin，需要在项目的哪些部分动手呢？**
- 在 `/plugin` 目录下仿照其他 plugin 的结构新建一个文件写逻辑代码
- 注意你应该至少需要阅读两个 plugin 并对比其差异（其中 food 和 homework 用到了数据库，所以会显得有些复杂哦）
- 在 `main.go` 的 `pluginRegister()` 中调用你 plugin 的注册函数

**Q2: 我可以通过哪些方法把我的 plugin 写的更加...格式化？**
- 把 bot 的语言部分统一写在 `/constant/index.go` 里
- 把某些配置写在 `../../config/default.yml` 里，并对应修改 `config/config.go` 中 `Plugin{}` 的结构，让它可以被你的插件读取
- 对于上一个选择，记得相应的去更新 `prod.yml`
- 把新增的结构体写到 `/param` 目录里

**Q3: 我需要调用数据库？**
- 首先，请再次确认你的 plugin 是否真的有数据库的需求，是否可以简单的被 map/slice 等程序内部结构所替代，前者在程序重启时依然能保留数据，而后者不能
- 你需要在本地环境下安装 MongoDB，运行它，它会默认开放 127.0.0.1:27017 的端口。然后，你需要新建一个名为 carrot 的数据库，在其中创建用户 carrot，密码为 123，此处参考配置文件
- 在配置文件中有一个 database-use 的选项，设置它为 true
- 在 `/model` 目录下仿照 `homework.go` 和 `food.go` 的结构写你需要的逻辑，再于同目录下 `init.go` 的 `Model{}` 中增加你的 Interface

**Q4: 我应该如何在本地测试我的 plugin?**
- pull 这个 repo，按照上面的指引写好 plugin，修改 `../../config/default.yml` 中 qqbot-qq 为你自己的 QQ 号，然后在本地跑起来!
- 在 https://docs.go-cqhttp.org/ 安装 cqhttp，修改 config.yml 的 uin 为自己的 QQ 号，修改 http 通信设置为
```yml
  - http:
      host: 0.0.0.0 # 服务端监听地址，这里是本地
      port: 5701      # 服务端监听端口，也就是 bug-carrot 给 cqhttp 发信的端口
      post:           # 反向 HTTP POST 地址列表，这是 bug-carrot 的 API
        - url: 'http://127.0.0.1:3456/api/reverse'
```
- 把 cqhttp 也跑起来，然后愉快的测试吧！

**Q5, 我应该如何向 cqhttp 通信？**
- 项目在 `/util/qq.go` 里已经写好了许多可以直接调用的函数
- 如果需要更多操作，可以阅读 https://docs.go-cqhttp.org/api 的说明，并在 `/util/qq_test.go` 里进行测试

**Q6, 我想发戳一戳 / 音乐分享 / 图片等特殊消息格式?**
- 阅读 https://docs.go-cqhttp.org/cqcode 的说明

**Q7, 我想监听除了群消息以外的更多事件？比如检测加好友请求？**
- 阅读 https://docs.go-cqhttp.org/event 的说明
- 看看 `/controller/message.go` 是怎么用 friendAddRequestHandler 函数处理加好友请求的，仿照

**Q8. `/util/word.go` 是干什么的？**
- 它调用了 gojieba 的库，实现了对一条消息的简单分词，分词的具体信息可以参阅 https://github.com/fxsjy/jieba
- 这使得我们可以调用 `param.GroupMessage.WordsMap.ExistWord("n", []string{"晚安"})` 来检查 msg 里有没有名词“晚安”
- 这在后续或许可以开发出有意思的功能(x

**Q9. 很好奇项目处理 time / private / group / listen 等多种需求的逻辑是什么？**
- 参考 `/controller/message.go` 和 `/controller/plugin.go`，前者接受 cqhttp 消息并处理，后者处理插件

**Q10. 你的 `param.GroupMessage{}` 怎么没有我想要的信息啊?**
- 阅读 https://docs.go-cqhttp.org/event 的说明
- 阅读 `/controller/message.go` 的 `groupMessageHandler()` 函数，看看它做了什么，然后修改 `param.GroupMessage{}` 吧

**Q11. CD 啥时候能用上啊？**
- 问就是马上！

**Q12. 我发现了项目~~bug~~特性！**
- 速速发 issue!

**Q13. 新的 plugin 写好啦，我应该做什么呢？**
- 戳戳 GC!(这里之后会写的更详细的)

**Q14. ???**
- 提提 issue 吧！