package async

// ErrorHandle refers to the interface that fetches the error from the async job.
// Calling the Err method is a blocking one so defer it for async behaviour.
type ErrorHandle interface {
	Err() error
}

type errorChannel chan error

func (e errorChannel) Err() error {
	return <-e
}
