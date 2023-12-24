package portals

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type ArtRequest struct {
	ID uint `uri:"id" binding:"required,numeric"`
}

type LastRequest struct {
	Last uint `uri:"last" binding:"required,numeric"`
}

// ArtPortal - канал связи, по которому Мидгард получает данные картин
type ArtPortal struct {
	artPile artPile
	harbour harbour
}

func NewArtPortal(artPile artPile, harbour harbour) *ArtPortal {
	return &ArtPortal{artPile: artPile, harbour: harbour}
}

func (ap *ArtPortal) HandleArt(c *gin.Context) {
	var request ArtRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}
	if request.ID < 1 {
		c.JSON(http.StatusBadRequest, wrapError(errors.Errorf("ID must be positive")))
		return
	}

	art, err := ap.artPile.GetArtRecursive(c, request.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, wrapError(errors.Errorf("not found ID=%d", request.ID)))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	c.JSON(http.StatusOK, art)
}

func (ap *ArtPortal) HandleArtFlat(c *gin.Context) {
	var request ArtRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}
	if request.ID < 1 {
		c.JSON(http.StatusBadRequest, wrapError(errors.Errorf("ID must be positive")))
		return
	}

	art, err := ap.artPile.GetArtRecursive(c, request.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, wrapError(errors.Errorf("not found ID=%d", request.ID)))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	c.JSON(http.StatusOK, model.MakeFlatArt(art))
}

func (ap *ArtPortal) HandleLast(c *gin.Context) {
	var request LastRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}
	if request.Last < 1 || request.Last > 100 {
		c.JSON(http.StatusBadRequest, wrapError(errors.Errorf("Active must be 0-100")))
		return
	}

	arts, err := ap.artPile.GetLastArts(c, request.Last)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	flatArts := model.MakeFlatArts(arts)
	c.JSON(http.StatusOK, flatArts)
}

func (ap *ArtPortal) HandleChosen(c *gin.Context) {
	cargo, err := ap.harbour.SendCrownWaitCargo(c, model.RequestGiveChosenArt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}
	var artID uint
	if err := json.Unmarshal([]byte(cargo), &artID); err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	art, err := ap.artPile.GetArtRecursive(c, artID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, "not found")
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	c.JSON(http.StatusOK, model.MakeFlatArt(art))
}

func (ap *ArtPortal) HandleMaxArt(c *gin.Context) {
	maxArtID, err := ap.artPile.GetMaxArtID(c)
	if err != nil {
		err = fmt.Errorf("[art_portal] failed to GetMaxArtID: %w", err)
		log.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, maxArtID)
}
