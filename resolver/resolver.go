package resolver

import (
	"fmt"
	"net/http"
)

//RequestURL is an impl of URLResolver interface
type RequestURL struct {
	r    http.Request
	Port int
	Host string
}

//SetRequest to implement `RequestAwareResolverInterface`
func (m *RequestURL) SetRequest(r http.Request) {
	m.r = r
}

//GetBaseURL implements `URLResolver` interface
func (m RequestURL) GetBaseURL() string {
	return fmt.Sprintf("http://%s:%d", m.Host, m.Port)
}
