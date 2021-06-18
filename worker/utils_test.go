package worker

import "testing"

func TestToUtf8(t *testing.T) {
	type args struct {
		gbkString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{gbkString: `\u4eb2\uff0c\u4f60\u4e4b\u524d\u5df2\u7ecf\u7b7e\u8fc7\u4e86`},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToUtf8(tt.args.gbkString); got != tt.want {
				t.Errorf("ToUtf8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnicode2ZHCN(t *testing.T) {
	type args struct {
		unicode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{unicode: `\u4eb2\uff0c\u4f60\u4e4b\u524d\u5df2\u7ecf\u7b7e\u8fc7\u4e86`},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unicode2ZHCN(tt.args.unicode); got != tt.want {
				t.Errorf("Unicode2ZHCN() = %v, want %v", got, tt.want)
			}
		})
	}
}
