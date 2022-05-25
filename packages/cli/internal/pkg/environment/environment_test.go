package environment

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsesWesAdapterConsistency(t *testing.T) {
	for _, engine := range AllEngines {
		_, found := UsesWesAdapter[engine]
		assert.True(t, found)
	}
	assert.Equal(t, len(UsesWesAdapter), len(AllEngines))
}

func TestDefaultRepositoriesConsistency(t *testing.T) {
	for _, component := range AllComponents {
		_, found := DefaultRepositories[component]
		assert.True(t, found)
	}
	assert.Equal(t, len(DefaultRepositories), len(AllComponents))
}

func TestDefaultTagsConsistency(t *testing.T) {
	for _, component := range AllComponents {
		_, found := DefaultTags[component]
		assert.True(t, found)
	}
	assert.Equal(t, len(DefaultTags), len(AllComponents))
}

func TestCommonImagesConsistency(t *testing.T) {
	for _, component := range AllComponents {
		_, found := CommonImages[component]
		assert.True(t, found)
	}
	assert.Equal(t, len(CommonImages), len(AllComponents))
}
