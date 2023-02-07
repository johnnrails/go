package middleware

import (
	"net/http"

	apperrors "github.com/johnnrails/ddd_go/boilerplate/errors"
	httpResponseJson "github.com/johnnrails/ddd_go/boilerplate/http/response/json"
	md "github.com/johnnrails/ddd_go/boilerplate/metadata"
	"github.com/vardius/gorouter/v4"
)

const requestMetadataKey = "m"

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP statusCode to be captured for metadata.
type responseWriter struct {
	http.ResponseWriter
	metadata    *md.Metadata
	wroteHeader bool
}

func (rw *responseWriter) WriteStatusCodeOnHeader(statusCode int) {
	rw.metadata.StatusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
	return
}

func WithMetadata() gorouter.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var metadata *md.Metadata
			if !r.URL.Query().Has(requestMetadataKey) {
				metadata = md.CreateMetadataFromRequest(r)
			} else {
				err := md.GetMetadataFromQuery(metadata, requestMetadataKey, r)
				if err != nil {
					httpResponseJson.JSONError(r.Context(), w, http.StatusInternalServerError, apperrors.NewAppErrorFromError(err))
				}
			}

			ctx := md.ContextWithMetadata(r.Context(), metadata)
			next.ServeHTTP(&responseWriter{
				ResponseWriter: w,
				metadata:       metadata,
			}, r.WithContext(ctx))
		})
	}
}
