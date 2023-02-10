package dnscommon

type Client struct {
	SecertId  string
	SecertKey string
}

type DNSRecord struct {
	Type      string
	RecordId  uint64
	SubDomain string
	Value     string
	Status    string
	Line      string
	LineId    string

	// 记录权重，用于负载均衡记录
	// 注意：此字段可能返回 null，表示取不到有效值。
	Weight *uint64

	// MX值，只有MX记录有
	// 注意：此字段可能返回 null，表示取不到有效值。
	MX *uint64 `json:"MX,omitempty" name:"MX"`

	// 记录备注说明
	Remark *string `json:"Remark,omitempty" name:"Remark"`

	// 记录缓存时间
	TTL *uint64 `json:"TTL,omitempty" name:"TTL"`
}
