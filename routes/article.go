package routes

import (
	"github.com/wawandx/rest-api-gin/config"
	"github.com/wawandx/rest-api-gin/models"
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
)

func GetHome(context *gin.Context) {
	items := []models.Article{}
	config.DB.Find(&items)

	context.JSON(200, gin.H {
		"status": "success",
		"data": items,
	})
}

func GetArticle(context *gin.Context) {
	slug := context.Param("slug")

	var item models.Article

	if config.DB.First(&item, "slug = ?", slug).RecordNotFound() {
		context.JSON(404, gin.H{"status": "error", "message": "record not found"})
		context.Abort() //membatalkan semua fungsi yang akan d jalankan di bawah
		return
	}

	context.JSON(200, gin.H {
		"status": "success",
		"data": item,
	})
}

func PostArticle(context *gin.Context) {
	item := models.Article {
		Title : context.PostForm("title"),
		Desc  : context.PostForm("desc"),
		Tag  : context.PostForm("tag"),
		Slug  : slug.Make(context.PostForm("title")),
		UserID: uint(context.MustGet("jwt_user_id").(float64)),
	}

	config.DB.Create(&item)

	context.JSON(200, gin.H {
		"status": "success",
		"data": item,
	})
}

func GetArticleByTag(context *gin.Context) {
	tag := context.Param("tag")
	items := []models.Article{}

	config.DB.Where("tag LIKE ?", "%" + tag + "%").Find(&items)

	context.JSON(200, gin.H {
		"data": items,
	})
}
