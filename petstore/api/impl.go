package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/evanlindsey/go-micro/petstore/ent"
)

type Server struct {
	client *ent.Client
}

func NewServer(client *ent.Client) *Server {
	return &Server{client: client}
}

// (GET /pets)
func (s *Server) ListPets(ctx context.Context, req ListPetsRequestObject) (ListPetsResponseObject, error) {
	log.Println("Received request to ListPets")

	var pets []*ent.Pet
	var err error

	if req.Params.Limit != nil {
		log.Printf("Limit provided: %d", *req.Params.Limit)
		pets, err = s.client.Pet.Query().Limit(int(*req.Params.Limit)).All(ctx)
	} else {
		pets, err = s.client.Pet.Query().All(ctx)
	}
	if err != nil {
		res := ListPetsdefaultJSONResponse{Body: Error{Code: http.StatusInternalServerError, Message: fmt.Sprintf("failed to retrieve pets: %s", err)}, StatusCode: http.StatusInternalServerError}
		return res, nil
	}

	log.Printf("Successfully retrieved %d pets", len(pets))

	res := make([]Pet, len(pets))
	for i, pet := range pets {
		res[i] = Pet{
			Id:   pet.ID,
			Name: pet.Name,
			Tag:  &pet.Tag,
		}
	}
	return ListPets200JSONResponse{Body: res}, nil
}

// (POST /pets)
func (s *Server) CreatePets(ctx context.Context, req CreatePetsRequestObject) (CreatePetsResponseObject, error) {
	log.Println("Received request to CreatePets")

	if req.Body.Name == "" {
		res := CreatePetsdefaultJSONResponse{Body: Error{Code: http.StatusBadRequest, Message: "name is required"}, StatusCode: http.StatusBadRequest}
		return res, nil
	}

	pet, err := s.client.Pet.Create().
		SetName(req.Body.Name).
		SetNillableTag(req.Body.Tag).
		Save(ctx)
	if err != nil {
		res := CreatePetsdefaultJSONResponse{Body: Error{Code: http.StatusInternalServerError, Message: fmt.Sprintf("failed to create pet: %s", err)}, StatusCode: http.StatusInternalServerError}
		return res, nil
	}

	log.Printf("Successfully created pet: ID: %d, Name: %s", pet.ID, pet.Name)

	res := CreatePets201Response{}
	return res, nil
}

// (GET /pets/{petId})
func (s *Server) ShowPetById(ctx context.Context, req ShowPetByIdRequestObject) (ShowPetByIdResponseObject, error) {
	log.Printf("Received request to ShowPetById with petId: %s", req.PetId)

	petId, err := strconv.ParseInt(req.PetId, 10, 64)
	if err != nil {
		res := ShowPetByIddefaultJSONResponse{Body: Error{Code: http.StatusBadRequest, Message: fmt.Sprintf("invalid petId: %s", err)}, StatusCode: http.StatusBadRequest}
		return res, nil
	}

	pet, err := s.client.Pet.Get(ctx, petId)
	if err != nil {
		res := ShowPetByIddefaultJSONResponse{Body: Error{Code: http.StatusNotFound, Message: fmt.Sprintf("failed to retrieve pet: %s", err)}, StatusCode: http.StatusNotFound}
		return res, nil
	}

	log.Printf("Successfully retrieved pet: ID: %d, Name: %s", pet.ID, pet.Name)

	res := ShowPetById200JSONResponse{
		Id:   pet.ID,
		Name: pet.Name,
		Tag:  &pet.Tag,
	}
	return res, nil
}
