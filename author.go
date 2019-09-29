package clippy

// Author represents the author of the software.
type Author struct {
	Name  string // Name of the author. For example, "Patrick McNamara".
	Email string // Email of the author. For example, "hello@patrickmcnamara.xyz".
}

func (a *Author) String() string {
	return a.Name + " " + "<" + a.Email + ">"
}
