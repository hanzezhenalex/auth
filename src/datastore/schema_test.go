package datastore

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestScopes(t *testing.T) {
	rq := require.New(t)
	origin := "\"" + "scope1" + delimiter + "scope2" + delimiter + "scope3" + "\"\n"
	scopes := Scopes{
		"scope1", "scope2", "scope3",
	}

	t.Run("Marshal", func(t *testing.T) {
		var p []byte
		buf := bytes.NewBuffer(p)
		rq.NoError(json.NewEncoder(buf).Encode(&scopes))
		rq.Equal(origin, buf.String())
	})

	t.Run("Unmarshal", func(t *testing.T) {
		var _scopes Scopes
		rq.NoError(json.NewDecoder(bytes.NewBuffer([]byte(origin))).Decode(&_scopes))
		rq.EqualValues(scopes, _scopes)
	})

}
