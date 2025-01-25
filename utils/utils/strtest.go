package utils

func getMaxStr(str string) string {
	if str == "" {
		return ""
	}
	_slice := []byte(str)
	_len := len(_slice)
	for {
		_start := 0
		if _len <= 255 {
			for {
				if (_start + _len) > len(_slice) {
					break
				}
				_testSlice := _slice[_start : _start+_len]
				_testChar := _testSlice[len(_testSlice)-1]
				if loopTest(_testSlice[0:len(_testSlice)-1], _testChar) {
					return string(_testSlice)
				}
				_start = _start + 1
			}
		}
		_len = _len - 1
	}
}

func loopTest(slice []byte, b byte) bool {
	for _, _v := range slice {
		if b == _v {
			return false
		}
	}
	_i := len(slice)
	if _i > 1 {
		return loopTest(slice[0:len(slice)-1], slice[len(slice)-1])
	}
	return true
}
