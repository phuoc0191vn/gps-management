package mongo

import (
	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const AccountMongoCollection = "account"

type AccountMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewAccountMongoRepository(provider *mongo.MongoProvider) *AccountMongoRepository {
	repo := &AccountMongoRepository{provider, AccountMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"username",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"email",
		},
		Unique: true,
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"userID",
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

func (repo *AccountMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *AccountMongoRepository) All() ([]model.Account, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Account, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByEmail(email string) (*model.Account, error) {
	collection, close := repo.collection()
	defer close()

	var result model.Account
	err := collection.Find(bson.M{"email": email}).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) Save(account model.Account) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(account))
}

func (repo *AccountMongoRepository) UpdateByUsername(username string, account model.Account) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Update(bson.M{"username": username}, account))
}

func (repo *AccountMongoRepository) RemoveByUsername(username string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"username": username}))
}
