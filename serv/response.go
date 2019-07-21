package serv

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/naminomare/gogutil/fileio"
)

// PostMultiFileFunc mult fileをpostされたときの対応
func PostMultiFileFunc(w http.ResponseWriter, r *http.Request, fileSaveDirectory string) {
	r.ParseForm()

	mr, err := r.MultipartReader()
	if err != nil {
		panic(err)
	}
	form, err := mr.ReadForm(r.ContentLength)
	if err != nil {
		panic(err)
	}
	for _, v := range form.File {
		for _, fheader := range v {
			f, err := fheader.Open()
			if err != nil {
				panic(err)
			}
			defer f.Close()

			dst, err := fileio.GetNonExistFileName(filepath.Join(fileSaveDirectory, fheader.Filename), 100)
			if err != nil {
				panic(err)
			}
			fhandle, err := os.OpenFile(dst, os.O_CREATE, os.ModePerm)
			if err != nil {
				panic(err)
			}
			defer fhandle.Close()

			io.Copy(fhandle, f)
		}
	}
}

// StaticResponseFunc 静的ファイル
func StaticResponseFunc(w http.ResponseWriter, r *http.Request) {
	extIndex := strings.LastIndex(r.URL.Path, ".")
	if extIndex == -1 {
	}
	ext := r.URL.Path[extIndex:]

	fhandle, err := os.OpenFile(r.URL.Path[1:], os.O_RDONLY, 0644)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404"))
		return
	}
	defer fhandle.Close()
	//debugPrint(ext)
	mimeType := mime.TypeByExtension(ext) //http.DetectContentType(bytes)
	if ext == ".ico" {
		mimeType = "image/x-icon"
	}
	info, err := fhandle.Stat()
	if err != nil {
		panic(err)
	}
	filesize := info.Size()
	w.Header().Set("Content-Length", strconv.FormatInt(filesize, 10))
	w.Header().Set("Content-Type", mimeType)

	w.WriteHeader(http.StatusOK)
	io.Copy(w, fhandle)
}
