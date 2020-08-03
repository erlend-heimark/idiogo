package extdadjokes

import (
	"testing"

	"github.com/erlend-heimark/idiogo/envvars"
	"github.com/stretchr/testify/require"
)

func TestFetcher_GetRandom(t *testing.T) {
	e := envvars.New()
	f := NewFetcher(e.DadJokeURL)
	joke, err := f.GetRandom()
	require.NoError(t, err)
	require.NotNil(t, joke)
	require.NotNil(t, joke.Joke)
}
