// Package async wraps golang library methods and provides a way to call them in a non blocking way.
package async

// ErrorHandle refers to the function that fetches the error from the async job.
// Calling the ErrorHandle function is a blocking one so defer it for async behaviour.
type ErrorHandle func() error
