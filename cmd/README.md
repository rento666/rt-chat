# cmd包
cmd 包是项目的主干，是编译构建的入口，main 文件通常放置在此处。一个典型的 cmd 包的目录结构如下所示：

- cmd
    - app1
        - main.go
    - app2
        - main.go

从上述例子可以看出，cmd 下可以允许挂载多个需要编译的应用，只需要在不同的包下编写 main 文件即可。需要注意的是，cmd 中的代码应该尽量「保持简洁」，main 函数中可能仅仅是参数初始化、配置加载、服务启动的操作。

---
我这个项目就一个main函数，反正截止到目前来说，是的，嗯嗯。