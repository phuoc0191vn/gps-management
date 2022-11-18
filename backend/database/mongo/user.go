package mongo

import (
	"fmt"

	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const UserMongoCollection = "user"

type UserMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewUserMongoRepository(provider *mongo.MongoProvider) *UserMongoRepository {
	repo := &UserMongoRepository{provider, UserMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"identityID",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"phoneNumber",
		},
		Unique: true,
	})

	return repo
}

func (repo *UserMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *UserMongoRepository) All() ([]model.User, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.User, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *UserMongoRepository) FindByID(id string) (*model.User, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	var result model.User
	err := collection.FindId(id).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *UserMongoRepository) FindByIdentifyID(identifyID string) (*model.User, error) {
	collection, close := repo.collection()
	defer close()

	var result model.User
	err := collection.Find(bson.M{"identifyID": identifyID}).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *UserMongoRepository) FindByPhoneNumber(phoneNumber string) (*model.User, error) {
	collection, close := repo.collection()
	defer close()

	var result model.User
	err := collection.Find(bson.M{"phoneNumber": phoneNumber}).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *UserMongoRepository) Save(user model.User) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(user))
}

func (repo *UserMongoRepository) UpdateByID(id string, user model.User) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	err := collection.UpdateId(id, user)
	return repo.provider.NewError(err)
}

func (repo *UserMongoRepository) UpdateByIdentifyID(identifyID string, user model.User) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Update(bson.M{"phoneNumber": identifyID}, user)
	return repo.provider.NewError(err)
}

func (repo *UserMongoRepository) UpdateByPhoneNumber(phoneNumber string, user model.User) error {
	collection, close := repo.collection()
	defer close()

	err := collection.Update(bson.M{"phoneNumber": phoneNumber}, user)
	return repo.provider.NewError(err)
}

func (repo *UserMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(id))
}

func (repo *UserMongoRepository) RemoveByIdentifyID(identifyID string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"identifyID": identifyID}))
}

func (repo *UserMongoRepository) RemoveByPhoneNumber(phoneNumber string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"phoneNumber": phoneNumber}))
}
