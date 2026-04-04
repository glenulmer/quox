package globals

import (
//	"errors"
    "fmt"
	"regexp"
    "strconv"
    "strings"
)

type Dec2_t int

func (d Dec2_t)String() string {
    isNeg := d < 0
    num := abs(int(d))
    
    decim := num % 100
    integ := num / 100
    
    decimStr := fmt.Sprintf("%02d", decim)
    integStr := fmt.Sprintf("%d", integ)

    length := len(integStr)
    if length > 3 {
        var result []byte
        for i := 0; i < length; i++ {
            if i > 0 && (length-i)%3 == 0 {
                result = append(result, '.')
            }
            result = append(result, integStr[i])
        }
        integStr = string(result)
    }
    
	sign := ""
	if isNeg { sign = "-" }
    return sign + integStr + "," + decimStr
}

func abs(n int) int {
    if n < 0 { return -n }
    return n
}

func CentInt(s string) int {
    if len(s) == 0 { return -1 }

    s = strings.Trim(s, "€")
    s = strings.Trim(s, ".")
    s = strings.TrimSpace(s)

    validChar := regexp.MustCompile(`^[0-9.,\s]+$`)
    if !validChar.MatchString(s) {
        return -1
    }

    var integ, decim string
    if strings.Contains(s, ",") {
        parts := strings.Split(s, ",")
        if len(parts) != 2 || len(parts[1]) > 2 {
            return -1
        }
        integ = parts[0]
        decim = parts[1]
        if len(decim) == 0 {
            decim = "00"
        } else if len(decim) == 1 {
            decim = decim + "0"
        }
    } else {
        integ = s
        decim = "00"
    }

    integInt, err := strconv.Atoi(integ)
    if err != nil {
        return -1
    }
    decimInt, err := strconv.Atoi(decim)
    if err != nil {
        return -1
    }

    return integInt*100 + decimInt
}
