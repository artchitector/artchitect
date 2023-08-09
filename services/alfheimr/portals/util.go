package portals

import (
	"github.com/artchitector/artchitect2/model"
	"github.com/gin-gonic/gin"
	"strings"
)

func wrapError(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func makeFlatArts(arts []model.Art) []FlatArt {
	fArts := make([]FlatArt, 0, len(arts))
	for _, a := range arts {
		fArts = append(fArts, makeFlatArt(a))
	}
	return fArts
}

func makeFlatArt(art model.Art) FlatArt {
	return FlatArt{
		ID:                 art.ID,
		CreatedAt:          art.CreatedAt,
		Version:            art.Version,
		IdeaSeed:           art.Idea.Seed,
		IdeaNumberOfWords:  uint(len(art.Idea.Words)),
		IdeaWords:          strings.Split(art.Idea.WordsStr, ","),
		SeedEntropyEncoded: art.Idea.SeedEntropy.Entropy.ImageEncoded,
		SeedChoiceEncoded:  art.Idea.SeedEntropy.Choice.ImageEncoded,
	}
}
