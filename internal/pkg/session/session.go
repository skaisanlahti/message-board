package session

import (
	"context"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/skaisanlahti/message-board/internal/pkg/middleware"
)

type Options struct {
	CookieName      string
	SessionDuration time.Duration
}

type Session struct {
	SessionID string
	UserID    int
	Expires   time.Time
}

func new(userID int, expires time.Time) Session {
	sessionID := uuid.New().String()
	sessionID = base64.URLEncoding.EncodeToString([]byte(sessionID))
	session := Session{
		SessionID: sessionID,
		UserID:    userID,
		Expires:   expires,
	}

	return session
}

type Store struct {
	sync.Mutex
	sessionsBySessionID map[string]Session
}

func NewStore() *Store {
	return &Store{
		sessionsBySessionID: make(map[string]Session),
	}
}

func (store *Store) Get(sessionID string) (Session, bool) {
	store.Lock()
	defer store.Unlock()

	session, ok := store.sessionsBySessionID[sessionID]
	if !ok {
		return session, false
	}

	if session.Expires.Before(time.Now()) {
		delete(store.sessionsBySessionID, sessionID)
		return session, false
	}

	return session, true
}

func (store *Store) Set(session Session) {
	store.Lock()
	defer store.Unlock()

	store.sessionsBySessionID[session.SessionID] = session
}

func (store *Store) Clear(sessionID string) {
	store.Lock()
	defer store.Unlock()

	delete(store.sessionsBySessionID, sessionID)
}

func (store *Store) ClearExpired() {
	store.Lock()
	defer store.Unlock()

	for _, session := range store.sessionsBySessionID {
		if session.Expires.Before(time.Now()) {
			delete(store.sessionsBySessionID, session.SessionID)
		}
	}
}

type Service struct {
	store   *Store
	options Options
}

func NewService(options Options) *Service {
	return &Service{
		store:   NewStore(),
		options: options,
	}
}

func (service *Service) Start(
	userID int,
	response http.ResponseWriter,
) {
	session := new(userID, time.Now().Add(service.options.SessionDuration))
	service.store.Set(session)
	cookie := &http.Cookie{
		Name:     service.options.CookieName,
		Value:    session.SessionID,
		HttpOnly: true,
		Expires:  session.Expires,
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(response, cookie)
}

func (service *Service) Stop(
	response http.ResponseWriter,
	request *http.Request,
) {
	cookie, err := request.Cookie(service.options.CookieName)
	if err != nil {
		return
	}

	service.store.Clear(cookie.Value)
	cookie = &http.Cookie{
		Name:     service.options.CookieName,
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now(),
		Path:     "/",
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(response, cookie)
}

func GetUserFromContext(request *http.Request) (int, bool) {
	userID, ok := request.Context().Value("userID").(int)
	return userID, ok
}

func AddUserToContext(service *Service) middleware.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			cookie, err := request.Cookie(service.options.CookieName)
			if err != nil {
				handler.ServeHTTP(response, request)
				return
			}

			session, ok := service.store.Get(cookie.Value)
			if !ok {
				handler.ServeHTTP(response, request)
				return
			}

			ctx := context.WithValue(
				request.Context(),
				"userID",
				session.UserID,
			)

			handler.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}

func RequireUser() middleware.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			_, ok := GetUserFromContext(request)
			if !ok {
				response.WriteHeader(http.StatusUnauthorized)
				return
			}

			handler.ServeHTTP(response, request)
		})
	}
}