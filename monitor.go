package gocache

type ChangeMonitor interface {
	HasChanged() bool
	UniqueId() string
	Dispose()
	NotifyOnChanged(f OnChangedCallback)
	OnChanged(data interface{})
}

type ChangeMonitorList struct {
}

func NewChangeMonitorList() *ChangeMonitorList {
	return &ChangeMonitorList{}
}
