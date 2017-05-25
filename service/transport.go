package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ory/hydra/sdk"
	"github.com/studiously/introspector"
)

var (
	ErrBadRequest = errors.New("the request is malformed or invalid")
)

func MakeHTTPHandler(s Service, logger log.Logger, client *sdk.Client) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerBefore(introspector.ToHTTPContext()),
	}

	// GET /classes/
	// Get a list of classes the user has access to.
	r.Methods("GET").Path("/classes/{class}").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.get")(e.GetClassEndpoint),
		decodeGetClassRequest,
		encodeResponse,
		append(options, httptransport.ServerBefore(introspector.ToHTTPContext()))...
	))

	r.Methods("POST").Path("/classes/").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.new")(e.CreateClassEndpoint),
		decodeCreateClassRequest,
		encodeResponse,
		options...
	))

	r.Methods("GET").Path("/classes/").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.list")(e.ListClassesEndpoint),
		decodeListClassesRequest,
		encodeResponse,
		options...
	))

	r.Methods("PATCH").Path("/classes/{class}").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.update")(e.UpdateClassEndpoint),
		decodeUpdateClassRequest,
		encodeResponse,
		options...
	))

	r.Methods("DELETE").Path("/classes/{class}").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.delete")(e.DeleteClassEndpoint),
		decodeDeleteClassRequest,
		encodeResponse,
		options...
	))

	r.Methods("GET").Path("/classes/{class}/members").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.list_members")(e.ListMembersEndpoint),
		decodeListMembersRequest,
		encodeResponse,
		options...
	))

	r.Methods("POST").Path("/classes/{class}/join").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.join")(e.JoinClassEndpoint),
		decodeJoinClassRequest,
		encodeResponse,
		options...
	))

	leaveClassServer := httptransport.NewServer(
		introspector.New(client.Introspection, "classes.leave")(e.LeaveClassEndpoint),
		decodeLeaveClassRequest,
		encodeResponse,
		options...
	)
	r.Methods("DELETE").Path("/classes/{class}/leave").Handler(leaveClassServer)
	r.Methods("DELETE").Path("/classes/{class}/leave/{user}").Handler(leaveClassServer)

	r.Methods("PATCH").Path("/classes/{class}/members/{user}").Handler(httptransport.NewServer(
		introspector.New(client.Introspection, "classes.members:update")(e.SetRoleEndpoint),
		decodeSetRoleRequest,
		encodeResponse,
		options...
	))

	return r
}

func encodeGetClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(getClassRequest)
	classID := url.QueryEscape(r.Id.String())
	req.Method, req.URL.Path = "GET", "/classes/"+classID
	return encodeRequest(ctx, req, request)
}

func decodeGetClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return getClassRequest{}, ErrBadRequest
	}
	return getClassRequest{class}, nil
}

func encodeCreateClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.Method, req.URL.Path = "POST", "/classes/"
	return encodeRequest(ctx, req, request)
}

func decodeCreateClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response createClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeCreateClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req createClassRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeListClassesRequest(ctx context.Context, req *http.Request, request interface{}) error {
	req.Method, req.URL.Path = "GET", "/classes/"
	return encodeRequest(ctx, req, request)
}

func decodeListClassesResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response listClassesResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeListClassesRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return nil, nil
}

func encodeUpdateClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(updateClassRequest)
	classID := url.QueryEscape(r.Class.String())
	req.Method, req.URL.Path = "PATCH", "/classes/"+classID
	return encodeRequest(ctx, req, request)
}

func decodeUpdateClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response updateClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeUpdateClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req updateClassRequest
	if e := json.NewDecoder(r.Body).Decode(&req); e != nil {
		return nil, e
	}
	return req, nil
}

func encodeDeleteClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(deleteClassRequest)
	classID := url.QueryEscape(r.Id.String())
	req.Method, req.URL.Path = "DELETE", "/classes/"+classID
	return encodeRequest(ctx, req, request)
}

func decodeDeleteClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return deleteClassRequest{}, ErrBadRequest
	}
	return deleteClassRequest{class}, nil
}

func encodeListMembersRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(listMembersRequest)
	classID := url.QueryEscape(r.Class.String())
	req.Method, req.URL.Path = "GET", "/classes/"+classID+"/members"
	return encodeRequest(ctx, req, request)
}

func decodeListMembersResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response listMembersResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeListMembersRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return listMembersRequest{}, ErrBadRequest
	}
	return listMembersRequest{class}, nil
}

func encodeJoinClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(joinClassRequest)
	classID := url.QueryEscape(r.Class.String())
	req.Method, req.URL.Path = "POST", "/classes/"+classID+"/join"
	return encodeRequest(ctx, req, request)
}

func decodeJoinClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response joinClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeJoinClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return joinClassRequest{}, ErrBadRequest
	}
	return joinClassRequest{
		Class: class,
	}, nil
}

func encodeLeaveClassRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(leaveClassRequest)
	classID := url.QueryEscape(r.Class.String())
	req.Method = "DELETE"
	if r.User == uuid.Nil {
		req.URL.Path = "/classes/" + classID + "/leave"
	} else {
		userID := url.QueryEscape(r.User.String())
		req.URL.Path = "/classes/" + classID + "/leave/" + userID
	}
	return encodeRequest(ctx, req, request)
}

func decodeLeaveClassResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response leaveClassResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeLeaveClassRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req leaveClassRequest
	vars := mux.Vars(r)
	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return leaveClassRequest{}, ErrBadRequest
	}
	req.Class = class
	userS, ok := vars["user"]
	user := uuid.Nil
	if ok {
		user, err = uuid.Parse(userS)
		if err != nil {
			return nil, ErrBadRequest
		}
	}
	req.User = user
	return req, nil
}

func encodeSetRoleRequest(ctx context.Context, req *http.Request, request interface{}) error {
	r := request.(setRoleRequest)
	classID := url.QueryEscape(r.Class.String())
	userID := url.QueryEscape(r.User.String())
	req.Method, req.URL.Path = "PATCH", "/classes/"+classID+"/members/"+userID
	return encodeRequest(ctx, req, request)
}

func decodeSetRoleResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response setRoleResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeSetRoleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req setRoleRequest
	vars := mux.Vars(r)

	class, err := uuid.Parse(vars["class"])
	if err != nil {
		return setRoleRequest{}, ErrBadRequest
	}
	req.Class = class
	user, err := uuid.Parse(vars["user"])
	if err != nil {
		return nil, ErrBadRequest
	}
	req.User = user
	if e := json.NewDecoder(r.Body).Decode(&req.Role); e != nil {
		return nil, e
	}
	return req, nil
}

// errorer is implemented by all concrete response types that may contain
// errors. It allows us to change the HTTP response code without needing to
// trigger an endpoint (transport-level) error. For more information, read the
// big comment in endpoints.go.
type errorer interface {
	error() error
}

// encodeResponse is the common method to encode all response types to the
// client. I chose to do it this way because, since we're using JSON, there's no
// reason to provide anything more specific. It's certainly possible to
// specialize on a per-response (per-method) basis.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// encodeRequest likewise JSON-encodes the request to the HTTP request body.
// Don't use it directly as a transport/http.Client EncodeRequestFunc:
// profilesvc endpoints require mutating the HTTP method and request path.
func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
