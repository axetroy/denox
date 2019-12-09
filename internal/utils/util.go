package utils

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/cheggaaa/pb/v3"
	"github.com/pkg/errors"
)

// Download file from URL to the filepath
func DownloadFile(filepath string, url string) (*pb.ProgressBar, error) {
	tmpl := fmt.Sprintf(`{{string . "prefix"}}{{ green "%s" }} {{counters . }} {{ bar . "[" "=" ">" "-" "]"}} {{percent . }} {{speed . }}{{string . "suffix"}}`, filepath)

	// Get the data
	response, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrapf(err, "Download `%s` fail", url)
	}

	if response.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(fmt.Sprintf("download file with status code %d", response.StatusCode))
	}

	defer response.Body.Close()

	// Create the file
	writer, err := os.Create(filepath)

	if err != nil {
		return nil, errors.Wrapf(err, "Create `%s` fail", filepath)
	}

	defer func() {
		err = writer.Close()

		if err != nil {
			_ = os.Remove(filepath)
		}
	}()

	bar := pb.ProgressBarTemplate(tmpl).Start64(response.ContentLength)

	bar.SetWriter(os.Stdout)

	barReader := bar.NewProxyReader(response.Body)

	_, err = io.Copy(writer, barReader)

	bar.Finish()

	if err != nil {
		err = errors.Wrap(err, "copy fail")
	}

	return bar, err
}

func DecompressZip(tarFile, dest string) (*string, error) {
	r, err := zip.OpenReader(tarFile)

	if err != nil {
		return nil, errors.Wrapf(err, "read zip file `%s` fail", tarFile)
	}

	defer r.Close()

	if len(r.File) > 1 {
		return nil, errors.New("window .zip file should only contain single file")
	}

	f := r.File[0]

	newFilepath := path.Join(dest, f.Name)

	src, err := f.Open()

	if err != nil {
		return nil, err
	}

	defer src.Close()

	dst, err := os.Create(newFilepath)

	if err != nil {
		return nil, err
	}

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	return &newFilepath, nil
}

func DecompressGz(tarFile, dest string) (*string, error) {
	fileReader, err := os.Open(tarFile)

	if err != nil {
		return nil, errors.Wrapf(err, "open file `%s` fail", tarFile)
	}

	defer fileReader.Close()

	gzipReader, err := gzip.NewReader(fileReader)

	if err != nil {
		return nil, errors.Wrapf(err, "gzip decode fail")
	}

	defer gzipReader.Close()

	newFilepath := path.Join(dest, "deno")

	fileWriter, err := os.Create(newFilepath)

	if err != nil {
		return nil, errors.Wrapf(err, "create file `%s` fail", newFilepath)
	}

	defer func() {
		err = fileWriter.Close()

		if err != nil {
			err = os.Remove(newFilepath)
		}
	}()

	if _, err = io.Copy(fileWriter, gzipReader); err != nil {
		return nil, err
	}

	if err := fileWriter.Chmod(os.FileMode(0755)); err != nil {
		return nil, errors.Wrap(err, "change file mod fail")
	}

	return &newFilepath, nil
}

// Decompress gzip file and return filepath
func Decompress(tarFile, dest string) (*string, error) {
	if path.Ext(tarFile) == ".zip" {
		return DecompressZip(tarFile, dest)
	} else {
		return DecompressGz(tarFile, dest)
	}
}
