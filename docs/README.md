# swagger接口文档

## 1. 在router包下挂载uri

示例：
```go
package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	//...
)

func Router() *gin.Engine {
	r := gin.New()
	// 全局中间件
	
	// v1版本,项目可多版本共存
	v1 := r.Group("/v1")
	{
		// 挂载swagger
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return r
}

```

#### 2. 在`objectview`包下的`objecthandler.go`中，为方法添加swagger注释。
在`pkg/tools/any/anyservice.go`文件中，添加注释。其中`any`为任意字符。

#### 3. 在命令行输入：`swag init`,就会更新`doc`目录下的文件，然后运行项目后，进入`/swagger/index.html`即可看到更新后的内容