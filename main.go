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
//具体删除某本书
func deleteBookHandler(c *gin.Context){
	//取query string参数
	idStr :=c.Query("id")
	idVal,err:= strconv.ParseInt(idStr,10,64)  //将字符串转换为int
	if err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":1,
			"msg":err,
		})
		return
	}
	//数据是个正常数字，去数据库删除具体的记录
	err=deleteBook(idVal)
	if err!=nil{
		c.JSON(http.StatusOK,gin.H{
			"code":1,
			"msg":err,
		})
		return
	}
	//删除成功还是跳转到书籍列表页
	c.Redirect(http.StatusMovedPermanently,"/book_list")

}
func updateBookHandler(c *gin.Context){
	//需要给模板渲染上原来的旧数据
	//1.取到用户编辑的内容  从querystring取到要编辑的书的ID值
	bookIdStr:=c.Query("id")
	if len(bookIdStr)==0{
		//请求没有携带参数，该请求无效
		c.String(http.StatusBadRequest,"无效的请求")
	}
	//HTTP请求传过来的参数通常都是string类型，要根据自己的需要去转换成相应的数据类型
	bookId,err:=strconv.ParseInt(bookIdStr,10,64)  //转换成64位10进制
	if err!=nil{
		//请求中没有携带要用的数据，该请求是无效的
		c.String(http.StatusBadRequest,"无效请求")
		return
	}
	if c.Request.Method == "POST"{
		//1.获取用户提交的数据
		titleVal:=c.PostForm("title")
		priceStr:=c.PostForm("price")
		priceVal,err:=strconv.ParseFloat(priceStr,64)
		if err!=nil{
			c.String(http.StatusBadRequest,"无效的价格信息")
			return
		}
		//2.去数据库更新对应的书籍
		err=editBook(titleVal,priceVal,bookId)
		if err!=nil{
			c.String(http.StatusInternalServerError,"更新数据失败")
			return
		}
		//3.跳转回book/list页面查看是否修改成功
		c.Redirect(http.StatusMovedPermanently,"/book_list")
		//相同网站可以写相对路径，不同网站需要写绝对路径


	}else{

		//2.根据id获取书籍信息
		bookObj,err:=queryBookById(bookId)
		if err!=nil{
			//请求中没有携带要用的数据，该请求是无效的
			c.String(http.StatusBadRequest,"无效的书籍id")
			return
		}
		//3.把数据渲染到页面上
		c.HTML(http.StatusOK,"book_edit.html",bookObj)
	}

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
	r.GET("/book/delete",deleteBookHandler)
	r.Any("book/edit",updateBookHandler)

	r.Run()

}
