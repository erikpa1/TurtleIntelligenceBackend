package db

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/erikpa1/TurtleIntelligenceBackend/credentials"
	"github.com/erikpa1/TurtleIntelligenceBackend/interfaces"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnyDBConnection struct {
	Mongo    *mongo.Client
	Context  *context.Context
	IsCosmos bool
}

var DB = NewAnyDBConnection()

func (self *AnyDBConnection) Col(name string) *mongo.Collection {
	return DB.Mongo.Database(credentials.GetDBName()).Collection(name)
}

func (self *AnyDBConnection) Db(name string) *mongo.Database {
	return DB.Mongo.Database(credentials.GetDBName())
}

func (self *AnyDBConnection) QueryEntity(collection string, query bson.M) *any {
	var elem any
	err := self.Col(collection).FindOne(context.TODO(), query).Decode(&elem)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		return nil
	}
	return &elem
}
func (self *AnyDBConnection) QueryEntitiesCursor(collectionName string, query bson.M, projection bson.M) (*mongo.Cursor, error) {
	collection := self.Col(collectionName)
	// Access the collection from the MongoDB client
	if projection == nil {
		projection = bson.M{}
	}

	findOptions := options.Find().SetProjection(projection)
	cursor, err := collection.Find(context.Background(), query, findOptions)
	if err != nil {
		return nil, err
	}
	return cursor, nil
}

func QueryEntity[T any](collection string, query bson.M) *T {
	var elem T
	err := DB.Col(collection).FindOne(context.TODO(), query).Decode(&elem)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			return nil
		} else {
			lg.LogStackTraceErr(err)
		}
		return nil
	}
	return &elem
}

func InsertEntity(collection string, entity any) {
	_, err := DB.Col(collection).InsertOne(context.TODO(), entity)
	if err != nil {
		lg.LogE(err)
	}
}

func InsertMany(collection string, entities []interface{}) {
	_, err := DB.Col(collection).InsertMany(context.TODO(), entities)
	if err != nil {
		lg.LogStackTraceErr(err.Error())
	}
}

func DeleteEntity(collection string, query bson.M) {
	_, err := DB.Col(collection).DeleteOne(context.TODO(), query)
	if err != nil {
		lg.LogE(err)
	}
}

func DeleteEntityWithUid(collection string, uid primitive.ObjectID) {
	_, err := DB.Col(collection).DeleteOne(context.TODO(), bson.M{
		"_id": uid,
	})
	if err != nil {
		lg.LogE(err)
	}
}

func DeleteEntitiesOfParent(collection string, parent primitive.ObjectID) {
	_, err := DB.Col(collection).DeleteMany(context.TODO(), bson.M{
		"parent": parent,
	})
	if err != nil {
		lg.LogE(err)
	}
}

func DeleteEntities(collection string, query bson.M) {
	_, err := DB.Col(collection).DeleteMany(context.TODO(), query)
	if err != nil {
		lg.LogE(err)
	}
}

func UpdateEntity(collection string, entity any) {
	jData, _ := json.Marshal(entity)
	uid := gjson.Get(string(jData), "uid").String()

	_, err := DB.Col(collection).UpdateOne(context.TODO(), bson.M{"uid": uid}, bson.M{"$set": entity})

	if err != nil {
		lg.LogStackTraceErr(err)
	}
}

func SetByOrgAndId(collection string, _id primitive.ObjectID, orgId primitive.ObjectID, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := DB.Col(collection).UpdateOne(context.TODO(), bson.M{
		"_id": _id,
		"org": orgId,
	}, bson.M{
		"$set": update,
	}, opts...)

	if err != nil {
		lg.LogStackTraceErr(err)
		return err
	}
	return nil
}

func SetById(collection string, _id primitive.ObjectID, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := DB.Col(collection).UpdateOne(context.TODO(), bson.M{
		"_id": _id,
	}, bson.M{
		"$set": update,
	}, opts...)

	if err != nil {
		lg.LogStackTraceErr(err)
		return err
	}
	return nil
}

func UpdateOneCustom(collection string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) error {
	_, err := DB.Col(collection).UpdateOne(context.TODO(), filter, update, opts...)

	if err != nil {
		lg.LogStackTraceErr(err)
		return err
	}
	return nil
}

func UpdateEntitiesWhere(collection string, filter bson.M, data any) {
	_, err := DB.Col(collection).UpdateMany(context.TODO(), filter, bson.M{"$set": data})

	if err != nil {
		lg.LogStackTraceErr(err)
	}
}

func HasEntityWithUid(collection string, uid string) bool {
	value, _ := DB.Col(collection).CountDocuments(context.TODO(), bson.M{"uid": uid})
	return value > 0
}

func COUEntity(collection string, entity interfaces.UidProvider) {

	exists := HasEntityWithUid(collection, entity.GetUid())

	if exists {
		UpdateEntity(collection, entity)
	} else {
		InsertEntity(collection, entity)
	}
}

func QueryEntitiesChannel[T any](ctx context.Context, collection string, query bson.M, opts ...*options.FindOptions) <-chan *T {
	resultsChan := make(chan *T)

	go func() {
		defer tools.Recover("Failed to do mongo thing")
		defer close(resultsChan) // Ensure the channel is closed when done

		cursor, err := DB.Col(collection).Find(ctx, query, opts...)

		if err != nil {
			lg.LogStackTraceErr(err.Error())
			return // Exit if there's an error
		}
		defer cursor.Close(context.TODO()) // Ensure cursor is closed after iteration

		for cursor.Next(context.TODO()) {
			var elem T // Changed to T instead of *T
			if err := cursor.Decode(&elem); err != nil {
				lg.LogE(err)
				return // Exit if there's an error
			}
			resultsChan <- &elem // Send the address of elem to the channel
		}

		if err := cursor.Err(); err != nil {
			lg.LogE(err) // Log any cursor errors
		}
	}()

	return resultsChan // Return the channel
}
func QueryEntitiesAsCopy[T any](collection string, query bson.M, opts ...*options.FindOptions) []T {
	cursor, err := DB.Col(collection).Find(context.TODO(), query, opts...)

	if err != nil {
		lg.LogE(err)
	}

	result := []T{}

	for cursor.Next(context.TODO()) {
		var elem T

		err := cursor.Decode(&elem)
		if err != nil {
			lg.LogE(err)
		}
		result = append(result, elem)
	}

	return result
}

func QueryEntities[T any](collection string, query bson.M, opts ...*options.FindOptions) []*T {

	cursor, err := DB.Col(collection).Find(context.TODO(), query, opts...)

	if err != nil {
		lg.LogE(err)
	}

	result := []*T{}

	for cursor.Next(context.TODO()) {
		var elem *T

		err := cursor.Decode(&elem)
		if err != nil {
			lg.LogE(err)
		}
		result = append(result, elem)
	}

	return result
}

func (self *AnyDBConnection) ListContainers() []string {

	names, err := self.Db(credentials.GetDBName()).ListCollectionNames(context.TODO(), bson.M{})

	if err != nil {
		lg.LogE(err.Error())
		return make([]string, 0)
	}

	return names

}

func (self *AnyDBConnection) DeleteEntity(collection string, query bson.M) {
	self.Col(collection).DeleteOne(context.TODO(), query)
}

func NewAnyDBConnection() AnyDBConnection {
	clientOptions := options.Client().ApplyURI(credentials.GetDBConnStr()) // Replace with your MongoDB URI

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		lg.LogE(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		lg.LogE(err)
	}

	lg.LogOk("Connected to MongoDB! ", credentials.GetDBConnStr())
	lg.LogOk("Database name: ", credentials.GetDBName())

	var conn = AnyDBConnection{}
	conn.SetIsCosmos(credentials.GetDBConnStr())
	conn.IsCosmos = false
	conn.Mongo = client

	return conn
}

func EntityExists(collection string, query bson.M) bool {
	var tmp bson.M
	err := DB.Col(collection).FindOne(context.TODO(), query).Decode(&tmp)

	if err != nil {
		return false
	}

	return true

}

func CountEntities(collection string, query bson.M) int64 {
	value, err := DB.Col(collection).CountDocuments(context.TODO(), query)

	if err != nil {
		lg.LogStackTraceErr(err)
		return 0
	}

	return value
}
func (self *AnyDBConnection) SetIsCosmos(connStr string) {
	self.IsCosmos = strings.Contains(connStr, "mongo.cosmos.azure.com")
}

func (self *AnyDBConnection) SafeSort(cursor *mongo.Cursor, sort interface{}) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	// Check if the database is Cosmos DB or MongoDB, and apply logic accordingly
	if self.IsCosmos {
		// For Cosmos DB, assume no need to sort; just return the cursor as-is
		for cursor.Next(context.Background()) {
			var result map[string]interface{}
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
	} else {
		// For MongoDB, sorting is handled at the query level, so we iterate over the sorted cursor
		for cursor.Next(context.Background()) {
			var result map[string]interface{}
			if err := cursor.Decode(&result); err != nil {
				return nil, err
			}
			results = append(results, result)
		}
		// Check for any cursor errors after iteration
		if err := cursor.Err(); err != nil {
			return nil, err
		}
	}

	return results, nil
}
