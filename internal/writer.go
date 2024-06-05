package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Writer interface {
	Write(r *http.Request) error
}

type stdWriter struct {
	wc io.Writer
}

func (w *stdWriter) Write(r *http.Request) error {
	if _, err := w.wc.Write(RequestFormat(r)); err != nil {
		return err
	}

	return nil
}

type fileWriter struct {
	dir os.FileInfo
}

func (w *fileWriter) Write(r *http.Request) error {
	name := filepath.Join(
		w.dir.Name(),
		fmt.Sprintf("%s_%s.log", time.Now().Format("20060102150405"), r.Context().Value(requestID)))
	return os.WriteFile(name, RequestFormat(r), 0600)
}

func RequestFormat(r *http.Request) []byte {
	buf := NewLineBuffer()
	buf.WriteString("********************************************************************************")
	buf.WriteString("RequestID: %s", r.Context().Value(requestID))
	buf.WriteString("Method: %s", r.Method)
	buf.WriteString("URL: %s", r.URL)
	buf.WriteString("Header:")
	for k, vs := range r.Header {
		buf.WriteString("  %s: %s", k, strings.Join(vs, " "))
	}
	if len(r.PostForm) > 0 {
		buf.WriteString("PostForm:")
		for k, vs := range r.PostForm {
			buf.WriteString("  %s=%s", k, strings.Join(vs, " "))
		}
	}
	if r.Body != nil {
		mediaType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
		buf.WriteString("Body: %s", mediaType)

		switch mediaType {
		case "text/plain",
			"text/csv":
			if b, err := io.ReadAll(r.Body); err == nil {
				buf.WriteString("")
				buf.Write(b)
			}
		case "application/json":
			if b, err := io.ReadAll(r.Body); err == nil {
				buf.WriteString("")
				if err := json.Indent(buf.buf, b, "", "  "); err == nil {
					buf.WriteString("")
				} else {
					buf.WriteString(string(b))
				}
			} else {
				buf.WriteString(err.Error())
			}
		default:
			buf.WriteString("Unsupported media type")
		}
	}
	buf.WriteString("")

	return buf.Bytes()
}

type LineBuffer struct {
	buf *bytes.Buffer
}

func NewLineBuffer() *LineBuffer {
	return &LineBuffer{&bytes.Buffer{}}
}

func (b *LineBuffer) WriteString(format string, a ...any) {
	b.buf.WriteString(fmt.Sprintf(format+"\n", a...))
}

func (b *LineBuffer) Write(p []byte) {
	b.buf.Write(append(p, '\n'))
}

func (b *LineBuffer) String() string {
	return b.buf.String()
}

func (b *LineBuffer) Bytes() []byte {
	return b.buf.Bytes()
}
