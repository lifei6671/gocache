package gocache

type ChangeMonitor interface {
	HasChanged() bool
	UniqueId() string
	Dispose()
	NotifyOnChanged(f OnChangedCallback)
	OnChanged(data interface{})
}

type OnChangedCallback func()
