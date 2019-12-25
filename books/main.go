package main

import (
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "os"
  "regexp"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/dynamodb"

  "github.com/aws/aws-lambda-go/events"
  "github.com/aws/aws-lambda-go/lambda"
)

var ( 
  isbnRegexp = regexp.MustCompile(`[0-9]{3}\-[0-9]{10}`)
  errorLogger = log.New(os.Stderr, "ERROR", log.Llongfile)
  svc = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1")
  db = &Store{
    Db: svc
  }
)

type book struct {
  ISBN  string `json:"isbn"`
}

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse,
error) {
  switch req.HTTPMethod {
    case "GET":
      return show(req)
    case "POST":
      return create(req)
    default:
      return clientError(http.StatusMethodNotAllowed)
  }
}

func show(req events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,
error){
  isbn := req.QueryStringParameters["isbn"]
  if !isbnRegexp.MatchString(isbn){
    return clientError(http.StatusBadRequest)
  }

  bk, err := db.getItem(isbn)
  if err != nil {return serverError(err)}
  if bk == nil {return clientError(http.StatusNotFound)}
  js, err := json.Marshal(bk)
  if err != nil {return serverError(err)}
  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusOK,
    Body: string(js),
  },nil
}

// request data must be validated
// see validate library instead of helper function
func create(req events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse,error){
  if req.Headers["content-type"] != "application/json" &&
    req.Headers["Content-Type"] != "application/json" {
      return clientError(http.StatusNotAcceptable)
  }
  bk := &book{}
  err := json.Unmarshal([]byte(req.Body), bk)
  if err != nil { 
    return clientError(http.StatusUnprocessableEntity)
  }
  if !isbnRegexp.MatchString(bk.ISBN) {
    return clientError(http.StatusBadRequest)
  }
  err = db.putItem(bk)
  if err != nil {return serverError(err)}
  return events.APIGatewayProxyResponse{
    StatusCode: 201,
    Headers:  map[string]string{"Location":fmt.Sprintf("/books?isbn=%s", bk.ISBN)},
  }
}

// helper function
func serverError(err error) (events.APIGatewayProxyResponse, error) {
  errorLogger.Println(err.Error())
  return events.APIGatewayProxyResponse{
    StatusCode: http.StatusInternalServerError,
    Body:       http.StatusText(http.StatusInternalServerError),
  },nil
}

// helper function
func clientError(err error) (events.APIGatewayProxyResponse, error) {
  return events.APIGatewayProxyResponse{
    StatusCode: status,
    Body:       http.StatusText(status),
  },nil
}

func main(){
  lambda.Start(router)
}
