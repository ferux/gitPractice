package hidstruct_test

import (
	"encoding/json"
	"fmt"
	"testing"

	h "github.com/ferux/gitPractice/hidstruct"
)

func TestString(t *testing.T) {
	var ss h.StringSafe
	ss = "Hello!"
	if x := fmt.Sprintf("%v", ss); x != h.FilterMessage {
		t.Errorf("expected %s\tgot %s\n", h.FilterMessage, x)
	}
	if x := fmt.Sprintf("%+v", ss); x != h.FilterMessage {
		t.Errorf("expected %s\tgot %s\n", h.FilterMessage, x)
	}
	if x := fmt.Sprintf("%#v", ss); x != h.FilterMessage {
		t.Errorf("expected %s\tgot %s\n", h.FilterMessage, x)
	}
	if x := fmt.Sprintf("%v", string(ss)); x != "Hello!" {
		t.Errorf("expected %s\tgot %s\n", "Hello!", x)
	}

	tt := struct{ Name h.StringSafe }{"Hello!"}
	tr := struct{ Name string }{h.FilterMessage}
	data, err := json.Marshal(&tt)
	if err != nil {
		t.Fatalf("can't unmarshal: %v\n", err)
	}
	data2, err := json.Marshal(&tr)
	if err != nil {
		t.Fatalf("can't unmarshal: %v\n", err)
	}
	if string(data) != string(data2) {
		t.Errorf("Data and Data2 are not equal: [%s]\t[%s]", string(data), string(data2))
	}
}
