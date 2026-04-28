package maxmind

import (
	"context"
	"fmt"
	"hash/crc64"
	"io"

	"pkg/database/maxmind/model"
	"pkg/errors"
	"pkg/log"

	"github.com/oschwald/maxminddb-golang"
)

const logStep = 500_000

func (mm *MaxMind) makeLocations(ctx context.Context, d *mmData, locations []string) error {
	networks := d.reader.Networks(maxminddb.SkipAliasedNetworks)

	locMap := make(map[string]struct{}, len(locations))
	for _, value := range locations {
		locMap[value] = struct{}{}
	}

	tmpMap := make(map[model.Location]*model.Location, 1000)

	var count int

	log.Info("loading geo data: started")

	for networks.Next() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var geodata model.Geodata
		if _, err := networks.Network(&geodata); err != nil {
			return errors.Default.Wrap(err)
		}

		if loc, iso := geodata.GetLocation(false); loc != nil {
			if _, ok := locMap[iso]; ok {
				tmpMap[*loc] = loc
			}
		}

		count++
		if count%logStep == 0 {
			log.Info(fmt.Sprintf("loading geo data: processed %d records", count))
		}
	}

	log.Info(fmt.Sprintf("loading geo data: done, processed %d records", count))

	hash := crc64.New(crc64.MakeTable(crc64.ECMA))
	paths := make(map[string]uint64, len(tmpMap))
	taxonomy := make(map[uint64]*model.Taxonomy, len(tmpMap))

	for _, value := range tmpMap {
		if (value.CountryRu == "" && value.CountryEn == "") ||
			((value.RegionRu == "" && value.RegionEn == "") && (value.CityRu != "" && value.CityEn != "")) {
			continue
		}

		var prevHashSum uint64

		for _, name := range value.Names() {
			if _, err := io.WriteString(hash, name[2]); err != nil {
				return errors.Default.Wrap(err)
			}

			hashSum := hash.Sum64()
			hash.Reset()

			if prevHashSum != 0 {
				taxonomy[prevHashSum].ParentId = hashSum
			}

			taxonomy[hashSum] = &model.Taxonomy{
				NameRu: name[0],
				NameEn: name[1],
			}

			paths[name[2]] = hashSum
			prevHashSum = hashSum
		}
	}

	for id, t := range taxonomy {
		if t.ParentId == 0 {
			taxonomy[id].Type = model.TaxonomyTypeCountry
		} else if taxonomy[t.ParentId] != nil && taxonomy[t.ParentId].ParentId == 0 {
			taxonomy[id].Type = model.TaxonomyTypeRegion
		} else {
			taxonomy[id].Type = model.TaxonomyTypeCity
		}
	}

	d.paths = paths
	d.taxonomy = taxonomy

	return nil
}
