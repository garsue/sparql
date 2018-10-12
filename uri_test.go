package sparql

import "testing"

func TestIRI_Ref(t *testing.T) {
	tests := []struct {
		name string
		i    URI
		want string
	}{
		{
			name: "replace <",
			i:    "<",
			want: "<%3C>",
		},
		{
			name: "replace >",
			i:    ">",
			want: "<%3E>",
		},
		{
			name: `replace "`,
			i:    `"`,
			want: "<%22>",
		},
		{
			name: "replace space",
			i:    " ",
			want: "<%20>",
		},
		{
			name: "replace {",
			i:    "{",
			want: "<%7B>",
		},
		{
			name: "replace }",
			i:    "}",
			want: "<%7D>",
		},
		{
			name: "replace |",
			i:    "|",
			want: "<%7C>",
		},
		{
			name: `replace \`,
			i:    `\`,
			want: "<%5C>",
		},
		{
			name: "replace ^",
			i:    "^",
			want: "<%5E>",
		},
		{
			name: "replace “",
			i:    "`",
			want: "<%60>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Ref(); got != tt.want {
				t.Errorf("URI.Ref() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrefixedName_Ref(t *testing.T) {
	p := PrefixedName("foo")
	want := "foo"
	if got := p.Ref(); got != want {
		t.Errorf("PrefixedName.Ref() = %v, want %v", got, want)
	}
}
