package main

import (
	"time"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
    "log"
    "os"
	"net/http"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"github.com/gorilla/mux"
)
var client *mongo.Client

//User Struct (Model)
type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"name,omitempty" bson:"name,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}
//Post Struct (Model)
type Post struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Caption  string `json:"caption,omitempty" bson:"caption,omitempty"`
	ImageURL string `json:"image,omitempty" bson:"image,omitempty"`
	PostedAt primitive.Timestamp `json:"postedAt,omitempty" bson:"postedAt,omitempty"` 
	User     primitive.ObjectID  `json:"user,omitempty" bson:"user,omitempty"`
}
// bcrypt hashing for password
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    return string(bytes), err
}
// function for password checking 
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

//GET all users
func getUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var users []User
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor,err:=collection.Find(ctx,bson.M{})
	if err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var user User
		cursor.Decode(&user)
		users = append(users,user)
	}
	if err:= cursor.Err() ; err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	json.NewEncoder(res).Encode(users)
}

//POST create new user
func createUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	user.Password, _ = HashPassword(user.Password)
	result,_:=collection.InsertOne(ctx,user)
	json.NewEncoder(res).Encode(result)

}

//GET user by id
func getUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req) // get the params
	id,_ := primitive.ObjectIDFromHex(params["id"]) 
	var user User
	collection:= client.Database("appointy").Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err:=collection.FindOne(ctx,User{Id:id}).Decode(&user)
	if err!= nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	json.NewEncoder(res).Encode(user)
}

//GET all posts
func getPosts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var posts []Post
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor,err:=collection.Find(ctx,bson.M{})
	if err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var post Post
		cursor.Decode(&post)
		posts = append(posts,post)
	}
	if err:=cursor.Err() ; err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	json.NewEncoder(res).Encode(posts)
}

//POST create new post
func createPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	post.PostedAt.T = uint32(time.Now().Unix())
	result,_:=collection.InsertOne(ctx,post)
	json.NewEncoder(res).Encode(result)
}

//GET post by id
func getPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req) // get the params
	id,_ := primitive.ObjectIDFromHex(params["id"]) 
	var post Post
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err:=collection.FindOne(ctx,Post{Id:id}).Decode(&post)
	if err!= nil {
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	json.NewEncoder(res).Encode(post)
}

//GET all posts of user
func getPostsOfUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "application/json")
	params := mux.Vars(req) // get the params
	id,_ := primitive.ObjectIDFromHex(params["id"]) 
	var posts []Post
	collection:= client.Database("appointy").Collection("posts")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor,err:=collection.Find(ctx,Post{User:id})
	if err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx){
		var post Post
		cursor.Decode(&post)
		posts = append(posts,post)
	}
	if err:=cursor.Err() ; err!=nil{
		res.WriteHeader(http.StatusInternalServerError)
		res.Write([]byte(`{"message" :"`+ err.Error() +`"}`))
		return
	}
	json.NewEncoder(res).Encode(posts)
}

func main() {
	//dotenv 
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	
	//Mongodb connect
	mongoUri:= os.Getenv("MONGO_URI")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	clientOptions := options.Client().ApplyURI(mongoUri)
	client, _ = mongo.Connect(ctx, clientOptions)
	
	// Init router
	router := mux.NewRouter()

	//route handlers
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.HandleFunc("/users", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	router.HandleFunc("/posts/users/{id}", getPostsOfUser).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))

}
