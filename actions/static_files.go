package actions

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/buffalo"
)

var assetsBoxSeparator = ""
var assetMaxAge = 31536000

func configureAssetsBoxSeparator() {
	for _, name := range assetsBox.List() {
		if strings.Contains(name, "\\") {
			assetsBoxSeparator = "\\"
			return
		}
		if strings.Contains(name, "/") {
			assetsBoxSeparator = "/"
			return
		}
	}
	assetsBoxSeparator = "/"
}

func staticFileNotFound() (http.File, error) {
	return nil, fmt.Errorf("not found")
}

var builtinMimeTypesLower = map[string]string{
	".css":  "text/css; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpg":  "image/jpeg",
	".js":   "application/javascript",
	".wasm": "application/wasm",
	".pdf":  "application/pdf",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".xml":  "text/xml; charset=utf-8",
}

func staticFileGetMimeType(ext string) string {
	if v, ok := builtinMimeTypesLower[ext]; ok {
		return v
	}
	ct := mime.TypeByExtension(ext)
	if ct == "" {
		ct = "application/octet-stream"
	}
	return ct
}

func resolveStaticFile(c buffalo.Context) (http.File, error) {
	// Normalize the request to only contain forward slashes
	path := strings.ReplaceAll(c.Param("asset"), "\\", "/")

	// Trim leading and trailing forward slashes
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimPrefix(path, "/")

	// Check for obviously bad things and reject them early
	if strings.Contains(path, "../") || strings.Contains(path, "\x00") {
		return staticFileNotFound()
	}

	// Switch to backslashes if necessary to resolve the file
	bpath := path
	if assetsBoxSeparator != "/" {
		bpath = strings.ReplaceAll(path, "/", "\\")
	}

	// Return the gzip content if the client supports it and its available
	aenc := c.Request().Header["Accept-Encoding"]

	// Shortcut since everything sending this header accepts gzip
	gz := len(aenc) > 0
	if gz {
		f, err := assetsBox.Open(bpath + ".gz")
		if err == nil {
			return f, nil
		}
	}

	// Return the uncompressed content instead
	f, err := assetsBox.Open(bpath)
	if err != nil {
		return f, err
	}

	return f, nil
}

// StaticFileGet returns the contents of a static asset
func StaticFileGet(c buffalo.Context) error {
	f, err := resolveStaticFile(c)
	if err != nil {
		return c.Error(404, fmt.Errorf("not found"))
	}
	s, err := f.Stat()
	if err != nil {
		return c.Error(404, fmt.Errorf("not found"))
	}

	flen := s.Size()
	name := s.Name()
	isgz := false

	if filepath.Ext(name) == ".gz" {
		bits := strings.Split(name, ".")
		if len(bits) > 2 {
			name = strings.Join(bits[0:len(bits)-1], ".")
			isgz = true
		}
	}

	ct := staticFileGetMimeType(filepath.Ext(name))

	buff := make([]byte, s.Size())
	n, err := f.Read(buff)
	if err != nil {
		return c.Error(400, fmt.Errorf("unreadable"))
	}
	if int64(n) != flen {
		return c.Error(400, fmt.Errorf("truncated"))
	}

	c.Response().Header().Add("ETag", fmt.Sprintf("%.8x", s.ModTime().Unix()))
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", assetMaxAge))
	c.Response().Header().Add("Content-Length", fmt.Sprintf("%d", flen))
	c.Response().Header().Add("Content-Type", ct)

	if isgz {
		c.Response().Header().Add("Content-Encoding", "gzip")
	}

	c.Response().WriteHeader(200)
	c.Response().Write(buff)
	return nil
}

// StaticFileHead returns the headers of a static asset
func StaticFileHead(c buffalo.Context) error {
	f, err := resolveStaticFile(c)
	if err != nil {
		return c.Error(404, fmt.Errorf("not found"))
	}
	s, err := f.Stat()
	if err != nil {
		return c.Error(404, fmt.Errorf("not found"))
	}

	c.Response().Header().Add("ETag", fmt.Sprintf("%.16x", s.ModTime().UnixNano()))
	c.Response().Header().Add("Cache-Control", fmt.Sprintf("max-age=%d", assetMaxAge))
	c.Response().Header().Add("Content-Length", fmt.Sprintf("%d", s.Size()))
	c.Response().WriteHeader(200)
	return nil
}
