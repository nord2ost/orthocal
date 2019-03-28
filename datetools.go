/*
	Orthocal provides tools for the Orthodox calendar.

	Most calculations are done with the Julian Day Number, which we abbreviate
	to JDN or the distance from Pascha in days, which we abbreviate pdist.
*/

package orthocal

import (
	"errors"
	"time"
)

const (
	Sunday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

const (
	NoFast = iota
	Fast
	LentenFast
	ApostlesFast
	DormitionFast
	NativityFast
)

var FastLevels = map[int]string{
	NoFast:        "Посту немає",
	Fast:          "Піст",
	LentenFast:    "Великий піст",
	ApostlesFast:  "Петрівський піст",
	DormitionFast: "Успенський піст",
	NativityFast:  "Різдвяний піст",
}

var FastExceptions = map[int]string{
	1:  "Дозволено вино та олію",
	2:  "Дозволено рибу, вино та олію",
	3:  "Дозволено вино та олію",
	4:  "Дозволено рибу, вино та олію",
	5:  "Дозволено вино",
	6:  "Дозволено вино, ікру та олію",
	7:  "Не споживаємо м'яса",
	8:  "Суворий піст (вино та олія)",
	9:  "Суворий піст",
	10: "Піст є",
	11: "Посту немає",
}

var FeastLevels = map[int]string{
	-1: "Немає літургії",
	0:  "Літургія",
	1:  "Ранішесвячених дарів",
	2:  "Black squigg (6-stich typikon symbol)",
	3:  "Red squigg (doxology typikon symbol)",
	4:  "Red cross (polyeleos typikon symbol)",
	5:  "Red cross half-circle (vigil typikon symbol)",
	6:  "Red cross circle (great feast typikon symbol)",
	7:  "Major feast Theotokos",
	8:  "Major feast Lord",
}

// Pascha functions

// Compute the Julian date of Pascha for the given year.
func ComputeJulianPascha(year int) (int, int) {
	// Use the Meeus Julian algorithm to calculate the Julian date
	// See https://en.wikipedia.org/wiki/Computus#Meeus'_Julian_algorithm
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7
	month := (d + e + 114) / 31
	day := (d+e+114)%31 + 1
	return month, day
}

// Compute the Julian day number of Pascha for the given year.
func ComputePaschaJDN(year int) int {
	month, day := ComputeJulianPascha(year)
	return JulianDateToJDN(year, month, day)
}

// Compute the Gregorian date of Pascha for the given year.
// The year must be between 2001 and 2099.
func ComputeGregorianPascha(year int) (time.Time, error) {
	month, day := ComputeJulianPascha(year)

	gregorianDate, e := JulianToGregorian(year, month, day)
	if e != nil {
		return time.Now(), e
	}

	return gregorianDate, nil
}

// Compute the distance of a given day from Pascha. Returns the distance and the year.
// If the distance is < -77, the returned year will be earlier than the one passed in.
func ComputePaschaDistance(year, month, day int) (int, int) {
	JDN := GregorianDateToJDN(year, month, day)
	distance := JDN - ComputePaschaJDN(year)

	if distance < -77 {
		year--
		distance = JDN - ComputePaschaJDN(year)
	}

	return distance, year
}

// Compute the distance of a given day from Pascha. Returns the distance and the year.
// If the distance is < -77, the returned year will be earlier than the one passed in.
func ComputeJulianPaschaDistance(year, month, day int) (int, int) {
	JDN := JulianDateToJDN(year, month, day)
	distance := JDN - ComputePaschaJDN(year)

	if distance < -77 {
		year--
		distance = JDN - ComputePaschaJDN(year)
	}

	return distance, year
}

// Return the day of the week given the distance from Pascha.
func WeekDayFromPDist(distance int) int {
	return (7 + distance%7) % 7
}

func SurroundingWeekends(distance int) (int, int, int, int) {
	weekday := WeekDayFromPDist(distance)

	saturdayBefore := distance - weekday - 1
	sundayBefore := distance - 7 + ((7 - weekday) % 7)
	saturdayAfter := distance + 7 - ((weekday + 1) % 7)
	sundayAfter := distance + 7 - weekday

	return saturdayBefore, sundayBefore, saturdayAfter, sundayAfter
}

// Conversion functions

// Convert a Julian date to a Gregorian date.
func JulianToGregorian(year, month, day int) (time.Time, error) {
	// This function will be incorrect outside the range 2001-2099 for 2 reasons:
	//
	// 1. The offset of 13 is incorrect outside the range 1900-2099.
	// 2. If the Julian date is in February and on a year that is divisible by
	//    100, the Go time module will incorrectly add the offset because these years
	//    are leap years on the Julian, but not on the Gregorian.
	//
	// Hopefully this code will no longer be running by 2100.

	if year < 2001 || year > 2099 {
		return time.Now(), errors.New("The year must be between 1900 and 2099")
	}

	// Add an offset of 13 to convert from Julian to Gregorian
	julianDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return julianDate.AddDate(0, 0, 13), nil
}

// Convert a Gregorian date to a Julian date.
func GregorianToJulian(year, month, day int) (time.Time, error) {
	// This function will be incorrect outside the range 2001-2099 for 2 reasons:
	//
	// 1. The offset of 13 is incorrect outside the range 1900-2099.
	// 2. If the Julian date is in February and on a year that is divisible by
	//    100, the Go time module will incorrectly add the offset because these years
	//    are leap years on the Julian, but not on the Gregorian.
	//
	// Hopefully this code will no longer be running by 2100.

	if year < 2001 || year > 2099 {
		return time.Now(), errors.New("The year must be between 1900 and 2099")
	}

	// Add an offset of 13 to convert from Julian to Gregorian
	gregorianDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return gregorianDate.AddDate(0, 0, -13), nil
}

// Convert a Julian date to a Julian day number.
func JulianDateToJDN(year, month, day int) int {
	// See https://en.wikipedia.org/wiki/Julian_day#Converting_Julian_calendar_date_to_Julian_Day_Number
	return 367*year - (7*(year+5001+(month-9)/7))/4 + (275*month)/9 + day + 1729777
}

// Convert a Gregorian date to a Julian day number.
// This function mimic's PHP's gregoriantojd().
func GregorianDateToJDN(year, month, day int) int {
	if month > 2 {
		month -= 3
	} else {
		month += 9
		year--
	}

	// break up the year into the leftmost 2 digits (century) and the rightmost 2 digits
	century := year / 100
	ya := year - 100*century

	return (146097*century)/4 + (1461*ya)/4 + (153*int(month)+2)/5 + day + 1721119
}
