package main

import (
    "fmt"
    "strconv"
)

type item struct {
    id            string
    source        string
    name          string
    description   string
    categories    []string
    lang          string
    urls          []string
    imageUrls     []string
    dimensions    map[string]string
    price         string
    priceCurrency string
}

type dimensions struct {
    height float64
    width  float64
    depth  float64
}

type normalizedItem struct {
    id            string
    source        string
    name          string
    description   string
    categories    []string
    lang          string
    urls          []string
    imageUrls     []string
    dimensions    *dimensions
    price         float64
    priceCurrency string
}

func normalize(item *item) (*normalizedItem, error) {
    fmt.Printf("%+v\n", item)

    dimensions, err := normalizeDimensions(item.dimensions)
    if err != nil {
        return nil, err
    }

    price, err := normalizePrice(item.price)
    if err != nil {
        return nil, err
    }

    return &normalizedItem{
        id:            item.id,
        source:        item.source,
        name:          item.name,
        description:   item.description,
        categories:    item.categories,
        lang:          item.lang,
        urls:          item.urls,
        imageUrls:     item.imageUrls,
        dimensions:    dimensions,
        price:         price,
        priceCurrency: item.priceCurrency,
    }, nil
}

func normalizeDimensions(dimensionsMap map[string]string) (*dimensions, error) {
    return &dimensions{
        height: 12.2,
        width:  13.3,
        depth:  14.4,
    }, nil
}

func normalizePrice(price string) (float64, error) {
    return strconv.ParseFloat(price, 64)
}
