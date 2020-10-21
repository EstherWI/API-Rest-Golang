package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"ConexaoSolar/helper"
	"ConexaoSolar/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = helper.ConnectDB()

//GetAll retorna todos os usuários
func GetAll(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// we created Usuario array
	var usuarios []models.Usuario

	// bson.M{},  we passed empty filter. So we want to get all data.
	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	/*A defer statement defers the execution of a function until the surrounding function returns.
	simply, run cur.Close() process but after cur.Next() finished.*/
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var usuario models.Usuario
		// & character returns the memory address of the following variable.
		err := cur.Decode(&usuario) // decode similar to deserialize process.
		if err != nil {
			log.Fatal(err)
		}

		// add item our array
		usuarios = append(usuarios, usuario)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(usuarios) // encode similar to serialize process.
}

//GetByID Retorna um usuário pelo seu ID
func GetByID(w http.ResponseWriter, r *http.Request) {

	// set header.
	w.Header().Set("Content-Type", "application/json")

	var usuario models.Usuario
	// we get params with mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// We create filter. If it is unnecessary to sort data for you, you can use bson.M{}
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(usuario)
}

//Create cria um novo usuário
func Create(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var usuario models.Usuario

	// we decode our body request params
	_ = json.NewDecoder(r.Body).Decode(&usuario)

	result, err := collection.InsertOne(context.TODO(), usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(result)
}

//Update atualiza informações de um usuário
func Update(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Get id from parameters
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var usuario models.Usuario

	// Create filter
	filter := bson.M{"_id": id}

	// Read update model from body request
	_ = json.NewDecoder(r.Body).Decode(&usuario)

	// prepare update model.
	update := bson.D{
		{"$set", bson.D{
			{"name", usuario.Name},
			{"email", usuario.Email},
			{"senha", usuario.Senha},
		}},
	}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	usuario.ID = id

	json.NewEncoder(w).Encode(usuario)
}

//Delete exclui um usuário
func Delete(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// get params
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// prepare filter.
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

// var client *mongo.Client

func main() {
	//Init Router
	r := mux.NewRouter()

	r.HandleFunc("/ConexaoSolar", GetAll).Methods("GET")
	r.HandleFunc("/ConexaoSolar", Create).Methods("POST")
	r.HandleFunc("/ConexaoSolar", Update).Methods("PUT")
	r.HandleFunc("/ConexaoSolar", Delete).Methods("DELETE")
	r.HandleFunc("/ConexaoSolar/{id}", GetByID).Methods("GET")

	// set our port address
	var port = ":8000"
	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(port, r))

}
