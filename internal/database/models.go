// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package database

import (
	"net/netip"

	"github.com/jackc/pgx/v5/pgtype"
)

type RefreshToken struct {
	ID        int32
	UserID    pgtype.UUID
	TokenHash string
	CreatedAt pgtype.Timestamp
}

type User struct {
	ID           pgtype.UUID
	Email        string
	PasswordHash string
	IpAddress    *netip.Addr
}
