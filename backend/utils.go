package main

import "strings"

// frontend recipes page.tsx
func normalizeTitle(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "duplicate", "")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

// chat gippity
func containsAll(slice []string, subslice []string) bool {
	set := make(map[string]struct{})

	for _, v := range slice {
		set[v] = struct{}{}
	}

	for _, v := range subslice {
		if _, ok := set[v]; !ok {
			return false
		}
	}

	return true
}

func containsAny(slice []string, subslice []string) bool {
	set := make(map[string]struct{})

	for _, v := range slice {
		set[v] = struct{}{}
	}

	for _, v := range subslice {
		if _, ok := set[v]; ok {
			return true
		}
	}

	return false
}
