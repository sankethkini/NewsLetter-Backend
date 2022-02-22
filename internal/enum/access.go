package enum

type Access int

const (
	ADMIN Access = iota
	USER
)

func (r Access) String() string {
	return []string{"admin", "user"}[r]
}
