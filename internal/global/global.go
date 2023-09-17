package global

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	ROUTER *gin.Engine
	REDIS  *redis.Client
)
