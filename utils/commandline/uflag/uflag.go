package uflag

import (
	"flag"
	"strconv"
	"sync"
	"time"
)

var flagNames = &FlagNames{m: make(map[string]int)}

type FlagNames struct {
	m    map[string]int
	lock sync.RWMutex
}

func (fl *FlagNames) Set(key string, val int) {
	fl.lock.Lock()
	defer fl.lock.Unlock()
	fl.m[key] = val
}

func (fl *FlagNames) Get(key string) (val int, isOk bool) {
	fl.lock.RLock()
	defer fl.lock.RUnlock()
	val, isOk = fl.m[key]
	return
}

func flagName(expect string) (actual string) {
	num, ok := flagNames.Get(expect)
	if !ok {
		flagNames.Set(expect, 1)
		return expect
	}

	flagNames.Set(expect, num+1)
	return expect + strconv.Itoa(num)
}

func BoolVar(p *bool, expectName string, defaultValue bool, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.BoolVar(p, actualName, defaultValue, usage)
	return
}

func Bool(expectName string, defaultValue bool, usage string) (p *bool, actualName string) {
	actualName = flagName(expectName)
	return flag.Bool(actualName, defaultValue, usage), actualName
}

func IntVar(p *int, expectName string, defaultValue int, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.IntVar(p, actualName, defaultValue, usage)
	return
}

func Int(expectName string, defaultValue int, usage string) (p *int, actualName string) {
	actualName = flagName(expectName)
	return flag.Int(actualName, defaultValue, usage), actualName
}

func Int64Var(p *int64, expectName string, defaultValue int64, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.Int64Var(p, actualName, defaultValue, usage)
	return
}

func Int64(expectName string, defaultValue int64, usage string) (p *int64, actualName string) {
	actualName = flagName(expectName)
	return flag.Int64(actualName, defaultValue, usage), actualName
}

func UintVar(p *uint, expectName string, defaultValue uint, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.UintVar(p, actualName, defaultValue, usage)
	return
}

func Uint(expectName string, defaultValue uint, usage string) (p *uint, actualName string) {
	actualName = flagName(expectName)
	return flag.Uint(actualName, defaultValue, usage), actualName
}

func Uint64Var(p *uint64, expectName string, defaultValue uint64, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.Uint64Var(p, actualName, defaultValue, usage)
	return
}

func Uint64(expectName string, defaultValue uint64, usage string) (p *uint64, actualName string) {
	actualName = flagName(expectName)
	return flag.Uint64(actualName, defaultValue, usage), actualName
}

func StringVar(p *string, expectName string, defaultValue string, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.StringVar(p, actualName, defaultValue, usage)
	return
}

func String(expectName string, defaultValue string, usage string) (p *string, actualName string) {
	actualName = flagName(expectName)
	return flag.String(actualName, defaultValue, usage), actualName
}

func Float64Var(p *float64, expectName string, defaultValue float64, usage string) (actualName string) {
	actualName = flagName(expectName)
	flag.Float64Var(p, actualName, defaultValue, usage)
	return
}

func Float64(expectName string, defaultValue float64, usage string) (p *float64, actualName string) {
	actualName = flagName(expectName)
	return flag.Float64(actualName, defaultValue, usage), actualName
}

func DurationVar(p *time.Duration, expectName string, defaultValue time.Duration,
	usage string) (actualName string) {

	actualName = flagName(expectName)
	flag.DurationVar(p, actualName, defaultValue, usage)
	return
}

func Duration(expectName string, defaultValue time.Duration, usage string) (p *time.Duration, actualName string) {
	actualName = flagName(expectName)
	return flag.Duration(actualName, defaultValue, usage), actualName
}

func Parse() {
	flag.Parse()
}
