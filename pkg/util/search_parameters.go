package util

type SearchParameters map[string]string

// Add inserts key-value pair into parameters if the key does not exist
// updates the key if it exists
func (sp SearchParameters) Add(key, value string) {
	sp[key] = value
}

func (sp SearchParameters) Remove(key string) {
	delete(sp, key)
}

func (sp SearchParameters) Get(key string) (string, bool) {
	val, ok := sp[key]
	return val, ok
}

func (sp SearchParameters) Clear() SearchParameters {
	return make(SearchParameters)
}
