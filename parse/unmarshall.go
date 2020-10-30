package parse

func Unmarshall(configuration interface{}, methods ...Unmarshaler) (err error) {
	for _, m := range methods {
		err = m.Unmarshal(configuration)
		if err != nil {
			return
		}
	}
	return
}
