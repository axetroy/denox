package deno

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/axetroy/denox/internal/fs"
	"github.com/axetroy/denox/internal/utils"
	"github.com/pkg/errors"
)

var (
	ErrNotSupport = errors.New("not support your platform")
)

type Os string
type Arch string

const (
	OsLinux   Os = "linux"
	OsWindows Os = "win"
	OsOsx     Os = "osx"

	Archx64 Arch = "x64"
)

type Deno struct {
	Version  string
	Os       Os
	Arch     Arch
	cacheDir string
	DenoDir  string
}

func New(version *string) (*Deno, error) {
	if version == nil {
		if v, err := getLatestVersion(); err != nil {
			return nil, err
		} else {
			version = v
		}
	}

	denoOs, err := getOs()

	if err != nil {
		return nil, err
	}

	denoArch, err := getArch()

	if err != nil {
		return nil, err
	}

	cacheDir, err := getCacheDir()

	if err != nil {
		return nil, err
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return nil, err
	}

	DenoDir := ""

	if s := os.Getenv("DENO_DIR"); s != "" {
		DenoDir = s
	} else {
		DenoDir = path.Join(homeDir, ".denox", "deno_"+*version)
	}

	if err := fs.EnsureDir(path.Join(DenoDir, "bin")); err != nil {
		return nil, err
	}

	return &Deno{
		Os:       *denoOs,
		Arch:     *denoArch,
		Version:  *version,
		cacheDir: cacheDir,
		DenoDir:  DenoDir,
	}, nil
}

// clear download cache dir
func (d *Deno) Clean() error {
	return os.RemoveAll(d.cacheDir)
}

// download Deno from remote and returns the path of the executable file
func (d *Deno) Download() (executablePath string, err error) {
	var (
		tarExtName        = ".gz"
		remoteTarFilename string
		localTarFilename  string
		dstDir            = path.Join(d.DenoDir, "bin")
	)

	executablePath = path.Join(dstDir, "deno")

	if d.Os == OsWindows {
		tarExtName = ".zip"
		executablePath += ".exe"
	}

	remoteTarFilename = fmt.Sprintf("deno_%s_%s%s", d.Os, d.Arch, tarExtName)
	localTarFilename = fmt.Sprintf("deno_%s_%s_%s%s", d.Version, d.Os, d.Arch, tarExtName)

	localGzFilepath := path.Join(d.cacheDir, localTarFilename)

	// download
	downloadURL := fmt.Sprintf("https://github.com/denoland/deno/releases/download/%s/%s", d.Version, remoteTarFilename)

	// if Deno executable file not exist
	// then we should download it
	if exit, err := fs.PathExists(executablePath); err != nil {
		return "", errors.Wrapf(err, "stat file `%s` fail", executablePath)
	} else if !exit {
		// download the file for current platform
		if _, err = utils.DownloadFile(localGzFilepath, downloadURL); err != nil {
			return "", errors.Wrap(err, "download file fail")
		} else {

		}

		if _, err := utils.Decompress(localGzFilepath, dstDir); err != nil {
			return "", errors.Wrapf(err, "decompress gzip file fail")
		} else {
			// make sure is it is a executable file
			if d.Os != OsWindows {
				mod := os.FileMode(0755)
				if err := os.Chmod(executablePath, mod); err != nil {
					err = errors.Wrap(err, "set permission fail")
				}
			}
		}
	}

	return executablePath, nil
}

// get deno OS
func getOs() (*Os, error) {
	var denoOS Os

	switch runtime.GOOS {
	case "darwin":
		denoOS = OsOsx
		break
	case "linux":
		denoOS = OsLinux
		break
	case "windows":
		denoOS = OsWindows
		break
	default:
		return nil, ErrNotSupport
	}

	return &denoOS, nil
}

// get deno Arch
func getArch() (*Arch, error) {
	var denoArch Arch

	switch runtime.GOARCH {
	case "amd64":
		denoArch = Archx64
		break
	default:
		return nil, ErrNotSupport
	}

	return &denoArch, nil
}

// get latest version of Deno from remote
func getLatestVersion() (*string, error) {
	r, err := http.Get("https://denolib.github.io/setup-deno/release.json")

	if err != nil {
		return nil, err
	}

	if r.StatusCode >= http.StatusBadRequest {
		return nil, errors.New(r.Status)
	}

	b, err := ioutil.ReadAll(r.Body)

	if err != nil {
		return nil, errors.Wrap(err, "read body fail")
	}

	type Response struct {
		Name string `json:"name"`
	}

	var res []Response

	if err := json.Unmarshal(b, &res); err != nil {
		return nil, errors.Wrap(err, "parse JSON fail")
	}

	if res == nil || len(res) == 0 {
		return nil, errors.New("can not found version")
	}

	latestVersion := res[0].Name

	return &latestVersion, nil
}

// get cache dir for deno
func getCacheDir() (string, error) {
	if userCacheDir, err := os.UserCacheDir(); err != nil {
		return "", errors.Wrap(err, "get user cache dir fail")
	} else {
		denoXCacheDir := path.Join(userCacheDir, "denox")

		if err = fs.EnsureDir(denoXCacheDir); err != nil {
			return "", errors.Wrap(err, "ensure denox cache dir fail")
		}

		return denoXCacheDir, nil
	}
}
