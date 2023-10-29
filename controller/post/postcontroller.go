package post

import (
	"beepbop/helper"
	"beepbop/models"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetStoredPost(user_id uint) []models.Post {
	Posts := []models.Post{}
	helper.Database.Db.Where("user_id = ?", user_id).Find(&Posts)
	return Posts
}

func CreatePost(c echo.Context) error {
	user_data, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, err.Error(), nil)
	}

	type newPost struct {
		Title      string `form:"title" `
		Desc       string `form:"desc"`
		Music      string `form:"music"`
		TypeOfPost string `form:"type_of_post"`
	}

	request := newPost{}
	if err := c.Bind(&request); err != nil {
		return helper.ErrorResponse(c, "Error Validation", nil)
	}

	form, err := c.MultipartForm()
	if err != nil {
		return err
	}

	post_id := helper.RandomString(50)

	err = helper.MakeDir("./assets/posts/" + post_id)

	if err != nil {
		helper.ErrorResponse(c, err.Error(), nil)
	}

	files := form.File["files"]

	for _, file := range files {
		// ext := strings.Split(file.Filename, ".")

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		// dst, err := os.Create("./assets/posts/" + post_id + "/" + strconv.Itoa(key) + "." + ext[len(ext)-1])
		dst, err := os.Create("./assets/posts/" + post_id + "/" + file.Filename)
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

	}

	user := models.User{}
	helper.Database.Db.Select("id").Where("token = ?", user_data.Token).Find(&user)

	music_sep := strings.Split(request.Music, "-")

	if len(music_sep) < 1 {
		return c.Redirect(http.StatusSeeOther, "/404")
	}
	music_id := music_sep[len(music_sep)-1]

	new_post := models.Post{
		Title:  request.Title,
		Music:  music_id,
		Desc:   request.Desc,
		Type:   request.TypeOfPost,
		Path:   post_id,
		UserId: user.Id,
	}

	helper.Database.Db.Create(&new_post)

	return c.Redirect(http.StatusSeeOther, "/home/posts")

}

func DeletePost(c echo.Context) error {
	user_data, err := helper.JWT(c)

	if err != nil {
		return helper.ErrorResponse(c, err.Error(), nil)
	}

	post_id64, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		return helper.ErrorResponse(c, err.Error(), nil)
	}

	post_id := uint(post_id64)

	user := models.User{
		Token: user_data.Token,
	}

	helper.Database.Db.Select("id").First(&user)

	post := models.Post{
		Id:     post_id,
		UserId: user.Id,
	}
	helper.Database.Db.Delete(post)

	return c.Redirect(http.StatusSeeOther, "/home/posts")
}
