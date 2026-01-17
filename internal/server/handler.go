package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/piyushdaiya/antigravity-connect/internal/cdp"
	"github.com/piyushdaiya/antigravity-connect/web"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// LogFilter hides "TLS handshake error" spam
type LogFilter struct{}
func (f *LogFilter) Write(p []byte) (n int, err error) {
	if strings.Contains(string(p), "TLS handshake error") { return len(p), nil }
	return log.Writer().Write(p)
}

func Start(port string, tlsConfig *http.Server) {
	if err := cdp.Init("ws://127.0.0.1:9222"); err != nil {
		log.Printf("⚠️  CDP Error: %v (Ensure IDE is started with --remote-debugging-port=9222)", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(web.Assets)))
	mux.HandleFunc("/ws", handleWebSocket)

	server := &http.Server{
		Addr:      ":" + port,
		Handler:   mux,
		TLSConfig: tlsConfig.TLSConfig,
		ErrorLog:  log.New(&LogFilter{}, "", 0),
	}

	log.Printf("🚀 Server running on https://<YOUR_IP>:%s", port)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil { return }

	log.Println("📱 Phone Connected")
	stopChan := make(chan struct{})

	// Sender Goroutine (Screenshots)
	go func() {
		defer conn.Close()
		ticker := time.NewTicker(200 * time.Millisecond) // 5 FPS
		defer ticker.Stop()

		for {
			select {
			case <-stopChan: return
			case <-ticker.C:
				screenshot, err := cdp.GetScreenshot()
				if err != nil { continue }
				
				conn.SetWriteDeadline(time.Now().Add(500 * time.Millisecond))
				if err := conn.WriteMessage(websocket.BinaryMessage, screenshot); err != nil { return }
			}
		}
	}()

	// Reader Loop (Scroll Events)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil { break }
		
		var m struct{ Y int `json:"y"` }
		if json.Unmarshal(msg, &m) == nil {
			go cdp.SyncScroll(m.Y)
		}
	}
	close(stopChan) // Signal sender to stop
	log.Println("📱 Phone Disconnected")
}