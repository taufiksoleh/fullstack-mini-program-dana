package request

import "github.com/taufiksoleh/backend-mp/internal/usecase"

type PropertyBody struct {
	DocumentID    string   `json:"document_id"`
	Title         string   `json:"title"          binding:"required"`
	Description   string   `json:"description"`
	BannerURL     string   `json:"banner_url"`
	Price         float64  `json:"price"          binding:"required"`
	Terms         string   `json:"terms"`
	Conditions    string   `json:"conditions"`
	ImageURLs     []string `json:"image_urls"`
	FacilityNames []string `json:"facility_names"`
}

func (b PropertyBody) ToUsecaseRequest() usecase.PropertyRequest {
	return usecase.PropertyRequest{
		DocumentID:    b.DocumentID,
		Title:         b.Title,
		Description:   b.Description,
		BannerURL:     b.BannerURL,
		Price:         b.Price,
		Terms:         b.Terms,
		Conditions:    b.Conditions,
		ImageURLs:     b.ImageURLs,
		FacilityNames: b.FacilityNames,
	}
}
