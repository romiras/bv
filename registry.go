package main

import (
	"strings"
)

type Registry struct {
	filters []string
}

func NewRegistry() *Registry {
	var filtersPtr *string

	filters := parseFilters(filtersPtr)

	return &Registry{
		filters: filters,
	}
}

func parseFilters(filtersPtr *string) []string {
	var result []string

	if filtersPtr == nil {
		return make([]string, 0)
	}

	result = strings.Split(*filtersPtr, ",")

	return result
}
