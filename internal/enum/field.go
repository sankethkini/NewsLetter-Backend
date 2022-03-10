package enum

type Field int

const (
	PRICE Field = iota
	DAYS
)

func (f Field) String() string {
	return []string{"price", "days"}[f]
}
