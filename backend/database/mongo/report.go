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

func (repo *ReportMongoRepository) FindByID(id string) (*model.Report, error) {
	if !bson.IsObjectIdHex(id) {
		fmt.Errorf("invalid id")
	}

	collection, close := repo.collection()
	defer close()

	result := &model.Report{}
	err := collection.FindId(bson.ObjectIdHex(id)).One(result)
	return result, repo.provider.NewError(err)
}

func (repo *ReportMongoRepository) Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Report, error) {
	collection, close := repo.collection()
	defer close()

	query := collection.Find(condition)

	total, err := query.Count()
	if err != nil {
		return 0, nil, repo.provider.NewError(err)
	}

	result := make([]model.Report, 0)

	query = query.Limit(limit)
	if page > 1 {
		query = query.Skip(limit * (page - 1))
	}

	err = query.All(&result)
	return total, result, repo.provider.NewError(err)
}

func (repo *ReportMongoRepository) GetReportByStatus(status int) ([]model.Report, error) {
	collection, close := repo.collection()
	defer close()

	result := make([]model.Report, 0)
	err := collection.Find(bson.M{"status": status}).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *ReportMongoRepository) Save(report model.Report) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(report))
}

func (repo *ReportMongoRepository) UpdateStatusReport(id string, status int) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.UpdateId(bson.ObjectIdHex(id), bson.M{"$set": bson.M{"status": status}}))
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
