package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"strconv"
	"time"

	"gohst/internal/config"

	"github.com/redis/go-redis/v9"
)

const REDIS_DB_DEFAULT = 0

const REDIS_HOST_DEFAULT = "localhost"

const REDIS_PASSWORD_DEFAULT = ""

const REDIS_PORT_DEFAULT = 6379


// RedisSessionManager handles Redis-based sessions
type RedisSessionManager struct {
	redisClient *redis.Client
}

// NewRedisSessionManager initializes Redis connection
func NewRedisSessionManager() (*RedisSessionManager, string) {
	redisConf := config.Session.Redis
	password := REDIS_PASSWORD_DEFAULT
	db := REDIS_DB_DEFAULT

	if redisConf != nil {
		if redisConf.Host != "" {
			password = redisConf.Password
		}
		if redisConf.Port != 0 {
			db = redisConf.DB
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     GetRedisHostAddr(),
		Password: password,
		DB:       db,
	})

	return &RedisSessionManager{redisClient: client}, "redis"
}

// GetRedisHostAddr returns the Redis host address
func GetRedisHostAddr() string {
	redis := config.Session.Redis
	host := REDIS_HOST_DEFAULT
	port := REDIS_PORT_DEFAULT

	if redis != nil {
		if redis.Host != "" {
			host = redis.Host
		}
		if redis.Port != 0 {
			port = redis.Port
		}
	}

	return host + ":" + strconv.Itoa(port)
}

// StartSession creates a session in Redis using Gob
func (rsm *RedisSessionManager) StartSession(w http.ResponseWriter, r *http.Request) (*SessionData, string) {
	sessionID := GenerateSessionID()
	ctx := context.Background()

	sessionData := &SessionData{
		ID:      sessionID,
		Values:  make(map[string]interface{}),
		Expires: time.Now().Add(30 * time.Minute),
		manager: rsm,
	}

	// Encode session to Gob format
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(sessionData)
	if err != nil {
		log.Println("Error encoding session:", err)
		return &SessionData{}, ""
	}

	// Store Gob-encoded session in Redis (expires in 30 minutes)
	sessionLength := GetSessionLength()
	err = rsm.redisClient.Set(ctx, sessionID, buf.Bytes(), sessionLength).Err()
	if err != nil {
		log.Println("Error storing session in Redis:", err)
	}

	// Set session ID in cookie
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_NAME,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
	})

	return sessionData, sessionID
}

// GetSession retrieves session from Redis and decodes Gob
func (rsm *RedisSessionManager) GetSession(r *http.Request) (*SessionData, string) {
	cookie, err := r.Cookie(SESSION_NAME)
	if err != nil {
		return nil, ""
	}

	// If session ID is not in the cookie, return nil
	ctx := context.Background()
	val, err := rsm.redisClient.Get(ctx, cookie.Value).Bytes()
	if err != nil {
		return nil, ""
	}

	// Decode Gob data
	var sessionData SessionData
	decoder := gob.NewDecoder(bytes.NewReader(val))
	err = decoder.Decode(&sessionData)
	if err != nil {
		log.Println("Error decoding session:", err)
		return nil, ""
	}

	return &sessionData, cookie.Value
}

// SetValue stores a value in the session and encodes it with Gob
func (rsm *RedisSessionManager) SetValue(sessionID string, key string, value interface{}) {
	ctx := context.Background()

	// Get current session data
	sessionData, err := rsm.GetSessionByID(ctx, sessionID)
	if err != nil {
		log.Println("Session not found:", sessionID)
		return
	}

	// Update session values
	sessionData.Values[key] = value

	// Encode session data using Gob
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err = encoder.Encode(sessionData)
	if err != nil {
		log.Println("Error encoding session data:", err)
		return
	}

	// Save the updated session back to Redis
	err = rsm.redisClient.Set(ctx, sessionID, buf.Bytes(), 30*time.Minute).Err()
	if err != nil {
		log.Println("Error storing updated session:", err)
	}
}

// GetValue retrieves a session value
func (rsm *RedisSessionManager) GetValue(sessionID string, key string) (interface{}, bool) {
	sessionData, err := rsm.GetSessionByID(context.Background(), sessionID)
	if err != nil {
		return nil, false
	}
	val, ok := sessionData.Values[key]
	return val, ok
}

// GetSessionByID fetches session data directly using session ID
func (rsm *RedisSessionManager) GetSessionByID(ctx context.Context, sessionID string) (*SessionData, error) {
	val, err := rsm.redisClient.Get(ctx, sessionID).Bytes() // Get session as bytes
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	decoder := gob.NewDecoder(bytes.NewReader(val))
	err = decoder.Decode(&sessionData)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

// RemoveValue deletes a key from the session data in Redis
func (rsm *RedisSessionManager) Remove(sessionID string, key string) error {
    ctx := context.Background()

    // Get current session data
    sessionData, err := rsm.GetSessionByID(ctx, sessionID)
    if err != nil {
        return err
    }

    delete(sessionData.Values, key)

    var buf bytes.Buffer
    encoder := gob.NewEncoder(&buf)
    err = encoder.Encode(sessionData)
    if err != nil {
        return err
    }

    // Save the updated session back to Redis with the same expiration
    ttl, err := rsm.redisClient.TTL(ctx, sessionID).Result()
    if err != nil {
        ttl = 30 * time.Minute // Default if TTL can't be retrieved
    }

    // Save the updated session back to Redis
    err = rsm.redisClient.Set(ctx, sessionID, buf.Bytes(), ttl).Err()
    if err != nil {
        return err
    }

    return nil
}

// Save saves the entire session
func (rsm *RedisSessionManager) Save(sessionID string, session *SessionData) error {
    ctx := context.Background()

    // Encode session data using Gob
    var buf bytes.Buffer
    encoder := gob.NewEncoder(&buf)
    err := encoder.Encode(session)
    if err != nil {
        return err
    }

    // Get current TTL if possible
    ttl, err := rsm.redisClient.TTL(ctx, sessionID).Result()
    if err != nil || ttl < 0 {
        ttl = GetSessionLength() // Default if TTL can't be retrieved
    }

    // Save to Redis with the same expiration
    return rsm.redisClient.Set(ctx, sessionID, buf.Bytes(), ttl).Err()
}

// Delete removes the entire session
func (rsm *RedisSessionManager) Delete(sessionID string) error {
    ctx := context.Background()
    return rsm.redisClient.Del(ctx, sessionID).Err()
}
