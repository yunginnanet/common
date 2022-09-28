package network

import ipa "inet.af/netaddr"

/*
IterateNetRange will ingest:

  - an inet.af/netaddr.Range

  - an inet.af/netaddr.Prefix

  - or a string to be parsed as either of the above options

  - valid subnet string example: 192.168.69.0/24

  - valid range string example: 192.168.69.0-192.168.69.254

    it then returns a channel that will stream all the individual netaddr.IP types within the given range or prefix.
    if the input is invalid this function will return nil.
*/
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

	ch := make(chan ipa.IP, 254)
	go func(ret chan ipa.IP) {
		for head := addrs.From(); head != addrs.To(); head = head.Next() {
			if !head.IsUnspecified() {
				ret <- head
			}
		}
		close(ret)
	}(ch)
	return ch
}
