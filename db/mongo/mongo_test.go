package mongo

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestInit(t *testing.T) {
	dbName := fmt.Sprintf("anydb_%d", time.Now().UnixMicro())
	assert.NoError(t, Init(dbName), "Unable to initialize DB - do you have MongoDB running on localhost?")
	table := GetDB().Collection("my_table")
	_, err := table.InsertOne(context.TODO(), bson.M{"foo": "bar"})
	assert.NoError(t, err, "Unable to insert into my_table - does your MongoDB installation work?")
	assert.NoError(t, GetDB().Drop(context.TODO()), "Unable to drop database - is your MongoDB sane?")
}
