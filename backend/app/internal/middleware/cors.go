package middleware

import (
	"net/http"
	"os"
	"strings"
)

// CORS libera as origens definidas em ALLOWED_ORIGINS (separadas por vírgula).
// Se a env estiver vazia (dev/local), libera qualquer origem ("*") para não
// atrapalhar o desenvolvimento. Em produção, defina o domínio real do frontend
// (ex: https://dosecerta.tech,https://www.dosecerta.tech).
func CORS(next http.Handler) http.Handler {
	allowed := map[string]bool{}
	for _, o := range strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",") {
		if o = strings.TrimSpace(o); o != "" {
			allowed[o] = true
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		switch {
		case len(allowed) == 0:
			w.Header().Set("Access-Control-Allow-Origin", "*")
		case allowed[origin]:
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
