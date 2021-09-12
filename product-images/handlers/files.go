package handlers

import (
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/shaun7pan/building-microservices-gin/product-images/files"
)

type Files struct {
	log   hclog.Logger
	store files.Storage
}

type params struct {
	ID       int    `uri:"id" binding:"required"`
	FileName string `uri:"filename" binding:"required"`
}

func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{l, s}
}

func (f *Files) SaveFile(c *gin.Context) {

	params := params{}

	err := c.ShouldBindUri(&params)
	if err != nil {
		f.log.Error("Invalid path", "uri", params)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err.Error(),
		})
		return
	}

	f.log.Info("Handle POST request", "id", params.ID, "filename", params.FileName)

	idstr := strconv.Itoa(params.ID)

	f.saveFile(idstr, params.FileName, c)
}

func (f *Files) saveFile(id, path string, c *gin.Context) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fullpath := filepath.Join(id, path)
	err := f.store.Save(fullpath, c.Request.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Message": err.Error(),
		})
		return
	}
}
