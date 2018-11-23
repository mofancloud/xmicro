package utils

// extention package of time
// adative milliseconds and seconds in uints
import (
	"time"
)

func Unix(sec int64, nsec int64) time.Time {
	//如果sec> (1970.1.1).Unix()*1000
	// 那么认为sec 是毫秒为单位
	//thres := time.Date(1970, 1, 1, 0, 0, 0, 0, nil)
	thres := time.Now()
	if sec > thres.Unix()*900 {
		sec = sec / 1000 //转换成秒
	}
	return time.Unix(sec, nsec)

}

//time.Time 转换成
func MilliSecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}
