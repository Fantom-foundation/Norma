package driver

type ServiceID int

type ServiceInfo struct {
	Name string
	Port uint16
}

var servicesByName = map[string]ServiceID{}
var servicesById = map[ServiceID]ServiceInfo{}

func RegisterService(name string, defaultPort uint16) ServiceID {
	id, ok := servicesByName[name]
	if ok {
		return id
	}
	info := ServiceInfo{
		Name: name,
		Port: defaultPort,
	}
	nextId := ServiceID(len(servicesByName))
	servicesByName[name] = nextId
	servicesById[nextId] = info
	return nextId
}

func GetServiceInfo(service ServiceID) *ServiceInfo {
	info, ok := servicesById[service]
	if !ok {
		return nil
	}
	return &info
}
