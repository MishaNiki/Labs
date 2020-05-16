package spamdetection

import (
	"fmt"
	"strings"
	"sync"
)

// Word ...
type Word struct {
	Word       string
	AmountOK   int
	AmountSpam int
}

// Dictionary
type Dictionary struct {
	Dict map[string]*Word
	mu   *sync.Mutex
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		Dict: make(map[string]*Word),
		mu:   &sync.Mutex{},
	}
}

// AddSpam ...
func (d *Dictionary) AddSpam(word string) {
	v, ok := d.Dict[word]
	if ok {
		v.AmountSpam++
	} else {
		d.mu.Lock()
		d.Dict[word] = &Word{Word: word, AmountSpam: 1}
		d.mu.Unlock()
	}
}

// AddOK ...
func (d *Dictionary) AddOK(word string) {
	v, ok := d.Dict[word]
	if ok {
		v.AmountOK++
	} else {
		d.mu.Lock()
		d.Dict[word] = &Word{Word: word, AmountOK: 1}
		d.mu.Unlock()
	}
}

// AddWord ...
func (d *Dictionary) AddWord(word Word) {
	v, ok := d.Dict[word.Word]
	if ok {
		v.AmountOK += word.AmountOK
		v.AmountSpam += word.AmountSpam
	} else {
		d.mu.Lock()
		d.Dict[word.Word] = &word
		d.mu.Unlock()
	}
}

func (d *Dictionary) Print() {
	for _, v := range d.Dict {
		fmt.Println(v.Word, "\t count OK: ", v.AmountOK, "\t count spam: ", v.AmountSpam)
	}
}

// String() -  map Dictionary.Dict convert to string.
// Template - "Word AmountOK AmountSpam Word AmountOK AmountSpam ..."
func (d *Dictionary) String() string {
	var str string
	for _, v := range d.Dict {
		str += fmt.Sprintf("%v %v %v ", v.Word, v.AmountOK, v.AmountSpam)
	}
	return strings.Trim(str, " ")
}

func (d *Dictionary) GetProbabilitySpam(w string) float64 {
	v, ok := d.Dict[w]
	if !ok {
		return 0.5
	} else {
		return float64(v.AmountSpam) / float64(v.AmountSpam+v.AmountOK)
	}
}

func (d *Dictionary) GetProbabilityOK(w string) float64 {
	v, ok := d.Dict[w]
	if !ok {
		return 0.5
	} else {
		return float64(v.AmountOK) / float64(v.AmountSpam+v.AmountOK)
	}
}
