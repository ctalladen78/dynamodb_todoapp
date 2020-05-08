
### operations sequence

start dynamodb local
```
java -Djava.library.path=./DynamoDBLocal_lib -jar ~/Downloads/dynamodb_local_latest/DynamoDBLocal.jar -sharedDb -dbPath .
```

start dashboard
```
export DYNAMO_ENDPOINT=http://localhost:8000
dynamodb-admin -p 8001

```
start serverless local
`sls offline -start`

set alias
`ifconfig lo0 alias 172.16.123.1`

### create pem file

```
$ openssl genpkey -algorithm RSA -out ./privkey_apr2020.pem
$ openssl rsa -in privkey_apr2020.pem -pubout > pubkey_apr2020.pub
```

### create table routine 
```
aws dynamodb create-table \
--endpoint-url http://localhost:8000 \
--key-schema AttributeName=BucketName,KeyType=HASH AttributeName=CreatedAt,KeyType=RANGE \
--attribute-definitions AttributeName=BucketName,AttributeType=S \
AttributeName=CreatedAt,AttributeType=S,AttributeName=Todo,AttributeType=S \
--global-secondary-indexes file://lsi.json
--billing-mode PAY_PER_REQUEST \
--table-name Test
```

### put item with cli
test todo item
```
aws dynamodb put-item --endpoint-url http://localhost:8000 \
--table-name Test --item file://item_data.json
```

### http requests

```
GET http://localhost:5000/user?userid=third&todo=action
GET http://localhost:5000/userlist
POST http://localhost:5000/user/edit --form-data {"todo":"newvalue"}
GET http://localhost:5000/query?userid=1BUCKET-username123
```
### table schema
```
// todo object schema
// notes: having secondary index attribute is more flexible in terms of querying
// notes: this is less efficient than querying by key attributes
{
  "AttributeDefinitions": [
    {
      "AttributeName": "object_id", // cannot be queryied, only selected
      "AttributeType": "S"
    },
    {
      "AttributeName": "is_done", // secondary index attribute
      "AttributeType": "S"
    }
  ],
  "TableName": "todotable2",
  "KeySchema": [
    {
      "AttributeName": "object_id",
      "KeyType": "HASH"
    }
  ],
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "is_done",
      "KeySchema": [
        {
          "AttributeName": "is_done",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 77,
      "ItemCount": 3,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/todotable2/index/is_done"
    }
  ]
}

// user
{
  "AttributeDefinitions": [
    {
      "AttributeName": "object_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "current_city",
      "AttributeType": "S"
    }
  ],
  "TableName": "usertable2",
  "KeySchema": [
    {
      "AttributeName": "object_id",
      "KeyType": "HASH"
    }
  ],
  "TableSizeBytes": 1306,
  "ItemCount": 4,
  "TableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/usertable2",
  "BillingModeSummary": {
    "BillingMode": "PROVISIONED",
    "LastUpdateToPayPerRequestDateTime": "1970-01-01T00:00:00.000Z"
  },
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "current_city",
      "KeySchema": [
        {
          "AttributeName": "current_city",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 0,
      "ItemCount": 0,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/usertable2/index/current_city"
    }
  ]
}
// trip

{
  "AttributeDefinitions": [
    {
      "AttributeName": "object_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "created_by",
      "AttributeType": "S"
    },
    {
      "AttributeName": "city_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "category_id",
      "AttributeType": "S"
    }
  ],
  "TableName": "triptable3",
  "KeySchema": [
    {
      "AttributeName": "object_id",
      "KeyType": "HASH"
    }
  ],
  "BillingModeSummary": {
    "BillingMode": "",
    "LastUpdateToPayPerRequestDateTime": "1970-01-01T00:00:00.000Z"
  },
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "city_ref",
      "KeySchema": [
        {
          "AttributeName": "city_id",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 4631,
      "ItemCount": 17,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/triptable3/index/city_ref"
    },
    {
      "IndexName": "creator_ref",
      "KeySchema": [
        {
          "AttributeName": "created_by",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 4631,
      "ItemCount": 17,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/triptable3/index/creator_ref"
    },
    {
      "IndexName": "category_ref",
      "KeySchema": [
        {
          "AttributeName": "category_id",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 0,
      "ItemCount": 0,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/triptable3/index/category_ref"
    }
  ]
}
// Location
{
  "AttributeDefinitions": [
    {
      "AttributeName": "object_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "city_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "trip_id",
      "AttributeType": "S"
    }
  ],
  "TableName": "locationtable2",
  "KeySchema": [
    {
      "AttributeName": "object_id",
      "KeyType": "HASH"
    }
  ],
  "TableStatus": "ACTIVE",
  "CreationDateTime": "2020-03-06T06:04:16.268Z",
  "ProvisionedThroughput": {
    "LastIncreaseDateTime": "1970-01-01T00:00:00.000Z",
    "LastDecreaseDateTime": "1970-01-01T00:00:00.000Z",
    "NumberOfDecreasesToday": 0,
    "ReadCapacityUnits": 3,
    "WriteCapacityUnits": 3
  },
  "TableSizeBytes": 2148,
  "ItemCount": 6,
  "TableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/locationtable2",
  "BillingModeSummary": {
    "BillingMode": "PROVISIONED",
    "LastUpdateToPayPerRequestDateTime": "1970-01-01T00:00:00.000Z"
  },
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "city_ref",
      "KeySchema": [
        {
          "AttributeName": "city_id",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 2148,
      "ItemCount": 6,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/locationtable2/index/city_ref"
    },
    {
      "IndexName": "trip_ref",
      "KeySchema": [
        {
          "AttributeName": "trip_id",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 2148,
      "ItemCount": 6,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/locationtable2/index/trip_ref"
    }
  ]
}
// city
{
  "AttributeDefinitions": [
    {
      "AttributeName": "object_id",
      "AttributeType": "S"
    },
    {
      "AttributeName": "country",
      "AttributeType": "S"
    }
  ],
  "TableName": "citytable2",
  "KeySchema": [
    {
      "AttributeName": "object_id",
      "KeyType": "HASH"
    }
  ],
  "TableStatus": "ACTIVE",
  "CreationDateTime": "2020-03-07T08:25:27.319Z",
  "ProvisionedThroughput": {
    "LastIncreaseDateTime": "1970-01-01T00:00:00.000Z",
    "LastDecreaseDateTime": "1970-01-01T00:00:00.000Z",
    "NumberOfDecreasesToday": 0,
    "ReadCapacityUnits": 3,
    "WriteCapacityUnits": 3
  },
  "TableSizeBytes": 1167,
  "ItemCount": 6,
  "TableArn": "arn:aws:dynamodb:ddblocal:000000000000:table/citytable2",
  "BillingModeSummary": {
    "BillingMode": "PROVISIONED",
    "LastUpdateToPayPerRequestDateTime": "1970-01-01T00:00:00.000Z"
  },
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "country",
      "KeySchema": [
        {
          "AttributeName": "country",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      },
      "IndexStatus": "ACTIVE",
      "ProvisionedThroughput": {
        "ReadCapacityUnits": 3,
        "WriteCapacityUnits": 3
      },
      "IndexSizeBytes": 1167,
      "ItemCount": 6,
      "IndexArn": "arn:aws:dynamodb:ddblocal:000000000000:table/citytable2/index/country"
    }
  ]
}
```