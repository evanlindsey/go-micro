// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// Error defines model for Error.
type Error struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

// Pet defines model for Pet.
type Pet struct {
	Id   int64   `json:"id"`
	Name string  `json:"name"`
	Tag  *string `json:"tag,omitempty"`
}

// Pets defines model for Pets.
type Pets = []Pet

// ListPetsParams defines parameters for ListPets.
type ListPetsParams struct {
	// Limit How many items to return at one time (max 100)
	Limit *int32 `form:"limit,omitempty" json:"limit,omitempty"`
}

// CreatePetsJSONRequestBody defines body for CreatePets for application/json ContentType.
type CreatePetsJSONRequestBody = Pet

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// List all pets
	// (GET /pets)
	ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams)
	// Create a pet
	// (POST /pets)
	CreatePets(w http.ResponseWriter, r *http.Request)
	// Info for a specific pet
	// (GET /pets/{petId})
	ShowPetById(w http.ResponseWriter, r *http.Request, petId string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// List all pets
// (GET /pets)
func (_ Unimplemented) ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create a pet
// (POST /pets)
func (_ Unimplemented) CreatePets(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Info for a specific pet
// (GET /pets/{petId})
func (_ Unimplemented) ShowPetById(w http.ResponseWriter, r *http.Request, petId string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// ListPets operation middleware
func (siw *ServerInterfaceWrapper) ListPets(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params ListPetsParams

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", r.URL.Query(), &params.Limit)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "limit", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ListPets(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreatePets operation middleware
func (siw *ServerInterfaceWrapper) CreatePets(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreatePets(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// ShowPetById operation middleware
func (siw *ServerInterfaceWrapper) ShowPetById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "petId" -------------
	var petId string

	err = runtime.BindStyledParameterWithOptions("simple", "petId", chi.URLParam(r, "petId"), &petId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "petId", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.ShowPetById(w, r, petId)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pets", wrapper.ListPets)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/pets", wrapper.CreatePets)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/pets/{petId}", wrapper.ShowPetById)
	})

	return r
}

type ListPetsRequestObject struct {
	Params ListPetsParams
}

type ListPetsResponseObject interface {
	VisitListPetsResponse(w http.ResponseWriter) error
}

type ListPets200ResponseHeaders struct {
	XNext string
}

type ListPets200JSONResponse struct {
	Body    Pets
	Headers ListPets200ResponseHeaders
}

func (response ListPets200JSONResponse) VisitListPetsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-next", fmt.Sprint(response.Headers.XNext))
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response.Body)
}

type ListPetsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response ListPetsdefaultJSONResponse) VisitListPetsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type CreatePetsRequestObject struct {
	Body *CreatePetsJSONRequestBody
}

type CreatePetsResponseObject interface {
	VisitCreatePetsResponse(w http.ResponseWriter) error
}

type CreatePets201Response struct {
}

func (response CreatePets201Response) VisitCreatePetsResponse(w http.ResponseWriter) error {
	w.WriteHeader(201)
	return nil
}

type CreatePetsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response CreatePetsdefaultJSONResponse) VisitCreatePetsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type ShowPetByIdRequestObject struct {
	PetId string `json:"petId"`
}

type ShowPetByIdResponseObject interface {
	VisitShowPetByIdResponse(w http.ResponseWriter) error
}

type ShowPetById200JSONResponse Pet

func (response ShowPetById200JSONResponse) VisitShowPetByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ShowPetByIddefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response ShowPetByIddefaultJSONResponse) VisitShowPetByIdResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// List all pets
	// (GET /pets)
	ListPets(ctx context.Context, request ListPetsRequestObject) (ListPetsResponseObject, error)
	// Create a pet
	// (POST /pets)
	CreatePets(ctx context.Context, request CreatePetsRequestObject) (CreatePetsResponseObject, error)
	// Info for a specific pet
	// (GET /pets/{petId})
	ShowPetById(ctx context.Context, request ShowPetByIdRequestObject) (ShowPetByIdResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// ListPets operation middleware
func (sh *strictHandler) ListPets(w http.ResponseWriter, r *http.Request, params ListPetsParams) {
	var request ListPetsRequestObject

	request.Params = params

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ListPets(ctx, request.(ListPetsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListPets")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ListPetsResponseObject); ok {
		if err := validResponse.VisitListPetsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// CreatePets operation middleware
func (sh *strictHandler) CreatePets(w http.ResponseWriter, r *http.Request) {
	var request CreatePetsRequestObject

	var body CreatePetsJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.CreatePets(ctx, request.(CreatePetsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "CreatePets")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(CreatePetsResponseObject); ok {
		if err := validResponse.VisitCreatePetsResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// ShowPetById operation middleware
func (sh *strictHandler) ShowPetById(w http.ResponseWriter, r *http.Request, petId string) {
	var request ShowPetByIdRequestObject

	request.PetId = petId

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.ShowPetById(ctx, request.(ShowPetByIdRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ShowPetById")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(ShowPetByIdResponseObject); ok {
		if err := validResponse.VisitShowPetByIdResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}
