package funit

// Settings for the trigger
type Settings struct {
}

// HandlerSettings for the trigger
type HandlerSettings struct {
}

// Output of the trigger
type Output struct {
}

// Reply from the trigger
type Reply struct {
}

// FromMap sets trigger output values from a map
func (o *Output) FromMap(values map[string]interface{}) error {

	return nil
}

// ToMap converts trigger output to a map
func (o *Output) ToMap() map[string]interface{} {
	return nil
}

// FromMap sets trigger reply values from a map
func (r *Reply) FromMap(values map[string]interface{}) error {
	return nil
}

// ToMap converts trigger reply to a map
func (r *Reply) ToMap() map[string]interface{} {
	return nil
}
