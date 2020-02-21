package main

// PageDescription is the description of the page.
type PageDescription struct {
	HTMLTitle   string `json:"html_title"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Form is a form in the structure.
type Form struct {
	Name     string                             `json:"name"`
	Type     string                             `json:"type"`
	Warning  *string                            `json:"warning"`
	Message  *string                            `json:"message"`
	Children *map[string]map[string]interface{} `json:"children"`
}

// BaseStructure is the structure at the base of the JSON.
type BaseStructure struct {
	PageDescription PageDescription `json:"page_description"`
	Forms           []Form          `json:"forms"`
	Redirect        string          `json:"redirect"`
}
