package game

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
)

var (
	minPhotos = map[string]int{
		"1": 2,
		"2": 3,
		"3": 1,
		"4": 3,
	}

	passMsg = map[string][]string{
		"0": {
			"Amazing work! You get one point for that.",
			"Once you get 4 points, you can get a reward from your parents.",
			"Want to keep the excitement going?",
		},
	}

	failMsg = map[string][]string{
		"1": {
			"No worries at all!",
			"Every attempt is one step closer to success.",
			"Would you like to have another try?",
		},
		"2": {
			"You are doing well.",
			"I really like the effort that you are putting into this.",
			"Would you like to have another try?",
		},
		"3": {
			"No worries at all!",
			"Please try to observe their eyes, nose, mouth, and hair.",
			"These may help you find the answer.",
			"Would you like to have another try?",
		},
		"4": {
			"You are doing well.",
			"First, you can focus on the face.",
			"Second, you can focus on the eyes.",
			"Would you like to have another try?",
		},
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
	userObj, err := dao.GetUserByID(c.Request.Context(), userid)
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
		UserID:      userid,
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
	return passMsg["0"]
}

func getFailMsg(userObj *db.User, level string) []string {
	return failMsg[level]
}

func getTitle(userObj *db.User, level string) string {
	return fmt.Sprintf("Level %v", level)
}

func getTiles(userObj *db.User, level string) []db.MapTile {
	// TODO: performance
	rand.Shuffle(len(userObj.Album), func(i, j int) {
		userObj.Album[i], userObj.Album[j] = userObj.Album[j], userObj.Album[i]
	})
	ret := []db.MapTile{
		{ImageID: userObj.Album[0].Key, ImageTag: "original", URL: config.OurAddr + userObj.Album[0].Original},
		{ImageID: userObj.Album[1].Key, ImageTag: "original", URL: config.OurAddr + userObj.Album[1].Original},
		{ImageID: userObj.Album[0].Key, ImageTag: "original", URL: config.OurAddr + userObj.Album[0].Original},
		{ImageID: userObj.Album[1].Key, ImageTag: "original", URL: config.OurAddr + userObj.Album[1].Original},
	}
	rand.Shuffle(len(ret), func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	return ret
}
