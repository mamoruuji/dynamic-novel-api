package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"log"
	"net/http"

	"github.com/amomon/dynamic_novel_api/models"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetDataTodo(ctx context.Context, c *gin.Context) {
	var b models.Todo
	if err := c.Bind(&b); err != nil {
		_ = fmt.Errorf("%#v", err)
	}
	b.Status = 0
	b.UpdatedAt = time.Date(2001, 5, 20, 23, 59, 59, 0, time.UTC)
	err := b.InsertG(ctx, boil.Infer())
	if err != nil {
		_ = fmt.Errorf("%#v", err)
	}

	todos, err := models.Todos().AllG(ctx)
	if err != nil {
		_ = fmt.Errorf("Get todo error: %v", err)
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"todo": todos,
	})
}

func GetDoneTodo(ctx context.Context, c *gin.Context) {
	var b models.Todo
	if err := c.Bind(&b); err != nil {
		_ = fmt.Errorf("%#v", err)
	}

	if b.Status == 0 {
		b.Status = 1
	} else {
		b.Status = 0
	}
	b.UpdatedAt = time.Date(2001, 5, 20, 23, 59, 59, 0, time.UTC)

	_, err := b.UpdateG(ctx, boil.Whitelist("status", "updated_at"))
	if err != nil {
		_ = fmt.Errorf("Get todo error: %v", err)
	}

	todos, err := models.Todos().AllG(ctx)
	if err != nil {
		_ = fmt.Errorf("Get todo error: %v", err)
	}
	c.HTML(http.StatusOK, "index.html", map[string]interface{}{
		"todo": todos,
	})

}

func main() {
	ctx := context.Background()
	db, err := sql.Open("mysql", "root:rootpass@tcp(container_db:3306)/boxers?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		log.Fatalf("Cannot connect database: %v", err)
	}
	boil.SetDB(db)

	r := gin.Default()
	r.LoadHTMLFiles("./tmpl/index.html")
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "world",
		})
	})
	r.GET("/todo", func(c *gin.Context) {

		todos, err := models.Todos(OrderBy("updated_at desc")).AllG(ctx)
		if err != nil {
			_ = fmt.Errorf("Get todo error: %v", err)
		}

		c.HTML(http.StatusOK, "index.html", map[string]interface{}{
			"todo": todos,
		})
	})
	r.GET("/yaru", func(c *gin.Context) {
		GetDataTodo(ctx, c)
	})
	r.GET("/done", func(c *gin.Context) {
		GetDoneTodo(ctx, c)
	})
	_ = r.Run() // listen and serve on 0.0.0.0:8080

}
