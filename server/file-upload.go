package main

import (
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// IFileManager exposes interfaces to manage upload files.
type IFileManager interface {
	Put(pid int64, name string, r io.Reader) error
	Delete(pid int64, name string) error
	List(pid int64) ([]string, error)
}

// FileUpload operates on upload files
type FileUpload struct {
	mgr IFileManager
}

// NewFileUpload returns a new instance of FileUpload
func NewFileUpload(mgr IFileManager) *FileUpload {
	return &FileUpload{
		mgr: mgr,
	}
}

// Upload does file saving
func (o *FileUpload) Upload(c *gin.Context) error {
	var err error

	parent := toInt64(c.Param("parent"))
	name := c.Param("name")

	if err = o.mgr.Put(parent, name, c.Request.Body); err != nil {
		return err
	}

	return nil
}

// List does file listing
func (o *FileUpload) List(c *gin.Context) ([]string, error) {
	parent := toInt64(c.Param("parent"))
	return o.mgr.List(parent)
}

// Delete does file deleting
func (o *FileUpload) Delete(c *gin.Context) error {
	parent := toInt64(c.Param("parent"))
	name := c.Param("name")

	if name == "" {
		return fmt.Errorf("bad name")
	}

	return o.mgr.Delete(parent, name)
}
