package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"time"
)

type Experiment struct {
	P0     []float64   `json:"p0"`
	MTX    [][]float64 `json:"mtx"`
	Amount int         `json:"amount"`
	End    int         `json:"end"`

	Average []int
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "./config.json", " path to config")
}

func main() {
	flag.Parse()

	exper, err := ParseExperiment(configPath)
	if err != nil {
		panic(err)
	}

	exper.Run()

	var out, k, pog float64
	for i, v := range exper.Average {
		if i == 0 {
			out += float64(v)
			k = out / float64(i+1)
		} else {
			out += float64(v)
			pog += math.Abs(k - out/float64(i+1))
			k = out / float64(i+1)
		}
	}
	out /= float64(len(exper.Average))
	fmt.Println("Мат ожидание ", out)
	fmt.Println("Погрешность: ", pog/float64(len(exper.Average)))

	var disp float64
	for _, v := range exper.Average {
		disp += (float64(v) - out) * (float64(v) - out)
	}
	disp /= float64(len(exper.Average))

	fmt.Println("Дисперсия :", disp, "Среднеквадратичное :", math.Sqrt(disp))

	fmt.Println("/////Данные для нахождения распределния/////")

	rasp := make([]int, 30)

	for _, v := range exper.Average {
		rasp[v]++
	}

	for _, v := range rasp {
		fmt.Println(v)
	}

}

// Run ...
func (ex *Experiment) Run() {
	ex.Average = make([]int, ex.Amount)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < ex.Amount; i++ {
		s0 := ex.InitialState(rand.Float64())
		for s0 != ex.End {
			s0 = ex.State(s0, rand.Float64())
			ex.Average[i]++
		}

		var out float64
		for j := 0; j <= i; j++ {
			out += float64(ex.Average[j])
		}
		out /= float64(i + 1)
		fmt.Println(i+1, out)
	}
}

// InitialState ...
func (ex *Experiment) InitialState(p float64) int {
	var sm1 float64
	for i := 0; i < len(ex.P0); i++ {
		sm1 = 0
		for j := 0; j <= i; j++ {
			sm1 += ex.P0[j]
		}
		if p < sm1 {
			//log.Println("InitialState", i, p, sm1)
			return i
		}
	}
	return 0
}

// State ...
func (ex *Experiment) State(s int, p float64) int {
	var sm1 float64
	for i := 0; i < len(ex.MTX[s]); i++ {
		sm1 = 0
		for j := 0; j <= i; j++ {
			sm1 += ex.MTX[s][j]
		}
		if p < sm1 {
			//log.Println("State\t", i, p)
			return i
		}
	}
	return 0
}

// ParseExperiment ...
func ParseExperiment(path string) (*Experiment, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 2048)

	var lenBuf int
	for {
		len, e := file.Read(data)
		lenBuf += len
		if e == io.EOF {
			break
		}
	}

	exper := &Experiment{}
	err = json.Unmarshal(data[:lenBuf], exper)
	if err != nil {
		return nil, err
	}
	return exper, nil
}
