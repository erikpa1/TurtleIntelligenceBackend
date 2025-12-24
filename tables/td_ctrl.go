package tables

import (
	"fmt"
	"turtle/core/users"

	"turtle/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetTablesContainerName(namespace string) string {
	return fmt.Sprintf("%s_tables", namespace)
}

func QueryTableData(user *users.User, namespace string, query bson.M) []*TableData {

	findOne := options.FindOptions{}
	findOne.Projection = bson.M{
		"theme":  -1,
		"schema": -1,
	}

	return db.QueryEntities[TableData](
		GetTablesContainerName(namespace),
		user.FillOrgQuery(bson.M{}),
		&findOne,
	)

}
