package biz

import "github.com/google/wire"

// ProviderSet for PaymentUsecase
var ProviderSet = wire.NewSet(NewPaymentUsecase)
