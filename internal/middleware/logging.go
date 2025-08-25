package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Yash-Watchguard/Tasknest/internal/logger"
)

func LoggingMiddleWare(next http.Handler)http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
		start:= time.Now()

	    logger.Info(fmt.Sprintf("Startted %v %v from %v",r.Method,r.URL.Path,r.RemoteAddr))

		next.ServeHTTP(w,r)

		logger.Info(fmt.Sprintf("Completed %v %v in %v",r.Method, r.URL.Path, time.Since(start)))
	})
	
}