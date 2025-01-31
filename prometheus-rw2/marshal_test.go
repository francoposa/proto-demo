// Copyright 2023 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package writev2

import (
	"testing"

	"github.com/CrowdStrike/csproto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var a proto.Message

//func TestMarshalUnmarshalZeroValuesConsistency(t *testing.T) {
//	for _, tt := range []struct {
//		name string
//		msg  *TimeSeries
//	}{
//		{
//			name: "nonempty-zero-value-flat-member", // passes
//			msg: &TimeSeries{
//				CreatedTimestamp: 0,
//			},
//		},
//		{
//			name: "nonempty-zero-value-nested-member", // fails
//			msg: &TimeSeries{
//				Samples: []*Sample{
//					{
//						Value:     0,
//						Timestamp: 0,
//					},
//				},
//			},
//		},
//	} {
//		t.Run(tt.name, func(t *testing.T) {
//			// upstream behavior
//			msgCopy := proto.Clone(tt.msg).(*TimeSeries)
//			upstreamProtoMarshaledBytes, err := proto.Marshal(msgCopy)
//			assert.NoError(t, err)
//
//			unmarshaledFromUpstreamBytes := &TimeSeries{}
//			assert.NoError(t, proto.Unmarshal(upstreamProtoMarshaledBytes, unmarshaledFromUpstreamBytes))
//			assert.EqualExportedValues(t, msgCopy, unmarshaledFromUpstreamBytes)
//
//			// csproto behavior
//			csProtoMarshaledBytes, err := tt.msg.Marshal()
//			assert.NoError(t, err)
//
//			unmarshaledFromCSProtoBytes := &TimeSeries{}
//			assert.NoError(t, unmarshaledFromCSProtoBytes.Unmarshal(csProtoMarshaledBytes))
//			// assert marshal->unmarshal consistency
//			assert.EqualExportedValues(t, tt.msg, unmarshaledFromCSProtoBytes) // fails on nested zero values
//
//			// compare behavior
//			// assert marshal behavior against upstream
//			assert.Equal(t, upstreamProtoMarshaledBytes, csProtoMarshaledBytes) // fails on nested zero values
//			// assert marshal-unmarshal behavior against upstream
//			assert.EqualExportedValues(t, unmarshaledFromUpstreamBytes, unmarshaledFromCSProtoBytes) // fails on nested zero values
//		})
//	}
//}

func TestMarshalZeroValues(t *testing.T) {
	for _, tt := range []struct {
		name   string
		msg    csproto.Marshaler
		panics bool
	}{
		{
			name: "nonempty-zero-value-single-field-member",
			// does not panic because total size for entire message reports as zero;
			// there is no buffer and no serialization is performed
			msg: &TimeSeries{
				Samples: []*Sample{
					{
						Value:     0,
						Timestamp: 0,
					},
				},
			},
			panics: false,
		},
		{
			name: "nonempty-zero-value-single-nested-member-with-member-before",
			// panics because total size for each Sample message is report as zero
			// so bytes were not accounted for in the initial buffer size calculation,
			// but EncodeNested will still encode the tag and length for each Sample message;
			// buffer overflows when the Sample tag is written to it
			msg: &TimeSeries{
				LabelsRefs: []uint32{0, 1},
				Samples: []*Sample{
					{
						Value:     0,
						Timestamp: 0,
					},
				},
			},
			panics: true,
		},
		{
			name: "nonempty-zero-value-single-nested-member-with-member-after",
			// panics because total size for each Sample message is report as zero
			// so bytes were not accounted for in the initial buffer size calculation,
			// but EncodeNested will still encode the tag and length for each Sample message;
			// buffer overflows when the next (Metadata) field is written to it
			msg: &TimeSeries{
				Samples: []*Sample{
					{
						Value:     0,
						Timestamp: 0,
					},
				},
				Metadata: &Metadata{
					Type: Metadata_METRIC_TYPE_UNSPECIFIED,
				},
			},
			panics: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// csproto marshal
			msg := tt.msg.(proto.Message)
			if tt.panics {
				assert.Panics(t, func() {
					_, _ = proto.Marshal(msg)
				})
			} else {
				_, err := proto.Marshal(msg)
				assert.NoError(t, err)
			}
		})
	}
}
