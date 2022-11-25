package mongo

import (
	"fmt"

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
			"licensePlateNumber",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
		},
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

func (repo *DeviceMongoRepository) Pagination(page int, limit int) (int, []model.Device, error) {
	collection, close := repo.collection()
	defer close()

	query := collection.Find(nil)

	total, err := query.Count()
	if err != nil {
		return 0, nil, repo.provider.NewError(err)
	}

	result := make([]model.Device, 0)

	query = query.Limit(limit)
	if page > 1 {
		query = query.Skip(limit * (page - 1))
	}

	err = query.All(&result)
	return total, result, repo.provider.NewError(err)
}

func (repo *DeviceMongoRepository) FindByID(id string) (*model.Device, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	result := model.Device{}
	err := collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *DeviceMongoRepository) FindByAccountID(accountID string) ([]model.Device, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Device, 0)
	err := collection.Find(bson.M{"accountID": accountID}).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *DeviceMongoRepository) Save(device model.Device) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(device))
}

func (repo *DeviceMongoRepository) UpdateByID(id string, device model.Device) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.UpdateId(bson.ObjectIdHex(id), device))
}

func (repo *DeviceMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(bson.ObjectIdHex(id)))
}
