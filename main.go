package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// 自定义拦截器
func myHandler() gin.HandlerFunc {

	return func(context *gin.Context) {
		//通过自定义中间件设置的值，后续处理中只要调用了此中间件就可以拿到这里的参数
		context.Set("userSession", "userid-1")
		context.Next() //放行

		//context.Abort() //阻止
	}
}
func main() {

	//1.创建服务
	ginServer := gin.Default()

	//2.加载静态页面
	ginServer.LoadHTMLGlob("templates/*")

	//3.加载资源文件
	ginServer.Static("/static", "./static")

	//4.响应页面给前端
	ginServer.GET("/hello", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", gin.H{
			"msg": "Hello World!!!",
		})
	})

	//5.接收前端传来参数
	//5.1.路径参数
	//http://localhost:8080/user/info?userId=67&userName=peixiaoze
	ginServer.GET("/user/info", myHandler(), func(context *gin.Context) {
		//去出中间件的值
		userSession := context.MustGet("userSession").(string)
		log.Println("=================>", userSession)

		userId := context.Query("userId")
		userName := context.Query("userName")
		context.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	//5.2.变量参数
	//http://localhost:8080/user/info/1/peixiaoze
	ginServer.GET("/user/info/:userId/:userName", func(context *gin.Context) {
		userId := context.Param("userId")
		userName := context.Param("userName")
		context.JSON(http.StatusOK, gin.H{
			"userId":   userId,
			"userName": userName,
		})
	})

	//5.3.JSON对象
	ginServer.POST("/json", func(context *gin.Context) {
		//从请求体里获取json对象
		data, _ := context.GetRawData()

		var m map[string]interface{}
		_ = json.Unmarshal(data, &m)
		context.JSON(http.StatusOK, m)
	})

	//5.4.form表单
	ginServer.POST("/user/add", func(context *gin.Context) {
		username := context.PostForm("username")
		password := context.PostForm("password")

		context.JSON(http.StatusOK, gin.H{
			"message":  "OK!!!",
			"username": username,
			"password": password,
		})
	})

	//6.路由
	//6.1.301
	ginServer.GET("/test", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	})

	//6.2.404
	ginServer.NoRoute(func(context *gin.Context) {
		context.HTML(http.StatusNotFound, "404.html", nil)
	})

	//7.路由组
	userGroup := ginServer.Group("/user")
	{
		userGroup.GET("add")
		userGroup.GET("login")
		userGroup.GET("logout")
	}

	//8.服务器端口
	ginServer.Run(":8080")

}
