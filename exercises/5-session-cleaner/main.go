// ////////////////////////////////////////////////////////////////////
//
//	//
//
// Given is a SessionManager that stores session information in		// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we		// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will		// keep on adding new sessions to the manager our program will
// eventually run out of memory.		// eventually run out of memory.
//
//	//
//
// Your task is to implement a session cleaner routine that runs		// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that		// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually		// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).		// session times are much longer).
//
//	//
//
// Note that we expect the session to be removed anytime between 5 and		// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be		// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.		// very careful in order to prevent race conditions.
//
//	//
package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

const SessionLifeCycle = time.Second * 5

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	sessions map[string]Session
	ml       sync.Mutex
}

// Session stores the session's data
type Session struct {
	Data         map[string]interface{}
	CreationTime time.Time
}

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
	}
	go m.RemoveExpiredSession()
	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}
	m.ml.Lock()
	defer m.ml.Unlock()
	m.sessions[sessionID] = Session{
		Data:         make(map[string]interface{}),
		CreationTime: time.Now(),
	}
	return sessionID, nil
}
func (m *SessionManager) RemoveExpiredSession() {
	tick := time.Tick(time.Second)
	for {
		<-tick
		m.ml.Lock()
		for sessionID, session := range m.sessions {
			if time.Since(session.CreationTime) >= SessionLifeCycle {
				delete(m.sessions, sessionID)
			}
		}
		m.ml.Unlock()
	}
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.ml.Lock()
	defer m.ml.Unlock()
	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.ml.Lock()
	defer m.ml.Unlock()
	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:         data,
		CreationTime: time.Now(),
	}
	return nil
}
func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Created new session with ID", sID)
	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"
	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(".update session data, set website to longhoang.de")
	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Get session data:", updatedData)
}
