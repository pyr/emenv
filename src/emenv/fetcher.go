package emenv

import (
	"archive/tar"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func (env *Env) FetchFilePackage(p InstallDef, body []byte) error {

	dir := fmt.Sprintf("%s/%s-%s", env.PackageDir, p.Name, p.Version)
	file := fmt.Sprintf("%s/%s.el", dir, p.Name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, body, 0644); err != nil {
		return err
	}
	return nil
}

func ReadTarEntry(hdr *tar.Header, rdr *tar.Reader) ([]byte, error) {

	remaining := hdr.Size
	outbuf := make([]byte, 0)
	for {

		b := make([]byte, remaining)
		br, err := rdr.Read(b)
		if err == io.EOF {
			return outbuf, nil
		}
		if err != nil {
			return outbuf, err
		}
		b = b[0:br]
		outbuf = append(outbuf, b...)
		remaining = remaining - int64(br)
	}
	return outbuf, UnreachableError
}

func (env *Env) FetchTarPackage(p InstallDef, rdr *tar.Reader) error {

	dir := fmt.Sprintf("%s/%s-%s", env.PackageDir, p.Name, p.Version)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for {
		hdr, err := rdr.Next()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch {
		case (hdr.Typeflag == tar.TypeReg || hdr.Typeflag == tar.TypeRegA):

			b, err := ReadTarEntry(hdr, rdr)
			if err != nil {
				return err
			}
			fpath := fmt.Sprintf("%s/%s", env.PackageDir, hdr.Name)
			dir := path.Dir(fpath)
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
			if err = ioutil.WriteFile(fpath, b, 0644); err != nil {
				return err
			}
			break
		case (hdr.Typeflag == tar.TypeDir):
			path := fmt.Sprintf("%s/%s", env.PackageDir, hdr.Name)
			if err = os.MkdirAll(path, 0755); err != nil {
				return err
			}
			break
		default:
			fmt.Printf("unhandled entry type: %c\n", hdr.Typeflag)
		}
	}
}

func (env *Env) FetchPackage(idef InstallDef) error {

	fmt.Printf("fetching from: %s\n", idef.URL)
	resp, err := http.Get(idef.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch {
	case idef.StoreType == FileStorage:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		err = env.FetchFilePackage(idef, body)
		if err != nil {
			return err
		}
	case idef.StoreType == TarStorage:
		rdr := tar.NewReader(resp.Body)
		err = env.FetchTarPackage(idef, rdr)
		if err != nil {
			return err
		}
	default:
		return UnreachableError
	}

	return nil
}
