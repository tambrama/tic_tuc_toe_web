package middleware

import "net/http"
//необходим, когда ваш фронтенд и бэкенд работают на разных портах/доменах
func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return 
		}
		next(w, r)
	}
}
//автоматически устанавливает заголовок для всех ответов
func ContentTypeJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next(w, r)
	}
}
//позволяет комбинировать несколько middleware в цепочку, 
// которая выполняется последовательно.
func Chain(handler http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range middlewares{
		handler = middleware(handler)
	}
	return handler
}

//Запрос → EnableCORS → ContentTypeJSON → Logging → Handler
//Ответ  ← EnableCORS ← ContentTypeJSON ← Logging ← Handler