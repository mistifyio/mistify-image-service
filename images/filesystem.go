package images

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

type (
	// FS is an image store using the filesystem
	FS struct {
		Config *FSConfig
	}

	// FSConfig contains necessary config options to set up the fs store
	FSConfig struct {
		Dir string
	}
)

// fsLogFields contain fields to include on all logs
var fsLogFields = log.Fields{
	"type":  "images",
	"store": "fs",
}

// Validate checks whether the config is valid
func (fsc *FSConfig) Validate() error {
	if fsc.Dir == "" {
		return errors.New("empty dir path")
	}
	return nil
}

// Init parses the config and ensures the directory exists
func (fs *FS) Init(configBytes []byte) error {
	config := &FSConfig{}

	// Parse the config json
	if err := json.Unmarshal(configBytes, config); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error": err,
			"json":  string(configBytes),
		}).Error("failed to unmarshal config json")
		return err
	}

	if err := config.Validate(); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error": err,
		}).Error("failed config validation")
		return err
	}

	fs.Config = config
	log.WithFields(fsLogFields).WithFields(log.Fields{
		"config": fs.Config,
	}).Info("config loaded")

	// Create the directory (if necessary)
	mode := os.FileMode(0755)
	if err := os.MkdirAll(fs.Config.Dir, mode); err != nil && !os.IsExist(err) {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error": err,
			"dir":   fs.Config.Dir,
		}).Error("failed to create directory")
		return err
	}

	// Make sure permissions are correct
	if err := os.Chmod(config.Dir, mode); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":       err,
			"dir":         config.Dir,
			"desiredMode": mode,
		}).Error("failed to set directory permissions")
		return err
	}

	return nil
}

// Shutdown is a noop
func (fs *FS) Shutdown() error {
	return nil
}

// filepath generates the full filepath from the image id
func (fs *FS) filepath(imageID string) string {
	return filepath.Join(fs.Config.Dir, imageID)
}

// Stat retrieves file information about an image
func (fs *FS) Stat(imageID string) (os.FileInfo, error) {
	filepath := fs.filepath(imageID)
	info, err := os.Stat(filepath)
	if err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to stat image")
	}
	return info, err
}

// Get retrieves an image from the filesystem
func (fs *FS) Get(imageID string, out io.Writer) error {
	filepath := fs.filepath(imageID)
	file, err := os.Open(filepath)
	if err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to open image")
		return err
	}
	defer file.Close()

	if _, err := io.Copy(out, file); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to copy image data to output stream")
		return err
	}

	return nil
}

// Put stores an image in the filesystem
func (fs *FS) Put(imageID string, in io.Reader) error {
	filepath := fs.filepath(imageID)
	file, err := os.Create(filepath)
	if err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to create image file")
		return err
	}

	mode := os.FileMode(0755)
	if err := file.Chmod(mode); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":       err,
			"imageID":     imageID,
			"filepath":    filepath,
			"desiredMode": mode,
		}).Error("failed to chmod image file")
		return err
	}

	if _, err := io.Copy(file, in); err != nil {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to create image file")
		return err
	}
	return nil
}

// Delete removes an image from the filesystem
func (fs *FS) Delete(imageID string) error {
	filepath := fs.filepath(imageID)
	if err := os.Remove(filepath); err != nil && !os.IsNotExist(err) {
		log.WithFields(fsLogFields).WithFields(log.Fields{
			"error":    err,
			"imageID":  imageID,
			"filepath": filepath,
		}).Error("failed to remove image")
		return err
	}
	return nil
}

func init() {
	Register("fs", func() Store {
		return &FS{}
	})
}
