# PolyJSON

A library to generate code to semi-efficiently consume JSON documents that contain lists of polymorphic objects and visit the individual elements.

Please excuse the sparse documentation, this library was extracted from a different project that is not public yet.

## What does it do?

Parsing a JSON document containing a list of events like the one in [[testdata/events/example.json]] can be easily implemented by having a struct with a field for every type, but every code that touches this struct must do a lot of nil-checks, error handling and has to be updated whenever a new type is added.

This library generates the struct, a visitor interface and a convenience implementation for the visitor interface.

## How do I use it?

Create a struct for all type of objects in the JSON to be parsed and mark them with `polyjson.Implements[...]`. The parameter passed as the argument to the generic marker must not exist, as it is generated.

```golang
type FailedLogin struct {
	polyjson.Implements[Event]

	IPAddress string `json:"ip_address"`
}
```

The struct name will be converted to snake case and used as the json field name. You can override the generated `json` field tag by tagging the `polyjson.Implements[...]` marker.

```golang
type FailedLogin struct {
	polyjson.Implements[Event] `polyjson:"login_failed"`

	IPAddress string `json:"ip_address"`
}
```
