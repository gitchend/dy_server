package tools

import (
	"bufio"
	"io"
	"os"
)

func ReadLine(path string, fn func(line string) error) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()

	buf := bufio.NewReader(fp)
	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		err = fn(string(line))
		if err != nil {
			return err
		}
	}
	return nil
}

func FileIsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
