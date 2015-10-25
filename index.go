package gostruct

import (
	"fmt"
	"log"
	"os"
)

type Index struct {
	Elems []string
	Ids   map[string]int
}

func NewIndex() Index {
	this := Index{
		Elems: make([]string, 0, 10000),
		Ids:   make(map[string]int),
	}
	return this
}

func (this *Index) AddElem(elem string) {
	eid := len(this.Ids)
	this.Elems = append(this.Elems, elem)
	this.Ids[elem] = eid
}

func (this *Index) GetIdAndAddElemIfNotExists(elem string) int {
	var eid int = -1
	if _, ok := this.Ids[elem]; !ok {
		eid = len(this.Ids)
		this.Elems = append(this.Elems, elem)
		this.Ids[elem] = eid
	} else {
		eid = this.Ids[elem]
	}
	return eid
}

func (this *Index) HasElem(elem string) bool {
	_, ok := this.Ids[elem]
	return ok
}

func (this *Index) GetId(elem string) int {
	//         return this.Ids[elem]
	id, ok := this.Ids[elem]
	if !ok {
		log.Println(fmt.Sprintf("%s doesn't exist", elem))
		os.Exit(1)
	}
	return id
}
