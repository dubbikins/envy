package types

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type EnvironmentReader struct {
	Source string
	Values map[string]string 
}

//UnmarshalText reads data, if data is a file path that it can open, it reads the file, expecting each line to be a key value pair in the form <key>=<value>/n
//Otherwise, it attempts to extract the raw kvs from data, otherwise it returns an error
func (r *EnvironmentReader) UnmarshalText(data []byte) (err error) {
	if len(data) == 0 {
		return
	}
	if data[0] == '@' {
		for i, kv := range os.Environ() {
			_kv := strings.SplitN(kv, "=", 2)
			if len(_kv) != 2 {
				return fmt.Errorf("error parsing os environmen: %v has len %d but expected 2", kv[i], len(_kv))
			}
			r.Values[_kv[0]]= _kv[1]
		}
		return
	}
	if r.Values == nil {
		r.Values = make(map[string]string)
	}
	slog.Debug("Unmarshalling Environment Reader", "source", data)
	var reader io.Reader
	if reader, err = os.Open(string(data)); err != nil {
		reader = bytes.NewReader(data)
	}else {
		r.Source = string(data)
	}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if err = scanner.Err(); err == io.EOF {
			return nil
		}else if err != nil {
			return
		}
		kv := strings.SplitN(scanner.Text(), "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("invalid format '%s', expected <key>=<value>", scanner.Text())
		}
		r.Values[kv[0]] = kv[1]
	}
	return
}

func (r *EnvironmentReader) MarshalText() ([]byte, error) {
	return []byte(r.Source), nil
}