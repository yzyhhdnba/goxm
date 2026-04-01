package media

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	maxVideoSize = 200 << 20
	maxImageSize = 10 << 20
)

type StoredObject struct {
	RelativePath string
	PublicURL    string
}

type LocalStorage struct {
	rootDir       string
	publicBaseURL string
}

func NewLocalStorage(rootDir string, publicBaseURL string) *LocalStorage {
	return &LocalStorage{
		rootDir:       strings.TrimSpace(rootDir),
		publicBaseURL: normalizePublicBaseURL(publicBaseURL),
	}
}

func (s *LocalStorage) SaveVideoSource(videoID uint64, fileHeader *multipart.FileHeader) (StoredObject, error) {
	return s.save(videoID, fileHeader, "source", maxVideoSize, map[string]struct{}{
		".mp4":  {},
		".m4v":  {},
		".mov":  {},
		".webm": {},
	})
}

func (s *LocalStorage) SaveVideoCover(videoID uint64, fileHeader *multipart.FileHeader) (StoredObject, error) {
	return s.save(videoID, fileHeader, "cover", maxImageSize, map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
		".png":  {},
		".webp": {},
	})
}

func (s *LocalStorage) save(videoID uint64, fileHeader *multipart.FileHeader, fileBase string, maxSize int64, allowedExt map[string]struct{}) (StoredObject, error) {
	if s == nil {
		return StoredObject{}, fmt.Errorf("media storage is unavailable")
	}
	if fileHeader == nil {
		return StoredObject{}, fmt.Errorf("upload file is required")
	}
	if fileHeader.Size > 0 && fileHeader.Size > maxSize {
		return StoredObject{}, fmt.Errorf("file too large")
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if _, ok := allowedExt[ext]; !ok {
		return StoredObject{}, fmt.Errorf("unsupported file type")
	}

	relativePath := path.Join("videos", strconv.FormatUint(videoID, 10), fileBase+ext)
	fullPath := filepath.Join(s.rootDir, filepath.FromSlash(relativePath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return StoredObject{}, fmt.Errorf("create media directory: %w", err)
	}

	src, err := fileHeader.Open()
	if err != nil {
		return StoredObject{}, fmt.Errorf("open upload file: %w", err)
	}
	defer src.Close()

	dst, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return StoredObject{}, fmt.Errorf("create media file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return StoredObject{}, fmt.Errorf("write media file: %w", err)
	}

	return StoredObject{
		RelativePath: relativePath,
		PublicURL:    path.Join(s.publicBaseURL, relativePath),
	}, nil
}

func normalizePublicBaseURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "/uploads"
	}
	normalized := "/" + strings.Trim(trimmed, "/")
	if normalized == "/" {
		return "/uploads"
	}
	return normalized
}
