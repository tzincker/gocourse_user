package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_user/internal/user"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {

	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser,
		encodeResponse,
		opts...,
	)).Methods("POST")

	router.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUsers,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser,
		encodeResponse,
		opts...,
	)).Methods("GET")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser,
		encodeResponse,
		opts...,
	)).Methods("PATCH")

	router.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser,
		encodeResponse,
		opts...,
	)).Methods("DELETE")

	return router
}

func decodeCreateUser(_ context.Context, r *http.Request) (any, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	return req, nil
}

func decodeGetAllUsers(_ context.Context, r *http.Request) (any, error) {
	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}

	return req, nil
}

func decodeGetUser(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}

	return req, nil
}

func decodeUpdateUser(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	id := p["id"]

	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}

	req.ID = id
	return req, nil
}

func decodeDeleteUser(_ context.Context, r *http.Request) (any, error) {
	p := mux.Vars(r)
	req := user.DeleteReq{
		ID: p["id"],
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp any) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utd-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utd-8")

	resp, ok := err.(response.Response)

	if !ok {
		newResponse := response.BadRequest("error parsing body")
		w.WriteHeader(newResponse.StatusCode())
		_ = json.NewEncoder(w).Encode(newResponse)
		return
	}

	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)

}
