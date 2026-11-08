package server

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/user/sharedrop/internal/auth"
	"github.com/user/sharedrop/internal/config"
)

// FileInfo is what we send to the HTML template
type FileInfo struct {
	Name    string
	Size    string
	RawSize int64
	IsDir   bool
	ModTime string
	Icon    string
	Path    string
}

// Server is the main HTTP file server
type Server struct {
	cfg      *config.Config
	rootDir  string
	sessions *auth.SessionStore
	mux      *http.ServeMux
}

func New(cfg *config.Config, rootDir string) *Server {
	s := &Server{
		cfg:      cfg,
		rootDir:  rootDir,
		sessions: auth.NewSessionStore(),
		mux:      http.NewServeMux(),
	}
	s.registerRoutes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/", s.authMiddleware(s.handleBrowse))
	s.mux.HandleFunc("/download", s.authMiddleware(s.handleDownload))
	s.mux.HandleFunc("/upload", s.authMiddleware(s.handleUpload))
	s.mux.HandleFunc("/delete", s.authMiddleware(s.handleDelete))
	s.mux.HandleFunc("/zip", s.authMiddleware(s.handleZip))
	s.mux.HandleFunc("/api/files", s.authMiddleware(s.handleAPIFiles))
	s.mux.HandleFunc("/login", s.handleLogin)
	s.mux.HandleFunc("/logout", s.handleLogout)
	s.mux.HandleFunc("/static/", s.handleStatic)
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !s.cfg.PasswordSet {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("sd_session")
		if err == nil && s.sessions.Valid(cookie.Value) {
			next(w, r)
			return
		}

		http.Redirect(w, r, "/login?next="+r.URL.RequestURI(), http.StatusFound)
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.PasswordSet {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	if r.Method == http.MethodPost {
		password := r.FormValue("password")
		if auth.CheckPassword(password, s.cfg.PasswordHash) {
			token, err := auth.GenerateToken()
			if err != nil {
				http.Error(w, "Internal error", 500)
				return
			}
			s.sessions.Add(token)
			http.SetCookie(w, &http.Cookie{
				Name:     "sd_session",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				MaxAge:   86400 * 7,
			})
			next := r.FormValue("next")
			if next == "" {
				next = "/"
			}
			http.Redirect(w, r, next, http.StatusFound)
			return
		}
		renderLogin(w, true)
		return
	}

	renderLogin(w, false)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sd_session")
	if err == nil {
		s.sessions.Remove(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "sd_session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *Server) handleBrowse(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}

	absPath, err := s.safePath(relPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if !info.IsDir() {
		http.ServeFile(w, r, absPath)
		return
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	files := make([]FileInfo, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()

		if !s.cfg.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fi, err := entry.Info()
		if err != nil {
			continue
		}

		childRelPath := filepath.Join(relPath, name)
		if relPath == "." {
			childRelPath = name
		}

		finfo := FileInfo{
			Name:    name,
			IsDir:   entry.IsDir(),
			ModTime: fi.ModTime().Format("2006-01-02 15:04"),
			Path:    childRelPath,
		}

		if !entry.IsDir() {
			finfo.RawSize = fi.Size()
			finfo.Size = auth.FormatSize(fi.Size())
			finfo.Icon = fileIcon(name)
		} else {
			finfo.Icon = "[dir]"
		}

		files = append(files, finfo)
	}

	sort.Slice(files, func(i, j int) bool {
		if files[i].IsDir != files[j].IsDir {
			return files[i].IsDir
		}
		return strings.ToLower(files[i].Name) < strings.ToLower(files[j].Name)
	})

	crumbs := buildBreadcrumbs(relPath)

	data := map[string]interface{}{
		"Files":       files,
		"CurrentPath": relPath,
		"Breadcrumbs": crumbs,
		"AllowUpload": s.cfg.AllowUpload,
		"AllowDelete": s.cfg.AllowDelete,
		"PasswordSet": s.cfg.PasswordSet,
		"RootName":    filepath.Base(s.rootDir),
	}

	renderTemplate(w, browseTmpl, data)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	absPath, err := s.safePath(relPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if info.IsDir() {
		http.Redirect(w, r, "/?path="+relPath, http.StatusFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, info.Name()))
	http.ServeFile(w, r, absPath)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.AllowUpload {
		http.Error(w, "Upload disabled", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	maxBytes := int64(s.cfg.MaxUploadMB) * 1024 * 1024
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	if err := r.ParseMultipartForm(maxBytes); err != nil {
		http.Error(w, "File too large", http.StatusRequestEntityTooLarge)
		return
	}

	destPath := r.FormValue("path")
	absDestDir, err := s.safePath(destPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	files := r.MultipartForm.File["files"]
	for _, fileHeader := range files {
		src, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer src.Close()

		cleanName := filepath.Base(fileHeader.Filename)
		destFile := filepath.Join(absDestDir, cleanName)

		dst, err := os.Create(destFile)
		if err != nil {
			continue
		}
		defer dst.Close()

		io.Copy(dst, src)
	}

	http.Redirect(w, r, "/?path="+destPath, http.StatusFound)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	if !s.cfg.AllowDelete {
		http.Error(w, "Delete disabled", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	relPath := r.FormValue("path")
	absPath, err := s.safePath(relPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	if absPath == s.rootDir {
		http.Error(w, "Cannot delete root", http.StatusForbidden)
		return
	}

	os.Remove(absPath)
	parentPath := filepath.Dir(relPath)
	http.Redirect(w, r, "/?path="+parentPath, http.StatusFound)
}

func (s *Server) handleZip(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	absPath, err := s.safePath(relPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	info, err := os.Stat(absPath)
	if err != nil || !info.IsDir() {
		http.Error(w, "Not a directory", http.StatusBadRequest)
		return
	}

	zipName := filepath.Base(absPath) + ".zip"
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, zipName))
	w.Header().Set("Content-Type", "application/zip")

	zw := zip.NewWriter(w)
	defer zw.Close()

	filepath.Walk(absPath, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(absPath, path)
		f, err := zw.Create(rel)
		if err != nil {
			return nil
		}
		src, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer src.Close()
		io.Copy(f, src)
		return nil
	})
}

func (s *Server) handleAPIFiles(w http.ResponseWriter, r *http.Request) {
	relPath := r.URL.Query().Get("path")
	if relPath == "" {
		relPath = "."
	}

	absPath, err := s.safePath(relPath)
	if err != nil {
		http.Error(w, "Access denied", http.StatusForbidden)
		return
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		http.Error(w, "Cannot read directory", http.StatusInternalServerError)
		return
	}

	type APIFile struct {
		Name    string `json:"name"`
		Size    int64  `json:"size"`
		IsDir   bool   `json:"is_dir"`
		ModTime string `json:"mod_time"`
	}

	result := make([]APIFile, 0)
	for _, e := range entries {
		fi, _ := e.Info()
		result = append(result, APIFile{
			Name:    e.Name(),
			Size:    func() int64 { if fi != nil { return fi.Size() }; return 0 }(),
			IsDir:   e.IsDir(),
			ModTime: func() string { if fi != nil { return fi.ModTime().Format(time.RFC3339) }; return "" }(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) safePath(relPath string) (string, error) {
	if relPath == "" || relPath == "." {
		return s.rootDir, nil
	}

	clean := filepath.Clean(relPath)
	abs := filepath.Join(s.rootDir, clean)

	rootAbs, _ := filepath.Abs(s.rootDir)
	absClean, _ := filepath.Abs(abs)

	if !strings.HasPrefix(absClean, rootAbs) {
		return "", fmt.Errorf("path escape attempt")
	}

	return abs, nil
}

type Breadcrumb struct {
	Name string
	Path string
}

func buildBreadcrumbs(relPath string) []Breadcrumb {
	if relPath == "." || relPath == "" {
		return nil
	}

	parts := strings.Split(filepath.ToSlash(relPath), "/")
	crumbs := make([]Breadcrumb, 0, len(parts))
	cumulative := ""

	for _, part := range parts {
		if part == "" || part == "." {
			continue
		}
		if cumulative == "" {
			cumulative = part
		} else {
			cumulative = cumulative + "/" + part
		}
		crumbs = append(crumbs, Breadcrumb{Name: part, Path: cumulative})
	}

	return crumbs
}

func fileIcon(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".svg", ".ico":
		return "[img]"
	case ".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return "[vid]"
	case ".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a":
		return "[aud]"
	case ".pdf":
		return "[pdf]"
	case ".doc", ".docx":
		return "[doc]"
	case ".xls", ".xlsx":
		return "[xls]"
	case ".ppt", ".pptx":
		return "[ppt]"
	case ".zip", ".rar", ".tar", ".gz", ".7z":
		return "[zip]"
	case ".go", ".py", ".js", ".ts", ".java", ".c", ".cpp", ".rs", ".rb":
		return "[src]"
	case ".json", ".yaml", ".yml", ".xml", ".toml":
		return "[cfg]"
	case ".txt", ".md", ".log":
		return "[txt]"
	case ".exe", ".msi", ".deb", ".rpm", ".dmg", ".app":
		return "[bin]"
	default:
		ext2 := mime.TypeByExtension(ext)
		if strings.HasPrefix(ext2, "image/") {
			return "[img]"
		}
		return "[file]"
	}
}

func renderTemplate(w http.ResponseWriter, tmplStr string, data interface{}) {
	tmpl, err := template.New("page").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(tmplStr)
	if err != nil {
		http.Error(w, "Template error: "+err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, data)
}

func renderLogin(w http.ResponseWriter, wrongPass bool) {
	data := map[string]interface{}{
		"WrongPassword": wrongPass,
	}
	renderTemplate(w, loginTmpl, data)
}
