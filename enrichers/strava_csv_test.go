package enrichers

import (
	"testing"

	"github.com/alevinval/trainer"
	"github.com/stretchr/testify/assert"
)

func TestStravaCsvEnrichFailOpen(t *testing.T) {
	_, err := NewStravaCsvEnricher("../tests/missing-file.csv")
	assert.NotNil(t, err)
}

func TestStravaCsvEnrichActivity(t *testing.T) {
	enricher, err := NewStravaCsvEnricher("../tests/activities.csv")
	assert.Nil(t, err)

	a := &trainer.Activity{}
	a.SetMetadata(&trainer.Metadata{
		Name: "",
		DataSource: trainer.DataSource{
			Type: trainer.FileDataSource,
			Name: "783319746.fit.gz",
		},
	})

	enricher.Enrich(a)

	assert.Equal(t, "Lunch Run", a.Metadata().Name)
}