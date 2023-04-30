package main

import (
	"log"
	"os"
	"server/pkg/database"
	"server/pkg/model"
	"strconv"
	"strings"

	geojson "github.com/paulmach/go.geojson"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalln("usage: ./mktraildb [input file]")
	}

	if err := database.Init(); err != nil {
		log.Fatalf("failed to initialize database: %v\n", err)
	}

	data, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("failed to open file: %v\n", err)
	}

	collection, err := geojson.UnmarshalFeatureCollection(data)
	if err != nil {
		log.Fatalf("failed to unmarshal feature collection: %v\n", err)
		os.Exit(1)
	}

	for _, feature := range collection.Features {
		if source, err := feature.PropertyString("source"); err == nil {
			if strings.Contains(source, "Muntii Nostri gpx") && feature.Geometry.IsLineString() {
				length, _ := strconv.ParseFloat(feature.PropertyMustString("distance", "0"), 64)
				trail := &model.Trail{
					Name:   feature.PropertyMustString("name", ""),
					From:   feature.PropertyMustString("from", ""),
					To:     feature.PropertyMustString("to", ""),
					Length: length,
					Symbol: feature.PropertyMustString("osmc:symbol", ""),
				}

				for _, coords := range feature.Geometry.LineString {
					point := &model.Point{
						Lat:  coords[0],
						Long: coords[1],
					}
					trail.Points = append(trail.Points, point)
				}

				database.Instance.Create(trail)
			}
		}
	}
}