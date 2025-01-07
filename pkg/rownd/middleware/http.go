package rowndmiddleware

import (
	"errors"
	"net/http"

	"github.com/rownd/client-go/pkg/rownd"
)

// TODO add error handler
func WithAuthentication(handler Handler) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := handler.TokenExtractor(r)
			if err != nil {
				handler.ErrorHandler(w, r, errors.New("Forbidden"))
				return
			}

			ctx := r.Context()
			// Pass validation options to Validate
			validated, err := handler.Validator.Validate(ctx, token, handler.ValidationOpts)
			if err != nil {
				handler.ErrorHandler(w, r, errors.New("Forbidden"))
				return
			}
			// embed validated token into context.
			next.ServeHTTP(w, r.WithContext(rownd.AddTokenToCtx(ctx, validated))) // Remove * since validated is already a *rownd.Token
		})
	}
}
