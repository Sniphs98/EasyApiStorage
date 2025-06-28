# 📁 SimpelWebFileBrowser

> Simple web-based file browser with folder navigation and modern interface

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A simple and user-friendly web file browser with modern interface, built with Go and responsive frontend design.

## ✨ Features

- 📁 **Folder Navigation** - Create and navigate through folder structures
- 🔍 **Search** - Search through all files and folders
- 📊 **Sorting** - Sort by name, date, or size
- 📱 **Responsive Design** - Works on desktop, tablet, and mobile
- 📊 **Swagger API Docs** - Complete REST API documentation
- 🐳 **Docker Ready** - Containerized and production-ready

## 🚀 Quick Start

### Run Locally

```bash
# Clone repository
git clone https://github.com/yourusername/SimpelWebFileBrowser.git
cd SimpelWebFileBrowser

# Install dependencies
go mod tidy

# Start application
go run main.go
```

Access the application:

- **Web Interface**: http://localhost:8080/
- **API Docs**: http://localhost:8080/swagger/

### Run with Docker

```bash
# Build and run
docker build -t simpel-web-file-browser .
docker run -d -p 8080:8080 -v $(pwd)/store:/app/store simpel-web-file-browser
```

## 📖 API Reference

| Method | Endpoint                         | Description        |
| ------ | -------------------------------- | ------------------ |
| `GET`  | `/`                              | Web Interface      |
| `GET`  | `/files?path=<path>`             | List files/folders |
| `POST` | `/upload/{filename}?path=<path>` | Upload file        |
| `GET`  | `/download/{filepath}`           | Download file      |
| `POST` | `/create-folder`                 | Create folder      |

### Example API Calls

```bash
# List files
curl http://localhost:8080/files

# Upload file
curl -X POST -T myfile.txt http://localhost:8080/upload/myfile.txt

# Create folder
curl -X POST -H "Content-Type: application/json" \
  -d '{"name":"NewFolder","path":""}' \
  http://localhost:8080/create-folder
```

## 🔒 Security Notice

⚠️ **Important**: This project is designed for development and internal use. For production environments on the internet, additional security measures should be implemented:

**Current Security:**

- ✅ Path-traversal protection
- ✅ Input validation
- ✅ Non-root Docker container

**Recommended for Production:**

- 🔐 Authentication & authorization
- 🚦 Rate limiting
- 🔒 HTTPS/TLS encryption
- 📝 Audit logging

## 🤝 Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**[⭐ Star this project](https://github.com/yourusername/SimpelWebFileBrowser)** if you like it!

Made with ❤️ and Go

</div>
