package metadata

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/johnnrails/ddd_go/boilerplate/reverseproxy"
)

const metadataKey int = 1

type Metadata struct {
	Now        time.Time `json:"-"`
	TraceID    string    `json:"trace_id,omitempty"`
	IPAddress  net.IP    `json:"ip_address,omitempty"`
	StatusCode int       `json:"http_status,omitempty"`
	UserAgent  string    `json:"http_user_agent,omitempty"`
	RemoteAddr string    `json:"http_remote_addr,omitempty"`
	Referer    string    `json:"http_referer,omitempty"`
	Err        error     `json:"-"`
}

func New() *Metadata {
	return &Metadata{
		TraceID: uuid.New().String(),
		Now:     time.Now(),
	}
}

func GetMetadataFromQuery(metadata *Metadata, requestMetadataKey string, r *http.Request) error {
	m := r.URL.Query().Get(requestMetadataKey)
	data, err := base64.RawURLEncoding.DecodeString(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return err
	}
	return nil
}

func CreateMetadataFromRequest(r *http.Request) *Metadata {
	metadata := New()
	metadata.RemoteAddr = r.RemoteAddr
	metadata.UserAgent = r.UserAgent()
	metadata.Referer = r.Referer()
	metadata.StatusCode = http.StatusOK
	if ip, err := reverseproxy.GetIpAddress(r); err == nil {
		metadata.IPAddress = ip
	}
	return metadata
}

func ContextWithMetadata(ctx context.Context, m *Metadata) context.Context {
	if ctx == nil {
		return nil
	}
	if m == nil {
		return ctx
	}
	return context.WithValue(ctx, metadataKey, m)
}

func FromContext(ctx context.Context) *Metadata {
	if ctx == nil {
		return nil
	}
	m := ctx.Value(metadataKey).(*Metadata)
	return m
}
