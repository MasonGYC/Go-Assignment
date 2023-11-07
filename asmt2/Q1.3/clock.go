package main

// check whether c1 = c2
func clockEqual(c1 []int, c2 []int) bool {
	if len(c1) == len(c2) {
		for i := 0; i < len(c1); i++ {
			if c1[i] != c2[i] {
				// if there is !=, not equal
				return false
			}
		}
		// if all =, equal
		return true
	} else {
		return false
	}
}

// check whether c1 < c2
func clockLess(c1 []int, c2 []int) bool {
	// check comparability
	if len(c1) == len(c2) {
		// loop thru to see if all <=
		for i := 0; i < len(c1); i++ {
			if c1[i] > c2[i] {
				// if there is greater than, false
				return false
			}
		}
		// loop thru to see if there exist <
		for i := 0; i < len(c1); i++ {
			if c1[i] < c2[i] {
				// since all <=, if one <, then true
				return true
			}
		}
		// since all <=, there's no <, then false (equal)
		return false
	} else {
		return false
	}
}
