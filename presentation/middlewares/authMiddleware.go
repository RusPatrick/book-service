package middlewares

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/ruspatrick/book-service/domain/models"
)

func AuthMiddleware(next http.Handler, db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie("session_id")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		query := `SELECT cookie, exp FROM sessions WHERE cookie=$1`
		row := db.QueryRow(query, cookie.String())
		var session models.Session
		if err := row.Scan(&session.ID, &session.Exp); err != nil {
			log.Println("can't scan from DB" + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		if session.Exp < time.Now().Unix() {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		next.ServeHTTP(w, req)
	})
}
