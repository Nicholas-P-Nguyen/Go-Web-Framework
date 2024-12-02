/*****************************************************************************
 * router.go
 * Name: Nicholas Nguyen
 * NetId: nn5029
 *****************************************************************************/

package http_router

import (
	"net/http"
	"net/url"
	"strings"
)

// RoutesField fields the parameters needed to add a route
type RoutesField struct {
	Method  string
	Pattern string
	Handler http.HandlerFunc
}

// HTTPRouter stores a list of routes that contains: method, pattern, and handler that has been added
type HTTPRouter struct {
	Routes []RoutesField
}

// NewRouter creates a new HTTP Router, with no initial routes
func NewRouter() *HTTPRouter {
	return &HTTPRouter{
		Routes: []RoutesField{},
	}
}

//----------------------------------------------------------------------------------------------------------------------

// AddRoute adds a new route to the router and maps a given method, path, and handler
func (router *HTTPRouter) AddRoute(method string, pattern string, handler http.HandlerFunc) {
	// Edge case: ignore leading and trailing '/'
	if pattern == "/" {
		pattern = ""
	}

	method = strings.ToUpper(method)
	pattern = strings.TrimPrefix(pattern, "/")
	pattern = strings.TrimSuffix(pattern, "/")

	for i := range router.Routes {
		// checking for existing static patterns, updating handler accordingly
		if router.Routes[i].Method == method && router.Routes[i].Pattern == pattern {
			router.Routes[i].Handler = handler
			return
		}
		// checking for existing dynamic patterns, updating handler and pattern accordingly
		if router.Routes[i].Method == method && IsExistingPath(pattern, router.Routes[i].Pattern) {
			router.Routes[i].Pattern = pattern
			router.Routes[i].Handler = handler
			return
		}
	}
	router.Routes = append(router.Routes, RoutesField{Method: method, Pattern: pattern, Handler: handler})
}

// IsExistingPath is a helper function for AddRoute that checks if an existing pattern in the router matches the new
// route being added
func IsExistingPath(newPattern, existingPattern string) bool {
	newPathSplit := strings.Split(newPattern, "/")
	existingPathSplit := strings.Split(existingPattern, "/")

	if len(newPathSplit) != len(existingPathSplit) {

		return false
	}
	if newPattern == "" && existingPattern == "" {
		return true
	}

	for i := range existingPathSplit {
		if strings.HasPrefix(existingPathSplit[i], ":") && strings.HasPrefix(newPathSplit[i], ":") {
			continue
		}
		if newPathSplit[i] != existingPathSplit[i] {
			return false
		}
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------

// ServeHTTP For the given request, finds the correct handler
func (router *HTTPRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	requestPattern := request.URL.Path
	// Edge case: ignore leading and trailing '/'
	if requestPattern == "/" {
		requestPattern = ""
	}
	requestPattern = strings.TrimPrefix(requestPattern, "/")
	requestPattern = strings.TrimSuffix(requestPattern, "/")

	requestMethod := strings.ToUpper(request.Method)

	var dynamicRoutes []RoutesField
	var bestRoute RoutesField
	bestRouteFound := false

	for _, route := range router.Routes {
		// Checking for static routes FIRST and invoke the handler to write HTTP response
		if route.Method == requestMethod && IsSameStaticPattern(requestPattern, route.Pattern) {
			route.Handler(response, request)
			return
		}
		// Checking for dynamic routes and appending to slice
		if route.Method == requestMethod && IsSameDynamicPattern(requestPattern, route.Pattern) {
			dynamicRoutes = append(dynamicRoutes, route)
		}
	}

	for _, route := range dynamicRoutes {
		if !bestRouteFound || IsHigherPrecedence(route.Pattern, bestRoute.Pattern) {
			bestRoute = route
			bestRouteFound = true
		}
	}

	if bestRouteFound {
		captureToValues := GetCapturesValues(requestPattern, bestRoute.Pattern)
		query := url.Values{}
		for capture, valueSlice := range captureToValues {
			for _, value := range valueSlice {
				query.Add(capture, value)
			}
		}
		request.URL.RawQuery = query.Encode()
		bestRoute.Handler(response, request)
		return
	}
	http.NotFound(response, request)
}

// IsHigherPrecedence is a helper function for ServeHTTP that compares two routes and finds out which pattern has
// largest number of non-capturing path components to the left of its first capture.
func IsHigherPrecedence(currentPattern, bestPattern string) bool {
	currentPatternSplit := strings.Split(currentPattern, "/")
	bestPatternSplit := strings.Split(bestPattern, "/")

	for i := 0; i < len(currentPatternSplit) && i < len(bestPatternSplit); i++ {
		if strings.HasPrefix(currentPatternSplit[i], ":") && !strings.HasPrefix(bestPatternSplit[i], ":") {
			return false
		}
		if !strings.HasPrefix(currentPatternSplit[i], ":") && strings.HasPrefix(bestPatternSplit[i], ":") {
			return true
		}
	}
	return len(currentPatternSplit) < len(bestPatternSplit)
}

// IsSameStaticPattern is a helper function for ServeHTTP that checks if two static paths are the same
func IsSameStaticPattern(requestPath string, routePath string) bool {
	requestPathSplit := strings.Split(requestPath, "/")
	routePathSplit := strings.Split(routePath, "/")

	if len(requestPathSplit) != len(routePathSplit) {
		return false
	}

	if requestPath == "" && routePath == "" {
		return true
	}

	for i := range routePathSplit {
		if requestPathSplit[i] != routePathSplit[i] {
			return false
		}
	}
	return true
}

// IsSameDynamicPattern is a helper function for ServeHTTP that checks if two dynamic paths are the same
func IsSameDynamicPattern(requestPath string, path string) bool {
	// users/nicholas/recent -> users nicholas recent
	requestPathSplit := strings.Split(requestPath, "/")
	// users/:user/recent -> users :user recent
	pathSplit := strings.Split(path, "/")

	if len(requestPathSplit) != len(pathSplit) {
		return false
	}

	if requestPath == "" && path == "" {
		return true
	}

	for i := range pathSplit {
		if strings.HasPrefix(pathSplit[i], ":") {
			continue
		}
		if requestPathSplit[i] != pathSplit[i] {
			return false
		}
	}
	return true
}

// GetCapturesValues is a helper function for ServeHTTP that maps query parameters to the value provided by the client
func GetCapturesValues(requestPath string, path string) map[string][]string {
	// users/nicholas/recent -> users nicholas recent
	requestPathSplit := strings.Split(requestPath, "/")
	// users/:user/recent -> users :user recent
	pathSplit := strings.Split(path, "/")

	captureToValue := make(map[string][]string)

	for i := range pathSplit {
		if strings.HasPrefix(pathSplit[i], ":") {
			capture := pathSplit[i][1:]
			captureToValue[capture] = append(captureToValue[capture], requestPathSplit[i])
		}
	}
	return captureToValue
}
