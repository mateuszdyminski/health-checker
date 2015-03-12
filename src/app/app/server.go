package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func LaunchServer(opts Options) {
	// Handle routes

	r := mux.NewRouter()

	r.HandleFunc("/wsapi/ws", serveWs)

	r.Handle("/{path:.*}", http.FileServer(http.Dir(opts.StaticDir))).Name("statics")

	// Listen on hostname:port
	glog.Infof("Listening on %s:%d...\n", opts.Hostname, opts.Port)
	http.Handle("/", &loggingHandler{r})
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", opts.Hostname, opts.Port), nil)
	if err != nil {
		glog.Errorf("Error: %s", err)
	}
}

type loggingHandler struct {
	http.Handler
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	h.Handler.ServeHTTP(w, req)

	elapsed := time.Since(t)
	glog.Infof("%s - [%s] \"%s %s %s\" \"%s\" \"%s\" \"Took: %s\"\n",
		strings.Split(req.RemoteAddr, ":")[0],
		t.Format("02/Jan/2006:15:04:05 -0700"), req.Method,
		req.RequestURI, req.Proto, req.Referer(), req.UserAgent(), elapsed)

	if elapsed > 200*time.Millisecond {
		glog.Errorf("Long run request: %s - \"%s %s %s\" Took: %s", strings.Split(req.RemoteAddr, ":")[0], req.Method, req.RequestURI, req.Proto, elapsed)
	}
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Registering client to WS")
	if req.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		fmt.Printf("Error %+v\n", err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	go c.writePump()
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}
