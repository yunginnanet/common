package network

import (
	"bufio"
	"strings"
	"testing"

	ipa "inet.af/netaddr"
)

var testdata29 string = `
192.168.69.240
192.168.69.241
192.168.69.242
192.168.69.243
192.168.69.244
192.168.69.245
192.168.69.246
192.168.69.247
`

var test29str = "192.168.69.240/29"

var test29rangestr = "192.168.69.240-192.168.69.247"

var test29 []*ipa.IP

func init() {
	test29 = nil
	xerox := bufio.NewScanner(strings.NewReader(testdata29))
	for xerox.Scan() {
		if line := xerox.Text(); len(line) < 1 {
			continue
		}
		ip := ipa.MustParseIP(xerox.Text())
		test29 = append(test29, &ip)
	}
}

func TestIterateNetRange(t *testing.T) {
	type args struct {
		ips interface{}
	}
	type test struct {
		name string
		args args
		want []*ipa.IP
	}

	var tests = []test{
		{
			name: "prefix",
			args: args{ips: ipa.MustParseIPPrefix(test29str)},
			want: test29,
		},
		{
			name: "range",
			args: args{ips: ipa.MustParseIPRange(test29rangestr)},
			want: test29,
		},
		{
			name: "string",
			args: args{ips: test29str},
			want: test29,
		},
		{
			name: "bogus",
			args: args{ips: "whatever, man. I'm just trynt'a vibe."},
			want: nil,
		},
	}
	for _, tt := range tests {
		index := 0
		retchan := IterateNetRange(tt.args.ips)
		if tt.want == nil && retchan != nil {
			t.Fatalf("return should have been nil, it was %v", retchan)
		}
		if retchan == nil {
			continue
		}
		t.Logf("test: %s", tt.name)
	mainloop:
		for {
			select {
			case ip := <-retchan:
				if ip.String() != test29[index].String() {
					t.Errorf("[%s] failed, wanted %s, got %s", tt.name, tt.want[index].String(), ip.String())
				} else {
					t.Logf("[%s] success (%s == %s)", tt.name, tt.want[index].String(), ip.String())
				}
				index++
			default:
				if index == 7 {
					break mainloop
				}
			}
		}
	}
}
