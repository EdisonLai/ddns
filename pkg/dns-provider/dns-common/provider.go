package dnscommon

import "context"

type Provider interface {
	NewClient(ctx context.Context, secretId, secretKey string) (Client, error)
	DescribeDomain(ctx context.Context, client Client, domain string, subdomainfilter []string) (records []DNSRecord, err error)
	ModifyDynamicDNS(ctx context.Context, client Client, recordId uint64, domain, subdomain, lineId, value string) (err error)
}
