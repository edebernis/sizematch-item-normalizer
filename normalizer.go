package main

import (
    "strconv"
)

// Item ...
type Item struct {
    ID            string
    Source        string
    Name          string
    Description   string
    Categories    []string
    Lang          string
    Urls          []string
    ImageUrls     []string
    Dimensions    map[string]string
    Price         string
    PriceCurrency string
}

// Dimensions ...
type Dimensions struct {
    Height float64
    Width  float64
    Depth  float64
}

// NormalizedItem ...
type NormalizedItem struct {
    ID            string
    Source        string
    Name          string
    Description   string
    Categories    []string
    Lang          string
    Urls          []string
    ImageUrls     []string
    Dimensions    *Dimensions
    Price         float64
    PriceCurrency string
}

func normalize(item *Item) (*NormalizedItem, error) {
    dimensions, err := normalizeDimensions(item.Dimensions)
    if err != nil {
        return nil, err
    }

    price, err := normalizePrice(item.Price)
    if err != nil {
        return nil, err
    }

    return &NormalizedItem{
        ID:            item.ID,
        Source:        item.Source,
        Name:          item.Name,
        Description:   item.Description,
        Categories:    item.Categories,
        Lang:          item.Lang,
        Urls:          item.Urls,
        ImageUrls:     item.ImageUrls,
        Dimensions:    dimensions,
        Price:         price,
        PriceCurrency: item.PriceCurrency,
    }, nil
}

func normalizeDimensions(dimensionsMap map[string]string) (*Dimensions, error) {
    return &Dimensions{
        Height: 0.0,
        Width:  0.0,
        Depth:  0.0,
    }, nil
}

func normalizePrice(price string) (float64, error) {
    return strconv.ParseFloat(price, 64)
}
