package stripe

import (
	"github.com/BTBurke/recur/pb"
	"github.com/stripe/stripe-go"
)

// convert an error response from Stripe to a pb.Error
func respToError(err *stripe.Error) *pb.Error {
	return &pb.Error{
		Type:           convertErrorType(err.Type),
		ChargeId:       err.ChargeID,
		Message:        err.Msg,
		HttpStatusCode: int32(err.HTTPStatusCode),
		Code:           convertCardError(err.Code),
		Param:          err.Param,
		RequestId:      err.RequestID,
	}
}

// map from Stripe error types to proto
func convertErrorType(t stripe.ErrorType) pb.ErrorType {
	lookup := map[stripe.ErrorType]pb.ErrorType{
		stripe.ErrorTypeAPI:            pb.ErrorType_API,
		stripe.ErrorTypeAPIConnection:  pb.ErrorType_APIConnection,
		stripe.ErrorTypeAuthentication: pb.ErrorType_Authentication,
		stripe.ErrorTypeCard:           pb.ErrorType_Card,
		stripe.ErrorTypeInvalidRequest: pb.ErrorType_InvalidRequest,
		stripe.ErrorTypePermission:     pb.ErrorType_Permission,
		stripe.ErrorTypeRateLimit:      pb.ErrorType_RateLimit,
	}
	return lookup[t]
}

// map from Stripe card errors to proto
func convertCardError(t stripe.ErrorCode) pb.CardErrors {
	lookup := map[stripe.ErrorCode]pb.CardErrors{
		stripe.IncorrectNum:  pb.CardErrors_IncorrectNumber,
		stripe.InvalidNum:    pb.CardErrors_InvalidNumber,
		stripe.InvalidExpM:   pb.CardErrors_InvalidExpirationMonth,
		stripe.InvalidExpY:   pb.CardErrors_InvalidExpirationYear,
		stripe.InvalidCvc:    pb.CardErrors_InvalidCvc,
		stripe.ExpiredCard:   pb.CardErrors_Expired,
		stripe.IncorrectCvc:  pb.CardErrors_IncorrectCvc,
		stripe.IncorrectZip:  pb.CardErrors_IncorrectZip,
		stripe.CardDeclined:  pb.CardErrors_Declined,
		stripe.Missing:       pb.CardErrors_Missing,
		stripe.ProcessingErr: pb.CardErrors_ProcessingError,
		stripe.RateLimit:     pb.CardErrors_RateLimited,
	}
	return lookup[t]
}
