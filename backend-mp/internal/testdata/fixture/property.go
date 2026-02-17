package fixture

import (
	"time"

	"github.com/google/uuid"
	"github.com/taufiksoleh/backend-mp/internal/domain/entity"
)

var (
	PropertyID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	FacilityID = uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")
	ImageID    = uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd")
	FixedTime  = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
)

func Property() *entity.Property {
	return &entity.Property{
		ID:          PropertyID,
		DocumentID:  "doc-001",
		Title:       "Test Property",
		Description: "A nice property",
		BannerURL:   "https://example.com/banner.jpg",
		Price:       500000,
		Terms:       "Monthly",
		Conditions:  "Good condition",
		Images: []entity.PropertyImage{
			{ID: ImageID, PropertyID: PropertyID, URL: "https://example.com/image.jpg", CreatedAt: FixedTime, UpdatedAt: FixedTime},
		},
		Facilities: []entity.Facility{
			{ID: FacilityID, Name: "Swimming Pool", CreatedAt: FixedTime, UpdatedAt: FixedTime},
		},
		CreatedAt: FixedTime,
		UpdatedAt: FixedTime,
	}
}

func PropertyList(n int) []entity.Property {
	list := make([]entity.Property, n)
	for i := range list {
		list[i] = entity.Property{
			ID:         uuid.New(),
			DocumentID: "doc-" + string(rune('0'+i+1)),
			Title:      "Property " + string(rune('A'+i)),
			Price:      float64((i + 1) * 100000),
			CreatedAt:  FixedTime,
			UpdatedAt:  FixedTime,
		}
	}
	return list
}
