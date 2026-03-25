package user

import (
	"context"
	"errors"

	"github.com/tzincker/go_lib_response/response"
	"github.com/tzincker/gocourse_meta/meta"
)

type (
	Controller func(ctx context.Context, request any) (any, error)

	Endpoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateReq struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	GetReq struct {
		ID string
	}

	GetAllReq struct {
		FirstName string
		LastName  string
		Limit     int
		Page      int
	}

	UpdateReq struct {
		ID        string
		FirstName *string `json:"first_name"`
		LastName  *string `json:"last_name"`
		Email     *string `json:"email"`
		Phone     *string `json:"phone"`
	}

	DeleteReq struct {
		ID string
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) Endpoints {
	return Endpoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(CreateReq)

		if req.FirstName == "" {
			return nil, response.BadRequest("firstname is required")
		}

		if req.LastName == "" {
			return nil, response.BadRequest("lastname is required")
		}

		if req.Email == "" {
			return nil, response.BadRequest("email is required")
		}

		if req.Phone == "" {
			return nil, response.BadRequest("phone is required")
		}

		user, err := s.Create(ctx, req.FirstName, req.LastName, req.Email, req.Phone)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.Created("success", user, nil), nil
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetReq)

		user, err := s.Get(ctx, req.ID)
		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", user, nil), nil
	}
}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetAllReq)
		filters := Filters{
			FirstName: req.FirstName,
			LastName:  req.LastName,
		}

		count, err := s.Count(ctx, filters)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		meta, err := meta.New(req.Page, req.Limit, count, config.LimPageDef)
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		users, err := s.GetAll(ctx, filters, meta.Offset(), meta.Limit())
		if err != nil {
			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", users, meta), nil
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {

		req := request.(UpdateReq)

		if req.FirstName != nil && *req.FirstName == "" {
			return nil, response.BadRequest(ErrFirstNameRequired.Error())
		}

		if req.LastName != nil && *req.LastName == "" {
			return nil, response.BadRequest(ErrLastNameRequired.Error())
		}

		err := s.Update(ctx, req.ID, req.FirstName, req.LastName, req.Email, req.Phone)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DeleteReq)
		err := s.Delete(ctx, req.ID)

		if err != nil {
			if _, ok := errors.AsType[*ErrNotFound](err); ok {
				return nil, response.NotFound(err.Error())
			}

			return nil, response.InternalServerError(err.Error())
		}

		return response.OK("success", nil, nil), nil
	}
}
