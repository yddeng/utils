package pselect

import (
	"errors"
	"reflect"
)

var errBadChannel = errors.New("event: Subscribe argument does not have sendable channel type")

type caseList []reflect.SelectCase

type RecvSelect struct {
	pList       []caseList
	defaultCase reflect.SelectCase
}

func MakeRecv(defaulted bool, chanSlice ...[]interface{}) *RecvSelect {
	if len(chanSlice) == 0 {
		return nil
	}

	s := new(RecvSelect)

	s.pList = make([]caseList, 0, len(chanSlice))
	for _, chans := range chanSlice {
		if len(chans) != 0 {
			list := make([]reflect.SelectCase, len(chans))
			for i, c := range chans {
				chanval := reflect.ValueOf(c)
				chantyp := chanval.Type()
				if chantyp.Kind() != reflect.Chan || chantyp.ChanDir()&reflect.RecvDir == 0 {
					panic(errBadChannel)
				}

				list[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: chanval}
			}
			s.pList = append(s.pList, list)
		}
	}

	if defaulted {
		s.defaultCase = reflect.SelectCase{Dir: reflect.SelectDefault}
	}

	return s
}

func (this *RecvSelect) Recv() (value interface{}, channel interface{}) {
	reflect.Select()
}
