//  Copyright 2022 Blockdaemon Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pyth

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttrsMap(t *testing.T) {
	caseMap := map[string]string{
		"1pythians":   "are",
		"2incredibly": "based",
	}
	attrs, err := NewAttrsMap(caseMap)
	require.NoError(t, err)
	assert.Equal(t, [][2]string{
		{"1pythians", "are"},
		{"2incredibly", "based"},
	}, attrs.Pairs)
	assert.Equal(t, caseMap, attrs.KVs())

	buf, err := attrs.MarshalBinary()
	require.NoError(t, err)
	assert.Equal(t, []byte(
		"\x09"+"1pythians"+"\x03"+"are"+
			"\x0b"+"2incredibly"+"\x05"+"based",
	), buf)

	var attrs2 AttrsMap
	require.NoError(t, attrs2.UnmarshalBinary(buf))
	assert.Equal(t, attrs, attrs2)
}

func TestAttrsMap_LongKey(t *testing.T) {
	longKey := strings.Repeat("A", 256)
	caseMap := map[string]string{
		longKey: ":)",
	}
	attrs, err := NewAttrsMap(caseMap)
	assert.EqualError(t, err, `key too long (256 > 0xFF): "`+longKey+`"`)
	assert.Len(t, attrs.Pairs, 0)
	assert.Len(t, attrs.KVs(), 0)
}

func TestAttrsMap_LongValue(t *testing.T) {
	caseMap := map[string]string{
		"bla": strings.Repeat("A", 256),
	}
	attrs, err := NewAttrsMap(caseMap)
	assert.EqualError(t, err, `value too long (256 > 0xFF): "`+caseMap["bla"]+`"`)
	assert.Len(t, attrs.Pairs, 0)
	assert.Len(t, attrs.KVs(), 0)
}
