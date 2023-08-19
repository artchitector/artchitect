package portals

import (
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"net/http"
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
	Unity    model.Unity     `json:"unity"`
	Arts     []model.FlatArt `json:"arts"`
	Children []model.Unity   `json:"children"`
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

	response := unityResponse{
		Unity:    parent,
		Arts:     []model.FlatArt{},
		Children: []model.Unity{},
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
