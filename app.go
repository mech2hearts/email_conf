package main
import (
	"context"
	"fmt"
	"net/http"
	"time"
	"encoding/json"
	"encoding/base64"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/smtp"
	"os"
)

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email string             	`json:"email,omitempty" bson:"email,omitempty"`
	Verified  bool
	ConfLink string 			`json:"conflink,omitempty" bson:"conflink,omitempty"`
}

type Configuration struct {
	Port								string
	Sender 							string
	SenderPass 					string
	SenderServer 				string
	SenderServerPort 		string
	MongoLink 					string
}

var client *mongo.Client
var config = ReadConfig()

func ReadConfig() Configuration {
	var config Configuration
	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
    fmt.Println(err.Error())
  }
  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)
  return config
}

//Email functions
func hashAddress(x string) string {
        data := []byte(x)
        message := base64.StdEncoding.EncodeToString(data)
        return message
}




//MongoDB functions
func CreateUser(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	collection := client.Database("email_conf").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var person Person
	var person1 Person
	person.Verified = false
	json.NewDecoder(request.Body).Decode(&person)
	//checks if email already exists and if it is verified
	result1 := collection.FindOne(ctx, bson.M{"email" : person.Email}).Decode(&person1)
	if(result1==nil){
		if(person1.Verified==true){
		json.NewEncoder(response).Encode(bson.M{"message":"This email is already verified"})
		return
		} else {
			json.NewEncoder(response).Encode(bson.M{"message":"A confirmation link has already been sent to this email address"})
			return
		}
	}
	person.ConfLink = hashAddress(person.Email)
	//result, _ := collection.InsertOne(ctx, person)
	collection.InsertOne(ctx, person)
	sendEmail(person.ConfLink, person.Email)
	json.NewEncoder(response).Encode(bson.M{"message":"A confirmation link has been sent to the email address"})
}


func VerifyAccount(response http.ResponseWriter, request *http.Request){
	response.Header().Set("content-type", "application/json")
	var person1 Person
	params := mux.Vars(request)
	confLink, _ := params["conflink"]
	collection := client.Database("email_conf").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	result1 := collection.FindOne(ctx,bson.M{"conflink": confLink}).Decode(&person1)
	verified := person1.Email + " is now verified."
	vCopy := person1.Email + " is already verified."
	if(result1==nil){
		if(person1.Verified==true){
		json.NewEncoder(response).Encode(vCopy)
		return
		}
		collection.UpdateOne(ctx,bson.M{"conflink": confLink}, bson.M{ "$set": bson.M{ "verified" : true }})
		json.NewEncoder(response).Encode(verified)
		return
	}
	json.NewEncoder(response).Encode("Invalid Link")
}

func sendEmail(confLink string, receiver string){
	sender := config.Sender
  // Set up authentication information.
  auth := smtp.PlainAuth("", sender, config.SenderPass, "smtp.gmail.com")

	autoMessage := "Confirm your account by going to the following link: https://radiant-refuge-49608.herokuapp.com/api/verify/"+confLink
  // Connect to the server, authenticate, set the sender and recipient,
  // and send the email all in one step.
  to := []string{receiver}
  msg := []byte("To: "+receiver+"\r\n" +
          "Subject: Email Confirmation\r\n" +
          "\r\n" + autoMessage + ".\r\n")
  smtp.SendMail("smtp.gmail.com:587", auth, sender, to, msg)
}

func main() {
	fmt.Println(config.SenderServerPort)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI(config.MongoLink))
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/", CreateUser).Methods("POST")
	api.HandleFunc("/verify/{conflink}", VerifyAccount).Methods("GET")
	port := os.Getenv("PORT")
	if(port==""){
		port=config.Port
	}
	buildHandler := http.FileServer(http.Dir("./client/build"))
	router.PathPrefix("/").Handler(buildHandler)
	http.ListenAndServe(":"+port, router)

}
