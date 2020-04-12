package normalizer

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/go/items"
    "regexp"
    "strconv"
    "strings"
)

var dimensions = []dimension{
    {
        Name:            items.Dimension_HEIGHT,
        NameRegexp:      regexp.MustCompile(`(?is)height`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_WIDTH,
        NameRegexp:      regexp.MustCompile(`(?is)width`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_DEPTH,
        NameRegexp:      regexp.MustCompile(`(?is)depth`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_LENGTH,
        NameRegexp:      regexp.MustCompile(`(?is)length`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_DIAMETER,
        NameRegexp:      regexp.MustCompile(`(?is)diameter`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_THICKNESS,
        NameRegexp:      regexp.MustCompile(`(?is)thickness`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mms?|cms?|ms?)`),
    },
    {
        Name:            items.Dimension_VOLUME,
        NameRegexp:      regexp.MustCompile(`(?is)volume`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>mm3|cm3|m3|l)`),
    },
    {
        Name:            items.Dimension_WEIGHT,
        NameRegexp:      regexp.MustCompile(`(?is)weight`),
        ValueUnitRegexp: regexp.MustCompile(`(?is)(?P<value>\d*\.?\d+)\s*(?P<unit>gs?|kgs?)`),
    },
}

type dimension struct {
    Name            items.Dimension_Name
    NameRegexp      *regexp.Regexp
    ValueUnitRegexp *regexp.Regexp
}

func (d *dimension) MatchName(s string) bool {
    return d.NameRegexp.MatchString(s)
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
    unit, err := d.ParseUnit(match[2])
    if err != nil {
        return 0, 0, err
    }
    return value, unit, nil
}

func (d *dimension) ParseUnit(s string) (items.Dimension_Unit, error) {
    switch strings.ToLower(s) {
    case "mm", "mms":
        return items.Dimension_MM, nil
    case "cm", "cms":
        return items.Dimension_CM, nil
    case "m", "ms":
        return items.Dimension_M, nil
    case "g", "gs":
        return items.Dimension_G, nil
    case "kg", "kgs":
        return items.Dimension_KG, nil
    case "mm2":
        return items.Dimension_MM2, nil
    case "cm2":
        return items.Dimension_CM2, nil
    case "m2":
        return items.Dimension_M2, nil
    case "mm3":
        return items.Dimension_MM3, nil
    case "cm3":
        return items.Dimension_CM3, nil
    case "m3":
        return items.Dimension_M3, nil
    case "l":
        return items.Dimension_L, nil
    default:
        return 0, fmt.Errorf("Unknown unit: %s", s)
    }
}

func findDimension(name string) (*dimension, error) {
    result := (*dimension)(nil)
    for i, d := range dimensions {
        if d.MatchName(name) {
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
