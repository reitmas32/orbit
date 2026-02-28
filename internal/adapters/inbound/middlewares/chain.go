package middlewares

import "net/http"

// Middleware es una función que toma un http.HandlerFunc y retorna otro http.HandlerFunc
type Middleware func(http.HandlerFunc) http.HandlerFunc

type Chain struct {
	middlewares []Middleware
}

// NewChain construye una cadena de middlewares.
// El orden importa: el primero definido será el primero en ejecutarse.
func NewChain(middlewares ...Middleware) Chain {
	return Chain{middlewares: middlewares}
}

// Then envuelve el manejador final con los middlewares.
// El 'handler' que se pasa aquí siempre será el último en ejecutarse en la cadena.
func (c Chain) Then(handler http.HandlerFunc) http.HandlerFunc {
	// Hacemos un loop inverso para que el primer middleware en la cadena sea la capa más externa
	for i := len(c.middlewares) - 1; i >= 0; i-- {
		handler = c.middlewares[i](handler)
	}
	return handler
}
