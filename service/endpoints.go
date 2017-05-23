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
	ListClassesEndpoint endpoint.Endpoint
	GetClassEndpoint    endpoint.Endpoint
	CreateClassEndpoint endpoint.Endpoint
	UpdateClassEndpoint endpoint.Endpoint
	DeleteClassEndpoint endpoint.Endpoint
	JoinClassEndpoint   endpoint.Endpoint
	SetRoleEndpoint     endpoint.Endpoint
	LeaveClassEndpoint  endpoint.Endpoint
	ListMembersEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		ListClassesEndpoint: MakeListClassesEndpoint(s),
		GetClassEndpoint:    MakeGetClassEndpoint(s),
		CreateClassEndpoint: MakeCreateClassEndpoint(s),
		UpdateClassEndpoint: MakeUpdateClassEndpoint(s),
		DeleteClassEndpoint: MakeDeleteClassEndpoint(s),
		//GetMemberEndpoint:   MakeGetMemberEndpoint(s),
		//GetMemberByUserEndpoint: MakeGetMemberByUserEndpoint(s),
		JoinClassEndpoint:   MakeJoinClassEndpoint(s),
		SetRoleEndpoint:     MakeSetRoleEndpoint(s),
		LeaveClassEndpoint:  MakeLeaveClassEndpoint(s),
		ListMembersEndpoint: MakeListMembersEndpoint(s),
	}
}

func MakeListClassesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		classes, err := s.ListClasses(ctx)
		return listClassesResponse{classes, err}, nil
	}
}

type listClassesResponse struct {
	Classes []uuid.UUID
	Error   error `json:"error,omitempty"`
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
		id, e := s.CreateClass(ctx, req.Name)
		return createClassResponse{id, e}, nil
	}
}

type createClassRequest struct {
	Name string `json:"name"`
}

type createClassResponse struct {
	Id    uuid.UUID `json:"id"`
	Error error `json:"error,omitempty"`
}

func MakeUpdateClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(updateClassRequest)
		e := s.UpdateClass(ctx, req.Class, req.Name, req.CurrentUnit)
		return updateClassResponse{e}, nil
	}
}

type updateClassRequest struct {
	Class       uuid.UUID
	Name        string `json:"class"`
	CurrentUnit uuid.UUID `json:"current_unit"`
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

//func MakeGetMemberEndpoint(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(getMemberRequest)
//		member, e := s.GetMember(ctx, req.Id)
//		return getMemberResponse{member, e}, nil
//	}
//}
//
//type getMemberRequest struct {
//	Id uuid.UUID `json:"id"`
//}
//
//type getMemberResponse struct {
//	*models.Member
//	Error error `json:"error,omitempty"`
//}

//func MakeGetMemberByUserEndpoint(s Service) endpoint.Endpoint {
//	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
//		req := request.(getMemberByUserRequest)
//		member, e := s.GetMemberByUser(ctx, req.User, req.Class)
//		return getMemberByUserResponse{member, e}, nil
//	}
//}
//
//type getMemberByUserRequest struct {
//	User  uuid.UUID `json:"user"`
//	Class uuid.UUID `json:"class"`
//}
//
//type getMemberByUserResponse struct {
//	*models.Member `json:"member"`
//	Error error `json:"error,omitempty"`
//}

func MakeJoinClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(joinClassRequest)
		e := s.JoinClass(ctx, req.Class)
		return joinClassResponse{e}, nil
	}
}

type joinClassRequest struct {
	Class uuid.UUID `json:"class"`
}

type joinClassResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeLeaveClassEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(leaveClassRequest)
		e := s.LeaveClass(ctx, req.User, req.Class)
		return leaveClassResponse{e}, nil
	}
}

type leaveClassRequest struct {
	User  uuid.UUID `json:"user,omitempty"`
	Class uuid.UUID `json:"class"`
}

type leaveClassResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeSetRoleEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(setRoleRequest)
		e := s.SetRole(ctx, req.User, req.Class, req.Role)
		return setRoleResponse{e}, nil
	}
}

type setRoleRequest struct {
	User  uuid.UUID `json:"user"`
	Class uuid.UUID `json:"class"`
	Role  models.UserRole `json:"role"`
}

type setRoleResponse struct {
	Error error `json:"error,omitempty"`
}

func MakeListMembersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listMembersRequest)
		members, e := s.ListMembers(ctx, req.Class)
		return getClassMembersResponse{members, e}, nil
	}
}

type listMembersRequest struct {
	Class uuid.UUID `json:"id"`
}

type getClassMembersResponse struct {
	Members []*models.Member `json:"members"`
	Error   error `json:"error,omitempty"`
}
