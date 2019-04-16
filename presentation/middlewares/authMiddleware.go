package middlewares

import (
	"database/sql"
	"log"
	"net/http"

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

		query := `SELECT session_id FROM sessions WHERE session_id=$1`

		row := db.QueryRow(query, cookie.String())

		var session models.Session

		if err := row.Scan(&session.ID); err != nil {
			log.Println("can't scan from DB" + err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized"))
			return
		}

		next.ServeHTTP(w, req)
	})
}
