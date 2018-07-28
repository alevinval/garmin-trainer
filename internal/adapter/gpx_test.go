package adapter

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/alevinval/trainer/internal/trainer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadGpx(t *testing.T) {
	data, err := ioutil.ReadFile("testdata/sample.gpx")
	require.Nil(t, err)

	g, err := Gpx(data)
	require.Nil(t, err)

	assert.Equal(t, "", g.Metadata().DataSource.Name)
	assert.Equal(t, trainer.DataSourceType(""), g.Metadata().DataSource.Type)
	assert.Equal(t, "Some activity name", g.Metadata().Name)
	assert.Equal(t, "2015-01-20 13:26:30 +0000 UTC", g.Metadata().Time.String())

	assert.Equal(t, 3, len(g.DataPoints()))

	// Compare first data point
	expectedTime, _ := time.Parse("2006-01-02T15:04:05.000Z", "2017-06-19T16:49:35.000Z")
	expected := &trainer.DataPoint{
		Time: expectedTime,
		Coords: trainer.Point{
			Lat:       1,
			Lon:       1,
			Elevation: 100,
		},
		Hr:  94,
		Cad: 170,
		N:   1,
	}
	assert.Equal(t, expected, g.DataPoints()[0])
}
