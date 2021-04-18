package main

import (
	"flag"
	"strings"
)

type Registry struct {
	filters []string
}

func NewRegistry() *Registry {
	filtersPtr := flag.String("f", "ae", "Comma-separated names of filters")
	flag.Parse()

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
