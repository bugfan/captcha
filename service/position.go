package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	DefaultExpire    time.Duration = 15 * time.Second
	DefaultCleanup   time.Duration = 1 * time.Second
	DefaultKeyLength int           = 20
	FontStr          string        = "天地玄黄宇宙洪荒日月盈昃辰宿列张寒来暑往秋收冬藏闰余成岁律吕调阳云腾致雨露结为霜金生丽水玉出昆冈剑号巨阙珠称夜光果珍李柰菜重芥姜海咸河淡鳞潜羽翔龙师火帝鸟官人皇始制文字乃服衣裳推位让国有虞陶唐吊民伐罪周发殷汤坐朝问道垂拱平章爱育黎首臣伏戎羌遐迩体率宾归王" //不重复的汉字

)

func NewPosition() *Position {
	p := &Position{}
	p.Data = cache.New(DefaultExpire, DefaultCleanup)
	return p
}

type Position struct {
	Data *cache.Cache
}

func (p *Position) New(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	key := RandStringRunes(DefaultKeyLength)
	msg := &Msg{
		Key:   key,
		Count: 4,
		Font:  []string{"昃", "辰", "宿", "列"},
		X:     []int{120, 240, 360, 480},
		Y:     []int{30, 130, 80, 180},
		Use:   []int{0, 2},
	}
	data, _ := json.Marshal(msg)
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
	check := &positionCheck{
		Point: make([]*Point, 0),
	}
	json.Unmarshal(data, &check)
	val, has := p.Data.Get(check.Key)
	if !has {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	msg, ok := val.(*Msg)
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	for i, p := range check.Point {
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
		Key:  check.Key,
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
type positionCheck struct {
	Key   string   `json:"key"`
	Point []*Point `json:"point"`
}
type Result struct {
	Key  string `json:"key"`
	Code int    `json:"code"`
}
type Msg struct {
	Key   string   `json:"key"`
	Count int      `json:"count"`
	X     []int    `json:"x"`
	Y     []int    `json:"y"`
	Font  []string `json:"font"`
	Use   []int    `json:"use"`
}
