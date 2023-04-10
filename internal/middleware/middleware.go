package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipReader struct {
	http.ResponseWriter
	Reader io.Reader
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (r gzipReader) Read(b []byte) (int, error) {
	return r.Reader.Read(b)
}

func CompressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				return
			}
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

func DecompressGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		reader, err := gzip.NewReader(r.Body)
		if err != nil {
			_, err := io.WriteString(w, err.Error())
			if err != nil {
				return
			}
			return
		}
		defer reader.Close()

		next.ServeHTTP(gzipReader{ResponseWriter: w, Reader: reader}, r)
	})
}
