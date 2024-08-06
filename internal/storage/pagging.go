package storage

const DefaultOffset int = 0
const DefaultLimit int = 30
const MaxLimit int = 1000

type Page struct {
	Limit  int
	Offset int
}

func Pagging(page, pageSize *int) Page {
	var limit int

	switch {
	case pageSize == nil:
		limit = DefaultLimit
	case *pageSize >= MaxLimit:
		limit = MaxLimit
	default:
		limit = *pageSize
	}

	offset := DefaultOffset
	if page != nil && *page > 0 {
		offset = (*page - 1) * limit
	}

	return Page{
		Limit:  limit,
		Offset: offset,
	}
}
