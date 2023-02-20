package identity

import (
	"context"

	"github.com/google/uuid"
)

type Identity struct {
	Token        string     `json:"token"`
	Permission   Permission `json:"permission"`
	UserID       uuid.UUID  `json:"user_id"`
	ClientID     uuid.UUID  `json:"client_id,omitempty"`
	ClientDomain string     `json:"client_domain,omitempty"`
}

type Permission uint8

func (p Permission) Add(flag Permission) Permission    { return p | flag }
func (p Permission) Remove(flag Permission) Permission { return p &^ flag }
func (p Permission) Has(flag Permission) bool          { return p&flag != 0 }

const (
	PermissionUserRead Permission = 1 << iota
	PermissionUserWrite
	PermissionClientWrite
	PermissionClientRead
	PermissionTokenRead
)

type key struct{}

func ContextWithIdentity(ctx context.Context, i *Identity) context.Context {
	if ctx == nil {
		return nil
	}
	if i == nil {
		return ctx
	}
	return context.WithValue(ctx, key{}, i)
}

func FromContext(ctx context.Context) (*Identity, bool) {
	if ctx == nil {
		return nil, false
	}
	i, ok := ctx.Value(key{}).(*Identity)
	return i, ok
}
