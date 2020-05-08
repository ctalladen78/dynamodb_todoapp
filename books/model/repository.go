package model 

// storage logic for this model (ie getters, setters)
import (
	// dao client for dynamodb
)

// GET /user?userid=one
func GetUser(c *gin.Context) {
	todo := &TodoObject{
		Id:   c.Query("userid"),
		Todo: c.Query("todo"),
	}
	res, err := ctrl.GetItem(todo, "Test2")
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": res})
}

func QueryUser(c *gin.Context) {
	q := c.Query("userid")
	// res, err := ctrl.GetItem(todo, "Test2")
	res, err := ctrl.QueryFilter("Test2", CREATED_BY, q)
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": res})
}
func GetAllUsers(c *gin.Context) {
	// http://github.com/gin-gonic/examples
	table := "Test2"
	resList, err := ctrl.Scan(table)
	if err != nil {
		c.AbortWithError(501, err)
	}
	fmt.Printf("RESULTS %s", resList)
	// c.String(http.StatusOK, string(result))
	// c.HTML(http.StatusOK, "template.tmpl", gin.H{"title": "helloworld"})
	// c.Stream()
	c.JSONP(http.StatusOK, gin.H{"data": resList})
}

func PutUser(c *gin.Context) {
	// c.GetPostForm()
	to := c.PostForm("todo")
	u := c.PostForm("userid")
	b := c.PostForm("bucket")
	t := &TodoObject{}
	t.CreatedAt = time.Now().Format(time.RFC3339) // uuid.New()
	t.Id = string([]byte(b + "-" + u))
	t.Todo = to
	o, err := ctrl.PutItem("Test2", t)
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": o})
}

type FormInput struct {
	Nt string
	Ot string
	Id string
}

func UpdateUser(c *gin.Context) {
	tt := &FormInput{
		Nt: c.PostForm("newtodo"),
		Ot: c.PostForm("oldtodo"),
		Id: c.PostForm("userid"),
	}
	u, err := ctrl.Update("Test2", tt)
	if err != nil {
		c.AbortWithError(501, err)
	}
	c.JSONP(200, gin.H{"data": u})
}

**/