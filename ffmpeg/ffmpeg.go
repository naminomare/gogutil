package ffmpeg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/naminomare/gogutil/fileio"
)

// VCodec vcodec
type VCodec string

var (
	// VCodecPNG png
	VCodecPNG VCodec = "png"
)

// Client クライアント
type Client struct {
	exePath         string
	outputFiles     []string
	outputDirectory string
}

// NewClient クライアントを返す
func NewClient(exepath string) *Client {
	return &Client{
		exePath: exepath,
	}
}

// PictOutput 画像出力
func (t *Client) PictOutput(
	movieFile string,
	outputDir string,
	vcodec VCodec,
	rate float64,
	endtime int,
	fileprefix string,
) error {
	outputDirectory, err := fileio.GetNonExistFileName(filepath.Join(outputDir, "tmp"), 1000)
	if err != nil {
		return err
	}

	// rateが1秒当たりの画像枚数なので、endtime時間でわかる
	// expectedOutputFilesNum := int(rate * float64(endtime))
	os.MkdirAll(outputDirectory, os.ModePerm)

	cmd := exec.Command(
		t.exePath,
		"-y",
		"-i", movieFile,
		"-r", fmt.Sprint(rate),
		"-vcodec", string(vcodec),
		"-t", strconv.Itoa(endtime),
		outputDirectory+"/"+fileprefix+"%05d.png",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	finfos, err := ioutil.ReadDir(outputDirectory)
	for _, finfo := range finfos {
		f, _ := filepath.Abs(filepath.Join(outputDirectory, finfo.Name()))
		t.outputFiles = append(t.outputFiles, f)
	}

	return nil
}

// EachFiles 出力したファイルに対して処理する
func (t *Client) EachFiles(fn func(outputfilepath string) interface{}) []interface{} {
	ret := make([]interface{}, len(t.outputFiles))
	for i, f := range t.outputFiles {
		ret[i] = fn(f)
	}
	return ret
}

// OutputFileDeleteAll 出力したファイルを消しておく
func (t *Client) OutputFileDeleteAll() {
	os.RemoveAll(t.outputDirectory)
}

// GetOutputFilePath ファイルパスのスライスを返す
func (t *Client) GetOutputFilePath() []string {
	return t.outputFiles
}
