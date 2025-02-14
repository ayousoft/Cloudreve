package middleware

import (
	"github.com/HFO4/cloudreve/pkg/conf"
	"github.com/HFO4/cloudreve/pkg/util"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// Store session存储
var Store memstore.Store

// Session 初始化session
func Session(secret string) gin.HandlerFunc {
	// Redis设置不为空，且非测试模式时使用Redis
	if conf.RedisConfig.Server != "" && gin.Mode() != gin.TestMode {
		var err error
		Store, err = redis.NewStoreWithDB(10, "tcp", conf.RedisConfig.Server, conf.RedisConfig.Password, conf.RedisConfig.DB, []byte(secret))
		if err != nil {
			util.Log().Panic("无法连接到 Redis：%s", err)
		}

		util.Log().Info("已连接到 Redis 服务器：%s", conf.RedisConfig.Server)
	} else {
		Store = memstore.NewStore([]byte(secret))
	}

	// Also set Secure: true if using SSL, you should though
	// TODO:same-site policy
	Store.Options(sessions.Options{HttpOnly: true, MaxAge: 7 * 86400, Path: "/"})
	return sessions.Sessions("cloudreve-session", Store)
}
