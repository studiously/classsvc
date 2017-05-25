package service

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
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

// MakeClientEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the remote instance, via a transport/http.Client.
// Useful in a profilesvc client.
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	// Note that the request encoders need to modify the request URL, changing
	// the path and method. That's fine: we simply need to provide specific
	// encoders for each endpoint.

	return Endpoints{
		ListClassesEndpoint: httptransport.NewClient("GET", tgt, encodeListClassesRequest, decodeListClassesResponse, options...).Endpoint(),
		GetClassEndpoint:    httptransport.NewClient("GET", tgt, encodeGetClassRequest, decodeGetClassResponse, options...).Endpoint(),
		CreateClassEndpoint: httptransport.NewClient("POST", tgt, encodeCreateClassRequest, decodeCreateClassResponse, options...).Endpoint(),
		UpdateClassEndpoint: httptransport.NewClient("PATCH", tgt, encodeUpdateClassRequest, decodeUpdateClassResponse, options...).Endpoint(),
		DeleteClassEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteClassRequest, decodeDeleteClassResponse, options...).Endpoint(),
		JoinClassEndpoint:   httptransport.NewClient("POST", tgt, encodeJoinClassRequest, decodeJoinClassResponse, options...).Endpoint(),
		SetRoleEndpoint:     httptransport.NewClient("PATCH", tgt, encodeSetRoleRequest, decodeSetRoleResponse, options...).Endpoint(),
		LeaveClassEndpoint:  httptransport.NewClient("DELETE", tgt, encodeLeaveClassRequest, decodeLeaveClassResponse, options...).Endpoint(),
		ListMembersEndpoint: httptransport.NewClient("GET", tgt, encodeListMembersRequest, decodeListMembersResponse, options...).Endpoint(),
	}, nil
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

func (r listClassesResponse) error() error {
	return r.Error
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

func (r getClassResponse) error() error {
	return r.Error
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

func (r createClassResponse) error() error {
	return r.Error
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

func (r updateClassResponse) error() error {
	return r.Error
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

func (r deleteClassResponse) error() error {
	return r.Error
}

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

func (r joinClassResponse) error() error {
	return r.Error
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

func (r leaveClassResponse) error() error {
	return r.Error
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

func (r setRoleResponse) error() error {
	return r.Error
}

func MakeListMembersEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(listMembersRequest)
		members, e := s.ListMembers(ctx, req.Class)
		return listMembersResponse{members, e}, nil
	}
}

type listMembersRequest struct {
	Class uuid.UUID `json:"id"`
}

type listMembersResponse struct {
	Members []*models.Member `json:"members"`
	Error   error `json:"error,omitempty"`
}

func (r listMembersResponse) error() error {
	return r.Error
}
