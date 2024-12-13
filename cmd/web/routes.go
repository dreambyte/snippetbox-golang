package main

import (
	"net/http"

	"github.com/justinas/alice"
	"snippetbox.victran/ui"
)

// The routes() method returns a servemux containing our application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	// fileServer := http.FileServer(http.Dir("./ui/static/"))
	// mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.Handle("GET /static/", http.FileServerFS(ui.Files)) // use embedded files

	// Add the authenticate() middleware to the chain.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippet/create", protected.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", protected.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", protected.ThenFunc(app.userLogoutPost))
	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standard.Then(mux)

	// mux.Handle("GET /{$}", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.home))))
	// mux.Handle("GET /snippet/view/{id}", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.snippetView))))
	// mux.Handle("GET /user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignup))))
	// mux.Handle("POST /user/signup", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userSignupPost))))
	// mux.Handle("GET /user/login", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userLogin))))
	// mux.Handle("POST /user/login", app.sessionManager.LoadAndSave(noSurf(http.HandlerFunc(app.userLoginPost))))

	// mux.Handle("GET /snippet/create", app.sessionManager.LoadAndSave(noSurf(app.requireAuthentication(http.HandlerFunc(app.snippetCreate)))))
	// mux.Handle("POST /snippet/create", app.sessionManager.LoadAndSave(noSurf(app.requireAuthentication(http.HandlerFunc(app.snippetCreatePost)))))
	// mux.Handle("POST /user/logout", app.sessionManager.LoadAndSave(noSurf(app.requireAuthentication(http.HandlerFunc(app.userLogoutPost)))))
	// return app.recoverPanic(app.logRequest(commonHeaders(mux)))
}
