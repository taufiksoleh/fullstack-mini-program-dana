package response

import (
	"fmt"
	"time"

	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
)

type BannerDTO struct {
	URL string `json:"url"`
}

type ImageDTO struct {
	URL string `json:"url"`
}

type PropertyListItem struct {
	DocumentID string    `json:"documentId"`
	Banner     BannerDTO `json:"banner"`
	Title      string    `json:"title"`
	Price      string    `json:"price"`
	CreatedAt  time.Time `json:"createdAt"`
}

type PropertyDetail struct {
	Banner      BannerDTO  `json:"banner"`
	Title       string     `json:"title"`
	DocumentID  string     `json:"documentId"`
	Description string     `json:"description"`
	Images      []ImageDTO `json:"images"`
	Facilities  []string   `json:"facilities"`
	Price       string     `json:"price"`
	Terms       string     `json:"terms"`
	Conditions  string     `json:"conditions"`
}

func ToPropertyListItem(p *entity.Property) PropertyListItem {
	return PropertyListItem{
		DocumentID: p.DocumentID,
		Banner:     BannerDTO{URL: p.BannerURL},
		Title:      p.Title,
		Price:      fmt.Sprintf("%.0f", p.Price),
		CreatedAt:  p.CreatedAt,
	}
}

func ToPropertyDetail(p *entity.Property) PropertyDetail {
	images := make([]ImageDTO, len(p.Images))
	for i, img := range p.Images {
		images[i] = ImageDTO{URL: img.URL}
	}
	facilities := make([]string, len(p.Facilities))
	for i, f := range p.Facilities {
		facilities[i] = f.Name
	}
	return PropertyDetail{
		Banner:      BannerDTO{URL: p.BannerURL},
		Title:       p.Title,
		DocumentID:  p.DocumentID,
		Description: p.Description,
		Images:      images,
		Facilities:  facilities,
		Price:       fmt.Sprintf("%.0f", p.Price),
		Terms:       p.Terms,
		Conditions:  p.Conditions,
	}
}
