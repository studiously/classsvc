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

func (e Endpoints) ListClasses(ctx context.Context) ([]uuid.UUID, error) {
	response, err := e.ListClassesEndpoint(ctx, nil)
	if err != nil {
		return nil, err
	}
	resp := response.(listClassesResponse)
	return resp.Classes, resp.Error
}

func (e Endpoints) GetClass(ctx context.Context, classID uuid.UUID) (*models.Class, error) {
	request := getClassRequest{ClassID: classID}
	response, err := e.GetClassEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(getClassResponse)
	return resp.Class, resp.Error
}

func (e Endpoints) CreateClass(ctx context.Context, name string) error {
	request := createClassRequest{Name: name}
	response, err := e.CreateClassEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(createClassResponse)
	return resp.Error
}

func (e Endpoints) UpdateClass(ctx context.Context, classID uuid.UUID, name string, currentUnit uuid.UUID) error {
	request := updateClassRequest{ClassID: classID, Name: name, CurrentUnit: currentUnit}
	response, err := e.UpdateClassEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(updateClassResponse)
	return resp.Error
}

func (e Endpoints) DeleteClass(ctx context.Context, classID uuid.UUID) error {
	request := deleteClassRequest{ClassID: classID}
	response, err := e.DeleteClassEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteClassResponse)
	return resp.Error
}

func (e Endpoints) JoinClass(ctx context.Context, classID uuid.UUID) error {
	request := joinClassRequest{ClassID: classID}
	response, err := e.JoinClassEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(joinClassResponse)
	return resp.Error
}

func (e Endpoints) SetRole(ctx context.Context, userID uuid.UUID, classID uuid.UUID, role models.UserRole) error {
	request := setRoleRequest{UserID: userID, ClassID: classID, Role: role}
	response, err := e.SetRoleEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(setRoleResponse)
	return resp.Error
}

func (e Endpoints) LeaveClass(ctx context.Context, classID uuid.UUID, userID uuid.UUID) error {
	request := leaveClassRequest{ClassID: classID, UserID: userID}
	response, err := e.LeaveClassEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(leaveClassResponse)
	return resp.Error
}

func (e Endpoints) ListMembers(ctx context.Context, classID uuid.UUID) ([]*models.Member, error) {
	request := listMembersRequest{ClassID: classID}
	response, err := e.ListMembersEndpoint(ctx, request)
	if err != nil {
		return nil, err
	}
	resp := response.(listMembersResponse)
	return resp.Members, resp.Error
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
		class, e := s.GetClass(ctx, req.ClassID)
		return getClassResponse{class, e}, nil
	}
}

type getClassRequest struct {
	ClassID uuid.UUID `json:"id,omitempty"`
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
		e := s.UpdateClass(ctx, req.ClassID, req.Name, req.CurrentUnit)
		return updateClassResponse{e}, nil
	}
}

type updateClassRequest struct {
	ClassID     uuid.UUID
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
		e := s.DeleteClass(ctx, req.ClassID)
		return deleteClassResponse{e}, nil
	}
}

type deleteClassRequest struct {
	ClassID uuid.UUID `json:"id"`
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
		e := s.JoinClass(ctx, req.ClassID)
		return joinClassResponse{e}, nil
	}
}

type joinClassRequest struct {
	ClassID uuid.UUID `json:"class"`
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
		e := s.LeaveClass(ctx, req.UserID, req.ClassID)
		return leaveClassResponse{e}, nil
	}
}

type leaveClassRequest struct {
	UserID  uuid.UUID `json:"user,omitempty"`
	ClassID uuid.UUID `json:"class"`
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
		e := s.SetRole(ctx, req.UserID, req.ClassID, req.Role)
		return setRoleResponse{e}, nil
	}
}

type setRoleRequest struct {
	UserID  uuid.UUID `json:"user"`
	ClassID uuid.UUID `json:"class"`
	Role    models.UserRole `json:"role"`
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
		members, e := s.ListMembers(ctx, req.ClassID)
		return listMembersResponse{members, e}, nil
	}
}

type listMembersRequest struct {
	ClassID uuid.UUID `json:"id"`
}

type listMembersResponse struct {
	Members []*models.Member `json:"members"`
	Error   error `json:"error,omitempty"`
}

func (r listMembersResponse) error() error {
	return r.Error
}
