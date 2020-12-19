package enum

/*
type WeekDay int

const (
	Sun WeekDay = iota
	Mon
	Tue
	Wed
	Thu
	Fri
	Sat
)

func (w WeekDay) String() string {
	switch w {
	case Sun:
		return "Sun"
	case Mon:
		return "Mon"
	case Tue:
		return "Tue"
	case Wed:
		return "Wed"
	case Thu:
		return "Thu"
	case Fri:
		return "Fri"
	case Sat:
		return "Sat"
	default:
		return "Unknown"
	}
}
*/

type Month int

const (
	_ Month = iota
	Jan
	Feb
	Mar
	Apr
	May
	Jun
	Jul
	Aug
	Sep
	Oct
	Nov
	Dec
)

func (m Month) String() string {
	switch m {
	case Jan:
		return "Jan"
	case Feb:
		return "Feb"
	case Mar:
		return "Mar"
	case Apr:
		return "Apr"
	case May:
		return "May"
	case Jul:
		return "Jul"
	case Jun:
		return "Jun"
	case Aug:
		return "Aug"
	case Sep:
		return "Sep"
	case Oct:
		return "Oct"
	case Nov:
		return "Nov"
	case Dec:
		return "Dec"
	default:
		return "Unknown"
	}
}
