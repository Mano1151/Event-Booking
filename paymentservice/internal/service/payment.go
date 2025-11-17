package service

import (
	"context"
	pb "paymentservice/api/paymentservice/v1"
	"paymentservice/internal/biz"
)

type PaymentService struct {
	pb.UnimplementedPaymentServiceServer
	uc *biz.PaymentUsecase
}

func NewPaymentService(uc *biz.PaymentUsecase) *PaymentService {
	return &PaymentService{uc: uc}
}

func (s *PaymentService) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentReply, error) {
	payment, err := s.uc.ProcessPayment(ctx, req.BookingId, req.PaymentMethod)
	if err != nil {
		return nil, err
	}

	return &pb.CreatePaymentReply{
		PaymentId: payment.ID,
		BookingId: payment.BookingID,
		Amount:    float32(payment.Amount),
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
