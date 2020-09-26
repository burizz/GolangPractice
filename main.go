package main

import (
	"context"

	// Built-in
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect" // get an object type
	"time"

	// JSON
	"encoding/json"

	// HTTP
	"github.com/gorilla/mux"

	// Mongo
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Person - team members
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	MemberID  string             `json:"memberid,omitempty" bson:"memberid,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Hours     int                `json:"hours,omitempty" bson:"hours,omitempty"`
	Month     string             `json:"month,omitempty" bson:"month,omitempty"`
}

var listOfPeople = []Person{
	// Aug 2019
	Person{MemberID: "1", Firstname: "Boris", Lastname: "Yakimov", Hours: 108},
	Person{MemberID: "2", Firstname: "Vname2", Lastname: "name2", Hours: 157},
	Person{MemberID: "3", Firstname: "Kname3", Lastname: "name3", Hours: 111},
	Person{MemberID: "4", Firstname: "Yname4", Lastname: "name4", Hours: 116},
	Person{MemberID: "5", Firstname: "Kaname5", Lastname: "name5", Hours: 162},
	Person{MemberID: "6", Firstname: "Krname6", Lastname: "name6", Hours: 0},
}

var client *mongo.Client

func main() {
	log.Println("Starting application ...")

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Cancel the context created for MongoDB, to prevent memory leaks - https://www.sohamkamani.com/golang/2018-06-17-golang-using-context-cancellation/
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, clientOptions)

	// Test if connected
	log.Println("Connecting to DB ...")
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to MongoDB!")
	}

	// https://kb.objectrocket.com/mongo-db/how-to-insert-mongodb-documents-from-json-using-the-golang-driver-457
	seedDB()

	// API - Mux Router - responses should be in JSON
	requestRouter := mux.NewRouter()
	log.Println("REST API - Mux Request Router")

	// Routes
	requestRouter.HandleFunc("/", homePage)
	requestRouter.HandleFunc("/team", getTeam)
	requestRouter.HandleFunc("/person/{ID}", getPerson)
	requestRouter.HandleFunc("/addToTeam", addToTeam).Methods("POST")
	requestRouter.HandleFunc("/delFromTeam/{ID}", deleteFromTeam).Methods("DELETE")
	requestRouter.HandleFunc("/updatePerson/{ID}", updatePerson).Methods("PUT")

	// Serve
	httpPort := "10000"
	log.Print("HTTP Listening on port [" + httpPort + "]")
	if err := http.ListenAndServe(":"+httpPort, requestRouter); err != nil {
		log.Fatal(err)
	}
}

//TODO: Update if already exist
func seedDB() {
	// Load input json with Seed data
	byteValues, err := ioutil.ReadFile("seed_data/input_data_latest.json")
	if err != nil {
		fmt.Println("ioutil.ReadFile ERROR:", err)
	} else {
		// Verify Seed Data
		fmt.Println("ioutil.ReadFile byteValues TYPE:", reflect.TypeOf(byteValues))
		fmt.Println("byteValues:", byteValues)
		fmt.Println("byteValues:", string(byteValues))
	}

	// Mongo config
	collection := client.Database("gofte").Collection("people")
	// Context to add a timeout of whatever is executed in it
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var people []Person
	if err := json.Unmarshal(byteValues, &people); err != nil {
		log.Fatal(err)
	}
	fmt.Println("People :", reflect.TypeOf(people))

	for i := range people {
		person := people[i]
		fmt.Println("nperson _id:", person.ID)
		fmt.Println("person Field Str:", person.ID)
		// Insert in Mongo
		result, insertErr := collection.InsertOne(ctx, person)

		// Check for insertion errors
		if insertErr != nil {
			fmt.Println("InsertOne ERROR:", insertErr)
		} else {
			fmt.Println("InstertOne() API result:", result)
		}
	}
}

// API //

func homePage(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, "Home page")
	log.Println("Endpoint hit: homePage")
}

func getTeam(response http.ResponseWriter, request *http.Request) {
	log.Println("Endpoint hit: getTeam")
	response.Header().Set("content-type", "application/json")

	// var team []Person
	var team []bson.M

	collection := client.Database("gofte").Collection("people")
	// Context to add a timeout of whatever is executed in it
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Cursor - pointer to collection - db.collection.find () by default
	// result of the query will be iterated automatically and returned
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}"`))
		return
	}

	// defer cursor.Close(ctx) // Not sure about this, IDE complains if cancel is not used in ctx

	// Iterate cursor object and build array of team
	for cursor.Next(ctx) {
		// var person Person
		var person bson.M
		cursor.Decode(&person)
		team = append(team, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `"}"`))
		return
	}

	// TODO:
	// sort team before returning it - https://golang.org/pkg/sort/
	json.NewEncoder(response).Encode(team)
}

// Retrieve member by either value - except hours
func getPerson(response http.ResponseWriter, request *http.Request) {
	log.Println("Endpoint hit: getPerson")
	response.Header().Set("content-type", "application/json")

	// Parse request vars
	vars := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	//memberID := vars["MemberID"]

	var person bson.M

	// Set DB and collection to use
	collection := client.Database("gofte").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use Person structure to find the proper field, doesn't work with bson.M()
	if err := collection.FindOne(ctx, Person{ID: id}).Decode(&person); err != nil {
		log.Printf("Cannot find teamMember with ID [%v]", id)
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + err.Error() + `" }`))
		return
	}

	json.NewEncoder(response).Encode(person)
}

// Add a new Team Member - HTTP POST
func addToTeam(response http.ResponseWriter, request *http.Request) {
	log.Println("Endpoint hit: AddToTeam")
	response.Header().Set("content-type", "application/json")

	var newMember Person
	var alreadyExist bool

	// Parse request body
	json.NewDecoder(request.Body).Decode(&newMember)

	// Set DB and collection to use
	collection := client.Database("gofte").Collection("people")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if newMember already exists
	// for _, teamMember := range listOfPeople {
	// 	if teamMember.ID == newMember.ID {
	// 		alreadyExist = true
	// 	} else {
	// 		alreadyExist = false
	// 	}
	// }

	alreadyExist = false

	if alreadyExist {
		log.Printf("A Team Member with ID [%v] already exists", newMember.ID)
	} else {
		// Add newMember
		log.Printf("Adding new Team Member: %v", newMember)
		listOfPeople = append(listOfPeople, newMember)

		// Insert into Mongo
		result, err := collection.InsertOne(ctx, newMember)
		if err != nil {
			log.Fatal(err)
		}

		// Response
		json.NewEncoder(response).Encode(result)
	}
}

// Delete Memberequest by ID - HTTP DELETE
func deleteFromTeam(response http.ResponseWriter, request *http.Request) {
	log.Println("Endpoint hit: delFromTeam")
	// Similar to getPerson() logic of finding person by mongo bson ID
	// but use delete equivalent of mongo lib
	vars := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(vars["id"])
	// id := vars["ID"]

	for index, teamMember := range listOfPeople {
		if teamMember.ID == id {
			// Do it in order
			//listOfPeople = append(listOfPeople[:index], listOfPeople[index+1:]...)
			// Do it faster but do not respect order
			log.Printf("Removing %v from Team", listOfPeople[index].Firstname)
			listOfPeople[index] = listOfPeople[len(listOfPeople)-1]
			listOfPeople = append(listOfPeople[:len(listOfPeople)-1])
		} else {
			log.Printf("Team member ID %v doesn't exist", teamMember.ID)
		}
	}
}

// Update existing member by ID - HTTP PUT
func updatePerson(response http.ResponseWriter, request *http.Request) {
	log.Println("Endpoint hit: updatePerson")
	// Similar to getPerson() logic of finding person by mongo bson ID
	// than update and store it using same logic as addToTeam()

	// Parse HTTP PUT body
	reqBody, _ := ioutil.ReadAll(request.Body)
	var updatePerson Person
	json.Unmarshal(reqBody, &updatePerson)

	for index, teamMember := range listOfPeople {
		if teamMember.ID == updatePerson.ID {
			// Works but doesn't keep the order of elements
			// will make a sort in the get team function
			listOfPeople = append(listOfPeople[:index], listOfPeople[index+1:]...)
			listOfPeople = append(listOfPeople, updatePerson)

			log.Printf("Updated member [%v] with new values [%v]", teamMember.Firstname, updatePerson)
		} else {
			log.Printf("No members with matching ID [%v] found", updatePerson.ID)
		}
	}
}

// END API //

// TODO:
// currentTime := time.Now()
// currentDate := currentTime.Format("01-Jan-2006")
//
// // Get hours of each member
// for i := 0; i < len(listOfPeople); i++ {
// 	// Idiomatic to Go - assign variable and error in same if ; while calling the actual function
// 	if percentUtil, err := calculateUtilization(listOfPeople[i].Hours); err != nil {
// 		log.Printf("input errror: %v\n", err)
// 	} else {
// 		printUtilization(listOfPeople[i].Firstname, listOfPeople[i].Hours, int(percentUtil), currentDate)
// 	}
// }

// TODO:
// Write tests - https://golang.org/doc/code.html#Testing

func calculateUtilization(trackedHours int) (float32, error) {
	if trackedHours < 0 {
		// check if negative number
		return float32(trackedHours), errors.New("Tracked hours should be a positive number")
	}

	// Calculate remaining percent to fullFte
	fullFte := 168
	percentUtil := (float32(trackedHours) / float32(fullFte)) * 100
	return percentUtil, nil
}

func printUtilization(name string, hours int, percent int, date string) {
	log.Printf("%v tracked [%v] hours ; [%v%%] utilization ; %v\n", name, hours, percent, date)
}
