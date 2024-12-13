package main

import (
	"io/fs"
	"path/filepath"
	"text/template"
	"time"

	"snippetbox.victran/internal/models"
	"snippetbox.victran/ui"
)

// Define a templateData type to act as the holding structure for // any dynamic data that we want to pass to our HTML templates. // At the moment it only contains one field, but we'll add more // to it as the build progresses.
type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	// Convert the time to UTC before formatting it.
	return t.UTC().Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	// Use fs.Glob() to get a slice of all filepaths in the ui.Files embedded
	// filesystem which match the pattern 'html/pages/*.tmpl'. This essentially // gives us a slice of all the 'page' templates for the application, just // like before.
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		// Create a slice containing the filepath patterns for the templates we
		// want to parse.
		patterns := []string{
			"html/base.tmpl", "html/partials/*.tmpl", page,
		}
		// Use ParseFS() instead of ParseFiles() to parse the template files
		// from the ui.Files embedded filesystem.
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

// not using embedded files
// func newTemplateCache() (map[string]*template.Template, error) { // Initialize a new map to act as the cache.
// 	cache := map[string]*template.Template{}
// 	// Use the filepath.Glob() function to get a slice of all filepaths that
// 	// match the pattern "./ui/html/pages/*.tmpl". This will essentially gives
// 	// us a slice of all the filepaths for our application 'page' templates
// 	// like: [ui/html/pages/home.tmpl ui/html/pages/view.tmpl]
// 	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Loop through the page filepaths one-by-one.
// 	for _, page := range pages {
// 		name := filepath.Base(page)
// 		// Parse the base template file into a template set.
// 		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Call ParseGlob() *on this template set* to add any partials.
// 		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Call ParseFiles() *on this template set* to add the page template.
// 		ts, err = ts.ParseFiles(page)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// Add the template set to the map as normal...
// 		cache[name] = ts
// 	}
// 	// Return the map.
// 	return cache, nil
// }
