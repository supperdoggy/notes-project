package main

// need to change a lot of code to implement this
// because structs.Map() transforms only exported fields
// and in result all keys in map start with uppercase
type apiResponse struct {
	Ok     bool        `json:"ok"`
	Error  string      `json:"error"`
	Answer interface{} `json:"answer"`
}
