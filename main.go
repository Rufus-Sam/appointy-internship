package main

import (
	"time"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
    "log"
    "os"
	"strconv"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
)
var client *mongo.Client
var pageSize int64 = 2

//User Struct (Model)
type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` 
	Name     string `json:"name,omitempty" bson:"name,omitempty"` 
	Email    string `json:"email,omitempty" bson:"email,omitempty"` 
	Password string `json:"password,omitempty" bson:"password,omitempty"` 
}
// When creating User in postman
// INPUT:
// {
// 	"name":"Rufus",
// 	"email":"sam2001rufus@gmail.com",
// 	"Password":"Welcome123"
// }
// OUTPUT:
// {
// 	"InsertedID": "6161d4b870181d4944c355e4"
// }
	
//Post Struct (Model)
type Post struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL string `json:"image,omitempty" bson:"image,omitempty"`
	PostedAt primitive.Timestamp `json:"postedAt,omitempty" bson:"postedAt,omitempty"` 
	User     primitive.ObjectID  `json:"user,omitempty" bson:"user,omitempty"`
}
// When creating post in postman
// INPUT:
// {
// 	"Caption":"Technology",
// 	"ImageURL":"http://some-domain/image1.jpg",
// 	"User":"6161d4b870181d4944c355e4"       => this is the user id generated(existing user)
// }
// OUTPUT:
// {
//     "InsertedID": "6161d55070181d4944c355e5"
// }

// bcrypt hashing for password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}
// function for password checking -> not used as no login route
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

//GET all users
func getUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	// needs a slice to store multiple values
	var users []User

	//refernce the collection
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// get the params to check if pagination route called
	params := mux.Vars(req)
	if params["page"]!="" {
		//pagination route called
		page,_:= strconv.Atoi(params["page"])
		p:= int64(page)
		options := options.Find()
		options.SetLimit(pageSize)
		options.SetSkip(pageSize * (p - 1))
		cursor,err:=collection.Find(ctx,bson.M{},options)

		//error handling
		if err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		defer cursor.Close(ctx)

		//add the data in slice
		for cursor.Next(ctx){
			var user User
			cursor.Decode(&user)
			users = append(users,user)
		}

		//error handling
		if err:= cursor.Err() ; err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		//send json data
		json.NewEncoder(res).Encode(users)
	}else{
		//no pagination
		cursor,err:=collection.Find(ctx,bson.M{})

		//error handling
		if err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		defer cursor.Close(ctx)

		//add the data in slice
		for cursor.Next(ctx){
			var user User
			cursor.Decode(&user)
			users = append(users,user)
		}

		//error handling
		if err:= cursor.Err() ; err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		//send json data
		json.NewEncoder(res).Encode(users)
	}
	
}

//POST create new user
func createUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	//declaring variable
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)

	//refernce the collection
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//hash the password using function declared above
	user.Password, _ = HashPassword(user.Password)
	result,_:=collection.InsertOne(ctx,user)

	//send id generated => successful creation
	json.NewEncoder(res).Encode(result)

}

//GET user by id
func getUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	// get the params to get id of user
	params := mux.Vars(req) 
	id,_ := primitive.ObjectIDFromHex(params["id"]) 

	//declaring variable
	var user User

	//refernce the collection
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//get the required data by matching id
	err:=collection.FindOne(ctx,User{Id:id}).Decode(&user)

	//error handling
	if err!= nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	//send json data
	json.NewEncoder(res).Encode(user)
}

//GET all posts
func getPosts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	// needs a slice to store multiple values
	var posts []Post

	//refernce the collection
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	// get the params to check whether pagination route called
	params := mux.Vars(req)
	if params["page"]!="" {
		//pagination route called 
		page,_:= strconv.Atoi(params["page"])
		p:= int64(page)
		options := options.Find()
		options.SetLimit(pageSize)
		options.SetSkip(pageSize * (p - 1))
		cursor,err:=collection.Find(ctx,bson.M{},options)

		//error handling
		if err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		defer cursor.Close(ctx)

		//add the data in slice
		for cursor.Next(ctx){
			var post Post
			cursor.Decode(&post)
			posts = append(posts,post)
		}

		//error handling
		if err:=cursor.Err() ; err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		//send json data
		json.NewEncoder(res).Encode(posts)
	}else{
		//no pagination 
		cursor,err:=collection.Find(ctx,bson.M{})

		//error handling
		if err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		defer cursor.Close(ctx)

		//add the data in slice
		for cursor.Next(ctx){
			var post Post
			cursor.Decode(&post)
			posts = append(posts,post)
		}

		//error handling
		if err:=cursor.Err() ; err!=nil{
			res.WriteHeader(http.StatusInternalServerError)
			res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
			return
		}
		//send json data
		json.NewEncoder(res).Encode(posts)
	}
}

//POST create new post
func createPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	//declaring variable
	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)

	//refernce the collection
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//add the timestamp and create post
	post.PostedAt.T = uint32(time.Now().Unix())
	result,_:=collection.InsertOne(ctx,post)

	//send the id generated => successful creation
	json.NewEncoder(res).Encode(result)
}

//GET post by id
func getPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	// get the params to get id of user
	params := mux.Vars(req) 
	id,_ := primitive.ObjectIDFromHex(params["id"]) 

	//declaring variable
	var post Post

	//refernce the collection
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//get the required data by matching id
	err:=collection.FindOne(ctx,Post{Id:id}).Decode(&post)	

	//error handling
	if err!= nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	//send json data
	json.NewEncoder(res).Encode(post)
}

//GET all posts of user
func getPostsOfUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")

	// get the params to get id of user
	params := mux.Vars(req) 
	id,_ := primitive.ObjectIDFromHex(params["id"]) 

	// needs a slice to store multiple values
	var posts []Post

	//refernce the collection
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	//get the required data by matching id
	cursor,err:=collection.Find(ctx,Post{User:id})

	//error handling
	if err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	defer cursor.Close(ctx)

	//add the data in slice
	for cursor.Next(ctx){
		var post Post
		cursor.Decode(&post)
		posts = append(posts,post)
	}

	//error handling
	if err:=cursor.Err() ; err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}

	//send json data
	json.NewEncoder(res).Encode(posts)
}

func main() {

	//dotenv for mongoatlas uri
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mongoUri:= os.Getenv("MONGO_URI")

	//Mongodb connect	
	//had to use <ctx> again and again in functions in order prevent server panic 
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) 
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, _ = mongo.Connect(ctx, clientOptions)
	
	// Init router
	router := mux.NewRouter()

	//route handlers

	//get all users with or without pagination 
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users/{page}", getUsers).Methods("GET")

	//Create new user
	router.HandleFunc("/users", createUser).Methods("POST")

	//get user by id
	router.HandleFunc("/users/{id}", getUser).Methods("GET")

	//get all posts with or without pagination 
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts/{page}", getPosts).Methods("GET")

	//create new post along with existing user id
	router.HandleFunc("/posts", createPost).Methods("POST")

	//get post by its id
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")

	//get posts of a particular user by their id
	router.HandleFunc("/posts/users/{id}", getPostsOfUser).Methods("GET")

	//Listens on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))

}
