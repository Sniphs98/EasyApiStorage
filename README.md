# API Storage

Eine moderne File Storage API mit Web-Interface, gebaut mit Go und modernem Frontend.

## Features

- 📁 **Ordner-Navigation** - Erstellen und navigieren durch Ordnerstrukturen
- 🔍 **Suche & Sortierung** - Durchsuchen und sortieren von Dateien
- 📤 **Drag & Drop Upload** - Einfaches Hochladen von Dateien
- 📱 **Responsive Design** - Funktioniert auf Desktop und Mobile
- 🚀 **Streaming Support** - Effizientes Handling großer Dateien
- 📊 **Swagger API Docs** - Vollständige API-Dokumentation

## Schnellstart

### Lokal ausführen

```bash
# Repository klonen
git clone <repository-url>
cd apiStorage

# Dependencies installieren
go mod tidy

# Anwendung starten
go run main.go
```

Die Anwendung ist dann verfügbar unter:
- **Web Interface**: http://localhost:8080/
- **Swagger API Docs**: http://localhost:8080/swagger/

### Mit Docker

```bash
# Image bauen und starten
docker-compose up --build

# Oder nur das Image bauen
docker build -t api-storage .

# Container starten
docker run -p 8080:8080 -v $(pwd)/store:/app/store api-storage
```

### Mit Docker Compose (empfohlen)

```bash
# Starten
docker-compose up -d

# Logs anzeigen
docker-compose logs -f

# Stoppen
docker-compose down
```

## API Endpoints

| Method | Endpoint | Beschreibung |
|--------|----------|--------------|
| GET | `/` | Web Interface |
| GET | `/files?path=<path>` | Dateien/Ordner auflisten |
| POST | `/upload/{filename}?path=<path>` | Datei hochladen |
| GET | `/download/{filepath}` | Datei herunterladen |
| POST | `/create-folder` | Ordner erstellen |
| GET | `/swagger/` | API Dokumentation |

## Konfiguration

### Umgebungsvariablen

- `GIN_MODE`: Gin Framework Modus (`debug`, `release`)
- `PORT`: Server Port (Standard: 8080)

### Docker Volumes

Das `store/` Verzeichnis wird als Volume gemountet, um Dateien persistent zu speichern:

```yaml
volumes:
  - ./store:/app/store
```

## Entwicklung

### Projekt-Struktur

```
apiStorage/
├── main.go              # Hauptanwendung
├── templates/           # HTML Templates
│   └── index.html
├── docs/               # Swagger Dokumentation
├── store/              # Datei-Speicher (ignoriert von Git)
├── Dockerfile          # Docker Image Definition
├── docker-compose.yml  # Docker Compose Konfiguration
└── README.md          # Diese Datei
```

### Build

```bash
# Lokaler Build
go build -o apiStorage main.go

# Cross-Platform Build
GOOS=linux GOARCH=amd64 go build -o apiStorage-linux main.go
GOOS=windows GOARCH=amd64 go build -o apiStorage-windows.exe main.go
```

## Sicherheit

- ✅ Path-Traversal Schutz (`..` wird blockiert)
- ✅ Eingabe-Validierung für alle Endpunkte
- ✅ Non-root Docker Container
- ✅ Sichere Pfad-Konstruktion
- ✅ CORS-Header für Web-Sicherheit

## Lizenz

MIT License - siehe LICENSE Datei für Details.