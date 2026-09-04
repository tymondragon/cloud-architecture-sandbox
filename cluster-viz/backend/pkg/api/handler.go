package api

import (
	"context"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/k8s"
	"github.com/tymondragon/cloud-architecture-sandbox/cluster-viz/pkg/models"
	"nhooyr.io/websocket"
)


var frontendFS embed.FS

// Handler manages HTTP and WebSocket connections
type Handler struct {
	graphBuilder *k8s.GraphBuilder
	watcher      *k8s.ResourceWatcher
	hub          *Hub
}

// Hub manages WebSocket clients
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan models.GraphEvent
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Client represents a WebSocket connection
type Client struct {
	conn *websocket.Conn
	send chan models.GraphEvent
}

// NewHandler creates a new HTTP/WebSocket handler
func NewHandler(graphBuilder *k8s.GraphBuilder, watcher *k8s.ResourceWatcher) *Handler {
	hub := &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan models.GraphEvent, 100),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	return &Handler{
		graphBuilder: graphBuilder,
		watcher:      watcher,
		hub:          hub,
	}
}

// Start begins the broadcast loop
func (h *Handler) Start(ctx context.Context) {
	go h.hub.run(ctx)
	go h.broadcastFromWatcher(ctx)
}

// ServeHTTP routes requests
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/graph":
		h.handleGraph(w, r)
	case "/api/ws":
		h.handleWebSocket(w, r)
	default:
		h.handleStatic(w, r)
	}
}

// handleGraph returns the full graph snapshot
func (h *Handler) handleGraph(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	graph := h.graphBuilder.Snapshot()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(graph); err != nil {
		log.Printf("Error encoding graph: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// handleWebSocket upgrades to WebSocket and manages the connection
func (h *Handler) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow cross-origin for local dev
	})
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		conn: conn,
		send: make(chan models.GraphEvent, 10),
	}

	h.hub.register <- client

	// Send initial SYNC event
	graph := h.graphBuilder.Snapshot()
	syncEvent := models.GraphEvent{
		Type:  models.EventTypeSync,
		Graph: graph,
	}

	select {
	case client.send <- syncEvent:
	case <-time.After(5 * time.Second):
		log.Println("Timeout sending SYNC event")
	}

	// Start read/write pumps
	go client.writePump(r.Context())
	client.readPump(r.Context(), h.hub)
}

// handleStatic serves the embedded frontend SPA
func (h *Handler) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Strip the embedded path prefix
	staticFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Printf("Error creating sub filesystem: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Serve index.html for all non-file routes (SPA routing)
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	// Check if file exists
	if _, err := staticFS.Open(path[1:]); err != nil {
		// File not found, serve index.html for SPA routing
		path = "/index.html"
	}

	http.FileServer(http.FS(staticFS)).ServeHTTP(w, r)
}

// broadcastFromWatcher forwards watcher events to hub
func (h *Handler) broadcastFromWatcher(ctx context.Context) {
	for {
		select {
		case event := <-h.watcher.Events():
			h.hub.broadcast <- event
		case <-ctx.Done():
			return
		}
	}
}

// Hub methods

func (hub *Hub) run(ctx context.Context) {
	for {
		select {
		case client := <-hub.register:
			hub.mu.Lock()
			hub.clients[client] = true
			hub.mu.Unlock()
			log.Printf("Client connected. Total: %d", len(hub.clients))

		case client := <-hub.unregister:
			hub.mu.Lock()
			if _, ok := hub.clients[client]; ok {
				delete(hub.clients, client)
				close(client.send)
			}
			hub.mu.Unlock()
			log.Printf("Client disconnected. Total: %d", len(hub.clients))

		case event := <-hub.broadcast:
			hub.mu.RLock()
			for client := range hub.clients {
				select {
				case client.send <- event:
				default:
					// Client buffer full, skip
				}
			}
			hub.mu.RUnlock()

		case <-ctx.Done():
			return
		}
	}
}

// Client methods

func (c *Client) readPump(ctx context.Context, hub *Hub) {
	defer func() {
		hub.unregister <- c
		c.conn.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, _, err := c.conn.Read(ctx)
		if err != nil {
			if websocket.CloseStatus(err) != websocket.StatusNormalClosure {
				log.Printf("WebSocket read error: %v", err)
			}
			return
		}
		// We don't expect messages from client, just keep connection alive
	}
}

func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case event, ok := <-c.send:
			if !ok {
				c.conn.Close(websocket.StatusNormalClosure, "")
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := c.conn.Write(writeCtx, websocket.MessageText, mustJSON(event))
			cancel()

			if err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			// Send ping to keep connection alive
			writeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			err := c.conn.Ping(writeCtx)
			cancel()

			if err != nil {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// mustJSON marshals v to JSON or panics
func mustJSON(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
