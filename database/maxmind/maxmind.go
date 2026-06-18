package maxmind

import (
	"context"
	"sync/atomic"

	"pkg/database/maxmind/model"
	"pkg/errors"

	"github.com/oschwald/maxminddb-golang"
)

type mmData struct {
	reader   *maxminddb.Reader
	taxonomy map[uint64]*model.Taxonomy
	paths    map[string]uint64
}

type MaxMind struct {
	file string
	data atomic.Pointer[mmData]
}

func New(file string) *MaxMind {
	return &MaxMind{
		file: file,
	}
}

func (mm *MaxMind) Run(ctx context.Context, countries []string) error {
	mmReader, err := maxminddb.Open(mm.file)
	if err != nil {
		return errors.Default.Wrap(err)
	}

	d := &mmData{reader: mmReader}

	if len(countries) != 0 {
		if err = mm.makeLocations(ctx, d, countries); err != nil {
			return err
		}
	}

	mm.data.Store(d)
	return nil
}

func (mm *MaxMind) GetTaxonomy() map[uint64]*model.Taxonomy {
	if d := mm.data.Load(); d != nil {
		return d.taxonomy
	}
	return nil
}
