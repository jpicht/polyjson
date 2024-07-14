package att

type Marker struct {
	Package       string
	InterfaceName string
}

func Mark(pkg, iface string) Marker {
	return Marker{pkg, iface}
}

func (m Marker) String() string {
	return m.Package + "." + m.InterfaceName
}
