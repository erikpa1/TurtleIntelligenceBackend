package tables

import (
	"context"

	"github.com/erikpa1/TurtleIntelligenceBackend/credentials"
	"github.com/erikpa1/TurtleIntelligenceBackend/db"
	"github.com/erikpa1/TurtleIntelligenceBackend/lg"
	"github.com/erikpa1/TurtleIntelligenceBackend/models"
	"github.com/erikpa1/TurtleIntelligenceBackend/tools"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const CT_TABLE_DATA = "tables_data"
const CT_TABLES = "tables"

func QueryTables(user *models.User, query bson.M) []*TurtleTable {
	return db.QueryEntities[TurtleTable](CT_TABLES, user.FillOrgQuery(query))
}

func GetTable(user *models.User, tableUid primitive.ObjectID) *TurtleTable {
	return db.GetByIdAndOrg[TurtleTable](CT_TABLES, tableUid, user.Org)
}

func COUTable(user *models.User, table *TurtleTable) {

	if table.Uid.IsZero() {
		table.Org = user.Org
		table.CreatedBy = user.Uid
		table.UpdatedBy = user.Uid
		table.CreatedAt = tools.GetNow()
		table.UpdatedAt = tools.GetNow()

		result, err := db.InsertEntity(CT_TABLES, table)

		if err == nil {
			table.Uid = result.InsertedID.(primitive.ObjectID)
			db.DB.Mongo.Database(credentials.GetDBName()).CreateCollection(context.Background(), table.GetTableMongoName())
		}

	} else {
		table.UpdatedAt = tools.GetNow()
		table.UpdatedBy = user.Uid
		db.SetByOrgAndId(CT_TABLES, table.Uid, table.Org, table)
	}
}

func DeleteTable(user *models.User, uid primitive.ObjectID) {

	table := GetTable(user, uid)

	deletionWasOk := false

	if table != nil {
		if table.HasDatabaseTable {
			err := db.DB.Mongo.Database(credentials.GetDBName()).Collection(table.GetTableMongoName()).Drop(context.Background())

			if err == nil {
				deletionWasOk = true
			} else {
				lg.LogStackTraceErr(err)
			}
		} else {
			err := db.DeleteEntities(CT_TABLE_DATA, bson.M{
				"org":    user.Org,
				"parent": uid,
			})
			if err == nil {
				deletionWasOk = true
			} else {
				lg.LogStackTraceErr(err)
			}
		}

		if deletionWasOk {

			db.DeleteByIdAndOrg(CT_TABLES, uid, user.Org)
		} else {
			lg.LogE("Not going to delete table because of previous error")
		}
	} else {
		lg.LogE("No table to delete: ", uid.Hex())
	}

}
