package users

import (
	"context"
	"log"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := doRegister(context.TODO(), &user)
	if err != nil {
		if _, ok := err.(*db.ErrDuplicateKey); ok {
			c.JSON(http.StatusConflict, gin.H{"code": 1001, "message": "User already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": result,
	})
}

func doRegister(ctx context.Context, user *UserRegisterReq) (*db.User, error) {
	ret, err := dao.CreateUser(ctx, user.Name)
	log.Println("doRegister", ret, err)
	return ret, err
}
