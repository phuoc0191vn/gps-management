package mongo

import (
	"fmt"

	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ActivityLogMongoCollection = "device"

type ActivityLogMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewActivityLogMongoRepository(provider *mongo.MongoProvider) *ActivityLogMongoRepository {
	repo := &ActivityLogMongoRepository{provider, ActivityLogMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"deviceID",
		},
		Unique: true,
	})

	return repo
}

func (repo *ActivityLogMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *ActivityLogMongoRepository) All() ([]model.ActivityLog, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.ActivityLog, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *ActivityLogMongoRepository) Save(log model.ActivityLog) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(log))
}

func (repo *ActivityLogMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(id))
}
