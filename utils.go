package main

func setPageBounds(limit int, page int, len int) (int, int) {
	if limit == 0 {
		return 0, len
	}
	lower := limit * page
	upper := lower + limit
	if len < lower {
		lower, upper = 0, 0
	} else if len < upper {
		upper = len
	}
	return lower, upper
}
