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
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		query := `SELECT session_id, exp FROM books.sessions WHERE session_id=$1`
		row := db.QueryRow(query, cookie.Value)
		var session models.Session
		if err := row.Scan(&session.ID, &session.Exp); err != nil {
			log.Println("can't scan from DB" + err.Error())
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if session.Exp < time.Now().Unix() {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, req)
	})
}
