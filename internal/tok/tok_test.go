package tok

import (
	"embed"
	"reflect"
	"testing"
)

// codeDir contains a dir of src code for testing
//
//go:embed test_data
var codeDir embed.FS

func TestClient_Tokenize(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want []Token
	}{
		{
			"example 1",
			args{
				data: func() string {
					data, err := codeDir.ReadFile("test_data/example1.bk")
					if err != nil {
						t.Fatalf("Tokenize() failed to setup example: %s", err)
					}
					return string(data)
				}(),
			},
			[]Token{"import", " ", "echo", " ", "as", " ", "print", "\n",
				"\n",
				"#", " ", "main", " ", "is", " ", "the", " ", "entry", " ", "point", " ", "for", " ", "any", " ", "bk", " ", "script", "\n",
				"main", " ", ":", "(", "args", " ", ":", "[", "]", "string", ")", ":", " ", "{", "\n",
				" ", " ", " ", " ", "$", "print", "[", "\"", "Hello", "\"", ",", " ", "\"", "World", "\"", "]", "\n",
				"}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Tokenize(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.Tokenize() = \n%#v, want \n%#v", got, tt.want)
			}
		})
	}
}
