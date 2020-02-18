package routes

import (
	"time"
	"strconv"
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
	var oldItem models.Article
	slug := slug.Make(context.PostForm("title"))

	if !config.DB.First(&oldItem, "slug = ?", slug).RecordNotFound() {
		slug = slug + strconv.FormatInt(time.Now().Unix(), 10)
	}

	item := models.Article {
		Title : context.PostForm("title"),
		Desc  : context.PostForm("desc"),
		Tag  : context.PostForm("tag"),
		Slug  : slug,
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

func UpdateArticle(context *gin.Context) {
	id := context.Param("id")

	var item models.Article

	if config.DB.First(&item, "id = ?", id).RecordNotFound() {
		context.JSON(404, gin.H{"status": "error", "message": "record not found"})
		context.Abort() //membatalkan semua fungsi yang akan d jalankan di bawah
		return
	}

	if uint(context.MustGet("jwt_user_id").(float64)) != item.UserID {
		context.JSON(403, gin.H{"status": "error", "message": "this data is forbidden"})
		context.Abort()
		return
	}

	config.DB.Model(&item).Where("id = ?", id).Updates(models.Article{
		Title: context.PostForm("title"),
		Desc: context.PostForm("desc"),
		Tag: context.PostForm("tag"),
	})

	context.JSON(200, gin.H {
		"status": "success",
		"data": item,
	})
}
