package stripe

func defaultInt(actual int, def int) int {
	switch {
	case actual == 0:
		return def
	default:
		return actual
	}
}

func defaultUint64(actual uint64, def uint64) uint64 {
	switch {
	case actual == 0:
		return def
	default:
		return actual
	}
}
