package service

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/studiously/classsvc/models"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	GetClassEndpoint        endpoint.Endpoint
	CreateClassEndpoint     endpoint.Endpoint
	UpdateClassEndpoint     endpoint.Endpoint
	DeleteClassEndpoint     endpoint.Endpoint
	GetMemberEndpoint       endpoint.Endpoint
	GetMemberByUserEndpoint endpoint.Endpoint
	CreateMemberEndpoint    endpoint.Endpoint
	UpdateMemberEndpoint    endpoint.Endpoint
	DeleteMemberEndpoint    endpoint.Endpoint
	GetClassMembers         endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		GetClassEndpoint:        MakeGetClassEndpoint(s),
		CreateClassEndpoint:     MakeCreateClassEndpoint(s),
		UpdateClassEndpoint:     MakeUpdateClassEndpoint(s),
		DeleteClassEndpoint:     MakeDeleteClassEndpoint(s),
		GetMemberEndpoint:       MakeGetMemberEndpoint(s),
		GetMemberByUserEndpoint: MakeGetMemberByUserEndpoint(s),
		CreateMemberEndpoint:    MakeCreateMemberEndpoint(s),
		UpdateMemberEndpoint:    MakeUpdateMemberEndpoint(s),
		DeleteMemberEndpoint:    MakeDeleteMemberEndpoint(s),
		GetClassMembers:         MakeGetClassMembersEndpoint(s),
	}
}

func MakeGetClassesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		s.GetClasses(ctx, )
		return getClassesResponse{}, nil
	}
}


type getClassesResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeGetClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getClassRequest)
		class, e := s.GetClass(ctx, req.Id)
		return getClassResponse{class, e}, nil
	}
}

type getClassRequest struct {
	Id uuid.UUID `json:"id,omitempty"`
}

type getClassResponse struct {
	*models.Class `json:"class,omitempty"`
	Error error `json:"error,omitempty"`
}

func MakeCreateClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createClassRequest)
		id, e := s.CreateClass(ctx, req.Class)
		return createClassResponse{id, e}, nil
	}
}

type createClassRequest struct {
	*models.Class `json:"class"`
}

type createClassResponse struct {
	Id    uuid.UUID `json:"id"`
	Error error `json:"error,omitempty"`
}

func MakeUpdateClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateClassRequest)
		e := s.UpdateClass(ctx, req.Class)
		return updateClassResponse{e}, nil
	}
}

type updateClassRequest struct {
	*models.Class `json:"class"`
}

type updateClassResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeDeleteClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteClassRequest)
		e := s.DeleteClass(ctx, req.Id)
		return deleteClassResponse{e}, nil
	}
}

type deleteClassRequest struct {
	Id uuid.UUID `json:"id"`
}

type deleteClassResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeGetMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getMemberRequest)
		member, e := s.GetMember(ctx, req.Id)
		return getMemberResponse{member, e}, nil
	}
}

type getMemberRequest struct {
	Id uuid.UUID `json:"id"`
}

type getMemberResponse struct {
	*models.Member
	Error error `json:"error,omitempty"`
}

func MakeGetMemberByUserEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getMemberByUserRequest)
		member, e := s.GetMemberByUser(ctx, req.User, req.Class)
		return getMemberByUserResponse{member, e}, nil
	}
}

type getMemberByUserRequest struct {
	User  uuid.UUID `json:"user"`
	Class uuid.UUID `json:"class"`
}

type getMemberByUserResponse struct {
	*models.Member `json:"member"`
	Error error `json:"error,omitempty"`
}

func MakeCreateMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createMemberRequest)
		id, e := s.CreateMember(ctx, req.Member)
		return createMemberResponse{id, e}, nil
	}
}

type createMemberRequest struct {
	*models.Member `json:"member"`
}

type createMemberResponse struct {
	Id    uuid.UUID `json:"id"`
	Error error `json:"error,omitempty"`
}

func MakeDeleteMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteMemberRequest)
		e := s.DeleteMember(ctx, req.Id)
		return deleteMemberResponse{e}, nil
	}
}

type deleteMemberRequest struct {
	Id uuid.UUID `json:"id"`
}

type deleteMemberResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeUpdateMemberEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateMemberRequest)
		e := s.UpdateMember(ctx, req.Member)
		return updateMemberResponse{e}, nil
	}
}

type updateMemberRequest struct {
	*models.Member `json:"member"`
}

type updateMemberResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeGetClassMembersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getClassMembersRequest)
		members, e := s.GetClassMembers(ctx, req.Class)
		return getClassMembersResponse{members, e}, nil
	}
}

type getClassMembersRequest struct {
	Class uuid.UUID `json:"id"`
}

type getClassMembersResponse struct {
	Members []*models.Member `json:"members"`
	Error   error `json:"error,omitempty"`
}
