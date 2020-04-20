package normalizer

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "regexp"
    "strconv"
    "strings"
)

type dimension struct {
    Name            items.Dimension_Name
    NameRegexp      map[items.Lang]*regexp.Regexp
    ValueUnitRegexp *regexp.Regexp
}

var dimensions = []dimension{
    {
        Name: items.Dimension_HEIGHT,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)height`),
            items.Lang_FR: regexp.MustCompile(`(?is)hauteur`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_WIDTH,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)width`),
            items.Lang_FR: regexp.MustCompile(`(?is)largeur`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_DEPTH,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)depth`),
            items.Lang_FR: regexp.MustCompile(`(?is)profondeur`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_LENGTH,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)length`),
            items.Lang_FR: regexp.MustCompile(`(?is)longueur`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_DIAMETER,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)diameter`),
            items.Lang_FR: regexp.MustCompile(`(?is)diam(è|e)tre`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_THICKNESS,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)thickness`),
            items.Lang_FR: regexp.MustCompile(`(?is)(é|e)paisseur`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name: items.Dimension_VOLUME,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)(volume|capacity)`),
            items.Lang_FR: regexp.MustCompile(`(?is)(volume|capacit(é|e))`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mm3|cm3|m3|ml|cl|l)`),
    },
    {
        Name: items.Dimension_WEIGHT,
        NameRegexp: map[items.Lang]*regexp.Regexp{
            items.Lang_EN: regexp.MustCompile(`(?is)(weight|mass)`),
            items.Lang_FR: regexp.MustCompile(`(?is)(poids|masse)`),
        },
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>gs?|kgs?)`),
    },
}

var units = map[string]items.Dimension_Unit{
    "mm":  items.Dimension_MM,
    "mms": items.Dimension_MM,
    "cm":  items.Dimension_CM,
    "cms": items.Dimension_CM,
    "m":   items.Dimension_M,
    "ms":  items.Dimension_M,
    "g":   items.Dimension_G,
    "gs":  items.Dimension_G,
    "kg":  items.Dimension_KG,
    "kgs": items.Dimension_KG,
    "mm2": items.Dimension_MM2,
    "cm2": items.Dimension_CM2,
    "m2":  items.Dimension_M2,
    "mm3": items.Dimension_MM3,
    "cm3": items.Dimension_CM3,
    "m3":  items.Dimension_M3,
    "ml":  items.Dimension_ML,
    "cl":  items.Dimension_CL,
    "l":   items.Dimension_L,
}

func (d *dimension) MatchName(s string, l items.Lang) bool {
    return d.NameRegexp[l].MatchString(s)
}

func (d *dimension) GetValueAndUnit(s string) (float64, items.Dimension_Unit, error) {
    match := d.ValueUnitRegexp.FindStringSubmatch(s)
    if match == nil {
        return 0, 0, fmt.Errorf("No match found for value and unit: %s", s)
    }
    value, err := strconv.ParseFloat(match[1], 64)
    if err != nil {
        return 0, 0, err
    }
    unit := units[strings.ToLower(match[2])]
    return value, unit, nil
}

func findDimension(name string, lang items.Lang) (*dimension, error) {
    result := (*dimension)(nil)
    for i, d := range dimensions {
        if d.MatchName(name, lang) {
            if result != nil {
                return nil, fmt.Errorf("Multiple matching dimensions for name: %s", name)
            }
            result = &dimensions[i]
        }
    }

    if result == nil {
        return nil, fmt.Errorf("No dimension found for name: %s", name)
    }

    return result, nil
}

func getMaxValueAndUnit(d *dimension, values []string) (float64, items.Dimension_Unit, error) {
    if len(values) == 0 {
        return 0, 0, fmt.Errorf("No values provided")
    }

    value := (float64)(0)
    unit := (items.Dimension_Unit)(0)
    for _, s := range values {
        v, u, err := d.GetValueAndUnit(s)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        if v > value {
            value = v
            unit = u
        }
    }
    return value, unit, nil
}
