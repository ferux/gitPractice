package rnd

func init() {
	InitVersion = "InitAloha"
}

//Settings contains app settings
type Settings struct {
	Name  string
	Value int
}

//Version here
var (
	InitVersion string
	Version     string
)

//Init the settings
func Init(name string, value int) *Settings {
	Version = "Aloha"
	return &Settings{name, value}
}
