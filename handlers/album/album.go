package album

import (
	"fmt"
	"net/http"
	"path"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
)

type AddPhotoReq struct {
	Desc string `json:"desc"`
	File string `json:"filename"`
}

func AddPhoto(c *gin.Context) {
	userid, err := c.Cookie("userid")
	if err != nil || userid == "" {
		userid = "wolfgang"
	}
	userObj, err := dao.GetUserByName(c.Request.Context(), userid)
	if err != nil {
		if err == db.ErrNotFound {
			c.JSON(http.StatusForbidden, gin.H{"code": 1006, "message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "Query failed"})
		}
		return
	}

	var req AddPhotoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	dao.IncrPhotoSeq(c.Request.Context(), userObj.ID)
	origFilePath := path.Join(config.PhotoDir, userObj.Name, req.File)
	origFileURL := config.OurAddr + path.Join("/api/v1/files", userObj.Name, req.File)
	picKey := fmt.Sprintf("%s-%03d", userObj.Name, userObj.NextPicSeq)
	var x = db.Photo{
		Seq:      userObj.NextPicSeq,
		Key:      picKey,
		Desc:     req.Desc,
		Original: origFilePath,
		UserID:   userObj.ID.Hex(),
		URLs: map[string]string{
			"normal": origFileURL,
		},
	}
	if err = dao.AddPhoto(c.Request.Context(), &x); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "AddPhoto() failed"})
		return
	}
	// TODO: Generate variants!!!
	c.JSON(http.StatusOK, gin.H{
		"id":       userObj.NextPicSeq,
		"key":      picKey,
		"progress": 0,    // Generate
		"message":  "OK", // Generate
	})
}
