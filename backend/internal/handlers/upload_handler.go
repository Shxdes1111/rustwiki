package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const uploadDir = "uploads/icons/weapons"
const maxUploadSize = 5 << 20

func (h *WeaponHandler) UploadIcon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, http.StatusBadRequest, "File too large or invalid form")
		return
	}

	file, header, err := r.FormFile("icon")
	if err != nil {
		writeError(w, http.StatusBadRequest, "Missing icon file")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext != ".avif" && ext != ".jpeg" && ext != ".jpg" && ext != ".png" && ext != ".webp" {
		writeError(w, http.StatusBadRequest, "Unsupported format. Use avif, jpeg, png, or webp")
		return
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		h.Logger.WithError(err).Error("UploadIcon: mkdir")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst, err := os.Create(filepath.Join(uploadDir, filename))
	if err != nil {
		h.Logger.WithError(err).Error("UploadIcon: create file")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		h.Logger.WithError(err).Error("UploadIcon: copy")
		writeError(w, http.StatusInternalServerError, "Server error")
		return
	}

	iconPath := fmt.Sprintf("/%s/%s", uploadDir, filename)
	json.NewEncoder(w).Encode(map[string]string{"icon": iconPath})
}
