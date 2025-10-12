package requests

type Validateable interface {
	Validate() error
}
