// file: main.go
package main

import (
	"fmt"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/pprof"
	iris_redis "github.com/kataras/iris/v12/sessions/sessiondb/redis"
	prometheusMiddleware "novel-server/prometheus"
	tools "novel-server/tools"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	routes "novel-server/web/routes"
)

func main() {
	config, err := tools.GetAppConfig("config.yaml")
	// 创建日志记录器
	customLogger := logger.New(logger.Config{
		Status: true, // 记录响应状态码
		IP:     true, // 记录请求IP地址
		Method: true, // 记录请求方法
		Path:   true, // 记录请求路径
		Query:  true, // 记录查询参数
	})

	if err != nil {
		panic(err)
	}
	tools.InitRedisClient()
	// 初始化 MySQL 客户端
	tools.InitMySQLClient()
	// Access configuration fields as needed
	fmt.Println("App Name:", config.AppName)
	fmt.Println("Port:", config.Port)
	app := iris.New()
	app.Use(prometheusMiddleware.New("novel", 0.3, 1.2, 5.0).ServeHTTP)
	// 将日志记录器添加到中间件
	app.Use(customLogger)
	iris_redis_db := iris_redis.New(iris_redis.Config{
		Network:   "tcp",
		Addr:      config.Redis.Addr + ":" + config.Redis.Port,
		Password:  config.Redis.Password, // Specify your Redis password if required
		Timeout:   time.Duration(30) * time.Second,
		MaxActive: 10,
		Prefix:    "novel_session_id-",
		Driver:    iris_redis.GoRedis(),
	})
	tools.InitSessionManager()
	sessionManager := tools.SessionClient()
	sessionManager.UseDatabase(iris_redis_db)
	app.Use(sessionManager.Handler())
	// 转换404 500
	//app.OnErrorCode(iris.StatusNotFound, controllers.NotFound)
	//app.OnErrorCode(iris.StatusInternalServerError, controllers.InternalServerError)

	app.Logger().SetLevel(config.LoggerLevel)
	// 注册路由
	routes.RegisterRoutes(app)
	pprof_type := pprof.New()
	app.Any("/debug/pprof", pprof_type)
	app.Any("/debug/pprof/{action:path}", pprof_type)
	app.Get("/metrics", iris.FromStd(promhttp.Handler()))
	app.HandleDir("/static", "./static")
	app.Listen(":"+config.Port, iris.WithOptimizations)
}
