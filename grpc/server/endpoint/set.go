package endpoint

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/partyzanex/esender/domain"
	"context"
)

type Set struct {
	CreateEndpoint endpoint.Endpoint
	UpdateEndpoint endpoint.Endpoint
	GetEndpoint    endpoint.Endpoint
	SearchEndpoint endpoint.Endpoint
	SendEndpoint   endpoint.Endpoint
}

func (s Set) Create(ctx context.Context, email domain.Email) (*domain.Email, error) {
	res, err := s.CreateEndpoint(ctx, email)
	if err != nil {
		return nil, err
	}

	result := res.(*domain.Email)
	return result, nil
}

func (s Set) Update(ctx context.Context, email domain.Email) (*domain.Email, error) {
	res, err := s.UpdateEndpoint(ctx, email)
	if err != nil {
		return nil, err
	}

	result := res.(*domain.Email)
	return result, nil
}

func (s Set) Get(ctx context.Context, id int64) (*domain.Email, error) {
	res, err := s.GetEndpoint(ctx, id)
	if err != nil {
		return nil, err
	}

	result := res.(*domain.Email)
	return result, nil
}

func (s Set) Search(ctx context.Context, filter *domain.Filter) ([]*domain.Email, error) {
	res, err := s.CreateEndpoint(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := res.([]*domain.Email)
	return result, nil
}

func (s Set) Send(email domain.Email) (bool, error) {
	res, err := s.SendEndpoint(context.Background(), email)
	return res.(bool), err
}

func (s Set) Name() string {
	return "set"
}

func (s Set) AgentConfig() domain.AgentConfig {
	return domain.AgentConfig{}
}

func New(svc domain.EmailService, logger log.Logger) Set {
	var createEndpoint endpoint.Endpoint
	{
		createEndpoint = MakeCreateEndpoint(svc)
		createEndpoint = LoggingMiddleware(log.With(logger, "method", "Create"))(createEndpoint)
	}
	var updateEndpoint endpoint.Endpoint
	{
		updateEndpoint = MakeUpdateEndpoint(svc)
		updateEndpoint = LoggingMiddleware(log.With(logger, "method", "Update"))(updateEndpoint)
	}
	var getEndpoint endpoint.Endpoint
	{
		getEndpoint = MakeGetEndpoint(svc)
		getEndpoint = LoggingMiddleware(log.With(logger, "method", "Get"))(getEndpoint)
	}
	var searchEndpoint endpoint.Endpoint
	{
		searchEndpoint = MakeSearchEndpoint(svc)
		searchEndpoint = LoggingMiddleware(log.With(logger, "method", "Search"))(searchEndpoint)
	}
	var sendEndpoint endpoint.Endpoint
	{
		sendEndpoint = MakeSendEndpoint(svc)
		sendEndpoint = LoggingMiddleware(log.With(logger, "method", "Send"))(sendEndpoint)
	}

	return Set{
		CreateEndpoint: createEndpoint,
		UpdateEndpoint: updateEndpoint,
		GetEndpoint:    getEndpoint,
		SendEndpoint:   sendEndpoint,
		SearchEndpoint: searchEndpoint,
	}
}

func MakeCreateEndpoint(svc domain.EmailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email := request.(domain.Email)
		return svc.Create(ctx, email)
	}
}

func MakeUpdateEndpoint(svc domain.EmailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email := request.(domain.Email)
		return svc.Update(ctx, email)
	}
}

func MakeGetEndpoint(svc domain.EmailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		id := request.(int64)
		return svc.Get(ctx, id)
	}
}

func MakeSearchEndpoint(svc domain.EmailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var filter *domain.Filter
		if request != nil {
			filter = request.(*domain.Filter)
		}

		return svc.Search(ctx, filter)
	}
}

func MakeSendEndpoint(svc domain.EmailService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		email := request.(domain.Email)
		return svc.Send(email)
	}
}
