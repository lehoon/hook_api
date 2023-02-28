package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/lehoon/hook_api/v2/library/logger"
)

func RequestLoggerFilter(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		buff, _ := io.ReadAll(r.Body)
		fmt.Printf("Request url=%s,method=%s,body=%s\n", r.RequestURI, r.Method, string(buff))
		logger.Log().Infof("Request url=%s,method=%s,body=%s\n", r.RequestURI, r.Method, string(buff))

		//把读出来的数据再写到request.body上
		r.Body = io.NopCloser(bytes.NewBuffer(buff))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
