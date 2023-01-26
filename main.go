package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" //一定记得导入数据库驱动,否则运行时会报错
	"net/http"
	"xorm.io/xorm"
)

type Todo struct {
	Id     int    `json:"id"` //默认0
	Title  string `json:"title"`
	Status bool   `json:"status"` //默认false
}

var (
	userName  string = "root"
	password  string = "123456789"
	ipAddress string = "127.0.0.1"
	port      int    = 3306
	dbName    string = "test"
	charset   string = "utf8mb4"
)

func main() {

	//数据库连接

	//构建数据库连接信息
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", userName, password, ipAddress, port, dbName, charset)

	engine, err := xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		println("数据库连接失败")
	}
	println(engine)
	ginServer := gin.Default()

	v1Gruop := ginServer.Group("/v1")
	{
		//查看所有待办事项
		v1Gruop.GET("/todo", func(context *gin.Context) {
			//通过find查询多条数据
			var todolists []Todo
			err := engine.Limit(10, 0).Find(&todolists)
			if err != nil {
				return
			} else {
				context.JSON(http.StatusOK, todolists)
			}
		})
		//查看某一个待办事项
		v1Gruop.GET("/todo/:id", func(context *gin.Context) {

		})
		//删除某一个待办事项
		v1Gruop.DELETE("/todo/:id", func(context *gin.Context) {
			todo := Todo{}
			id := context.Param("id")
			_, err := engine.Where("id=?", id).Delete(&todo)
			if err != nil {
				context.JSON(http.StatusBadRequest, gin.H{
					"code":    500,
					"message": "删除失败",
				})
				return
			} else {
				context.JSON(http.StatusOK, gin.H{
					"code":    200,
					"message": "删除成功",
				})
			}
		})
		//添加待办事项
		v1Gruop.POST("/todo", func(context *gin.Context) {
			var to Todo
			//将前端发送来的数据与结构体绑定
			context.BindJSON(&to)
			//执行插入操作
			result, err := engine.Insert(&to)
			println(result)
			//判断插入是否成功
			if result >= 1 {
				fmt.Println("插入成功")
			} else {
				fmt.Println(err)
				fmt.Println("插入失败")
			}

		})
		//更新待办事项
		v1Gruop.PUT("/todo/:id", func(context *gin.Context) {
			//获取需要更新的id
			id := context.Param("id")
			var to Todo
			context.BindJSON(&to)
			//数据库更新
			_, err := engine.Exec("update todo set status=? where id=?", to.Status, id)
			if err != nil {
				fmt.Printf("%v", err)
				context.JSON(http.StatusOK, gin.H{
					"message": "修改失败",
				})
				return
			} else {
				context.JSON(http.StatusOK, gin.H{
					"message": "修改成功",
				})
			}

		})
	}
	ginServer.Run(":9000")
}
