package portals

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"net/http"
)

type LikePortal struct {
	authService authService
	likePile    likePile
}

func NewLikePortal(authService authService, likePile likePile) *LikePortal {
	return &LikePortal{authService: authService, likePile: likePile}
}

type LikedRequest struct {
	ArtID uint `uri:"art_id" json:"art_id" binding:"required,numeric"`
}

func (lp *LikePortal) HandleLikedList(c *gin.Context) {
	userId := lp.authService.GetUserIdFromContext(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	likes, err := lp.likePile.GetList(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	}

	liked := make([]uint, 0)
	for _, like := range likes {
		liked = append(liked, like.ArtID)
	}

	c.JSON(http.StatusOK, liked)
}

func (lp *LikePortal) HandleLikedArt(c *gin.Context) {
	userId := lp.authService.GetUserIdFromContext(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	var req LikedRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	like, err := lp.likePile.Get(c, userId, req.ArtID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{"liked": false})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"liked": true, "liked_at": like.CreatedAt})
		return
	}
}

func (lp *LikePortal) HandleLike(c *gin.Context) {
	userId := lp.authService.GetUserIdFromContext(c)
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not authorized"})
		return
	}

	var req LikedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, wrapError(err))
		return
	}

	var newLiked bool
	_, err := lp.likePile.Get(c, userId, req.ArtID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		newLiked = true
	} else {
		newLiked = false
	}

	err = lp.likePile.Set(c, userId, req.ArtID, newLiked)
	if err != nil {
		c.JSON(http.StatusInternalServerError, wrapError(err))
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"liked": newLiked})
	}
}
