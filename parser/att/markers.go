package att

type Marker struct {
	Package       string
	Name          string
	InterfaceName string
}

func Mark(pkg, name, iface string) Marker {
	return Marker{pkg, name, iface}
}

func (m Marker) String() string {
	return m.Name
}
