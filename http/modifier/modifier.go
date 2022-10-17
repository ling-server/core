package modifier

import "net/http"

type Modifier interface {
	Modify(*http.Request) error
}
