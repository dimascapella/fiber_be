package request

type NoteRequest struct {
	Code string `form:"code" validate:"required"`
	Note string `form:"note" validate:"required"`
	Qty  int    `form:"qty" validate:"required"`
}
