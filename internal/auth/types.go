package auth

import "time"

type AuthData struct {
    UserID     uint64
    Email      string
    Name       string
    IsAdmin    bool
    LoggedInAt time.Time
}
