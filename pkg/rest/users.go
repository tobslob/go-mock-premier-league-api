package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"tobslob.com/go-mock-premier-league-api/pkg/config"
	"tobslob.com/go-mock-premier-league-api/pkg/users"
)

func Users(ctx context.Context, r *chi.Mux, app *config.App) error {
	userRepo, err := users.NewRepository(ctx, app.DB)

	if err != nil {
		log.Fatal(err.Error())
	}

	r.Route("/users", func(r chi.Router) {
		r.Post("/", CreateUser(userRepo))
	})
	return nil
}

func CreateUser(userRepo *users.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var dto users.UserDTO
		config.ReadJSON(r, &dto)

		user, err := userRepo.Create(r.Context(), dto)
		if err != nil {
			mailErr, ok := err.(users.ErrEmail)
			if ok {
				panic(config.JSendError{
					Code:    http.StatusConflict,
					Message: "This user already exist.",
					Err:     mailErr,
				})
			}
			panic(err)
		}
		config.SendSuccess(w, r, user)
	}
}
