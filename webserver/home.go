// Copyright 2016 Graeme Dykes.

// Based on https://golang.org/doc/articles/wiki/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"strconv"
	"time"

	//	"log"
	"net/http"
	// "nowbridge/nbLog"
	// "nowbridge/nowbridgeComms"
	// "nowbridge/nowbridgeMQ"
	"os"
	"runtime"

	"github.com/gorilla/websocket"
	// "time"
)

type HomeTemplateParams struct {
	UserName string
	PageHTML template.HTML
}

// var messageQueue nowbridgeComms.MessageQueue
var sessionList SessionList

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	sessionToken, err := r.Cookie("sessionTokenStaff")
	_ /*sessionIsValid*/, userName, _ := sessionList.CheckSessionToken(sessionToken, err)

	/*
		if !sessionIsValid {

			http.Redirect(w, r, "/login/", http.StatusFound)
			return
		}
	*/

	switch r.Method {
	case "GET":

		pageHTML := ""

		w.Header().Set("Content-Type", "text/html")
		t, _ := template.ParseFiles("MainTemplate.htmx")
		t.Execute(w, &HomeTemplateParams{UserName: userName, PageHTML: template.HTML(pageHTML)})
	case "POST":
		return
	}

}

func logFileName() string {
	if runtime.GOOS == "windows" {
		return "C:\\nowbridgelogs\\nowbridge.dynWebStaff"
	}
	// Assume linux else.
	return "/var/log/nowbridge.dynWebStaff"
}

func main() {

	filename = "test.txt"

	// nbLog.Init(logFileName())

	// nbLog.Log.Println("Starting dynWebStaff")

	loadConfiguration()
	// messageQueue.Init(configuration.CentralDBServerAddress, strconv.Itoa(configuration.CentralDBServerPort))

	sessionList.init()
	// http.HandleFunc("/", defaultHandler)

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("."))))

	http.HandleFunc("/ws", handlerWs)

	// http.HandleFunc("/login/favicon.ico", loginHandler_favicon)

	// http.HandleFunc("/canSee/", handler_canSee)
	// http.HandleFunc("/clientTaskServers/", handler_clientTaskServers)
	// http.HandleFunc("/clientUploadLog/", handler_uploadLog)
	// http.HandleFunc("/commands/", handler_commands)
	// http.HandleFunc("/emailTemplates/", handler_emailTemplates)
	// http.HandleFunc("/groups/", handler_groups)
	// http.HandleFunc("/login/", loginHandler)
	// http.HandleFunc("/netInfoServers/", handler_netInfoServers)
	// http.HandleFunc("/options/", handler_options)
	// http.HandleFunc("/restartClient/", handler_restartClient)
	// http.HandleFunc("/sendClientCommand/", handler_sendClientCommand)
	// http.HandleFunc("/showUserLog/", handler_showUserLog)
	// http.HandleFunc("/signups/", handler_signups)
	// http.HandleFunc("/updateClient/", handler_updateClient)
	// http.HandleFunc("/userEvents/", handler_userEvents)
	// http.HandleFunc("/users/", handler_users)

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))

	// http.ListenAndServe(":9092", nil)

	fmt.Println("Listening on port", configuration.HttpListenAddressPort)

	http.ListenAndServe(configuration.HttpListenAddressPort, nil)

}

var (
	addr = flag.String("addr", ":8080", "http service address")
	// homeTempl = template.Must(template.New("").Parse(homeHTML))
	filename string
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

const (
	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second

	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second
)

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func readFileIfModified(lastMod time.Time) ([]byte, time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return nil, lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return nil, lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fi.ModTime(), err
	}
	return p, fi.ModTime(), nil
}

func writer(ws *websocket.Conn, lastMod time.Time) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	fileTicker := time.NewTicker(filePeriod)
	defer func() {
		pingTicker.Stop()
		fileTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-fileTicker.C:
			var p []byte
			var err error

			p, lastMod, err = readFileIfModified(lastMod)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, p); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func handlerWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			fmt.Println(err)
		}
		return
	}

	var lastMod time.Time
	if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		lastMod = time.Unix(0, n)
	}

	go writer(ws, lastMod)
	reader(ws)
}

/*
func encodeAndSend(message interface{}, messageId int, dbSvrConn *nowbridgeComms.LongTcp) {

	// -------------------------------------------------

	// First, JSON encode and add command id to the beginning

	var encodeBuf bytes.Buffer

	encodedBytes, err := json.Marshal(message)

	if err != nil {
		return
	}

	// First 2 bytes are message type. Other server uses these to decode the message.
	encodeBuf.WriteByte(byte(messageId & 0xff))
	encodeBuf.WriteByte(byte(messageId >> 8))

	// Next 2 bytes are encoded message length.
	// encodeBuf.WriteByte(byte(len(encodedBytes) & 0xff))
	// encodeBuf.WriteByte(byte(len(encodedBytes) >> 8))

	encodeBuf.Write(encodedBytes)

	// -------------------------------------------------

	// Next, wrap encoded JSON up in another wrapper

	var messageBuf bytes.Buffer

	// Message Marker
	messageBuf.WriteByte(0xFF)

	// Message Length
	len := encodeBuf.Len()
	messageBuf.WriteByte(byte(len))
	messageBuf.WriteByte(byte(len >> 8))
	messageBuf.WriteByte(byte(len >> 16))
	messageBuf.WriteByte(byte(len >> 24))

	messageBuf.Write(encodeBuf.Bytes())

	n, err2 := dbSvrConn.Write(messageBuf.Bytes())

	if err2 != nil {
		fmt.Println("Error writing.", err2.Error())
	} else {
		if false {
			fmt.Println("n =", n)
		}
	}
}
*/

type Configuration struct {
	CentralDBServerAddress string
	CentralDBServerPort    int
	HttpListenAddressPort  string
}

func configFile() string {
	if runtime.GOOS == "windows" {
		return "conf.json"
	}
	// Assume linux else.
	// return "/etc/nowbridge/dynWebStaff/conf.json"
	return "conf.json"
}

var configuration Configuration

func loadConfiguration() {

	// Defaults
	configuration.CentralDBServerAddress = "127.0.0.1"
	// configuration.CentralDBServerAddress = "185.114.224.29"
	configuration.CentralDBServerPort = 9091
	configuration.HttpListenAddressPort = ":80"

	file, err := os.Open(configFile())

	if err != nil {
		// nbLog.Log.Println("Error in loadConfiguration -", err)
		fmt.Println("Error in loadConfiguration -", err)
		return
	}

	decoder := json.NewDecoder(file)
	configuration = Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// nbLog.Log.Println("CentralDBServerAddress:", configuration.CentralDBServerAddress)
	// nbLog.Log.Println("CentralDBServerPort:", configuration.CentralDBServerPort)
	// nbLog.Log.Println("HttpListenAddressPort:", configuration.HttpListenAddressPort)
}

func handleMessageLoginResponse() {
	fmt.Println("handleMessage_LoginResponse")
}

func returnErrorMessageHtml(w http.ResponseWriter, r *http.Request, errorMessage string, userName string) {
	w.Header().Set("Content-Type", "text/html")
	t, _ := template.ParseFiles("MainTemplate.htmx")
	t.Execute(w, &HomeTemplateParams{UserName: userName, PageHTML: template.HTML(errorMessage)})
}
