package service

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	DefaultExpire    time.Duration = 15 * time.Second
	DefaultCleanup   time.Duration = 1 * time.Second
	DefaultKeyLength int           = 20
)

var FontStr []string = []string{"天", "往", "秋", "收", "冬", "藏", "闰", "余", "成", "岁", "律", "吕", "调", "阳", "云", "腾", "致", "雨", "露", "结"} //不重复的汉字

func NewPosition() *Position {
	p := &Position{}
	p.Data = cache.New(DefaultExpire, DefaultCleanup)
	return p
}

type Position struct {
	Data *cache.Cache
}

func generateRandomNumber(start int, end int, count int) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		num := r.Intn((end - start)) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}

		if !exist {
			nums = append(nums, num)
		}
	}

	return nums
}
func (p *Position) New(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	key := RandStringRunes(DefaultKeyLength)
	msg := &Msg{
		Count: 4,
		X:     []int{120, 240, 360, 480},
		Y:     []int{30, 130, 80, 180},
	}
	fontIndex := generateRandomNumber(1, len(FontStr), 4)
	for _, f := range fontIndex {
		msg.Font = append(msg.Font, FontStr[f])
	}
	msg.Use = generateRandomNumber(0, msg.Count, 2)
	msg.Y = generateRandomNumber(30, 180, 4)

	data, _ := json.Marshal(msg)
	data, _ = json.Marshal(&Wrapper{
		Key:  key,
		Data: EncodeStr(string(data)),
	})
	p.Data.SetDefault(key, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(data)
}
func (p *Position) Post(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	data, _ := ioutil.ReadAll(r.Body)
	wrapper := &Wrapper{}
	json.Unmarshal(data, wrapper)
	val, has := p.Data.Get(wrapper.Key)
	if !has {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var points []Point
	err := json.Unmarshal([]byte(DecodeStr(wrapper.Data)), &points)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	msg, ok := val.(*Msg)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	for i, p := range points {
		tmpX := float64(msg.X[msg.Use[i]])
		tmpY := float64(msg.Y[msg.Use[i]])
		// fmt.Println("E:", tmpX, tmpY, p.X, p.Y)
		if p.X > tmpX && p.X < tmpX+20 && p.Y > tmpY-20 && p.Y < tmpY {
			continue
		}
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	data, _ = json.Marshal(&Result{
		Key:  wrapper.Key,
		Code: 1000,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(data)
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
type Result struct {
	Key  string `json:"key"`
	Code int    `json:"code"`
}
type Wrapper struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}
type Msg struct {
	Count int      `json:"count"`
	X     []int    `json:"x"`
	Y     []int    `json:"y"`
	Font  []string `json:"font"`
	Use   []int    `json:"use"`
}
