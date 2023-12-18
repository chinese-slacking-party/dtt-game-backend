package game

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"path"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
)

const (
	imgTagNormal  = "normal"
	imgTagChanged = "changed"
)

var (
	minPhotos = map[string]int{
		"1": 2,
		"2": 3,
		"3": 1,
		"4": 3,
	}
)

func Start(c *gin.Context) {
	// TODO: Abstraction
	// (This is copied from album.AddPhoto())
	userid, err := c.Cookie("userid")
	if err != nil || userid == "" {
		c.JSON(http.StatusForbidden, gin.H{"code": 1003, "message": "Not logged in"})
		return
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
	if userObj.Album, err = dao.LoadUserPhotos(c.Request.Context(), userObj.ID.Hex()); err != nil {
		log.Printf("Error loading user %s photos: %+v", userObj.ID.String(), err)
		return
	}

	level := c.Param("level")
	if _, ok := minPhotos[level]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"code": 1000, "message": "No such level"})
		return
	}
	if !hasEnoughPhotos(userObj, level) {
		c.JSON(http.StatusPreconditionFailed, gin.H{
			"code":    1007,
			"message": fmt.Sprintf("You need to upload at least %d photos before starting this level", minPhotos[level]),
		})
		return
	}

	var ret = db.Game{
		Title:       getTitle(userObj, level),
		SuccessText: getPassMsg(userObj, level),
		FailureText: getFailMsg(userObj, level),
		Width:       2, // Hard-coded: all our 4 levels have 2x2 maps
		Height:      2,
		LevelID:     level, // TODO: Use ObjectID from collection `levels`
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Status:      "initial",
		UserID:      userObj.Name,
	}
	tileList := getTiles(userObj, level)
	for i := 0; i < 2; i++ {
		ret.Map = append(ret.Map, nil)
		for j := 0; j < 2; j++ {
			ret.Map[i] = append(ret.Map[i], tileList[i*2+j])
		}
	}
	// TODO: Save to DB

	c.JSON(http.StatusOK, &ret)
}

func Finish(c *gin.Context) {
	// TODO
}

func hasEnoughPhotos(userObj *db.User, level string) bool {
	return len(userObj.Album) >= minPhotos[level]
}

func getPassMsg(userObj *db.User, level string) []string {
	return config.PassMsg["0"]
}

func getFailMsg(userObj *db.User, level string) []string {
	return config.FailMsg[level]
}

func getTitle(userObj *db.User, level string) string {
	return fmt.Sprintf("Level %v", level)
}

func getTiles(userObj *db.User, level string) []db.MapTile {
	// TODO: performance
	rand.Shuffle(len(userObj.Album), func(i, j int) {
		userObj.Album[i], userObj.Album[j] = userObj.Album[j], userObj.Album[i]
	})
	var ret []db.MapTile
	switch level {
	case "1":
		ret = []db.MapTile{
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[1].Key, ImageTag: imgTagNormal, URL: userObj.Album[1].URLs[imgTagNormal]},
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[1].Key, ImageTag: imgTagNormal, URL: userObj.Album[1].URLs[imgTagNormal]},
		}
	case "2":
		ret = []db.MapTile{
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[1].Key, ImageTag: imgTagNormal, URL: userObj.Album[1].URLs[imgTagNormal]},
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[2].Key, ImageTag: imgTagNormal, URL: userObj.Album[2].URLs[imgTagNormal]},
		}
	case "3":
		ret = []db.MapTile{
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			*getPreloadedMapTile(0, 5),
			*getPreloadedMapTile(5, 5),
		}
	case "4":
		ret = []db.MapTile{
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagNormal, URL: userObj.Album[0].URLs[imgTagNormal]},
			{ImageID: userObj.Album[0].Key, ImageTag: imgTagChanged, URL: userObj.Album[0].URLs[imgTagChanged]},
			{ImageID: userObj.Album[1].Key, ImageTag: imgTagNormal, URL: userObj.Album[1].URLs[imgTagNormal]},
			{ImageID: userObj.Album[2].Key, ImageTag: imgTagNormal, URL: userObj.Album[2].URLs[imgTagNormal]},
		}
	}
	rand.Shuffle(len(ret), func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	return ret
}

func getPreloadedMapTile(skip int, limit int) *db.MapTile {
	idx := rand.Intn(limit) + skip
	imageID := fmt.Sprintf("00system-%03d", idx)
	return &db.MapTile{
		ImageID:  imageID,
		ImageTag: imgTagNormal,
		URL:      config.OurAddr + path.Join(config.APIFilesFull, imageID+".png"),
	}
}
