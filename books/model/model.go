package model

// https://github.com/EwanValentine/invoicely/blob/master/functions/sprints/model/repository.go


type TodoObject struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Todo      string `json:"todo"`
}


// https://yourbasic.org/golang/iota/
// enum for querying item attributes
type QueryCondition int

const (
	CREATED_AT = iota
	CREATED_BY
)

func (q QueryCondition) String() string {
	return [...]string{"CREATED_BY", "CREATED_AT"}[q]
}

type SessionObject struct {
	ObjectId string `json:"object_id"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
	UserObjectId string `json:"user_object_id"` 
	SessionToken string `json:"session_token"` 
}

type Pointer struct {
	Type string `json:"type"` 
	ClassName string `json:"classname"`
	ObjectRefId string `json:"objectref_id"`
}

type City struct {
	ObjectId string `json:"object_id"`
	CreatedAt string `json:"created_at"` 
}

//"currentLocation": {
	// "type": "Pointer",
	// "className": "City",
	// "objectId": "9FVeDdsiZR"
// }

type UserObject struct {
	ObjectId string `json:"object_id"`
	CreatedAt string `json:"created_at"` 
	UpdatedAt string `json:"updated_at"`
	Username string `json:"updated_at"` 
	Email string `json:"email"` 
	Bio string `json:"bio"` 
	Whatsapp string `json:"whatsapp"` 
	Interests []*Pointer
	FirstName string `json:"first_name"`
	CurrentLocation *Pointer
}
