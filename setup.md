# Sparsh Installation & Running Guide

This document outlines the steps required to set up the development environment, install dependencies, and run both the client and server components of the Sparsh Invoice Generator.

## Prerequisites

Regardless of your operating system, you must have **Go** installed.
- **Minimum Version**: Go 1.26.1 (as specified in [go.mod](file:///Users/nikhil/Desktop/sparsh/go.mod))
- **Download**: [golang.org/dl](https://golang.org/dl/)

---

##  macOS Setup

macOS usually comes with some tools, but for Go GUI development with Fyne, you need a C compiler.

1. **Install Xcode Command Line Tools**:
   Open your terminal and run:
   ```bash
   xcode-select --install
   ```

2. **Clone the Project**:
   ```bash
   git clone <repository-url>
   cd sparsh
   ```

3. **Install Go Dependencies**:
   ```bash
   go mod download
   ```

4. **Run the Project**:
   - **Start the Server**:
     ```bash
     go run cmd/server/main.go
     ```
   - **Start the Client (New Terminal)**:
     ```bash
     go run cmd/client/main.go
     ```

---

## ⊞ Windows Setup

Windows requires a C compiler (MinGW) to handle the graphics and native dialog dependencies used by Fyne and `sqweek/dialog`.

1. **Install a C Compiler (MSYS2)**:
   - Download and install [MSYS2](https://www.msys2.org/).
   - Open the **MSYS2 MINGW64** terminal.
   - Run the following command to install the MinGW-w64 toolchain:
     ```bash
     pacman -S mingw-w64-x86_64-toolchain
     ```
   - Add the resulting `bin` folder (usually `C:\msys64\mingw64\bin`) to your Windows **System PATH**.

2. **Clone the Project**:
   ```powershell
   git clone <repository-url>
   cd sparsh
   ```

3. **Install Go Dependencies**:
   ```powershell
   go mod download
   ```

4. **Run the Project**:
   - **Start the Server**:
     ```powershell
     go run cmd/server/main.go
     ```
   - **Start the Client (New Terminal)**:
     ```powershell
     go run cmd/client/main.go
     ```

---

## 📦 Key Dependencies Used

- **GUI Framework**: [Fyne v2](https://fyne.io/)
- **Native Dialogs**: [sqweek/dialog](https://github.com/sqweek/dialog) (Used for native OS Save/Open windows)
- **PDF Generation**: [gofpdf](https://github.com/jung-kurt/gofpdf)
- **Web Server**: [Gin Gonic](https://gin-gonic.com/)
- **Database**: [Modernc SQLite](https://modernc.org/sqlite) (CGO-free SQLite)

## 🛠️ Common Troubleshooting

### "CGO_ENABLED" Error
If you see errors related to `gcc` or `cgo` during `go run`:
- **macOS**: Ensure Xcode tools are installed.
- **Windows**: Ensure MinGW is in your PATH and `go env CGO_ENABLED` is set to `1`.

### Fyne Graphics Issues
Ensure your graphics drivers are up to date, as Fyne relies on OpenGL.
