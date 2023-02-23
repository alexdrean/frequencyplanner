package snmp

type OidKind int

const (
	OidGet OidKind = iota
	OidCount
)

type Oid struct {
	Oid  string
	Kind OidKind
}

type OidResult struct {
	Oid    Oid
	Result string
}
