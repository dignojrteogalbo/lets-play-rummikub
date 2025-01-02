package event

type Listener interface {
	Notify(messages ...string)
}
