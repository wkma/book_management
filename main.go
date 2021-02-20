package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//BookManagementSystem
func bookListHandler(c *gin.Context){
	//连数据库  查数据  返回给浏览器
	booklist,err:=queryAllBook()
	if err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":0,
			"msg":err,
		})
		return
	}
	//返回数据
	c.HTML(http.StatusOK,"book_list.html",gin.H{
		"code":1,
		"data":booklist,
	})
}

func newBookHandler(c *gin.Context){
	//给用户返回一个添加书籍的处理函数
	c.HTML(http.StatusOK,"book/new_book.html",nil)
}

func createBookHandler(c *gin.Context){
	//创建书籍的处理函数
	//从form表单提取数据
	var msg string
	titleVal:=c.PostForm("title")
	priceVal:=c.PostForm("price")
	price,err:=strconv.ParseFloat(priceVal,64)
	if err!=nil{
		msg="无效的参数"
		c.JSON(http.StatusOK,gin.H{
			"msg":msg,
		})
		return
	}
	fmt.Printf("%T %T\n",titleVal,priceVal)
	//拿到数据 插入数据库
	err=insertBook(titleVal,price)
	if err!=nil{
		msg="插入数据失败，请重试"
	c.JSON(http.StatusOK,gin.H{
		"msg":msg,
	})
		return
	}
	//插入数据库成功
	c.Redirect(http.StatusMovedPermanently,"/book_list")

}

func main(){
	//程序启动就应该连接数据库
	err:=initDB()
	if err!=nil{
		panic(err)
	}
	r:=gin.Default()
	r.LoadHTMLGlob("templates/*")
	//查看所有书籍
	r.GET("/book_list",bookListHandler)
	//返回一个页面给用户填写新增的书籍信息
	r.GET("/book/new",newBookHandler)
	r.POST("/book/new",createBookHandler)
	r.Run()

}
