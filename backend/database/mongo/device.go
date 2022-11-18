package mongo

import (
	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const DeviceMongoCollection = "device"

type DeviceMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewDeviceMongoRepository(provider *mongo.MongoProvider) *DeviceMongoRepository {
	repo := &DeviceMongoRepository{provider, DeviceMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"identityID",
		},
		Unique: true,
	})

	return repo
}

func (repo *DeviceMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *DeviceMongoRepository) All() ([]model.Device, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Device, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *DeviceMongoRepository) Save(device model.Device) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(device))
}

func (repo *DeviceMongoRepository) UpdateByIdentifyID(identifyID string, device model.Device) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Update(bson.M{"identifyID": identifyID}, device))
}

func (repo *DeviceMongoRepository) RemoveByIdentifyID(identifyID string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"identifyID": identifyID}))
}
