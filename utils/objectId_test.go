package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"time"
)

func TestObjectID_Hex(t *testing.T) {
	id := NewObjectID()
	hex := id.Hex()
	t.Log(hex)
	fromHex, err := ObjectIDFromHex(hex)
	require.Equal(t, nil, err)
	assert.Equal(t, id, fromHex)
	assert.Equal(t, false, id.IsZero())
	assert.Equal(t, id.String(), fmt.Sprintf("%+v", id))

	marshal, err := json.Marshal(id)
	require.Equal(t, nil, err)
	assert.Equal(t, hex, strings.Trim(string(marshal), "\""))
	var nid ObjectID
	err = json.Unmarshal(marshal, &nid)
	require.Equal(t, nil, err)
	assert.Equal(t, id, nid)

	marshal2, err := json.Marshal(map[string]interface{}{"$oid": id})
	require.Equal(t, nil, err)
	err = json.Unmarshal(marshal2, &nid)
	require.Equal(t, nil, err)
	assert.Equal(t, id, nid)
}

func TestObjectID_Timestamp(t *testing.T) {
	now := time.Now()
	id := NewObjectIDFromTimestamp(now)
	assert.Equal(t, now.Unix(), id.Timestamp().Unix())
}

func BenchmarkObjectID_Hex(b *testing.B) {
	b.Log(NewObjectID().Hex())
}
