# 📁 API Storage

> Eine moderne File Storage API mit Web-Interface und Ordner-Navigation

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

Eine benutzerfreundliche File Storage Lösung mit modernem Web-Interface, gebaut mit Go und responsivem Frontend-Design.

![API Storage Interface](https://via.placeholder.com/800x400/6366f1/ffffff?text=API+Storage+Interface)

## ✨ Features

- 📁 **Ordner-Navigation** - Erstellen und navigieren durch Ordnerstrukturen wie in einem Desktop-Dateimanager
- 🔍 **Intelligente Suche** - Echtzeit-Suche durch alle Dateien und Ordner
- 📊 **Flexible Sortierung** - Sortieren nach Name, Datum oder Größe (aufsteigend/absteigend)
- 📤 **Drag & Drop Upload** - Einfaches Hochladen von Dateien mit Progress-Anzeige
- 📱 **Responsive Design** - Optimiert für Desktop, Tablet und Mobile
- 🚀 **Streaming Support** - Effizientes Handling großer Dateien ohne Memory-Probleme
- 📊 **Swagger API Docs** - Vollständige REST API-Dokumentation
- 🔒 **Sicherheit** - Path-Traversal Schutz und sichere Datei-Operationen
- 🐳 **Docker Ready** - Containerisiert und produktionsbereit

## 🚀 Schnellstart

### Lokal ausführen

```bash
# Repository klonen
git clone https://github.com/yourusername/apiStorage.git
cd apiStorage

# Dependencies installieren
go mod tidy

# Anwendung starten
go run main.go
```

Die Anwendung ist dann verfügbar unter:
- **🌐 Web Interface**: http://localhost:8080/
- **📚 API Docs**: http://localhost:8080/swagger/

### Mit Docker (empfohlen)

```bash
# Image bauen
docker build -t api-storage .

# Container starten mit persistentem Storage
docker run -d \
  --name api-storage \
  -p 8080:8080 \
  -v $(pwd)/store:/app/store \
  api-storage

# Logs anzeigen
docker logs -f api-storage
```

## 📖 API Referenz

| Method | Endpoint | Beschreibung | Parameter |
|--------|----------|--------------|-----------|
| `GET` | `/` | Web Interface | - |
| `GET` | `/files` | Dateien/Ordner auflisten | `?path=<pfad>` |
| `POST` | `/upload/{filename}` | Datei hochladen | `?path=<pfad>` |
| `GET` | `/download/{filepath}` | Datei herunterladen | - |
| `POST` | `/create-folder` | Ordner erstellen | JSON Body |
| `GET` | `/swagger/` | API Dokumentation | - |

### Beispiel API-Calls

```bash
# Dateien im Root-Verzeichnis auflisten
curl http://localhost:8080/files

# Dateien in einem Unterordner auflisten
curl http://localhost:8080/files?path=documents

# Datei hochladen
curl -X POST -T myfile.txt http://localhost:8080/upload/myfile.txt

# Ordner erstellen
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"NewFolder","path":""}' \
  http://localhost:8080/create-folder
```

## ⚙️ Konfiguration

### Umgebungsvariablen

| Variable | Beschreibung | Standard |
|----------|--------------|----------|
| `PORT` | Server Port | `8080` |
| `GIN_MODE` | Gin Framework Modus | `debug` |

### Docker Volumes

```bash
# Persistente Datei-Speicherung
-v $(pwd)/store:/app/store

# Oder mit benanntem Volume
docker volume create api-storage-data
-v api-storage-data:/app/store
```

## 🛠️ Entwicklung

### Projekt-Struktur

```
apiStorage/
├── main.go              # 🚀 Hauptanwendung & API Server
├── templates/           # 🎨 HTML Templates
│   └── index.html       #   └── Web Interface
├── docs/               # 📚 Swagger API Dokumentation
├── store/              # 💾 Datei-Speicher (Git-ignoriert)
├── Dockerfile          # 🐳 Container Definition
├── .gitignore          # 🚫 Git Ignore Regeln
└── README.md           # 📖 Diese Dokumentation
```

### Lokale Entwicklung

```bash
# Development Server mit Hot Reload
go run main.go

# Tests ausführen (falls vorhanden)
go test ./...

# Code formatieren
go fmt ./...

# Dependencies aufräumen
go mod tidy
```

### Build für verschiedene Plattformen

```bash
# Linux (64-bit)
GOOS=linux GOARCH=amd64 go build -o apiStorage-linux main.go

# Windows (64-bit)
GOOS=windows GOARCH=amd64 go build -o apiStorage-windows.exe main.go

# macOS (64-bit)
GOOS=darwin GOARCH=amd64 go build -o apiStorage-macos main.go

# ARM64 (z.B. Apple Silicon, Raspberry Pi)
GOOS=linux GOARCH=arm64 go build -o apiStorage-arm64 main.go
```

## 🔒 Sicherheitshinweise

⚠️ **Wichtig**: Dieses Projekt ist für Entwicklung und interne Nutzung konzipiert. Für Produktionsumgebungen im Internet sollten zusätzliche Sicherheitsmaßnahmen implementiert werden:

**Implementierte Basis-Sicherheit:**
- ✅ **Path-Traversal Schutz** - `..` Sequenzen werden blockiert
- ✅ **Eingabe-Validierung** - Pfade und Dateinamen werden bereinigt
- ✅ **Non-root Container** - Docker läuft mit unprivilegiertem User
- ✅ **Sichere Pfad-Konstruktion** - Verhindert Zugriff außerhalb des Store-Verzeichnisses

**Für Produktionsumgebungen empfohlen:**
- 🔐 Authentifizierung und Autorisierung
- 🚦 Rate Limiting und DDoS-Schutz
- 🔒 HTTPS/TLS Verschlüsselung
- 📝 Audit Logging und Monitoring
- 📏 Dateigrößen- und Typ-Beschränkungen

## 🤝 Contributing

Beiträge sind willkommen! Bitte:

1. **Fork** das Repository
2. **Branch** erstellen (`git checkout -b feature/amazing-feature`)
3. **Commit** Änderungen (`git commit -m 'Add amazing feature'`)
4. **Push** zum Branch (`git push origin feature/amazing-feature`)
5. **Pull Request** öffnen

## 📝 Lizenz

Dieses Projekt steht unter der MIT License - siehe [LICENSE](LICENSE) Datei für Details.

## 🙏 Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux) - HTTP Router
- [Swagger](https://swagger.io/) - API Dokumentation
- [Alpine Linux](https://alpinelinux.org/) - Minimales Docker Base Image

---

<div align="center">

**[⭐ Star dieses Projekt](https://github.com/yourusername/apiStorage)** wenn es dir gefällt!

Made with ❤️ and Go

</div>