package event

type Listener interface {
	Notify()
}