package normalizer

import (
    "fmt"
    "github.com/edebernis/sizematch-protobuf/build/go/items"
    "reflect"
    "strings"
)

// Normalize items
func Normalize(item *items.Item) (*items.NormalizedItem, error) {
    n := normalizer{
        Item:           item,
        NormalizedItem: &items.NormalizedItem{},
    }

    err := n.normalize()
    if err != nil {
        return nil, err
    }

    return n.NormalizedItem, nil
}

func isProtoField(field string) bool {
    return strings.HasPrefix(field, "XXX_")
}

type normalizer struct {
    Item           *items.Item
    NormalizedItem *items.NormalizedItem
}

func (n *normalizer) normalize() error {
    normalizedItemType := reflect.ValueOf(n.NormalizedItem).Elem()
    for i := 0; i < normalizedItemType.NumField(); i++ {
        fieldName := normalizedItemType.Type().Field(i).Name
        if isProtoField(fieldName) {
            continue
        }
        value, err := n.normalizeField(fieldName)
        if err != nil {
            return err
        }
        normalizedItemType.Field(i).Set(value)
    }
    return nil
}

func (n *normalizer) normalizeField(fieldName string) (reflect.Value, error) {
    if len(fieldName) < 3 {
        fieldName = strings.ToUpper(fieldName)
    }

    method, found := reflect.TypeOf(n).MethodByName(fieldName)
    if !found {
        return reflect.ValueOf(nil), fmt.Errorf("method not found for %s", fieldName)
    }

    in := []reflect.Value{reflect.ValueOf(n)}
    result := method.Func.Call(in)
    if result[1].IsNil() {
        return result[0], nil
    }
    return result[0], result[1].Interface().(error)
}

func (n *normalizer) ID() (string, error) {
    return n.Item.Id, nil
}

func (n *normalizer) Source() (string, error) {
    return n.Item.Source, nil
}

func (n *normalizer) Lang() (items.NormalizedItem_Lang, error) {
    return items.NormalizedItem_EN, nil
}

func (n *normalizer) Urls() ([]string, error) {
    return n.Item.Urls, nil
}

func (n *normalizer) Name() (string, error) {
    return n.Item.Name, nil
}

func (n *normalizer) Description() (string, error) {
    return n.Item.Description, nil
}

func (n *normalizer) Categories() ([]string, error) {
    return n.Item.Categories, nil
}

func (n *normalizer) ImageUrls() ([]string, error) {
    return n.Item.ImageUrls, nil
}

func (n *normalizer) Dimensions() (*items.Dimensions, error) {
    return &items.Dimensions{}, nil
}

func (n *normalizer) Price() (*items.Price, error) {
    return &items.Price{}, nil
}
