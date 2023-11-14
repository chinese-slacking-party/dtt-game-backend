package users

import (
	"fmt"
	"testing"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"

	"github.com/stretchr/testify/assert"
)

func init() {
	dbName := fmt.Sprintf("dtt_test_%d", time.Now().UnixMicro())
	mongo.Init(dbName)
}

func TestRegister(t *testing.T) {
	_, err := doRegister(&User{Name: "alice"})
	assert.NoError(t, err, "should register successfully")
	_, err = doRegister(&User{Name: "alice"})
	assert.Error(t, err, "should not register the same user twice")
}
