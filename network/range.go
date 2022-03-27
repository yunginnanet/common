package network

import ipa "inet.af/netaddr"

// IterateNetRange will ingest either a netaddr range or a netaddr prefix from the inet.af/netaddr package;
// returning a channel that will stream all the individual netaddr IPs within the given range or prefix.
// Alternatively, feed it a string in prefix or range format. (192.168.69.0/24) (192.168.69.0-192.168.69.254)
// Will return nil value if input is invalid.
func IterateNetRange(ips interface{}) chan ipa.IP {
	var addrs ipa.IPRange

	switch ips.(type) {
	case string:
		strefix, prefixErr := ipa.ParseIPPrefix(ips.(string))
		strange, rangeErr := ipa.ParseIPRange(ips.(string))
		switch {
		case rangeErr == nil:
			addrs = strange
		case prefixErr == nil:
			addrs = strefix.Range()
		default:
			return nil
		}
	case ipa.IPRange:
		addrs = ips.(ipa.IPRange)
	case ipa.IPPrefix:
		addrs = ips.(ipa.IPPrefix).Range()
	default:
		return nil
	}

	ch := make(chan ipa.IP)
	go func(ret chan ipa.IP) {
		for head := addrs.From(); head != addrs.To(); head = head.Next() {
			if !head.IsUnspecified() {
				ret <- head
			}
		}
	}(ch)
	return ch
}
