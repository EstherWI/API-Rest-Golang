package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"UserCRUD/helper"
	"UserCRUD/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Response estrutura
type Response struct {
	Method  string `json:"method"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// New Response
func NewResponse(method, message string, status int) Response {

	return Response{Method: method, Message: message, Status: status}

}

//GetAll : retorna todos os usuários
func GetAll(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	w.Header().Set("Content-Type", "application/json")
	HttpInfo(r)
	// array para Usuario
	var usuarios []models.Usuario

	cur, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		helper.GetError(err, w)
		return
	}

	// Close the cursor once finished
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		// cria um valor no qual o único documento possa ser decodificado
		var usuario models.Usuario
		// & retorna o endereço de memória da seguinte variável.
		err := cur.Decode(&usuario) // decode semelhante ao processo de desserialização.
		if err != nil {
			log.Fatal(err)
		}

		// add item no array
		usuarios = append(usuarios, usuario)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	json.NewEncoder(w).Encode(usuarios) // encode semelhante ao processo de serialização.
}

//GetByID : Retorna um usuário pelo seu ID
func GetByID(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	// set header.
	w.Header().Set("Content-Type", "application/json")

	var usuario models.Usuario
	// resgatando parametros com mux.
	var params = mux.Vars(r)

	// string to primitive.ObjectID
	id, _ := primitive.ObjectIDFromHex(params["id"])

	// Filtro pelo id.
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(usuario)
}

//GetByEmail : retorna um usuário pelo email
func GetByEmail(email string) (*models.Usuario, error) {
	var collection = helper.ConnectDB()
	var usuario models.Usuario
	fmt.Println(email)
	filter := bson.M{"email": email}
	err := collection.FindOne(context.TODO(), filter).Decode(&usuario)

	if err != nil {
		return nil, err
	}
	return &usuario, nil
}

//Create : cria um novo usuário
func Create(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	var output *models.Usuario

	w.Header().Set("Content-Type", "application/json")

	var usuario models.Usuario
	_ = json.NewDecoder(r.Body).Decode(&usuario)
	output, _ = GetByEmail(usuario.Email)

	if output != nil {
		json.NewEncoder(w).Encode(NewResponse(r.Method, "failed", 400))
		return
	}

	// inserindo usuario
	result, err := collection.InsertOne(context.TODO(), usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(NewResponse(r.Method, "sucess", 201))
	fmt.Println(result)
}

//Update : atualiza informações de um usuário
func Update(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	//Pegando id dos parametros
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var usuario models.Usuario

	// Filtro pelo id
	filter := bson.M{"_id": id}

	_ = json.NewDecoder(r.Body).Decode(&usuario)

	// Atualizando model.
	update := bson.D{
		{"$set", bson.D{
			{"name", usuario.Name},
			{"email", usuario.Email},
			{"senha", usuario.Senha},
		},
		}}

	err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&usuario)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	usuario.ID = id

	json.NewEncoder(w).Encode(usuario)
}

//Delete : exclui um usuário
func Delete(w http.ResponseWriter, r *http.Request) {
	var collection = helper.ConnectDB()
	// Set header
	w.Header().Set("Content-Type", "application/json")

	// Resgata parâmetros
	var params = mux.Vars(r)

	// string to primitve.ObjectID
	id, err := primitive.ObjectIDFromHex(params["id"])

	// Filtro pelo id
	filter := bson.M{"_id": id}

	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	json.NewEncoder(w).Encode(deleteResult)
}

//HttpInfo
func HttpInfo(r *http.Request) {
	fmt.Printf("%s\t %s\t %s%s\r\n", r.Method, r.Proto, r.Host, r.URL)
}
func main() {
	//Init Router
	r := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Request", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})
	//r.HandleFunc("/", handler)
	r.HandleFunc("/users", GetAll).Methods("GET")
	r.HandleFunc("/users", Create).Methods("POST")
	r.HandleFunc("/users/{id}", Update).Methods("PUT")
	r.HandleFunc("/users/{id}", Delete).Methods("DELETE")
	r.HandleFunc("/users/{id}", GetByID).Methods("GET")
	var port = ":8000"
	fmt.Println("Server running in port:", port)
	log.Fatal(http.ListenAndServe(":8000", handlers.CORS(headers, methods, origins)(r)))

	//log.Fatal(http.ListenAndServe(port, r))

}
