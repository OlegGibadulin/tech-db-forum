package uniq

func RemoveDuplicates(elements []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, elem := range elements {
		if encountered[elem] == false {
			encountered[elem] = true
			result = append(result, elem)
		}
	}
	return result
}
