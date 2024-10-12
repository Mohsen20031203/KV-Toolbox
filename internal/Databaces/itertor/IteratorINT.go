package itertor

type IterDB interface {
	Next() bool
	Key() string
	First() bool
	Value() string
	Prev() bool
}
