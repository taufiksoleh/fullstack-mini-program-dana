package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	pgmodel "github.com/taufiksoleh/backend-mp/internal/repository/postgres/model"
	"github.com/taufiksoleh/backend-mp/pkg/config"
	"github.com/taufiksoleh/backend-mp/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	adjectives = []string{
		"Luxury", "Modern", "Cozy", "Spacious", "Elegant",
		"Affordable", "Prime", "Exclusive", "Charming", "Strategic",
	}
	propTypes = []string{
		"Villa", "Apartment", "House", "Townhouse", "Studio",
		"Loft", "Penthouse", "Cottage", "Shophouse", "Ruko",
	}
	locations = []string{
		"Bali", "Jakarta Selatan", "Surabaya", "Bandung", "Yogyakarta",
		"Medan", "Semarang", "Makassar", "Depok", "Tangerang",
	}
	descriptions = []string{
		"A stunning property located in a prime area with easy access to schools, shops, and public transport.",
		"Well-maintained property surrounded by a quiet residential neighborhood. Perfect for families.",
		"Modern design with high-quality finishes throughout. Ready for immediate occupancy.",
		"Strategically located near commercial centers. Ideal for investment or personal use.",
		"A peaceful retreat offering privacy and comfort in a secure gated community.",
		"Bright and airy interiors with open-plan living spaces. Close to major highways and public transit.",
		"Nestled in a lush green area, this property offers a serene environment away from the city noise.",
		"Newly renovated with premium materials. Equipped with modern kitchen and bathroom fittings.",
		"Corner unit with panoramic city views. High-rise building with full amenities on every floor.",
		"Freehold property in a rapidly developing area with high investment potential.",
	}
	terms = []string{
		"Payment can be made in installments over 12–36 months. Bank financing available.",
		"Full payment required within 30 days of signing. KPR and cash hard accepted.",
		"Flexible payment terms available. Contact agent for custom payment arrangements.",
		"Down payment 30%, remainder within 90 days. Notarial deed issued upon full payment.",
		"KPR BTN and BCA available. Minimum down payment 20%. Maximum tenor 20 years.",
	}
	conditions = []string{
		"Property is sold as-is. Buyers are encouraged to conduct their own inspection.",
		"Certificate of ownership (SHM) is clean and free of encumbrances. Taxes are up to date.",
		"Building permit (IMB) and certificate available. No outstanding utilities bills.",
		"AJB and SHM ready. All legal documents completed. Ready for transfer.",
		"PBB paid up to date. No liens or encumbrances. Suitable for KPR application.",
	}
	facilityNames = []string{
		"Swimming Pool", "Gym", "Parking Area", "24-Hour Security",
		"WiFi", "CCTV", "Playground", "Garden", "Jogging Track",
		"Clubhouse", "BBQ Area", "Laundry Room", "Storage Room",
		"Elevator", "Power Backup",
	}
)

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cfg := config.Load()
	db, err := database.NewPostgres(cfg)
	if err != nil {
		slog.Error("seed: failed to connect to database", "err", err)
		os.Exit(1)
	}

	facilities := seedFacilities(db)
	seedProperties(db, rng, facilities)

	slog.Info("seed: completed successfully")
}

func seedFacilities(db *gorm.DB) []pgmodel.Facility {
	facilities := make([]pgmodel.Facility, len(facilityNames))
	for i, name := range facilityNames {
		f := pgmodel.Facility{Name: name}
		db.Clauses(clause.OnConflict{DoNothing: true}).FirstOrCreate(&f, pgmodel.Facility{Name: name})
		facilities[i] = f
	}
	slog.Info("seed: facilities ready", "count", len(facilities))
	return facilities
}

func seedProperties(db *gorm.DB, rng *rand.Rand, facilities []pgmodel.Facility) {
	inserted := 0
	for i := 1; i <= 100; i++ {
		docID := fmt.Sprintf("PROP-%03d", i)

		var count int64
		db.Model(&pgmodel.Property{}).Where("document_id = ?", docID).Count(&count)
		if count > 0 {
			continue
		}

		title := fmt.Sprintf("%s %s in %s",
			adjectives[rng.Intn(len(adjectives))],
			propTypes[rng.Intn(len(propTypes))],
			locations[rng.Intn(len(locations))],
		)
		price := float64(500_000_000 + rng.Intn(14_500)*1_000_000)

		imgCount := 2 + rng.Intn(4)
		images := make([]pgmodel.PropertyImage, imgCount)
		for j := 0; j < imgCount; j++ {
			images[j] = pgmodel.PropertyImage{
				ID:  uuid.New(),
				URL: fmt.Sprintf("https://picsum.photos/seed/%d-%d/800/600", i, j),
			}
		}

		facCount := 2 + rng.Intn(3)
		perm := rng.Perm(len(facilities))
		propFacilities := make([]pgmodel.Facility, facCount)
		for k := 0; k < facCount; k++ {
			propFacilities[k] = facilities[perm[k]]
		}

		prop := pgmodel.Property{
			ID:          uuid.New(),
			DocumentID:  docID,
			Title:       title,
			Description: descriptions[rng.Intn(len(descriptions))],
			BannerURL:   fmt.Sprintf("https://picsum.photos/seed/%d/800/600", i),
			Price:       price,
			Terms:       terms[rng.Intn(len(terms))],
			Conditions:  conditions[rng.Intn(len(conditions))],
			Images:      images,
			Facilities:  propFacilities,
		}

		if err := db.Create(&prop).Error; err != nil {
			slog.Error("seed: failed to insert property", "documentId", docID, "err", err)
			continue
		}
		inserted++
	}
	slog.Info("seed: properties inserted", "count", inserted)
}
