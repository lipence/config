package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fileUtils "github.com/lipence/utils/file"
)

var (
	ErrLoaderUndefined = fmt.Errorf("loader undefined")
	ErrPathUndefined   = fmt.Errorf("path undefined")
	ErrPathNotFound    = fmt.Errorf("path not found")
	ErrPathNoContent   = fmt.Errorf("path has no content")
	ErrPathUnavailable = fmt.Errorf("path unavailable")
	ErrPathNotFile     = fmt.Errorf("only accepts single file, but got directory")
)

type OnLoaderRegister interface {
	OnRegister() (err error)
}

var loader Loader

func Use(l Loader) error {
	loader = l
	if evt, ok := loader.(OnLoaderRegister); ok {
		if err := evt.OnRegister(); err != nil {
			return err
		}
	}
	return nil
}

func wrapErrorWithPath(err error, path string) error {
	return fmt.Errorf("%w (path: %s)", err, path)
}

func wrapErrorWithMsgAndPath(err error, path string, msg error) error {
	return fmt.Errorf("%w: %v (path: %s)", msg, err, path)
}

func LoadFromPath(path string) (instance Value, err error) {
	var _loader = loader
	if _loader == nil {
		return nil, ErrLoaderUndefined
	}
	// empty config path
	if path = filepath.Clean(strings.TrimSpace(path)); path == "" {
		return nil, ErrPathUndefined
	}
	// check config status
	var pathStat os.FileInfo
	if path, pathStat, err = fileUtils.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, wrapErrorWithMsgAndPath(err, path, ErrPathNotFound)
		}
		return nil, wrapErrorWithMsgAndPath(err, path, ErrPathUnavailable)
	}
	var files = map[string][]byte{}
	if pathStat.IsDir() {
		if !_loader.AllowDir() {
			return nil, wrapErrorWithPath(ErrPathNotFile, path)
		}
		var pathList []string
		if pathList, err = fileUtils.List(path, _loader.PathPattern(), true); err != nil {
			return nil, wrapErrorWithPath(err, path)
		}
		if len(pathList) == 0 {
			return nil, wrapErrorWithPath(ErrPathNoContent, path)
		}
		for _, subPath := range pathList {
			if files[subPath], err = os.ReadFile(subPath); err != nil {
				return nil, wrapErrorWithPath(err, subPath)
			}
		}
	} else {
		if files[path], err = os.ReadFile(path); err != nil {
			return nil, wrapErrorWithPath(err, path)
		} else if files[path] = bytes.TrimSpace(files[path]); len(files[path]) == 0 {
			return nil, wrapErrorWithPath(ErrPathNoContent, path)
		}
	}
	if instance, err = _loader.Load(path, files); err != nil {
		return nil, wrapErrorWithPath(err, path)
	}
	return instance, nil
}

func LoadConfigs(configPath string) error {
	instance, err := LoadFromPath(configPath)
	if err != nil {
		return err
	}
	configRWLocker.Lock()
	configInstance = instance
	configRWLocker.Unlock()
	return nil
}
