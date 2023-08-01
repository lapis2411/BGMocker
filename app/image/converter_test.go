package image

import (
	"testing"
)

func TestToPNGImageResponse(t *testing.T) {
	type args struct {
		img ImageBase64
	}
	tests := []struct {
		name    string
		args    args
		want    ImageResponse
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				img: ImageBase64{
					image: "test_base_64",
				},
			},
			want: ImageResponse{
				ImageData: "data:image/png;base64,test_base_64",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.args.img.ToPNGImageResponse()
			if got != tt.want {
				t.Errorf("ToPNGImageResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
