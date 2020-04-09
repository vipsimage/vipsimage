package operation

import (
	"testing"

	"github.com/vipsimage/vips"
)

func TestOperation_Execute(t *testing.T) {
	type fields struct {
		Thumbnail *Thumbnail
		Resize    *Resize
		Crop      *Crop
		SmartCrop *SmartCrop
		Watermark *Watermark
		Rotate    *Rotate
	}
	type args struct {
		img *vips.Image
	}

	img, _ := vips.NewFromFile("/Users/mo/Go/vipsimage/vipsimage/data/images/test.webp")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				Thumbnail: &Thumbnail{Width: 20},
				Resize:    nil,
				Crop:      nil,
				SmartCrop: nil,
				Watermark: nil,
				Rotate:    nil,
			},
			args: args{
				img: img,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			th := Operation{
				Thumbnail: tt.fields.Thumbnail,
				Resize:    tt.fields.Resize,
				Crop:      tt.fields.Crop,
				SmartCrop: tt.fields.SmartCrop,
				Watermark: tt.fields.Watermark,
				Rotate:    tt.fields.Rotate,
			}
			if err := th.Execute(tt.args.img); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
