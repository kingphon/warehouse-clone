package service

import (
	"testing"

	"git.selly.red/Selly-Modules/3pl/partnerapi/tnc"
)

func Test_outboundRequestTNC_getTrackingCode(t *testing.T) {
	type fields struct {
		auth string
	}
	type args struct {
		trackingCode string
		tplCode      string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "ghtk code has dot char",
			fields: fields{},
			args: args{
				trackingCode: "S18893960.BO.DN2.S14.300357090",
				tplCode:      tnc.TPLCodeGHTK,
			},
			want: "300357090",
		},
		{
			name:   "ghtk code has no dot char",
			fields: fields{},
			args: args{
				trackingCode: "300357090",
				tplCode:      tnc.TPLCodeGHTK,
			},
			want: "300357090",
		},
		{
			name:   "other tpl code",
			fields: fields{},
			args: args{
				trackingCode: "300357090",
				tplCode:      tnc.TPLCodeGHN,
			},
			want: "300357090",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := outboundRequestTNC{
				auth: tt.fields.auth,
			}
			if got := s.getTrackingCode(tt.args.trackingCode, tt.args.tplCode); got != tt.want {
				t.Errorf("getTrackingCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
