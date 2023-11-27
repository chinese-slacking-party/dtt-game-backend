package users

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/chinese-slacking-party/dtt-game-backend/db/mongo"

	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestRegister(t *testing.T) {
	dbName := fmt.Sprintf("dtt_test_register_%d", time.Now().UnixMicro())
	mongo.Init(dbName)
	_, err := doRegister(&UserRegisterReq{Name: "alice"})
	assert.NoError(t, err, "should register successfully")
	_, err = doRegister(&UserRegisterReq{Name: "alice"})
	assert.Error(t, err, "should not register the same user twice")
	assert.NoError(t, mongo.GetDB().Drop(context.TODO()), "Unable to drop database - is your MongoDB sane?")
}
