package main

import (
	"errors"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DbController struct {
	conn *dynamodb.DynamoDB
}

// uses localhost only
func InitDbConnection(h string) *DbController {
	return &DbController{
		conn: dynamodb.New(session.New(&aws.Config{
			Region:   aws.String("us-east-1"),
			Endpoint: aws.String(h),
		})),
	}
}

// get item by key attributes as per table schema
func (ctrl *DbController) GetItem(t *TodoObject, table string) (interface{}, error) {
	// https://github.com/ace-teknologi/memzy
	// https://github.com/nullseed/lesshomeless-backend/blob/master/services/user/dynamodb/dynamodb.go
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	// building pkey for search query
	var pkey = map[string]*dynamodb.AttributeValue{
		"id": {
			S: aws.String(t.Id),
		},
		"todo": {
			S: aws.String(t.Todo),
		},
	}

	// TodoObject and table key attributes do not match because of extra "cratedat" field
	// pkey, err := dynamodbattribute.MarshalMap(t)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(table),
		Key:       pkey,
	}
	res, err := ctrl.conn.GetItem(input)
	log.Println("GET ITEM output", res)
	if err != nil {
		return nil, err
	}
	var out *TodoObject
	err = dynamodbattribute.UnmarshalMap(res.Item, &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ensure item follows attribute value schema
func (ctrl *DbController) PutItem(table string, todo interface{}) (interface{}, error) {
	// https://stackoverflow.com/questions/38151687/dynamodb-adding-non-key-attributes/56177142
	newTodoAV, err := dynamodbattribute.MarshalMap(todo) // conver todo item to av map
	log.Printf("AV Map %v", newTodoAV)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.PutItemInput{
		Item:      newTodoAV,
		TableName: aws.String(table),
	}
	log.Printf("Put Input %v", input)
	o, err := ctrl.conn.PutItem(input)
	if err != nil {
		return nil, err
	}
	// var out map[string]interface{}
	// log.Printf("Put output %s", o.Attributes)
	// dynamodbattribute.UnmarshalMap(o.Attributes, &out)
	return o.Attributes, err
}

// pass in an empty attribute value struct which will be populated as a result
func (ctrl *DbController) Scan(table string) (interface{}, error) {
	if ctrl.conn == nil {
		return nil, errors.New("db connection error")
	}
	// get all items in table
	scanOutput, err := ctrl.conn.Scan(&dynamodb.ScanInput{
		TableName: aws.String(table),
	})
	if err != nil {
		return nil, err
	}
	var castTo []*TodoObject
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go
	err = dynamodbattribute.UnmarshalListOfMaps(scanOutput.Items, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil
}

// using dynamodb.Query as opposed to dynamodb.Scan
// query by enums CREATED_AT | CREATED_BY
//     --key-condition-expression 'Artist = :a AND SongTitle BETWEEN :t1 AND :t2' \
func (ctrl *DbController) QueryFilter(table string, qc QueryCondition, val string) (interface{}, error) {
	condition := ""
	switch qc {
	case CREATED_AT:
		condition = "CREATED_AT = :val"
	case CREATED_BY: // return items created by
		condition = "id = :val"
	default: // return all items
		condition = ""
	}
	qInput := &dynamodb.QueryInput{
		TableName: aws.String(table),
		// AttributesToGet : , // select only certain attribute values
		KeyConditionExpression: aws.String(condition),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":val": {S: aws.String(val)}, // use this value in an expression
		},
		// KeyConditionExpression: "",
	}
	res, err := ctrl.conn.Query(qInput)
	castTo := []*TodoObject{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &castTo)
	if err != nil {
		return nil, err
	}
	return castTo, nil

}
// https://stackoverflow.com/questions/43727426/can-anyone-provide-an-example-of-a-dynamodb-document-client-upsert
// https://stackoverflow.com/questions/33847477/querying-a-global-secondary-index-in-dynamodb-local
func (ctrl *DbController) QueryTodoByDoneIndex(isDone string) ([]interface{}, error) {

	qInput := &dynamodb.QueryInput{
		TableName:              aws.String("todotable2"),
		IndexName:              aws.String("is_done"),
		KeyConditionExpression: aws.String("is_done = :isdone"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isdone": {S: aws.String(isDone)},
		},
		ProjectionExpression: aws.String(""), // return values of only selected attributes
	}
	fmt.Println("QUERY TODO INPUT", qInput)
	res, err := ctrl.conn.Query(qInput)
	// castTo := []*TodoObject{}
	var retlist []interface{}
	err = dynamodbattribute.UnmarshalListOfMaps(res.Items, &retlist)
	if err != nil {
		return nil, err
	}
	return retlist, nil

}
// Update(table string, f *FormInput) (interface{}, err)
// returns updated object
func (ctrl *DbController) Update(table string, formInput *FormInput) (interface{}, error) {
	var err error
	// var keyMapAV2 map[string]*dynamodb.AttributeValue
	// var toUpdate map[string]*dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// keyMapAV, err := dynamodbattribute.MarshalMap(oldItem)
	oldItemKeys := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String(formInput.Id)},
		"todo": {S: aws.String(formInput.Ot)},
	}
	if err != nil {
		return nil, errors.New("itemkey error")
	}
	if err != nil {
		return nil, errors.New("newItem error")
	}
	// https://aws.amazon.com/blogs/developer/introducing-amazon-dynamodb-expression-builder-in-the-aws-sdk-for-go/
	// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go#L33
	itemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key:       oldItemKeys, // match key attributes per table definition
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(formInput.Nt)}, // set new value
		},
		// attribute being updated must not be part of the key
		ExpressionAttributeNames: map[string]*string{
			"#T": aws.String("createdat"),
		},
		// ConditionExpression:	"attribute_exists("#T"),
		// https://gist.github.com/doncicuto/d623ec0e74bf6ea0db7c364d88507393#file-main-go-L63
		ReturnValues:     aws.String("ALL_NEW"),     // enum of ReturnValue class UPDATED_NEW ALL_NEW ALL_OLD
		UpdateExpression: aws.String("set #T = :t"), // SET,REMOVE the attribute to update

	}
	result, err := ctrl.conn.UpdateItem(itemInput)
	if err != nil {
		return nil, err
	}

	// TODO print resulting updated attributes
	// var u *TodoObject
	out := &TodoObject{}

	// convert db result into inmemory struct
	err = dynamodbattribute.UnmarshalMap(result.Attributes, out)
	if err != nil {
		return nil, err
	}

	log.Printf("UPDATE RESULT %s", result.Attributes)
	return out, nil
}
// updates scalar non-primary attribute using primary key
func (ctrl *DbController) UpdateTodo(todoId string, attributeName string, newValue string) (interface{}, error) {

	primaryKey := map[string]*dynamodb.AttributeValue{
		"object_id": {S: aws.String(todoId)}, // hash key
	}
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String("todotable2"),
		Key:       primaryKey, // match key attributes per table definition
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":AV": {S: aws.String(newValue)}, // set new value
		},
		// attribute being updated must not be part of the key
		ExpressionAttributeNames: map[string]*string{
			"#AN": aws.String(attributeName),
		},
		// ConditionExpression:	"attribute_exists("#T"),
		// https://gist.github.com/doncicuto/d623ec0e74bf6ea0db7c364d88507393#file-main-go-L63
		ReturnValues:     aws.String("ALL_NEW"),       // enum of ReturnValue class UPDATED_NEW ALL_NEW ALL_OLD
		UpdateExpression: aws.String("set #AN = :AV"), // SET,REMOVE the attribute to update

	}
	result, err := ctrl.conn.UpdateItem(updateInput)
	if err != nil {
		return nil, err
	}
	var out interface{}
	err = dynamodbattribute.UnmarshalMap(result.Attributes, &out)
	if err != nil {
		return nil, err
	}
	// TODO print resulting updated attributes
	log.Printf("UPDATE RESULT %s", result.Attributes)
	return out, nil
}

// update has reference to old data
// returns new updated object
func (ctrl *DbController) Update(table string, formInput *FormInput) (interface{}, error) {
	var err error
	// var keyMapAV2 map[string]*dynamodb.AttributeValue
	// var toUpdate map[string]*dynamodb.AttributeValue
	// http://gist.github.com/doncicuto
	// keyMapAV, err := dynamodbattribute.MarshalMap(oldItem)
	todoPrimaryKeys := map[string]*dynamodb.AttributeValue{
		"id":   {S: aws.String(formInput.Id)},
		"todo": {S: aws.String(formInput.Ot)},
	}
	if err != nil {
		return nil, errors.New("itemkey error")
	}
	if err != nil {
		return nil, errors.New("newItem error")
	}
	// https://aws.amazon.com/blogs/developer/introducing-amazon-dynamodb-expression-builder-in-the-aws-sdk-for-go/
	// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.UpdateExpressions.html
	// https://github.com/mczal/go-gellato-membership/blob/master/service/UserService.go#L33
	itemInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(table),
		Key:       todoPrimaryKeys, // match key attributes per table definition
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":t": {S: aws.String(formInput.Nt)}, // set new value
		},
		// attribute being updated must not be part of the key
		ExpressionAttributeNames: map[string]*string{
			"#T": aws.String("createdat"),
		},
		// ConditionExpression:	"attribute_exists("#T"),
		// https://gist.github.com/doncicuto/d623ec0e74bf6ea0db7c364d88507393#file-main-go-L63
		ReturnValues:     aws.String("ALL_NEW"),     // enum of ReturnValue class UPDATED_NEW ALL_NEW ALL_OLD
		UpdateExpression: aws.String("set #T = :t"), // SET,REMOVE the attribute to update

	}
	result, err := ctrl.conn.UpdateItem(itemInput)
	if err != nil {
		return nil, err
	}

	// TODO print resulting updated attributes
	// var u *TodoObject
	out := &TodoObject{}

	// convert db result into inmemory struct
	err = dynamodbattribute.UnmarshalMap(result.Attributes, out)
	if err != nil {
		return nil, err
	}

	log.Printf("UPDATE RESULT %s", result.Attributes)
	return out, nil
}