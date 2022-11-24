package mongo

import (
	"fmt"

	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ReportMongoCollection = "report"

type ReportMongoRepository struct {
	provider       *mongo.MongoProvider
	collectionName string
}

func NewReportMongoRepository(provider *mongo.MongoProvider) *ReportMongoRepository {
	repo := &ReportMongoRepository{provider, ReportMongoCollection}
	collection, close := repo.collection()
	defer close()

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
		},
	})

	return repo
}

func (repo *ReportMongoRepository) collection() (collection *mgo.Collection, close func()) {
	session := repo.provider.MongoClient().GetCopySession()
	close = session.Close

	return session.DB(repo.provider.MongoClient().Database()).C(repo.collectionName), close
}

func (repo *ReportMongoRepository) All() ([]model.Report, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Report, 0)
	err := collection.Find(nil).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *ReportMongoRepository) Save(report model.Report) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(report))
}

func (repo *ReportMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(bson.ObjectIdHex(id)))
}

func (repo *ReportMongoRepository) RemoveByAccountID(accountID string) error {
	collection, close := repo.collection()
	defer close()

	_, err := collection.RemoveAll(bson.M{"accountID": accountID})

	return repo.provider.NewError(err)
}
