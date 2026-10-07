package goemit

// Adapted from dop251/goja date_parser.go; native parsing only.
// Upstream revision 481fdb442bb4.
const dateParserRuntime = `
/*
Copyright (c) 2016 Dmitry Panov

Copyright (c) 2012 Robert Krimen

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the "Software"), to deal in the Software without restriction, including without limitation
the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE
WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
type tsDPDate struct {
	year, month, day     int
	hour, min, sec, msec int
	timeZoneOffset       int // time zone offset in minutes
	isLocal              bool
}

func tsDPSkip(s string, c byte) (string, bool) {
	if len(s) > 0 && s[0] == c {
		return s[1:], true
	}
	return s, false
}

func tsDPSkipSpaces(s string) string {
	for len(s) > 0 && s[0] == ' ' {
		s = s[1:]
	}
	return s
}

func tsDPSkipUntil(s string, stopList string) string {
	for len(s) > 0 && !strings.ContainsRune(stopList, rune(s[0])) {
		s = s[1:]
	}
	return s
}

func tsDPMatch(s string, lower string) (string, bool) {
	if len(s) < len(lower) {
		return s, false
	}
	for i := 0; i < len(lower); i++ {
		c1 := s[i]
		c2 := lower[i]
		if c1 != c2 {
			// switch to lower-case; 'a'-'A' is known to be a single bit
			c1 |= 'a' - 'A'
			if c1 != c2 || c1 < 'a' || c1 > 'z' {
				return s, false
			}
		}
	}
	return s[len(lower):], true
}

func tsDPGetDigits(s string, minDigits, maxDigits int) (int, string, bool) {
	var i, v int
	for i < len(s) && i < maxDigits && s[i] >= '0' && s[i] <= '9' {
		v = v*10 + int(s[i]-'0')
		i++
	}
	if i < minDigits {
		return 0, s, false
	}
	return v, s[i:], true
}

func tsDPGetMilliseconds(s string) (int, string) {
	mul, v := 100, 0
	if len(s) > 0 && (s[0] == '.' || s[0] == ',') {
		const I_START = 1
		i := I_START
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			v += int(s[i]-'0') * mul
			mul /= 10
			i++
		}
		if i > I_START {
			// only consume the separator if digits are present
			return v, s[i:]
		}
	}
	return 0, s
}

// [+-]HH:mm or [+-]HHmm or Z
func tsDPGetTimeZoneOffset(s string, strict bool) (int, string, bool) {
	if len(s) == 0 {
		return 0, s, false
	}
	sign := s[0]
	if sign == '+' || sign == '-' {
		var hh, mm, v int
		var ok bool
		t := s[1:]
		n := len(t)
		if hh, t, ok = tsDPGetDigits(t, 1, 9); !ok {
			return 0, s, false
		}
		n -= len(t)
		if strict && n != 2 && n != 4 {
			return 0, s, false
		}
		for n > 4 {
			n -= 2
			hh /= 100
		}
		if n > 2 {
			mm = hh % 100
			hh = hh / 100
		} else if t, ok = tsDPSkip(t, ':'); ok {
			if mm, t, ok = tsDPGetDigits(t, 2, 2); !ok {
				return 0, s, false
			}
		}
		if hh > 23 || mm > 59 {
			return 0, s, false
		}
		v = hh*60 + mm
		if sign == '-' {
			v = -v
		}
		return v, t, true
	} else if sign == 'Z' {
		return 0, s[1:], true
	}
	return 0, s, false
}

var tsDPTzAbbrs = []struct {
	nameLower string
	offset    int
}{
	{"gmt", 0},        // Greenwich Mean Time
	{"utc", 0},        // Coordinated Universal Time
	{"ut", 0},         // Universal Time
	{"z", 0},          // Zulu Time
	{"edt", -4 * 60},  // Eastern Daylight Time
	{"est", -5 * 60},  // Eastern Standard Time
	{"cdt", -5 * 60},  // Central Daylight Time
	{"cst", -6 * 60},  // Central Standard Time
	{"mdt", -6 * 60},  // Mountain Daylight Time
	{"mst", -7 * 60},  // Mountain Standard Time
	{"pdt", -7 * 60},  // Pacific Daylight Time
	{"pst", -8 * 60},  // Pacific Standard Time
	{"wet", +0 * 60},  // Western European Time
	{"west", +1 * 60}, // Western European Summer Time
	{"cet", +1 * 60},  // Central European Time
	{"cest", +2 * 60}, // Central European Summer Time
	{"eet", +2 * 60},  // Eastern European Time
	{"eest", +3 * 60}, // Eastern European Summer Time
}

func tsDPGetTimeZoneAbbr(s string) (int, string, bool) {
	for _, tzAbbr := range tsDPTzAbbrs {
		if s, ok := tsDPMatch(s, tzAbbr.nameLower); ok {
			return tzAbbr.offset, s, true
		}
	}
	return 0, s, false
}

var tsDPMonthNamesLower = []string{
	"jan",
	"feb",
	"mar",
	"apr",
	"may",
	"jun",
	"jul",
	"aug",
	"sep",
	"oct",
	"nov",
	"dec",
}

func tsDPGetMonth(s string) (int, string, bool) {
	for i, monthNameLower := range tsDPMonthNamesLower {
		if s, ok := tsDPMatch(s, monthNameLower); ok {
			return i + 1, s, true
		}
	}
	return 0, s, false
}

func tsDPParseDateISOString(s string) (tsDPDate, bool) {
	if strings.Contains(s,","){return tsDPDate{},false}
	if len(s) == 0 {
		return tsDPDate{}, false
	}
	var d = tsDPDate{month: 1, day: 1}
	var ok bool

	// year is either yyyy digits or [+-]yyyyyy
	sign := s[0]
	if sign == '-' || sign == '+' {
		s = s[1:]
		if d.year, s, ok = tsDPGetDigits(s, 6, 6); !ok {
			return tsDPDate{}, false
		}
		if sign == '-' {
			if d.year == 0 {
				// reject -000000
				return tsDPDate{}, false
			}
			d.year = -d.year
		}
	} else if d.year, s, ok = tsDPGetDigits(s, 4, 4); !ok {
		return tsDPDate{}, false
	}
	if s, ok = tsDPSkip(s, '-'); ok {
		if d.month, s, ok = tsDPGetDigits(s, 2, 2); !ok || d.month < 1 {
			return tsDPDate{}, false
		}
		if s, ok = tsDPSkip(s, '-'); ok {
			if d.day, s, ok = tsDPGetDigits(s, 2, 2); !ok || d.day < 1 {
				return tsDPDate{}, false
			}
		}
	}
	if s, ok = tsDPSkip(s, 'T'); ok {
		if d.hour, s, ok = tsDPGetDigits(s, 2, 2); !ok {
			return tsDPDate{}, false
		}
		if s, ok = tsDPSkip(s, ':'); !ok {
			return tsDPDate{}, false
		}
		if d.min, s, ok = tsDPGetDigits(s, 2, 2); !ok {
			return tsDPDate{}, false
		}
		if s, ok = tsDPSkip(s, ':'); ok {
			if d.sec, s, ok = tsDPGetDigits(s, 2, 2); !ok {
				return tsDPDate{}, false
			}
			d.msec, s = tsDPGetMilliseconds(s)
		}
		d.isLocal = true
	}
	// parse the time zone offset if present
	if len(s) > 0 {
		if d.timeZoneOffset, s, ok = tsDPGetTimeZoneOffset(s, true); !ok {
			return tsDPDate{}, false
		}
		d.isLocal = false
	}
	// error if extraneous characters
	return d, len(s) == 0
}

func tsDPParseDateOtherString(s string) (tsDPDate, bool) {
	var d = tsDPDate{
		year:    2001,
		month:   1,
		day:     1,
		isLocal: true,
	}
	var nums [3]int
	var numIndex int
	var hasYear, hasMon, hasTime, ok bool
	for {
		s = tsDPSkipSpaces(s)
		if len(s) == 0 {
			break
		}
		c := s[0]
		n, val := len(s), 0
		if c == '+' || c == '-' {
			if hasTime {
				if val, s, ok = tsDPGetTimeZoneOffset(s, false); ok {
					d.timeZoneOffset = val
					d.isLocal = false
				}
			}
			if !hasTime || !ok {
				s = s[1:]
				if val, s, ok = tsDPGetDigits(s, 1, 9); ok {
					d.year = val
					if c == '-' {
						if d.year == 0 {
							return tsDPDate{}, false
						}
						d.year = -d.year
					}
					hasYear = true
				}
			}
		} else if val, s, ok = tsDPGetDigits(s, 1, 9); ok {
			if s, ok = tsDPSkip(s, ':'); ok {
				// time part
				d.hour = val
				if d.min, s, ok = tsDPGetDigits(s, 1, 2); !ok {
					return tsDPDate{}, false
				}
				if s, ok = tsDPSkip(s, ':'); ok {
					if d.sec, s, ok = tsDPGetDigits(s, 1, 2); !ok {
						return tsDPDate{}, false
					}
					d.msec, s = tsDPGetMilliseconds(s)
				}
				hasTime = true
				if t := tsDPSkipSpaces(s); len(t) > 0 {
					if t, ok = tsDPMatch(t, "pm"); ok {
						if d.hour < 12 {
							d.hour += 12
						}
						s = t
						continue
					} else if t, ok = tsDPMatch(t, "am"); ok {
						if d.hour == 12 {
							d.hour = 0
						}
						s = t
						continue
					}
				}
			} else if n-len(s) > 2 {
				d.year = val
				hasYear = true
			} else if val < 1 || val > 31 {
				d.year = val
				if val < 100 {
					d.year += 1900
				}
				if val < 50 {
					d.year += 100
				}
				hasYear = true
			} else {
				if numIndex == 3 {
					return tsDPDate{}, false
				}
				nums[numIndex] = val
				numIndex++
			}
		} else if val, s, ok = tsDPGetMonth(s); ok {
			d.month = val
			hasMon = true
			s = tsDPSkipUntil(s, "0123456789 -/(")
		} else if val, s, ok = tsDPGetTimeZoneAbbr(s); ok {
			d.timeZoneOffset = val
			if len(s) > 0 {
				if c := s[0]; (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
					return tsDPDate{}, false
				}
			}
			d.isLocal = false
			continue
		} else if c == '(' {
			// tsDPSkip parenthesized phrase
			level := 1
			s = s[1:]
			for len(s) > 0 && level != 0 {
				if s[0] == '(' {
					level++
				} else if s[0] == ')' {
					level--
				}
				s = s[1:]
			}
			if level > 0 {
				return tsDPDate{}, false
			}
		} else if c == ')' {
			return tsDPDate{}, false
		} else {
			if hasYear || hasMon || hasTime || numIndex > 0 {
				return tsDPDate{}, false
			}
			// tsDPSkip a word
			s = tsDPSkipUntil(s, " -/(")
		}
		for len(s) > 0 && strings.ContainsRune("-/.,", rune(s[0])) {
			s = s[1:]
		}
	}
	n := numIndex
	if hasYear {
		n++
	}
	if hasMon {
		n++
	}
	if n > 3 {
		return tsDPDate{}, false
	}

	switch numIndex {
	case 0:
		if !hasYear {
			return tsDPDate{}, false
		}
	case 1:
		if hasMon {
			d.day = nums[0]
		} else {
			d.month = nums[0]
		}
	case 2:
		if hasYear {
			d.month = nums[0]
			d.day = nums[1]
		} else if hasMon {
			d.year = nums[1]
			if nums[1] < 100 {
				d.year += 1900
			}
			if nums[1] < 50 {
				d.year += 100
			}
			d.day = nums[0]
		} else {
			d.month = nums[0]
			d.day = nums[1]
		}
	case 3:
		d.year = nums[2]
		if nums[2] < 100 {
			d.year += 1900
		}
		if nums[2] < 50 {
			d.year += 100
		}
		d.month = nums[0]
		d.day = nums[1]
	default:
		return tsDPDate{}, false
	}
	return d, d.month > 0 && d.day > 0
}

`
