package main

import (
	"time"
	"context"
	"encoding/json"
	"github.com/joho/godotenv"
    "log"
    "os"
	"math/rand"
	"net/http"
	"strconv"
	"go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gorilla/mux"
)

//Post Struct (Model)
type Post struct {
	Id       string `json:"_id"`
	Caption  string `json:"caption"`
	ImageURL string `json:"image"`
	User     *User  `json:"user"`
}

//User Struct (Model)
type User struct {
	Id       string `json:"_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//init posts as a slice Post struct
var posts []Post
var users []User

//GET all users
func getUsers(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(users)
}

//POST create new user
func createUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	user.Id = strconv.Itoa(rand.Intn(10000000)) // Mock Id -not safe
	users = append(users, user)
	json.NewEncoder(res).Encode(user)

}

//GET user by id
func getUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // get the params
	//loop through posts and find by id
	for _, item := range users {
		if item.Id == params["id"] {
			json.NewEncoder(res).Encode(item)
			return
		}
	}
	json.NewEncoder(res).Encode(&User{})
}

//GET all posts
func getPosts(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(posts)
}

//POST create new post
func createPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var post Post
	_ = json.NewDecoder(req.Body).Decode(&post)
	post.Id = strconv.Itoa(rand.Intn(10000000)) // Mock Id -not safe
	posts = append(posts, post)
	json.NewEncoder(res).Encode(post)
}

//GET post by id
func getPost(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // get the params
	//loop through posts and find by id
	for _, item := range posts {
		if item.Id == params["id"] {
			json.NewEncoder(res).Encode(item)
			return
		}
	}
	json.NewEncoder(res).Encode(&Post{})
}

//GET all posts of user
func getPostsOfUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) // get the params
	//loop through posts and find by id
	var postsOfUser []Post
	for _, item := range posts {
		if item.User.Id == params["id"] {
			postsOfUser = append(postsOfUser, item)
		}
	}
	json.NewEncoder(res).Encode(postsOfUser)
}

func main() {
	//dotenv for mongodb uri
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//Mongodb connect
	mongoUri:= os.Getenv("MONGO_URI")
	ctx,_:= context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
	defer client.Disconnect(ctx)

	// Init router
	router := mux.NewRouter()

	//Mock data - @todo - implement database
	posts = append(posts, Post{Id: "1", Caption: "watch", ImageURL: "https://www.youtube.com/watch?v=SonwZ6MF5BE", User: &User{Id: "1", Name: "Rufus", Email: "sam2001rufus@gmail.com", Password: "helloworld"}})
	users = append(users, User{Id: "1", Name: "Rufus", Email: "sam2001rufus@gmail.com", Password: "helloworld"})
	posts = append(posts, Post{Id: "2", Caption: "listen", ImageURL: "https://www.youtube.com/watch?v=SonwZ6MF5BE", User: &User{Id: "2", Name: "Sam", Email: "sam2001@gmail.com", Password: "hello"}})
	posts = append(posts, Post{Id: "3", Caption: "read", ImageURL: "https://www.youtube.com/watch?v=SonwZ6MF5BE", User: &User{Id: "2", Name: "Sam", Email: "sam2001@gmail.com", Password: "hello"}})
	users = append(users, User{Id: "2", Name: "Sam", Email: "sam2001@gmail.com", Password: "hello"})

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
