package ctxutils

import "context"

const requestID = "request_id"

func RequestID(ctx context.Context) string {
	return ctx.Value(requestID).(string)
}

func SetRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, requestID, reqID)
}
