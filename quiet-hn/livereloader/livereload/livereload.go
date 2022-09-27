package livereload

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func StartLiveReloadSocket() error {
	return StartLiveReloadSocketOnPort(9090)
}

func StartLiveReloadSocketOnPort(port int) error {

	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGHUP)

	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		if err != nil {
			log.Fatal(err)
		}

		for {
			s := <-sigs
			fmt.Println("Reloading using", s)
			conn.WriteMessage(websocket.BinaryMessage, []byte("reload"))
		}
	})

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("echo")
		conn, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

		if err != nil {
			fmt.Fprintf(w, "%+v", err)
			fmt.Println(err)
			return
		}

		for {
			// Read message from browser
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Print the message to the console
			fmt.Printf("%s sent: %s\n", conn.RemoteAddr(), string(msg))

			// Write message back to browser
			if err = conn.WriteMessage(msgType, msg); err != nil {
				return
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})

	fmt.Printf("Running Live Reloader on port: %d", port)
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", port), nil)
}
