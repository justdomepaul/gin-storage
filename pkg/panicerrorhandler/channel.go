package panicerrorhandler

// ErrorToChannel method
func ErrorToChannel(err error, ch chan<- error) {
	if err != nil {
		ch <- err
	}
}
