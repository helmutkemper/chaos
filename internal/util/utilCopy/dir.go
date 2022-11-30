package utilCopy

import (
	"io/ioutil"
	"os"
	"path"
)

// Dir copies a whole directory recursively
func Dir(dst string, src string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = Dir(dstFp, srcFp); err != nil {
				return err
			}
		} else {
			if err = File(dstFp, srcFp); err != nil {
				return err
			}
		}
	}
	return nil
}
