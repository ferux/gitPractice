package hidstruct

import (
	"encoding/json"
	"log"
)

// Hidden struct hides password field
type Hidden struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Hint     string `json:"hint,omitempty"`
}

// Init creates a new password
func Init(user, pass, hint string) *Hidden {
	return &Hidden{
		user,
		pass,
		hint,
	}
}

// MarshalJSON modified to hide password
func (h Hidden) MarshalJSON() ([]byte, error) {
	type hidden Hidden
	hh := hidden(h)
	hh.Password = "hidden"
	return json.Marshal(((*hidden)(&hh)))
}

// SetPassword updates password
func (h *Hidden) SetPassword(p string) bool {
	if len(p) == 0 {
		log.Println("New password has zero length. Can't assign.")
		return false
	}
	h.Password = p
	log.Println("New password applied successfuly")
	return true
}

// GetPassword from struct
func (h *Hidden) GetPassword() string {
	return h.Password
}

// UpdatePointer sets new address for struct
func (h *Hidden) UpdatePointer() {
	h = &Hidden{h.Name, h.Password, h.Hint}
}
