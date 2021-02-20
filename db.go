package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// 跟数据库相关的操作
var db *sqlx.DB

func initDB()(err error){
	dsn:="wang:123@tcp(127.0.0.1:3306)/go_shuji"
	db,err=sqlx.Connect("mysql",dsn)
	if err!=nil{
		return err
	}
	db.SetConnMaxLifetime(100)
	db.SetMaxIdleConns(16)
	return
}

//查数据
func queryAllBook()(booklist []*Book,err error){
	sqlStr:="select id,title,price from book;"
	err=db.Select(&booklist,sqlStr)
	if err!=nil{
		fmt.Println("查询所有书籍信息失败")
		return
	}
	return
}

//插入数据
func insertBook(title string, price float64) (err error) {
	sqlStr:="insert into book(title,price) value (?,?)"
	_,err=db.Exec(sqlStr,title,price)
	if err!=nil{
		fmt.Println("插入书籍信息失败")
		return
	}
	return
}