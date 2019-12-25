package main

import (
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/aws/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/service/dynamodb/dynamodbattribute"

)

// TODO check session variables with localstack
// var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1"))

type Store struct {
  Db  *dynamodb.DynamoDB
  }

// setup dynamodb with ISBN(S) field
// https://stackoverflow.com/questions/49129534/unmarshal-mapstringdynamodbattributevalue-into-a-struct
func (s *store) GetItem(s string) (*book, error) {
  input := &dynamodb.GetItemInput{
    TableName: aws.String("Books"),
    Key: map[string]*dynamodb.AttributeValue{
      "ISBN":{S: aws.String(s)},
    },
  }
  result, err := s.Db.GetItem(input)
  if err != nil{return nil,err}
  if result.Item == nil {return nil,nil}

  // since result.Item is of type map[string]*dynamodb.AttributeValue
  // we use the dynamodbattribute helper function to unmarshall into struct
  bk := &book{}
  err = dynamodbattribute.UnmarshalMap(result.Itme, bk)
  if err != nil {return nil, err}
  return bk, nil
}

// see book model
func (s *store) PutItem(b *book) error {
  input := &dynamodb.PutItemInput{
    TableName: aws.String("Books"),
    Item: map[string]*dynamodb.AttributeValue{
      "ISBN":{S: aws.String(b.ISBN)},
    },
  }
  _,err := s.Db.PutItem(input)
  return err
}


