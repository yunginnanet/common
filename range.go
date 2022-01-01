package common

import ipa "inet.af/netaddr"

func RangeIterate(ips interface{}) chan *ipa.IP {
	var addrs ipa.IPRange

	switch ips.(type) {
	case ipa.IPRange:
		addrs = ips.(ipa.IPRange)
	case ipa.IPPrefix:
		addrs = ips.(ipa.IPPrefix).Range()
	}

	ch := make(chan *ipa.IP)
	go func() {
		var head ipa.IP
		head = addrs.From()
		end := addrs.To()
		for head != end {
			if !head.IsUnspecified() {
				ch <- &head
			}
			head = head.Next()
		}
		close(ch)
	}()
	return ch
}
