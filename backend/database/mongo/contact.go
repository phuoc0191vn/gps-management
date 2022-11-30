package mongo

import (
	"fmt"

	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ContactMongoCollection = "contact"

type ContactMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewContactMongoRepository(provider *mongo.MongoProvider) *ContactMongoRepository {
	repo := &ContactMongoRepository{provider, ContactMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"email",
		},
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
		},
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"adminEmail",
		},
	})

	return repo
}

func (repo *ContactMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *ContactMongoRepository) All() ([]model.Contact, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Contact, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *ContactMongoRepository) Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Contact, error) {
	collection, close := repo.collection()
	defer close()

	query := collection.Find(condition)

	total, err := query.Count()
	if err != nil {
		return 0, nil, repo.provider.NewError(err)
	}

	result := make([]model.Contact, 0)

	query = query.Limit(limit)
	if page > 1 {
		query = query.Skip(limit * (page - 1))
	}

	err = query.All(&result)
	return total, result, repo.provider.NewError(err)
}

func (repo *ContactMongoRepository) FindByID(id string) (*model.Contact, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	var result model.Contact
	err := collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *ContactMongoRepository) Save(account model.Contact) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(account))
}

func (repo *ContactMongoRepository) UpdateStatus(id string, status int) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"status": status}}))
}

func (repo *ContactMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(bson.ObjectIdHex(id)))
}
