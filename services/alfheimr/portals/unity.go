package portals

import (
	"net/http"
	"time"

	"github.com/artchitector/artchitect/model"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UnityPortal struct {
	unityPile unityPile
	artPile   artPile
}

func NewUnityPortal(unityPile unityPile, artPile artPile) *UnityPortal {
	return &UnityPortal{unityPile: unityPile, artPile: artPile}
}

type unityRequest struct {
	Mask string `uri:"mask" binding:"required"`
}

type unityResponse struct {
	Unity        model.Unity     `json:"unity"`
	FirstArtTime *time.Time      `json:"firstArtTime"`
	LastArtTime  *time.Time      `json:"lastArtTime"`
	Arts         []model.FlatArt `json:"arts"`
	Children     []model.Unity   `json:"children"`
}

func (up *UnityPortal) HandleMain(c *gin.Context) {
	unities, err := up.unityPile.GetRoot(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
	}
	c.JSON(http.StatusOK, unities)
}

func (up *UnityPortal) HandleUnity(c *gin.Context) {
	var request unityRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	parent, err := up.unityPile.Get(c, request.Mask)
	if err != nil {
		c.JSON(http.StatusNotFound, wrapError(err))
		return
	}

	firstArtTime, lastArtTime, err := up.getFirstLastTimes(c, parent)
	response := unityResponse{
		Unity:        parent,
		Arts:         []model.FlatArt{},
		Children:     []model.Unity{},
		FirstArtTime: firstArtTime,
		LastArtTime:  lastArtTime,
	}

	if parent.Rank == model.Unity100 {
		// Для сотенного единства нужно получить список всех его карточек
		arts, err := up.artPile.GetArtsInterval(c, parent.MinID, parent.MaxID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, wrapError(err))
			return
		}
		flatArts := model.MakeFlatArts(arts)
		response.Arts = flatArts

	} else {
		// Для остальных единств нужно получать список дочерних единств
		children, err := up.unityPile.GetChildren(c, parent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, wrapError(err))
			return
		}
		response.Children = children
	}

	c.JSON(http.StatusOK, response)
}

func (up *UnityPortal) getFirstLastTimes(c *gin.Context, parent model.Unity) (*time.Time, *time.Time, error) {
	var firstArtTime, lastArtTime *time.Time
	minID := parent.MinID
	if minID == 0 {
		minID = 1
	}
	first, err := up.artPile.GetArt(c, minID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errors.Wrapf(err, "[unity_portal] ОШИБКА")
	} else if err == nil {
		firstArtTime = &first.CreatedAt
	}

	last, err := up.artPile.GetArt(c, parent.MaxID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, errors.Wrapf(err, "[unity_portal] ОШИБКА")
	} else if err == nil {
		lastArtTime = &last.CreatedAt
	}
	return firstArtTime, lastArtTime, nil
}
