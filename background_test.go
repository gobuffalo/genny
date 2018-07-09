package genny

import "context"

func (r *Suite) Test_Background_Context() {
	b := Background()
	r.Equal(b.Context(), context.Background())
}

func (r *Suite) Test_background_Context() {
	b := background{}
	r.Nil(b.Context())
}

func (r *Suite) Test_background_Parent() {
	b := background{}
	r.Nil(b.Parent())
}

func (r *Suite) Test_background_Cmd() {
	b := background{}
	r.Nil(b.Cmd())
}

func (r *Suite) Test_background_File() {
	b := background{}
	r.Nil(b.File())
}

func (r *Suite) Test_background_Logger() {
	b := background{}
	r.Nil(b.Logger())
}

func (r *Suite) Test_background_Run() {
	b := background{}
	r.NoError(b.Run())
}
