package interceptors

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"grpc-auth/internal/core/services"
)

var (
	invariantViolationError *services.InvariantViolationError
)

func ErrorHandlingAndLogging(logger *zap.SugaredLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, next grpc.UnaryHandler) (any, error) {
		method := info.FullMethod
		requestUuid := uuid.New().String()

		logger.Infow("begin", "requestUuid", requestUuid, "method", method, "param", req)

		ret, err := next(ctx, req)
		if err != nil {
			var st *status.Status
			if errors.As(err, &invariantViolationError) {
				st = status.New(codes.InvalidArgument, err.Error())

				logger.Infow("end", "requestUuid", requestUuid, "errorCode", st.Code(), "errorMessage", st.Message())
			} else {
				st = status.New(codes.Internal, fmt.Sprintf("Request UUID: %s. Please send this message to technical support.", requestUuid))

				logger.Errorw("end", "requestUuid", requestUuid, "errorCode", st.Code(), "errorMessage", st.Message(), "internalErrorDetail", err)
			}

			return nil, st.Err()
		} else {
			logger.Infow("end", "requestUuid", requestUuid, "return", ret)

			return ret, nil
		}
	}
}
