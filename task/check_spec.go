// @Author: Perry
// @Date  : 2020/10/18
// @Desc  : 检查任务的时间参数是否合法

package task

import (
	"errors"
	"regexp"
	"strings"
)

var (
	regCronExpression    = regexp.MustCompile(`^[\S]+ [\S]+ [\S]+ [\S]+ [\S]+$`)
	regCronPredefined    = []string{"@yearly", "@annually", "@monthly", "@weekly", "@daily", "@midnight", "@hourly"}
	regCronIntervalTime  = regexp.MustCompile(`^([\d]+h)?([\d]+m)?([\d]+s)?$`)
	regCronMinHoursMonth = regexp.MustCompile(`^[\d]{1,2}|\*|/|,|-$`)
	regCronDOMW          = regexp.MustCompile(`^[\d]{1,2}|\*|/|,|-|\?$`)
)

// 检查传入的时间参数是否合法
func checkSpec(spec string) error {
	if len(spec) == 0 {
		return errors.New("check spec failed,spec is empty")
	}
	if strings.HasPrefix(spec, "@every ") {
		interval := spec[6:]
		interval = strings.Trim(interval, " ")
		if regCronIntervalTime.MatchString(interval) && interval != "" {
			return nil
		}
		return errors.New("check spec failed,interval value =" + spec)
	}
	if strings.HasPrefix(spec, "@") {
		for _, v := range regCronPredefined {
			if v == spec {
				return nil
			}
		}
		return errors.New("check spec failed,predefined value =" + spec)
	}
	if regCronExpression.MatchString(spec) {
		expList := strings.Split(spec, " ")
		minutes := expList[0]
		hours := expList[1]
		dayOfMonth := expList[2]
		month := expList[3]
		dayOfWeek := expList[4]
		if regCronMinHoursMonth.MatchString(minutes) &&
			regCronMinHoursMonth.MatchString(hours) &&
			regCronMinHoursMonth.MatchString(month) &&
			regCronDOMW.MatchString(dayOfMonth) &&
			regCronDOMW.MatchString(dayOfWeek) {
			return nil
		}
		return errors.New("check spec failed,cron value =" + spec)
	}
	return errors.New("check spec failed,not matched value =" + spec)
}
