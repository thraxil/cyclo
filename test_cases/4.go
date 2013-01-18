package four

func a() {
	for i := range 10 {
		if i % 2 {
			i + 1
		}
	}
}
