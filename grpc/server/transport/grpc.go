package transport

import (
	"github.com/go-kit/kit/transport/grpc"
	"github.com/partyzanex/esender/grpc/server/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/partyzanex/esender/grpc/server/pb"
	"github.com/go-kit/kit/transport"
	"context"
	"fmt"
	"github.com/partyzanex/esender/domain"
	"time"
)

type grpcServer struct {
	create grpc.Handler
	update grpc.Handler
	get    grpc.Handler
	search grpc.Handler
	send   grpc.Handler
}

func (s *grpcServer) Create(ctx context.Context, email *pb.Message) (*pb.Message, error) {
	_, rep, err := s.create.ServeGRPC(ctx, email)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.Message), nil
}

func (s *grpcServer) Update(ctx context.Context, email *pb.Message) (*pb.Message, error) {
	_, rep, err := s.update.ServeGRPC(ctx, email)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.Message), nil
}

func (s *grpcServer) Get(ctx context.Context, email *pb.Message) (*pb.Message, error) {
	if email == nil || email.ID == 0 {
		return nil, fmt.Errorf("required email id")
	}

	_, rep, err := s.get.ServeGRPC(ctx, email.ID)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.Message), nil
}

func (s *grpcServer) Send(ctx context.Context, email *pb.Message) (*pb.SendResult, error) {
	_, rep, err := s.send.ServeGRPC(ctx, email)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.SendResult), nil
}

func (s *grpcServer) Search(ctx context.Context, filter *pb.SearchRequest) (*pb.SearchResponse, error) {
	_, rep, err := s.search.ServeGRPC(ctx, filter)
	if err != nil {
		return nil, err
	}

	return rep.(*pb.SearchResponse), nil
}

func NewGRPCServer(endpoints endpoint.Set, logger log.Logger) pb.EmailServer {
	options := []grpc.ServerOption{
		grpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		create: grpc.NewServer(
			endpoints.CreateEndpoint,
			decodeEmail,
			encodeEmail,
			options...
		),
		update: grpc.NewServer(
			endpoints.UpdateEndpoint,
			decodeEmail,
			encodeEmail,
			options...
		),
		get: grpc.NewServer(
			endpoints.GetEndpoint,
			decodeEmail,
			encodeEmail,
			options...
		),
		send: grpc.NewServer(
			endpoints.SendEndpoint,
			decodeEmail,
			encodeSendResponse,
			options...
		),
		search: grpc.NewServer(
			endpoints.GetEndpoint,
			decodeSearchRequest,
			encodeSearchResponse,
			options...
		),
	}
}

func encodeEmail(_ context.Context, response interface{}) (interface{}, error) {
	var email *domain.Email
	switch response.(type) {
	case domain.Email:
		resp := response.(domain.Email)
		email = &resp
	case *domain.Email:
		email = response.(*domain.Email)
	}

	var DTUpdated, DTSent time.Time
	{
		if email.DTUpdated != nil {
			DTUpdated = *email.DTUpdated
		}
		if email.DTSent != nil {
			DTSent = *email.DTSent
		}
	}

	return &pb.Message{
		ID:         email.ID,
		Recipients: encodeAddresses(email.Recipients...),
		BCC:        encodeAddresses(email.BCC...),
		CC:         encodeAddresses(email.CC...),
		Sender:     encodeAddress(email.Sender),
		Subject:    email.Subject,
		Body:       email.Body,
		MimeType:   pb.MimeType(pb.MimeType_value[email.MimeType.String()]),
		Status:     pb.EmailStatus(pb.EmailStatus_value[email.Status.String()]),
		Error:      *email.Error,
		DTCreated:  email.DTCreated.Unix(),
		DTUpdated:  DTUpdated.Unix(),
		DTSent:     DTSent.Unix(),
	}, nil
}

func encodeAddresses(addresses ...domain.Address) []*pb.Address {
	results := make([]*pb.Address, len(addresses))

	for i, address := range addresses {
		address := address
		results[i] = encodeAddress(address)
	}

	return results
}

func encodeAddress(address domain.Address) *pb.Address {
	return &pb.Address{
		Address: address.Address,
		Name:    address.Name,
	}
}

func decodeEmail(_ context.Context, grpcReply interface{}) (interface{}, error) {
	email := grpcReply.(*pb.Message)

	var DTUpdated, DTSent *time.Time
	{
		if email.DTUpdated > 0 {
			*DTUpdated = time.Unix(email.DTUpdated, 0)
		}
		if email.DTSent > 0 {
			*DTSent = time.Unix(email.DTSent, 0)
		}
	}

	return &domain.Email{
		ID:         email.ID,
		Recipients: decodeAddresses(email.Recipients...),
		BCC:        decodeAddresses(email.BCC...),
		CC:         decodeAddresses(email.CC...),
		Sender:     decodeAddress(email.Sender),
		Subject:    email.Subject,
		Body:       email.Body,
		MimeType:   domain.MimeTypeAlias(pb.MimeType_name[int32(email.MimeType)]),
		Status:     domain.EmailStatus(pb.EmailStatus_name[int32(email.Status)]),
		Error:      &email.Error,
		DTCreated:  time.Unix(email.DTCreated, 0),
		DTUpdated:  DTUpdated,
		DTSent:     DTSent,
	}, nil
}

func decodeAddresses(addresses ...*pb.Address) []domain.Address {
	results := make([]domain.Address, len(addresses))

	for i, address := range addresses {
		results[i] = decodeAddress(address)
	}

	return results
}

func decodeAddress(address *pb.Address) domain.Address {
	return domain.Address{
		Address: address.Address,
		Name:    address.Name,
	}
}

func encodeSendResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(bool)

	return &pb.SendResult{
		Success: resp,
	}, nil
}

func encodeSearchResponse(ctx context.Context, response interface{}) (interface{}, error) {
	resp := response.([]*domain.Email)

	emails := make([]*pb.Message, len(resp))

	for i, email := range resp {
		res, err := encodeEmail(ctx, email)
		if err != nil {
			return nil, err
		}

		emails[i] = res.(*pb.Message)
	}

	return &pb.SearchResponse{
		Emails: emails,
	}, nil
}

func decodeSearchRequest(ctx context.Context, grpcReply interface{}) (interface{}, error) {
	var filter *domain.Filter
	{
		if grpcReply != nil {
			reply := grpcReply.(*pb.SearchRequest)

			var timeRange domain.TimeRange
			if reply.Till > 0 && reply.Since > 0 {
				timeRange.SetTill(time.Unix(reply.Till, 0))
				timeRange.SetSince(time.Unix(reply.Since, 0))
			}

			filter = &domain.Filter{
				Recipient: reply.Recipient,
				Sender:    reply.Sender,
				Status:    domain.EmailStatus(pb.EmailStatus_name[int32(reply.Status)]),
				TimeRange: timeRange,
				Limit:     int(reply.Limit),
				Offset:    int(reply.Offset),
			}
		}
	}

	return filter, nil
}
