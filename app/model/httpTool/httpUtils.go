package httpTool

func Map2Str(m map[string]string) string {
	result := ""
	for k, v := range m {
		result += "&" + k + "=" + v
	}
	return result[1:]
}
func MakeHeader() map[string]string {
	m := map[string]string{}
	m["Content-Type"] = "application/x-www-form-urlencoded; charset=UTF-8"
	m["X-Requested-With"] = "XMLHttpRequest"
	return m
}
