# Antigravity Connect

**Antigravity Connect** is a lightweight, high-performance tool that allows you to control the scrolling of your **Google Antigravity IDE** (or VS Code/Chrome) directly from your smartphone.

It mirrors your IDE's screen to your phone and syncs scroll events in real-time, allowing for a "lean back" reading experience while reviewing code.

> Note:** This project has been completely rewritten in **Go (Golang)**. It no longer requires Node.js, NPM, or OpenSSL. It runs as a single, static binary on Linux, Windows, and macOS.This project is a refined fork/extension based on the original [Antigravity Shit-Chat](https://github.com/gherghett/Antigravity-Shit-Chat) by gherghett and [Antigravity Phone-Chat](https://github.com/krishnakanthb13/antigravity_phone_chat) by krishnakanthb13

## 🚀 Features

* **Zero Dependencies:** No `node_modules`, no Python, no external libraries required on the target machine.
* **Cross-Platform:** Native support for Linux (Ubuntu/Debian/Arch), Windows, and macOS.
* **Automatic HTTPS:** Generates valid self-signed certificates on-the-fly to ensure secure context (required for mobile browsers).
* **Smart Discovery:** Automatically detects and attaches to the correct IDE window, ignoring background processes.
* **Auto-Reconnect:** The mobile client automatically reconnects if the phone locks or the network drops.
* **High Performance:** Uses raw binary WebSockets for screen mirroring (low latency) instead of Base64 encoding.

## 🛠 Architecture

This tool bridges your phone and your IDE using the **Chrome DevTools Protocol (CDP)**.

1. **The Server (Go):** Starts a secure HTTPS/WebSocket server on your laptop (Port 3000).
2. **The Controller (CDP):** Connects to the IDE's debug port (`9222`) to capture screenshots and inject scroll commands.
3. **The Client (Mobile):** A lightweight HTML/JS frontend (embedded in the binary) that displays the stream and sends touch events.

## 📋 Prerequisites

1. **Google Antigravity IDE** (or VS Code / Chrome).
2. The IDE **must** be launched with the remote debugging port open.

## 📦 Installation & Build

### Option 1: Build from Source (Recommended)

You need [Go installed](https://go.dev/dl/) on your machine.

```bash
# 1. Clone the repository
git clone [https://github.com/piyushdaiya/antigravity-connect.git](https://github.com/piyushdaiya/antigravity-connect.git
cd antigravity-connect

# 2. Build the binary (Detects your OS automatically)
go build -o antigravity-connect ./cmd/server

# 3. Make executable (Linux/Mac only)
chmod +x antigravity-connect
```

## 🚦 Usage Guide

### Step 1: Launch the IDE with Debugging Enabled

The tool requires the IDE to "listen" for connections on port 9222. You cannot simply click the app icon unless you have modified the shortcut.

**Linux (Ubuntu/Debian)**

* **Temporary (Terminal):**
  **Bash**
  
  ```
  antigravity --remote-debugging-port=9222
  ```
* **Permanent (Desktop Shortcut):** Edit `/usr/share/applications/antigravity.desktop` and update the `Exec` line:
  **Ini, TOML**
  
  ```
  Exec=/usr/bin/antigravity --remote-debugging-port=9222 --args %F
  ```

**Windows**

* **PowerShell:**
  **PowerShell**
  
  ```
  & "C:\Users\$env:USERNAME\AppData\Local\Programs\Google Antigravity\Antigravity.exe" --remote-debugging-port=9222
  ```
* **Shortcut:** Right-click your desktop shortcut -> Properties -> Target. Add ` --remote-debugging-port=9222` to the end.

**macOS**

* **Terminal:**
  **Bash**
  
  ```
  /Applications/Google\ Antigravity.app/Contents/MacOS/Antigravity --remote-debugging-port=9222
  ```

### Step 2: Run the Connect Tool

Navigate to the folder where you built the tool and run the binary.

**Bash**

```
# Linux / macOS
./antigravity-connect

# Windows
.\antigravity-connect.exe
```

*Success Output:*

**Plaintext**

```
Found IDE Target ID: A1B2C3...
✅ Attached to IDE Target: Antigravity - Agent Manager
🚀 Server running on https://192.168.1.50:3000
```

### Step 3: Connect Your Phone

1. Ensure your phone is connected to the **same Wi-Fi network** as your computer.
2. Open Chrome (Android) or Safari (iOS).
3. Type the URL shown in the terminal (e.g., `https://192.168.1.50:3000`).
4. **Important:** You will see a security warning ("Your connection is not private").
   * **Android:** Click **Advanced** > **Proceed to... (unsafe)**.
   * **iOS:** Click **Show Details** > **visit this website**.
5. The screen should mirror immediately. Scroll to test!

---

## 🔧 Troubleshooting

### 1. "Waiting for snapshot..." or Black Screen

* **Cause:** The tool is connected to the phone, but the IDE is not sending images.
* **Solution:**
  * Check if the IDE window is **minimized**. Restore it (it can be in the background, but not minimized to the dock/taskbar).
  * Ensure you launched the IDE with `--remote-debugging-port=9222`.
  * Kill the `antigravity-connect` process and restart it.

### 2. "Connection Refused" (in Terminal)

* **Cause:** The tool cannot find the IDE's debugger.
* **Solution:**
  * Close **ALL** instances of Antigravity/Chrome/VS Code.
  * Run the launch command from Step 1 again.
  * Open `http://localhost:9222` in your browser to verify the debugger is active.

### 3. Phone says "Disconnected. Retrying..."

* **Cause:** The phone lost connection to the server.
* **Solution:**
  * This is normal if you locked your screen or switched apps.
  * Unlock the phone; the page will auto-reconnect within 2 seconds.
  * If it loops forever, refresh the page manually.

### 4. "TLS handshake error" in Terminal Logs

* **Cause:** The phone is rejecting the self-signed certificate.
* **Solution:**
  * This is expected behavior until you accept the certificate on the phone.
  * Follow the steps in **Step 3** to click "Proceed" or "Visit Website" on your mobile browser.
  * Once accepted, these errors will stop appearing.

