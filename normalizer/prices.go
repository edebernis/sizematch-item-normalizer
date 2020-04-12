package normalizer

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "regexp"
    "strings"
)

var currencyRegexp = regexp.MustCompile(`(?is)(?P<currency>eur|gbp)`)

func parseCurrency(s string) (items.Price_Currency, error) {
    match := currencyRegexp.FindStringSubmatch(s)
    if match == nil {
        return 0, fmt.Errorf("No match found for currency: %s", s)
    }

    switch strings.ToLower(match[1]) {
    case "eur":
        return items.Price_EUR, nil
    case "gbp":
        return items.Price_GBP, nil
    default:
        return 0, fmt.Errorf("Unknown unit: %s", s)
    }
}
