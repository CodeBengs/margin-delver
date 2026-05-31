package pagination

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100
)

type Result struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

type Meta struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

func Normalize(page int, limit int) (int, int) {
	if page < 1 {
		page = DefaultPage
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return page, limit
}

func NewResult(data interface{}, page int, limit int, total int64) Result {
	return Result{
		Data: data,
		Meta: Meta{
			Page:  page,
			Limit: limit,
			Total: total,
		},
	}
}
