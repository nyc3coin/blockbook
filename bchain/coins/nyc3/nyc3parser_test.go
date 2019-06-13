// +build unittest

package nyc3

import (
	"blockbook/bchain/coins/btc"
	"encoding/hex"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
)

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func Test_GetAddrDescFromAddress_Mainnet(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "P2PKH1",
			args:    args{address: "NSRgg8eXnJnQEyyjAm8SGMAAXYsdqxSt5v"},
			want:    "76a914476f7379f044adc34cfcce5b2eba7c015815180788ac",
			wantErr: false,
		}
	}
	parser := NewNYC3Parser(GetChainParams("main"), &btc.Configuration{})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.GetAddrDescFromAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAddrDescFromAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("GetAddrDescFromAddress() = %v, want %v", h, tt.want)
			}
		})
	}
}
