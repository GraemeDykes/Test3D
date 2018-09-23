package main

import (
	// "encoding/json"
	// "errors"
	"fmt"
	// "log"
	// "math/rand"
	// "net"
	// "nowbridge/binEncDec"
	// "nowbridge/nowbridgeComms"
	// "strconv"
	"net/http"
	"sync"
	"time"
)

type WebUser struct {
	userId           int64
	userName         string
	lastMessageInUtc time.Time
	sessionToken     string
}

type SessionList struct {
	sync.Mutex
	webUsers map[string]*WebUser
}

func (sessionList *SessionList) init() {
	sessionList.webUsers = make(map[string]*WebUser, 1000)
}

func (sessionList *SessionList) successfulLogin(userId int64, userName string, sessionToken string) {

	newWebUser := new(WebUser)
	newWebUser.userId = userId
	newWebUser.userName = userName
	newWebUser.lastMessageInUtc = time.Now().UTC()
	newWebUser.sessionToken = sessionToken

	sessionList.Lock()

	sessionList.webUsers[sessionToken] = newWebUser

	sessionList.Unlock()
}

func (sessionList *SessionList) invalidateSessionToken(sessionToken string) {
	sessionList.Lock()

	delete(sessionList.webUsers, sessionToken)

	defer sessionList.Unlock()
}

func (sessionList *SessionList) CheckSessionToken(sessionToken *http.Cookie, err error) (isValid bool, userName string, userId int64) {
	if err != nil {
		return false, "", 0
	}
	if sessionToken == nil {
		return false, "", 0
	}
	sessionList.Lock()
	defer sessionList.Unlock()

	webUser := sessionList.webUsers[sessionToken.Value]
	if webUser == nil {
		return false, "", 0
	}

	// Hard coded 15 minute timeout
	timeOutTimeUtc := time.Now().UTC().Add(-15 * time.Minute)

	if webUser.lastMessageInUtc.Before(timeOutTimeUtc) {
		return false, "", 0
	}

	webUser.lastMessageInUtc = time.Now().UTC()

	return true, webUser.userName, webUser.userId
}

func (sessionList *SessionList) updateMessageInTime(sessionToken string, userId int64) {

	sessionList.Lock()

	webUser := sessionList.webUsers[sessionToken]

	if webUser == nil {
		fmt.Println("Session not found, session=", sessionToken)
	} else {
		webUser.lastMessageInUtc = time.Now().UTC()
	}

	sessionList.Unlock()
}
