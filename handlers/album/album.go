package album

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"

	"github.com/chinese-slacking-party/dtt-game-backend/config"
	"github.com/chinese-slacking-party/dtt-game-backend/db"
	"github.com/chinese-slacking-party/dtt-game-backend/db/dao"

	"github.com/gin-gonic/gin"
	repl "github.com/replicate/replicate-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddPhotoReq struct {
	Desc string `json:"desc"`
	File string `json:"filename"`
}

func AddPhoto(c *gin.Context) {
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

	var req AddPhotoReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1000, "message": err.Error()})
		return
	}

	dao.IncrPhotoSeq(c.Request.Context(), userObj.ID)
	picKey := fmt.Sprintf("%s-%03d", userObj.Name, userObj.NextPicSeq)
	origFilePath := path.Join(config.PhotoDir, userObj.Name, req.File)
	newFilePath := path.Join(config.PhotoDir, userObj.Name, picKey+config.AIGenSuffix)
	origFileURL := config.OurAddr + path.Join(config.APIFilesFull, userObj.Name, req.File)
	newFileURL := config.OurAddr + path.Join(config.APIFilesFull, userObj.Name, picKey+config.AIGenSuffix)

	var x = db.Photo{
		Seq:      userObj.NextPicSeq,
		Key:      picKey,
		Desc:     req.Desc,
		Original: origFilePath,
		UserID:   userObj.ID.Hex(),
		URLs: map[string]string{
			"normal":  origFileURL,
			"changed": newFileURL,
		},
	}
	if err = dao.AddPhoto(c.Request.Context(), &x); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "AddPhoto() failed"})
		return
	}
	if err = changeClothes(c.Request.Context(), x.ID, origFileURL, newFilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 1000, "message": "changeClothes() failed " + err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"id":       userObj.NextPicSeq,
		"key":      picKey,
		"progress": 0,    // Generate
		"message":  "OK", // Generate
	})
}

func changeClothes(ctx context.Context, objectID primitive.ObjectID, inURL string, outFile string) error {
	// TODO: global client in `db` package
	client, err := repl.NewClient(repl.WithToken(config.ReplicateAPIKey))
	if err != nil {
		return err
	}
	prediction, err := client.CreatePredictionWithDeployment(ctx,
		config.ReplicateDeploymentOwner, config.ReplicateDeploymentName,
		repl.PredictionInput{
			"image":  inURL,
			"prompt": "a person wearing " + getOutfit(),
		},
		nil,   // No webhook for now
		false, // This model does not support streaming
	)
	if err != nil {
		return err
	}
	log.Println("The prediction is", mustMarshalJSONString(prediction))
	go waitForClothes(context.Background(), client, objectID, prediction, outFile)
	return nil
}

func waitForClothes(ctx context.Context, client *repl.Client, objectID primitive.ObjectID, prediction *repl.Prediction, outFile string) {
	predFinish, predError := client.WaitAsync(context.TODO(), prediction)
	for predFinish != nil || predError != nil {
		select {
		case pred, ok := <-predFinish:
			if !ok {
				predFinish = nil
				break
			}
			if pred == nil {
				continue
			}
			switch pred.Status {
			case repl.Starting:
				log.Println("Model still starting for", outFile)
			case repl.Processing:
				progress := pred.Progress()
				if progress == nil {
					log.Println("Progress not yet available for", outFile)
					continue
				}
				log.Println("Progress for", outFile, "is", progress.Percentage)
				dao.UpdatePhotoInitProgress(ctx, objectID, pred.Status.String(), int(100.0*progress.Percentage))
			case repl.Succeeded:
				log.Println("Downloading", pred.Output.([]interface{})[3])
				if err := downloadFile(pred.Output.([]interface{})[3].(string), outFile); err != nil {
					log.Println("ERROR! Unable to download result for", outFile, "from", pred.Output.([]interface{})[3], "with error", err)
					return
				}
				dao.UpdatePhotoInitProgress(ctx, objectID, pred.Status.String(), 100)
			default:
				log.Println("Unexpected prediction status", pred.Status, "for", outFile)
			}
		case err, ok := <-predError:
			if !ok {
				predError = nil
			}
			if err != nil {
				log.Println("ERROR!", err)
			}
		}
		if predFinish == nil && predError == nil {
			break
		}
	}
	log.Println("Prediction complete with", mustMarshalJSONString(prediction))
}

func getOutfit() string {
	return config.Outfits[rand.Intn(len(config.Outfits))]
}

func mustMarshalJSONString(obj any) string {
	bts, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return string(bts)
}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
