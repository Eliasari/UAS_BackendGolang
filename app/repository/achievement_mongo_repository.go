package repository

import (
	"context"
	"time"

	"uas-prestasi/app/model"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type AchievementMongoRepository struct {
	Collection *mongo.Collection
}

func NewAchievementMongoRepository(db *mongo.Database) *AchievementMongoRepository {
	return &AchievementMongoRepository{
		Collection: db.Collection("achievements"),
	}
}

func (r *AchievementMongoRepository) InsertDraft(achievement *model.Achievement) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	achievement.CreatedAt = time.Now()
	achievement.UpdatedAt = time.Now()
	achievement.Attachments = []model.Attachment{}

	result, err := r.Collection.InsertOne(ctx, achievement)
	if err != nil {
		return "", err
	}

	objectID := result.InsertedID.(primitive.ObjectID)

	return objectID.Hex(), nil
}

func (r *AchievementMongoRepository) FindByID(id string) (bson.M, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	var result bson.M
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = r.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}


func (r *AchievementMongoRepository) UpdateByID(id string, payload map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": payload},
	)

	return err
}

func (r *AchievementMongoRepository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.Collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (r *AchievementMongoRepository) AddAttachment(mongoID string, attachment model.Attachment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(mongoID)
	if err != nil {
		return err
	}

	_, err = r.Collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$push": bson.M{"attachments": attachment}},
	)

	return err
}

func (r *AchievementMongoRepository) UpdatePoints(id string, points int) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return err
    }

    _, err = r.Collection.UpdateOne(
        ctx,
        bson.M{"_id": objID},
        bson.M{"$set": bson.M{"points": points, "updated_at": time.Now()}},
    )
    return err
}


