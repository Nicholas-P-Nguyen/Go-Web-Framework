# COS316, Assignment 2: HTTP Routing Framework

# HTTP Routing Framework

As discussed in lecture, naming schemes are central to system design.
In this project, you'll build a general library, called an HTTP Routing Framework,
to help structure web applications based on patterns in end-user requests.
This general library will support a naming scheme for clients accessing resources
provided by a web application.

## Getting Started

Before you begin working on this assignment, you will need to set up your
development environment.

Namely, `git clone` your assignment 2 repository
from GitHub into whichever local directory you are using to store your
assignment files. (See the *General Assignment Instructions* section on the website's assignments tab for more detailed instructions).

## API

Your solution must implement the following API:

```go
// HTTPRouter: stores the information necessary to route HTTP requests
type HttpRouter struct {
    // this can include whatever fields you want
}

// AddRoute: adds a new route to the router, associates/maps a given method and path
// pattern to the designated http handler.
//
// Details:
// Routes any request that matches the given `method` and `pattern` to a
// `handler`.
//
// `method`: should support arbitrary method strings, and at least each of "POST",
//           "GET", "PUT", "DELETE". Method strings are case insensitive.
//
// `pattern`: patterns on the request _path_. Patterns can include arbitrary
//           "directories". Directories may be "captures" of the form
//           `:variable_name` that capture the actual directory value into an
//            HTTP query paramter. Leading and trailing '/' (slash) characters
//            are ignored.
//
// Example:
//
//   AddRoute("GET", "/users/:user/recent", GivenHandlerFunc)
//
// should map all GET requests with a path of the form "/users/*/recent" to the
// `GivenHandlerFunc` handler. 
//
func (router *HTTPRouter) AddRoute(method string, pattern string, handler http.HandlerFunc)

// ServeHTTP: For the given request, finds the correct handler
// associated with the route that is appropriate for the provided request's path. 
// Then, invokes this handler to allow it to write an HTTP response to the provided response writer.
// Note: the specific handler's implementation is NOT in the scope of this API.
//
// Continuing from the same example from AddRoute, for a request of the form 
//    "GET /users/cesar/recent HTTP/1.1",
//
// this function should call the GivenHandlerFunc handler function with an `http.Request` 
// with a raw query in the form of "user=cesar" (`URL.RawQuery = "user=cesar"`). As mentioned in the specs, 
// Go's url package and the "Values" type (and its Encode function) should be useful for this step.
func (router *HTTPRouter) ServeHTTP(response http.ResponseWriter, request *http.Request)
```

## Additional Specifications

Be sure that your implementation of the above basic API takes the following into account:

* Your router must support arbitrary paths and HTTP methods, not just those
  required by the microblog client discussed below. You need not (and should not)
  validate that a client-provided method is part of the [official HTTP spec](https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods#Specifications).

* You may assume that your router will only be provided with paths that are
  well-formed. That is, paths will be a (possibly empty) list of directory names,
  separated by `/`. Directory names may contain the characters `a-zA-Z0-9_-.`

* HTTP Methods are case-insensitive, whereas paths are case-sensitive. That is,
  `GET` and `get` are equivalent, but `/path/to/file` and `/Path/To/File` are not.

* If (and only if) the request provided to `ServeHTTP` has no associated route,
  your router must write an HTTP "404 Not Found" error as its response. Any response
  with a `Status` field equal to 404 will do.
  You may find Go's [http package](https://golang.org/pkg/net/http/) useful.

* Calling `AddRoute` with a `method`, `pattern` pair that is already associated
  with a handler should update the associated handler to the newly provided value.
  Further, if `pattern` contains a capture, the relevant query parameters
  should be updated. (e.g. `/path/to/:file --> /path/to/:dir`)

* Captures should only be overwritten when necessary. For example, `/photos/:year/july`
  and `/photos/:place/winter` can coexist while `/photos/:year/picture` and `photos/:place/picture`
  cannot.

* A single directory within a path may contain at most one capture. `/path/to/:file`
  is valid, but `/path/to/:file1:file2` is not.

* A capture must always comprise the entire directory in which it appears. It is
  *not valid* to "prefix" captures, such as in `/path/to/user:id/photos`.

* The empty capture is not valid, either when providing a pattern to AddRoute
  (e.g. `/path/:/file`), or when providing a request path to ServeHTTP
  (e.g. `/path//file`). You may assume your code will never encounter these cases.

* A pattern may include a capture for several values using the same name, as in
  `/path/to/:file/:file`. In this case, the http Response should have a
  `URL.RawQuery` of `file=<value>&file=<value>`. We impose no restrictions on the
  order in which key-value pairs appear, except that the values for a particular
  key must appear in the same order as they did in the request path.
  You may find Go's [url package](https://golang.org/pkg/net/url/)
  useful, particularly the [Values type](https://golang.org/pkg/net/url/#Values).

Be aware of the following notable edge case:
* Your router may be asked to route a request for which there are multiple
  applicable routes. Consider a router with a non-capturing route
  `GET /path/to/file` and a capturing route `GET /path/to/:filename`, that has
  been asked to route a request for `GET /path/to/file`.

  In such cases, routes without captures should be given precedence over routes
  with captures. The order in which the routes were added is not considered.

  To be explicit, the route used by the router should be the matching route with
  the largest number of non-capturing path components to the left of its first
  capture. This applies recursively, so that ties are broken by finding the
  matching route with the largest number of non-capturing path components between
  the first and second capture, and so on. If you like, you may informally think
  of this as finding the route which "waits the longest" before capturing each
  of its values.

Finally, you should *not* have to iterate through every single route in the
router to find a matching route. Finding a route in this way will significantly
hurt performance when your router contains a large number of routes, and you
should aim to make your implementation as quick as possible in these cases.
While **this does not affect your overall grade for this assignment**, it is a
useful exercise in reasoning about performance trade-offs.

## Unit Testing

Go provides the [testing package](https://golang.org/pkg/testing/), which is a
convenient framework that allows you to write unit tests for your code to ensure
that it is working correctly and providing the expected results.

For this assignment, you are provided with the file `router_test.go`, which
contains a couple of very simple tests for your router implementation and
demonstrates how to use the `testing` package.

You are encouraged to extend this file to create your own unit tests for your
`http_router` implementation.

You can run your unit tests with the command `go test`, which simply reports the
result of the test, and the reason for failure, if any, or you may add the `-v`
flag to see the verbose output of the unit tests.

For example, run the following from your top-level assignment 2 directory:
```bash
$ go test -v ./http_router
```
Equivalently, you may `cd` into the http_router directory and run the following:
```bash
$ go test -v
```

You will not be graded directly on the quality of your unit tests for this
assignment, but good unit tests will help you debug and understand your
program. We **highly** recommend writing *at least* three or four simple unit tests,
both to familiarize yourself with the API, and to help identify tricky corner
cases or debug unexpected results.


## Submission & Grading

Your assignment will be automatically submitted every time you push your changes
to your GitHub repository. Within a couple minutes of your submission, the
autograder will make a comment on your commit listing the output of our testing
suite when run against your code. **Note that you will be graded only on your
changes to the `router.go` file**, and not on your tests in the `router_test.go` file.

Remember to fill out the `readme_contributions` file as well! 
For more information, please consult the Grading and Policies tab of our website.
