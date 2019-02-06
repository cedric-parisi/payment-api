package docs

import "net/http"

// Handler ...
func Handler() http.Handler {
	fs := http.FileServer(http.Dir("./docs/swaggerui/"))
	return http.StripPrefix("/swaggerui/", fs)

}
