package hidstruct

import (
	"strings"
)

// Regular struct for copy and modifying
type Regular struct {
	Name   string
	Params []Param
	IDs    map[string]int
}

// Param for storing key-value params
type Param struct {
	Key   string
	Value string
}

// NewRegular creates new regular object
func NewRegular(name string) *Regular {
	return &Regular{
		Name:   name,
		Params: make([]Param, 0),
		IDs:    make(map[string]int),
	}
}

// Add new key-value pair to params
func (r *Regular) Add(key, value string) {
	for _, v := range r.Params {
		if v.Key == key {
			return
		}
	}
	r.Params = append(r.Params, Param{key, value})
}

// Set new value to param
func (r *Regular) Set(key, value string) {
	for i, v := range r.Params {
		if v.Key == key {
			r.Params[i].Value = value
			return
		}
	}
}

// Delete key-value from params
func (r *Regular) Delete(key string) {
	for i, v := range r.Params {
		if v.Key == key {
			r.Params = append(r.Params[:i], r.Params[i+1:]...)
			return
		}
	}
}

//DeepCopySafe returns new struct with hided password fields
func (r *Regular) DeepCopySafe() *Regular {
	newr := *r
	newr.Params = make([]Param, 0)
	newr.IDs = make(map[string]int)
	for _, v := range r.Params {
		if strings.ToLower(v.Key) == "password" {
			newr.Add(v.Key, "[protected]")
			continue
		}
		newr.Add(v.Key, v.Value)
	}
	for k, v := range r.IDs {
		newr.IDs[k] = v
	}
	return &newr
}
