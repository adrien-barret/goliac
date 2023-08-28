// Code generated by go-swagger; DO NOT EDIT.

package app

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// PostFlushCacheHandlerFunc turns a function with the right signature into a post flush cache handler
type PostFlushCacheHandlerFunc func(PostFlushCacheParams) middleware.Responder

// Handle executing the request and returning a response
func (fn PostFlushCacheHandlerFunc) Handle(params PostFlushCacheParams) middleware.Responder {
	return fn(params)
}

// PostFlushCacheHandler interface for that can handle valid post flush cache params
type PostFlushCacheHandler interface {
	Handle(PostFlushCacheParams) middleware.Responder
}

// NewPostFlushCache creates a new http.Handler for the post flush cache operation
func NewPostFlushCache(ctx *middleware.Context, handler PostFlushCacheHandler) *PostFlushCache {
	return &PostFlushCache{Context: ctx, Handler: handler}
}

/*
	PostFlushCache swagger:route POST /flushcache app postFlushCache

Flush the Github remote cache
*/
type PostFlushCache struct {
	Context *middleware.Context
	Handler PostFlushCacheHandler
}

func (o *PostFlushCache) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		*r = *rCtx
	}
	var Params = NewPostFlushCacheParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
