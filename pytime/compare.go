package pytime

import "time"

// 判断t_ahead+delay_sec是否大于t_after，即过期
func IsExpired(t_ahead time.Time, t_after time.Time, delay_sec int) bool {
	min_diff := t_after.Sub(t_ahead).Seconds()
	if min_diff < 0 {
		//t_ahead>t_after
		return false
	}
	bret := float64(delay_sec)-min_diff >= 0
	return !bret
}
