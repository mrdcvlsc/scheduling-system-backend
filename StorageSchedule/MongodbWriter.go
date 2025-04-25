package StorageSchedule

import (
	"context"
	"log"

	"github.com/mrdcvlsc/scheduling-system-backend/Schedule"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongodbWriter struct {
	Mongo *MongoDB
}

type MongodbSchedule struct {
	Semester int    `bson:"Semester"`
	Schedule []byte `bson:"Schedule"`
}

func (s *MongodbWriter) SaveSchedules(university_schedule Schedule.UniTimeTables, semester int) error {

	if s.Mongo.Schedules == nil {
		s.Mongo.Schedules = s.Mongo.Client.Database("gass").Collection("schedules")
	}

	schedule_collection := s.Mongo.Schedules

	opt := options.Replace().SetUpsert(true)

	result, err := schedule_collection.ReplaceOne(
		context.TODO(),

		bson.D{{
			Key:   "Semester",
			Value: semester,
		}},

		&MongodbSchedule{
			Semester: semester,
			Schedule: Schedule.SerializeUniversitySchedule(
				university_schedule,
			),
		},

		opt,
	)

	log.Println("SaveSchedules ReplaceOne result:", result)

	return err
}
