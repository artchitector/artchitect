package portals

import (
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

type ArtImageRequest struct {
	//ID   uint   `uri:"id" binding:"required,numeric"`
	//Size string `uri:"size" binding:"required"`
	Name string `uri:"name" binding:"required"`
}

type UnityImageRequest struct {
	Mask    string `uri:"mask" binding:"required"`
	Version uint   `uri:"version" binding:"numeric"`
	Size    string `uri:"size" binding:"required"`
}

type ImagePortal struct {
	warehouse warehouse
}

func NewImagePortal(warehouse warehouse) *ImagePortal {
	return &ImagePortal{warehouse: warehouse}
}

func (ip *ImagePortal) HandleArtImage(c *gin.Context) {
	var request ArtImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	var err error
	id, size, err := ip.parseName(request.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	var img []byte
	if size == model.SizeOrigin {
		img, err = ip.warehouse.DownloadArtOrigin(c, id)
	} else {
		img, err = ip.warehouse.DownloadArtImage(c, id, size)
	}

	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}
	c.Data(http.StatusOK, "image/jpeg", img)
}

func (ip *ImagePortal) HandleUnityImage(c *gin.Context) {
	var request UnityImageRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	img, err := ip.warehouse.DownloadUnityImage(c, request.Mask, request.Version, request.Size)
	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	c.Data(http.StatusOK, "image/jpeg", img)
}

func (ip *ImagePortal) parseName(name string) (uint, string, error) {
	// name в формате artchitect-ID-SIZE, например, artchitect-19500-f, artchitect-800-s...
	pieces := strings.Split(name, "-")
	if len(pieces) != 3 {
		return 0, "", errors.Errorf("имя файла должно быть в формате artchitect-<id>-<size>")
	}
	if pieces[0] != "artchitect" {
		return 0, "", errors.Errorf("имя файла должно быть в формате artchitect-<id>-<size>")
	}
	id, err := strconv.ParseUint(pieces[1], 10, 64)
	if err != nil {
		return 0, "", err
	}
	if id == 0 || id > 1000000000 {
		// не может быть миллиард картин. artchitect вообще рассчитан на 1 миллион, тут с запасом в x1000
		return 0, "", errors.Errorf("id должен быть от 1 до 1млрд")
	}

	size := pieces[2]
	if size != model.SizeF && size != model.SizeM && size != model.SizeS && size != model.SizeXS && size != model.SizeOrigin {
		return 0, "", errors.Errorf("некорректный размер картинки")
	}

	return uint(id), size, nil
}
