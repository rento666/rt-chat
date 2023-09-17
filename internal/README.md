# internal 包

我的代码分层：

- internal
- - global (全局变量挂载)
- - object (结构体的定义到业务)
- - router (多版本api挂载)
- 

未完待续。。。

---

internal 包主要用处在于提供一个项目级别的代码保护方式，存放在其中的代码仅供**项目内部**使用。具体使用的规则是：.../a/b/c/**internal**/d/e/f 仅仅可以被.../a/b/c下的目录导入，.../a/b/g则不允许。

在 internal 内部可以继续通过命名对目录的共享范围做区分，例如 internal/myapp 表示该目录下的代码是供 myapp 应用使用的；internal/pkg 表示该目录下的代码是可以供项目内多个应用使用的。

---

存放私有应用和库代码。

如果一些代码，你不希望被其他项目/库导入，可以将这部分代码放至/internal目录下。一般存储一些比较专属于当前项目的代码包。这是在代码编译阶段就会被限制的，该目录下的代码不可被外部访问到。一般有以下子目录：

- /router 路由
- /application 存放命令与查询
- - /command
- - query
- /middleware 中间件
- /model 模型定义
- /repository 仓储层，封装数据库操作
- /response 响应
- /errmsg 错误处理

在/internal目录下应存放每个组件的源码目录，当项目变大、组件增多时，扔可以将新增的组件代码存放到/internal目录下

internal目录并不局限在根目录，在各级子目录中也可以有internal子目录，也会同样起到作用。

