package main

import (
	"encoding/json"
	"os"
	"sort"
)

func writeJSONReport(pages map[string]PageData, filename string) error {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sorted := make([]PageData, 0, len(pages))
	for _, k := range keys {
		sorted = append(sorted, pages[k])
	}

	data, err := json.MarshalIndent(sorted, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}
