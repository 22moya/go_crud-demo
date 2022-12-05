package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strconv"
	"time"
	_ "time"
)

func main() {
	dsn := "root:cengliu0106.@tcp(127.0.0.1:3306)/go-crud-list?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			//解决查表的时候会自动添加复数的问题
			SingularTable: true,
		},
	})
	fmt.Println(err)
	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
	//结构体
	type List struct {
		gorm.Model        //处理添加主键
		Name       string `gorm:"type:varchar(20);not null" json:"name" binding:"required"`
		State      string `gorm:"type:varchar(20);not null" json:"state" binding:"required"`
		Phone      string `gorm:"type:varchar(20);not null" json:"Phone" binding:"required"`
		Email      string `gorm:"type:varchar(40);not null" json:"email" binding:"required"`
		Address    string `gorm:"type:varchar(200);not null" json:"address" binding:"required"`
	}
	db.AutoMigrate(&List{}) //迁移
	r := gin.Default()
	//测试
	/*r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "请求成功",
		})
	})*/
	//增
	r.POST("/user/add", func(context *gin.Context) {
		var data List
		err := context.ShouldBindJSON(&data)
		if err != nil {
			context.JSON(200, gin.H{
				"msg":  "添加失败",
				"data": gin.H{},
				"code": 400,
			},
			)
		} else {
			//数据库操作
			db.Create(&data) //创建一条数据
			context.JSON(200, gin.H{
				"msg":  "添加成功",
				"data": data,
				"code": 200,
			})
		}
	})
	//删
	r.DELETE("/user/delete/:id", func(context *gin.Context) {
		var data []List

		//1.接受id

		id := context.Param("id")
		//2.判断id是否存在
		db.Where("id=?", id).Find(&data)
		if len(data) == 0 {
			context.JSON(200, gin.H{
				"msg":  "id不存在",
				"code": 400,
			})

		} else {
			//操作数据库删除
			db.Where("id = ?", id).Delete(&data)
			context.JSON(200, gin.H{
				"msg":  "删除成功",
				"code": 200,
			})
		}

	})
	//改
	r.PUT("/user/update/:id", func(context *gin.Context) {
		var data List

		//1.接受id
		id := context.Param("id")
		//2.判断id是否存在
		db.Select("id").Where("id=?", id).Find(&data)
		if data.ID == 0 {
			context.JSON(200, gin.H{
				"msg":  "id不存在",
				"code": 400,
			})

		} else {
			err := context.ShouldBindJSON(&data)
			if err != nil {
				context.JSON(400, gin.H{
					"msg":  "更新数据失败",
					"code": 400,
				})
			} else {
				//操作数据库更改
				db.Where("id = ?", id).Updates(&data)
				context.JSON(200, gin.H{
					"msg":  "更新数据成功",
					"code": 200,
				})
			}

		}
	})
	//查 （1.条件查询   2.分页查询）
	r.GET("/user/list/:name", func(context *gin.Context) {
		name := context.Param("name")
		var datalist []List
		//查询数据库
		db.Where("name=?", name).Find(&datalist)
		//判断是否查询到数据
		if len(datalist) == 0 {
			context.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			context.JSON(200, gin.H{
				"msg":  "查询数据成功",
				"code": 200,
				"data": datalist,
			})
		}
	})
	//全部查询
	r.GET("/user/list", func(context *gin.Context) {
		var dataList []List
		pageSize, _ := strconv.Atoi(context.Query("pageSize"))
		pageNum, _ := strconv.Atoi(context.Query("pageNum"))
		//判断是否需要分页
		if pageSize == 0 {
			pageSize = -1
		}
		if pageNum == 0 {
			pageNum = -1
		}
		offsetVal := (pageNum - 1) * pageSize
		if pageNum == -1 && pageSize == -1 {
			offsetVal = -1
		}
		var total int64
		db.Model(dataList).Count(&total).Limit(pageSize).Offset(offsetVal).Find(&dataList)
		if len(dataList) == 0 {
			context.JSON(200, gin.H{
				"msg":  "没有查询到数据",
				"code": 400,
				"data": gin.H{},
			})
		} else {
			context.JSON(200, gin.H{
				"msg":  "查询成功",
				"code": 200,
				"data": gin.H{
					"list":     dataList,
					"total":    total,
					"pageNum":  pageNum,
					"pageSize": pageSize,
				},
			})

		}
	})
	r.Run(":8181")
}

//func main() {
//	type b map[int]int
//	var a map[int]int = b{
//		1: 3,
//		3: 4,
//	}
//	a[4] = 5
//	fmt.Printf("%v", a)
//	fmt.Printf("%v", a[1])
//	fmt.Printf("%v", a[3])
//	fmt.Printf("%v", a[4])
//	type c map[int]int
//	c2:=c{1: 2}
//	fmt.Printf("%v\n", c2)
//
//}
