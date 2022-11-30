package mongo

import (
	"fmt"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/model"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const ActivityLogMongoCollection = "activityLog"

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
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"deviceID",
		},
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
			"deviceID",
		},
	})

	collection.EnsureIndex(mgo.Index{
		Key: []string{
			"accountID",
			"deviceID",
			"date",
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

func (repo *ActivityLogMongoRepository) GetInDay(deviceID, accountID string, date time.Time) (*model.ActivityLog, error) {
	collection, close := repo.collection()
	defer close()

	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	result := model.ActivityLog{}
	err := collection.Find(bson.M{
		"accountID": accountID,
		"deviceID":  deviceID,
		"date":      date,
	}).One(&result)
	return &result, repo.provider.NewError(err)
}

func (repo *ActivityLogMongoRepository) GetRange(deviceID, accountID string, start, end time.Time) ([]model.ActivityLog, error) {
	collection, close := repo.collection()
	defer close()

	start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	end = time.Date(end.Year(), end.Month(), end.Day(), 0, 0, 0, 0, end.Location())

	result := make([]model.ActivityLog, 0)
	err := collection.Find(bson.M{
		"accountID": accountID,
		"deviceID":  deviceID,
		"date":      "date",
	}).All(&result)
	return result, repo.provider.NewError(err)
}

func (repo *ActivityLogMongoRepository) FilterReport(filter model.Filter) (chan model.ActivityLogData, chan error) {
	collection, closeCollection := repo.collection()

	filter.StartTime = time.Date(filter.StartTime.Year(), filter.StartTime.Month(), filter.StartTime.Day(), 0, 0, 0, 0, filter.StartTime.Location())
	filter.EndTime = time.Date(filter.EndTime.Year(), filter.EndTime.Month(), filter.EndTime.Day(), 0, 0, 0, 0, filter.EndTime.Location())

	match := bson.M{
		"accountID": filter.AccountID,
		"deviceID":  filter.DeviceID,
		"date":      bson.M{"$gte": filter.StartTime, "$lte": filter.EndTime},
	}

	pipeline := []bson.M{
		{"$match": match},
		{"$unwind": "$locations"},
		{"$project": bson.M{
			"date":      "$date",
			"longitude": "$locations.longitude",
			"latitude":  "$locations.latitude",
			"address":   "$locations.address",
			"timestamp": "$locations.timestamp",
			"speed":     "$locations.speed",
		}},
		{"$sort": bson.M{"date": 1}},
	}

	dataChan := make(chan model.ActivityLogData)
	errChan := make(chan error)

	go func() {
		defer close(dataChan)
		defer close(errChan)
		defer closeCollection()

		var data model.ActivityLogData
		iter := collection.Pipe(pipeline).Batch(1000).AllowDiskUse().Iter()
		for iter.Next(&data) {
			dataChan <- data
		}

		errChan <- iter.Err()
	}()

	return dataChan, errChan
}

func (repo *ActivityLogMongoRepository) Save(log model.ActivityLog) error {
	collection, close := repo.collection()
	defer close()

	log.Date = time.Date(log.Date.Year(), log.Date.Month(), log.Date.Day(), 0, 0, 0, 0, log.Date.Location())

	var oldLog model.ActivityLog
	err := collection.Find(bson.M{
		"accountID": log.AccountID,
		"deviceID":  log.DeviceID,
		"date":      log.Date,
	}).One(&oldLog)

	if err == nil {
		oldLog.Locations = append(oldLog.Locations, log.Locations...)
		return repo.provider.NewError(collection.UpdateId(oldLog.ID.Hex(), oldLog))
	}

	return repo.provider.NewError(collection.Insert(log))
}

func (repo *ActivityLogMongoRepository) UpdateByID(id string, log model.ActivityLog) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}

	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.UpdateId(bson.ObjectIdHex(id), log))
}

func (repo *ActivityLogMongoRepository) RemoveByID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return fmt.Errorf("invalid id: %s", id)
	}
	collection, close := repo.collection()
	defer close()

	return repo.provider.NewError(collection.RemoveId(bson.ObjectIdHex(id)))
}

func (repo *ActivityLogMongoRepository) RemoveByAccountID(accountID string) error {
	collection, close := repo.collection()
	defer close()

	_, err := collection.RemoveAll(bson.M{"accountID": accountID})

	return repo.provider.NewError(err)
}
