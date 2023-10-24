package http2

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEbpfTx_Path(t *testing.T) {
	type fields struct {
		Tup            connTuple
		Request_method uint8
		Path_size      uint8
		Request_path   [160]uint8
	}
	type args struct {
		buffer []byte
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []byte
		want1  bool
	}{
		{
			name: "test empty string",
			fields: fields{
				Request_method: 1,
				Path_size:      0,
				Request_path:   [160]uint8{},
			},
		},
	},
		for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &EbpfTx{
				Request_method: tt.fields.Request_method,
				Path_size:      tt.fields.Path_size,
				Request_path:   tt.fields.Request_path,
			}
			got, got1 := tx.Path(tt.args.buffer)
			assert.Equalf(t, tt.want, got, "Path(%v)", tt.args.buffer)
			assert.Equalf(t, tt.want1, got1, "Path(%v)", tt.args.buffer)
		})
	}
}
