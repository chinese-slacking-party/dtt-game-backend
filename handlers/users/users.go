package users

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
)

type UserRegisterReq struct {
	Name string `json:"name"`
}

func Register(c *gin.Context) {
	var user UserRegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	result, err := dao.CreateUser(c.Request.Context(), user.Name)
	if err != nil {
		if _, ok := err.(*db.ErrDuplicateKey); ok {
			c.JSON(http.StatusConflict, gin.H{"code": 1001, "message": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	c.SetCookie("userid", result.Name, int(db.CookieLife.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"profile": result,
	})
}

func UploadFile(c *gin.Context) {
	name := c.Param("name")
	filename := c.Param("filename")
	// TODO: check if Cookie matches user name
	userid, err := c.Cookie("userid")
	if err != nil || userid == "" {
		c.JSON(http.StatusForbidden, gin.H{"code": 1003, "message": "Not logged in"})
		return
	}
	if name == "" {
		name = userid
	}

	formData, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1004, "message": "No file in body"})
		return
	}

	dirPath := path.Join(config.PhotoDir, name)
	// TODO: Secure permissions
	if err = os.MkdirAll(dirPath, os.ModePerm); err != nil {
		log.Println("Error executing MkdirAll():", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "Unable to create directory"})
		return
	}

	filePath := path.Join(dirPath, filename)
	if _, err = os.Stat(filePath); err == nil {
		c.JSON(http.StatusConflict, gin.H{"code": 1005, "message": "File exists"})
		return
	}
	if err = c.SaveUploadedFile(formData, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "Unable to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
