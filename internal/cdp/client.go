package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/chromedp/cdproto/target"
	"github.com/chromedp/chromedp"
)

var (
	globalCtx    context.Context
	globalCancel context.CancelFunc
)

type Target struct {
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"` // Added URL for better detection
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func Init(debuggerPortURL string) error {
	if debuggerPortURL == "" {
		debuggerPortURL = "ws://127.0.0.1:9222"
	}

	targetID, err := findMainPageTarget(debuggerPortURL)
	if err != nil {
		return fmt.Errorf("IDE Discovery Failed: %v", err)
	}

	log.Printf("✅ Attached to Target: %s", targetID)

	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), debuggerPortURL)
	globalCtx, globalCancel = chromedp.NewContext(allocCtx, chromedp.WithTargetID(target.ID(targetID)))

	if err := chromedp.Run(globalCtx, chromedp.Evaluate("console.log('Antigravity Remote Connected')", nil)); err != nil {
		cancel()
		return fmt.Errorf("connection lost: %v", err)
	}
	return nil
}

func findMainPageTarget(baseURL string) (string, error) {
	// 1. Get List of Targets
	queryURL := strings.Replace(baseURL, "ws://", "http://", 1)
	resp, err := http.Get(queryURL + "/json")
	if err != nil {
		return "", fmt.Errorf("could not connect to %s. Is the IDE running? (Error: %v)", queryURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil { return "", err }

	var targets []Target
	if err := json.Unmarshal(body, &targets); err != nil { return "", err }

	// --- DIAGNOSTIC LOGS (Check your terminal!) ---
	log.Println("🔍 Scanning for IDE Windows...")
	for _, t := range targets {
		if t.WebSocketDebuggerURL != "" {
			log.Printf("   Found: [%s] '%s' (URL: %s)", t.Type, t.Title, t.URL)
		}
	}
	// ---------------------------------------------

	// STRATEGY 1: Strict match for Antigravity App
	for _, t := range targets {
		title := strings.ToLower(t.Title)
		if strings.Contains(title, "antigravity") || strings.Contains(title, "agent manager") {
			log.Printf("🎯 Match Found (By Title): %s", t.Title)
			return t.ID, nil
		}
	}

	// STRATEGY 2: Electron "App" Type (Very specific to IDEs)
	// Chrome tabs are usually type "page", Electron apps are often "app"
	for _, t := range targets {
		if t.WebSocketDebuggerURL != "" && t.Type == "app" {
			log.Printf("🎯 Match Found (By Type 'app'): %s", t.Title)
			return t.ID, nil
		}
	}

	// STRATEGY 3: Exclude Google Search / New Tab
	// If we are forced to pick a random page, ensure it's NOT a standard Chrome tab
	for _, t := range targets {
		if t.WebSocketDebuggerURL != "" && t.Type == "page" {
			// Skip known "Junk" targets
			if strings.Contains(t.Title, "SharedWorker") || 
			   strings.Contains(t.Title, "Extension") || 
			   strings.Contains(t.Title, "New Tab") || 
			   strings.Contains(t.Title, "Google Chrome") {
				continue
			}
			
			log.Printf("⚠️ Fallback Selection: %s", t.Title)
			return t.ID, nil
		}
	}

	return "", fmt.Errorf("no valid IDE window found. Close Chrome and restart the IDE")
}
// GetScreenshot captures a specific area of the screen
func GetScreenshot() ([]byte, error) {
	if globalCtx == nil {
		return nil, fmt.Errorf("CDP not ready")
	}

	var buf []byte

	// =====================================================================
	// 🎯 SELECTOR CONFIGURATION (Check ui_inspector.js in your repo!)
	// =====================================================================
	
	// OPTION A: Main Editor Area (Default for Agent Workspaces)
	// Use this if your Agent Manager is open as a main tab in the center.
	//selector := ".part.editor" 

	// OPTION B: Right-Side Chat Panel (Auxiliary Bar)
	// Use this if your Agent is in the right-hand sidebar (like Copilot/Chat).
	selector := ".part.auxiliarybar"

	// OPTION C: Left-Side Explorer Panel
	// Use this if your Agent is in the left sidebar.
	// selector := ".part.sidebar"
	
	// =====================================================================

	// 1. Try to capture the specific cropped area
	// 'chromedp.NodeVisible' ensures we wait until it's actually rendered
	err := chromedp.Run(globalCtx, chromedp.Screenshot(selector, &buf, chromedp.NodeVisible))
	
	// 2. Fallback: If the selector isn't found, capture the full screen
	if err != nil || len(buf) == 0 {
		// log.Printf("⚠️ Crop failed (Selector '%s' not found). Sending full screen.", selector)
		if err := chromedp.Run(globalCtx, chromedp.CaptureScreenshot(&buf)); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

func SyncScroll(scrollY int) error {
	if globalCtx == nil { return fmt.Errorf("CDP not ready") }
	script := fmt.Sprintf("window.scrollTo({top: %d, behavior: 'auto'});", scrollY)
	return chromedp.Run(globalCtx, chromedp.Evaluate(script, nil))
}