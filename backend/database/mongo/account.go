package mongo

import (
	"fmt"

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

func (repo *AccountMongoRepository) Pagination(page int, limit int, condition map[string]interface{}) (int, []model.Account, error) {
	collection, close := repo.collection()
	defer close()

	query := collection.Find(condition)

	total, err := query.Count()
	if err != nil {
		return 0, nil, repo.provider.NewError(err)
	}

	result := make([]model.Account, 0)

	query = query.Limit(limit)
	if page > 1 {
		query = query.Skip(limit * (page - 1))
	}

	err = query.All(&result)
	return total, result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByID(id string) (*model.Account, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	var result model.Account
	err := collection.FindId(bson.ObjectIdHex(id)).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) FindByEmail(email string) (*model.Account, error) {
	collection, close := repo.collection()
	defer close()

	var result model.Account
	err := collection.Find(bson.M{"email": email}).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) GetDeviceIDsByID(id string) ([]string, error) {
	if !bson.IsObjectIdHex(id) {
		return nil, fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	var result model.Account
	err := collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return nil, repo.provider.NewError(err)
	}

	return result.DeviceIDs, nil
}

func (repo *AccountMongoRepository) GetChildAccounts(email string) ([]model.Account, error) {
	collection, close := repo.collection()
	defer close()

	var result []model.Account
	err := collection.Find(bson.M{"createdBy": email}).Select(bson.M{"_id": 1, "email": 1}).All(&result)
	if err != nil {
		return nil, repo.provider.NewError(err)
	}

	return result, nil
}

func (repo *AccountMongoRepository) Save(account model.Account) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Insert(account))
}

func (repo *AccountMongoRepository) UpdateByEmail(email string, account model.Account) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Update(bson.M{"email": email}, account))
}

func (repo *AccountMongoRepository) UpdateDeviceIDs(id string, deviceIDs []string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	err := collection.UpdateId(bson.ObjectIdHex(id), bson.M{
		"$set": bson.M{
			"deviceIDs": deviceIDs,
		},
	})

	return repo.provider.NewError(err)
}

func (repo *AccountMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id")
	}
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(bson.ObjectIdHex(id)))
}

func (repo *AccountMongoRepository) RemoveByEmail(email string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"email": email}))
}

func (repo *AccountMongoRepository) RemoveByUserID(userID string) error {
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.Remove(bson.M{"userID": userID}))
}
