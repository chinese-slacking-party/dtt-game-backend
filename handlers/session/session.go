package session

import (
	"context"
	"log"
	"net/http"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
)

type UserLoginReq struct {
	Name string `json:"name"`
}

func Login(c *gin.Context) {
	var user UserLoginReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	result, err := doLogin(context.TODO(), &user)
	if err != nil {
		if err == db.ErrNotFound {
			if result, err = dao.CreateUser(c.Request.Context(), user.Name); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": err.Error()})
				return
			}
			c.SetCookie("userid", result.Name, int(config.CookieLife.Seconds()), "/", "", false, true)
			c.JSON(http.StatusOK, gin.H{
				"profile": result,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	c.SetCookie("userid", result.Name, int(config.CookieLife.Seconds()), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"profile": result,
	})
}

func doLogin(ctx context.Context, user *UserLoginReq) (*db.User, error) {
	userObj, err := dao.GetUserByName(ctx, user.Name)
	log.Println("Loaded user", userObj, err)
	if err != nil {
		return nil, err
	}
	if userObj.Album, err = dao.LoadUserPhotos(ctx, userObj.ID.Hex()); err != nil {
		log.Printf("Error loading user %s photos: %+v", userObj.ID.String(), err)
		return nil, err
	}
	return userObj, nil
}
