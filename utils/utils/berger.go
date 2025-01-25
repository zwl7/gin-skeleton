package utils

import (
	"fmt"
	"log"
	"math/rand"
	"reflect"
	"time"
)

// 编排算法工具函数
type AlgorithmSort struct {
	DataArray [][]map[int64]int64 //对打：贝格尔编排后的数据源、逆时针数据源
	SortArray [][]int64           //抽签：贝格尔编排、随机编排的数据源
	IsOddNum  bool                //是否奇数，默认不是，false
}

// 操作实例
func NewAlgorithmSort() *AlgorithmSort {
	return &AlgorithmSort{}
}

// 贝格尔经典算法编排，n为队伍数
func (that *AlgorithmSort) Berger(n int64) *AlgorithmSort {
	//遇单数队,最后一位数字补为O成为偶数
	m := int64(0)
	//是否奇数，存在轮空位
	if n%2 == 0 {
		m = n
	} else {
		m = n + 1
		that.IsOddNum = true
	}
	//环形起始位置，对阵A
	a := int64(1)
	//环形起始位置，对阵B
	b := int64(1)
	//队伍偶数量，包括轮空位
	c := m
	//索引分成均等两边
	index := int64(1)
	//轮数
	loop := 0
	//结果
	if that.IsOddNum {
		that.DataArray = make([][]map[int64]int64, n)
	} else {
		that.DataArray = make([][]map[int64]int64, n-1)
	}
	for i := int64(1); i <= (m-1)*(m/2); i++ { // 0<15
		if a >= m { //1>=6
			a = 1
		}
		if index > m/2 { //1>3
			index = 1
		}
		if index == 1 { //1==1
			loop++
			if i == 1 {
				b = m
			} else {
				b = a
			}
			if that.IsOddNum == true {
				c = 0
			} else {
				c = m
			}
			that.DataArray[loop-1] = []map[int64]int64{}
			println("===========第" + ToString(loop) + "轮===========")
			if (((i - 1) / (m / 2)) % 2) == 0 {
				that.DataArray[loop-1] = append(that.DataArray[loop-1], map[int64]int64{a: c})
				log.Printf(ToString(a) + " - " + ToString(c))
			} else {
				that.DataArray[loop-1] = append(that.DataArray[loop-1], map[int64]int64{c: a})
				log.Printf(ToString(c) + " - " + ToString(a))
			}
		} else if index > 1 && index < m/2 {
			if b > 1 {
				b--
			} else {
				b = m - 1
			}
			that.DataArray[loop-1] = append(that.DataArray[loop-1], map[int64]int64{a: b})
			log.Printf(ToString(a) + " - " + ToString(b))
		} else {
			if b-1 == 0 {
				c = m - 1
			} else {
				c = b - 1
			}
			that.DataArray[loop-1] = append(that.DataArray[loop-1], map[int64]int64{a: c})
			log.Printf(ToString(a) + " - " + ToString(c))
		}
		index++
		a++
	}
	return that
}

// 单循环赛贝格尔编排法
func (that *AlgorithmSort) BergerArrangement(nAmount int) *AlgorithmSort {
	if nAmount < 2 || nAmount > 90 {
		return nil
	}
	// 队伍数量
	nFixAmount := nAmount
	// 最后一支队伍的编号
	nLastPlayerNo := nAmount
	// 奇数队伍，补上一支虚拟的队伍，最后一支队伍的编号为0
	if nAmount%2 != 0 {
		nFixAmount++
		nLastPlayerNo = 0
		that.IsOddNum = true
		that.DataArray = make([][]map[int64]int64, nAmount)
	} else {
		that.DataArray = make([][]map[int64]int64, nAmount-1)
	}
	//轮数
	nMaxRound := nFixAmount - 1
	nHalfAmount := nFixAmount / 2
	// 移动的间隔
	nStep := 0
	if nFixAmount <= 4 {
		nStep = 1
	} else {
		nStep = (nFixAmount-4)/2 + 1
	}
	nRound := 1
	nFirstPlayerPos := 1
	nLastPlayerPos := 1
	result := [100][200]int{
		{0, 0},
	}
	for nRound <= nMaxRound {
		// 每次最后一个玩家的位置需要左右对调
		nLastPlayerPos = nFixAmount + 1 - nLastPlayerPos
		if nRound == 1 {
			nFirstPlayerPos = 1
		} else {
			nFirstPlayerPos = (nFirstPlayerPos + nStep) % (nFixAmount - 1)
			if nFirstPlayerPos == 0 {
				nFirstPlayerPos = nFixAmount - 1
			}
			if nFirstPlayerPos == nLastPlayerPos {
				nFirstPlayerPos = nFixAmount + 1 - nLastPlayerPos
			}
		}
		for i := 1; i <= nHalfAmount; i++ {
			nPos := [2]int{i, nFixAmount - i + 1}
			nPlayer := [2]int{0, 0}
			for j := 0; j < 2; j++ {
				if nPos[j] == nLastPlayerPos {
					nPlayer[j] = nLastPlayerNo
				} else if nPos[j] < nFirstPlayerPos {
					nPlayer[j] = nFixAmount - nFirstPlayerPos + nPos[j]
				} else {
					nPlayer[j] = nPos[j] - nFirstPlayerPos + 1
				}
				(result)[i-1][(nRound-1)*2+j] = nPlayer[j]
			}
		}
		nRound++
	}
	_mp := make([][]map[int]int, nAmount+1)
	fmt.Printf("%6s\n", ToString(nAmount)+"个队编排如下：")
	for i := 1; i <= nMaxRound; i++ {
		if i == 1 {
			fmt.Printf("%3s%-3d|", "r", i)
		} else {
			fmt.Printf("%4s%-3d|", "r", i)
		}
		_mp[i] = []map[int]int{}
		that.DataArray[i-1] = []map[int64]int64{}
	}
	fmt.Printf("\n")
	for i := 0; i < nHalfAmount; i++ {
		for j := 0; j < nMaxRound; j++ {
			_mp[i] = append(_mp[i], map[int]int{(result)[i][j*2]: (result)[i][j*2+1]})
			fmt.Printf("%-2d-%2d | ", (result)[i][j*2], (result)[i][j*2+1])
		}
		fmt.Printf("\n")
	}
	//
	for i := 0; i < nMaxRound; i++ {
		for j := 0; j < nHalfAmount; j++ {
			that.DataArray[i] = append(that.DataArray[i], map[int64]int64{ToInt64((result)[j][i*2]): ToInt64((result)[j][i*2+1])})
		}
	}
	fmt.Printf("\n\n")
	return that
}

// 返回一个数最小2次方的幂和指数
func (that *AlgorithmSort) GetMin2Power(n int64) (int64, int) {
	count := 0
	index := int64(1)
	if n == 1 {
		return 0, 0
	}
	for n > index {
		count++
		index = index * 2
	}
	return index, count
}

// 抽签：贝格尔算法
func (that *AlgorithmSort) BergerSort(n int64) *AlgorithmSort {
	if n < 3 {
		that.SortArray = [][]int64{}
		return that
	}
	h := int64(0)
	if n%2 == 1 {
		h = 0
	} else {
		h = n
	}
	num := int64(0)
	if h > 0 {
		num = n
	} else {
		num = n + 1
	}
	average := num / 2
	res1 := make([][]int64, 0)
	for i := int64(0); i < average; i++ {
		top := i + 1
		temp := make([]int64, 0)
		for j := int64(0); j < num-1; j++ {
			if top > num-1 {
				top = 1
			}
			temp = append(temp, top)
			top++
		}
		temp = append(temp, h)
		res1 = append(res1, temp)
	}
	res2 := make([][]int64, 0)
	for i := average; i < n; i++ {
		top := int64(0)
		if top+2 > num-1 {
			top = 1
		} else {
			top = top + 2
		}
		temp := make([]int64, 0)
		temp = append(temp, h)
		for j := int64(0); j < num-1; j++ {
			if top > num-1 {
				top = 1
			}
			temp = append(temp, top)
			top++
		}
		res2 = append(res2, temp)
	}
	that.SortArray = make([][]int64, 0)
	for _k, _v := range res1 {
		that.SortArray = append(that.SortArray, _v)
		if _k <= len(res2) {
			that.SortArray = append(that.SortArray, res2[_k])
		}
	}
	return that
}

// 抽签：随机，将n个人数分k组,随机生成,例如32个人分成4组，每组随机排列8人
func (that *AlgorithmSort) RandSort(n int64, k int64) *AlgorithmSort {
	arr := make([]int64, 0)
	for i := int64(0); i < n; i++ {
		arr = append(arr, i+1)
	}
	that.Shuffle(arr)
	that.SortArray = make([][]int64, k)
	for _k, _v := range arr {
		if that.SortArray[ToInt64(_k)%k] == nil {
			that.SortArray[ToInt64(_k)%k] = make([]int64, 0)
		}
		that.SortArray[ToInt64(_k)%k] = append(that.SortArray[ToInt64(_k)%k], _v)
	}
	return that
}

// 洗牌
func (that *AlgorithmSort) Shuffle(slice []int64) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

/*
	  锁定组内的"种子位"，其他数据随机排列组合
	  调用方法：LockSeedRandOrder(24, 3, [][]int{
					{2, 6},  //第一组的1、2号位置排种子位(2, 6)
					{4},     //第二组的1号位置排种子位(4)
	                {0,3},   //第二组的2号位置种子位排(3)
				})
*/
func (that *AlgorithmSort) LockSeedSort(seedData ...[][]int64) *AlgorithmSort {
	dataTemp := make([]int64, 0)
	seedTemp := make([]int64, 0)
	if len(seedData) > 0 {
		_seedData := seedData[0]
		//先随机编排
		for i := 0; i < len(that.SortArray); i++ {
			for j := 0; j < len(that.SortArray[i]); j++ {
				if that.SortArray[i][j] != 0 {
					dataTemp = append(dataTemp, that.SortArray[i][j])
				}
			}
		}
		//制作坑位，标识为种子位
		for i := 0; i < len(_seedData); i++ {
			seedDataItemLength := len(_seedData[i])
			for j := 0; j < seedDataItemLength; j++ {
				//过滤0，即跳过该位置
				if _seedData[i][j] != 0 {
					that.SortArray[i][j] = 0
					seedTemp = append(seedTemp, _seedData[i][j])
				}
			}
		}
		//调补
		for i := 0; i < len(dataTemp); i++ {
			if that.Contains(seedTemp, dataTemp[i]) != -1 {
				dataTemp = append(dataTemp[:i], dataTemp[i+1:]...)
				i--
			}
		}
		//填补
		for i := 0; i < len(that.SortArray); i++ {
			for j := 0; j < len(that.SortArray[i]); j++ {
				if that.SortArray[i][j] != 0 {
					that.SortArray[i][j] = dataTemp[0]
					if len(dataTemp) > 0 {
						dataTemp = append(dataTemp[:0], dataTemp[1:]...)
					}
				} else {
					that.SortArray[i][j] = _seedData[i][j]
				}
			}
		}
	}
	return that
}

// 逆时针轮转法
func (that *AlgorithmSort) RotaryGet(n int) *AlgorithmSort {
	i := 0
	j := int64(1)
	//先判断队伍是奇数还是偶数
	if (n & 1) == 1 {
		n++
		i = 1
		that.IsOddNum = true
	}
	//偶数循环N-1次,奇数循环N次
	array := make([]int64, n)
	//队伍赋值，若队伍为奇数，首位赋值就跳过，且冗余值为0
	for ; i < n; i++ {
		array[i] = j
		j++
	}
	//循环编排开始
	return that.move(array, n/2, n)
}

// 移动索引位置
func (that *AlgorithmSort) move(array []int64, t int, n int) *AlgorithmSort {
	that.DataArray = make([][]map[int64]int64, n-1)
	//t为圆圈的中间位置下标
	length := n
	for i := 0; i < length-1; i++ {
		test1 := t
		test2 := 0
		data := make([]int64, length)
		data[0] = array[0]
		that.DataArray[i] = []map[int64]int64{}
		println("第" + ToString(i+1) + "轮比赛：")
		for test1 <= length-1 {
			if array[test2] != 0 && test2 == 0 {
				that.DataArray[i] = append(that.DataArray[i], map[int64]int64{array[test2]: array[test1]})
				println(ToString(array[test2]) + " VS " + ToString(array[test1]))
			} else if test2 != 0 {
				that.DataArray[i] = append(that.DataArray[i], map[int64]int64{array[t+test2]: array[t-test2]})
				println(ToString(array[t+test2]) + " VS " + ToString(array[t-test2]))
			} else {
				//轮空位
				that.DataArray[i] = append(that.DataArray[i], map[int64]int64{array[t+test2]: 0})
			}
			//下面算法是为下一轮编排所用圆圈赋值
			if test1 == length-1 {
				data[1] = array[test1]
				data[test2+1] = array[test2]
			} else {
				data[test1+1] = array[test1]
				if test2 != 0 {
					data[test2+1] = array[test2]
				}
			}
			test1++
			test2++
		}
		array = data
	}
	return that
}

// 对打：将索引位置替换成实际位置
func (that *AlgorithmSort) ToSequence(sequence []int64, isDontZero ...bool) (sequenceArray []map[int64]int64) {
	nAmount := len(that.DataArray)
	if !that.IsOddNum && ((len(sequence) - 1) != nAmount) {
		return []map[int64]int64{}
	}
	if that.IsOddNum && len(sequence) != nAmount {
		return []map[int64]int64{}
	}
	//是否不需要轮空位
	zeroPostion := false
	if len(isDontZero) > 0 && isDontZero[0] == true {
		zeroPostion = true
	}
	sequenceArray = make([]map[int64]int64, nAmount)
	for _k, _v := range that.DataArray {
		for _, _v2 := range _v {
			for _k3, _v3 := range _v2 {
				if sequenceArray[_k] == nil {
					sequenceArray[_k] = make(map[int64]int64)
				}
				if _index := _k3 - 1; _index < 0 {
					if zeroPostion == false {
						sequenceArray[_k][0] = sequence[_v3-1]
					}
					continue
				}
				if _index := _v3 - 1; _index < 0 {
					if zeroPostion == false {
						sequenceArray[_k][sequence[_k3-1]] = 0
					}
					continue
				}
				sequenceArray[_k][sequence[_k3-1]] = sequence[_v3-1]
			}
		}
	}
	return
}

// 抽签：将索引位置替换成实际位置
func (that *AlgorithmSort) ToSortSequence(sequence []int64, groupNum int64) [][]int64 {
	nAmount := len(that.SortArray)
	if groupNum != ToInt64(nAmount) {
		return [][]int64{}
	}
	for _k, _v := range that.SortArray {
		for _k2, _v2 := range _v {
			if _v2-1 < 0 {
				continue
			}
			if len(sequence) < ToInt(_v2) {
				continue
			}
			that.SortArray[_k][_k2] = sequence[_v2-1]
		}
	}
	return that.SortArray
}

// Contains Returns the index position of the val in array
func (that *AlgorithmSort) Contains(array interface{}, val interface{}) (index int) {
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		{
			s := reflect.ValueOf(array)
			for i := 0; i < s.Len(); i++ {
				if reflect.DeepEqual(val, s.Index(i).Interface()) {
					index = i
					return
				}
			}
		}
	}
	return
}

// ContainsString Returns the index position of the string val in array
func (that *AlgorithmSort) ContainsString(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// 根据人数、组数，返回每轮的对打的人所在的轮次
func (that *AlgorithmSort) GetByRoundIds(peopleNum int64, groupNum int64) []int64 {
	//轮次
	_roundNum := peopleNum - 1
	//每轮
	_everyRound := peopleNum / 2
	_slice := make([]int64, 0)
	if (peopleNum & 1) == 1 {
		_roundNum = peopleNum
		_everyRound = (peopleNum + 1) / 2
		that.IsOddNum = true
	}
	//组数
	for i := int64(0); i < groupNum; i++ {
		//轮次
		for j := int64(0); j < _roundNum; j++ {
			//每小组在第几轮
			for k := int64(0); k < _everyRound; k++ {
				_slice = append(_slice, j+1)
			}
		}
	}
	return _slice
}

// 对阵图
type pvp struct {
	Letter   string `json:"letter"`   //编号
	Position int64  `json:"position"` //位置
}

func (that *AlgorithmSort) MatchChart(x [][]int64) Mp {
	if len(x) == 0 {
		x = [][]int64{
			{1, 2, 3, 4}, //第一组的1-4名
			{1, 2, 3, 4}, //第二组的1-4名
			{1, 2, 3, 4}, //第三组的1-4名
			{1, 2, 3, 4}, //第四组的1-4名
		}
	}
	//分组
	abcdIndex := 'A'
	abcd := map[string][]int64{} //字母
	for _, v := range x {
		abcd[fmt.Sprintf("%c", abcdIndex)] = v
		abcdIndex = abcdIndex + 1
	}
	upperHalf := map[string][]int64{}
	lowerHalf := map[string][]int64{}
	length := len(abcd)
	_az := '0'
	for i, j := 'A', 0; j < length; j++ {
		//上半区
		if j < length/2 {
			upperHalf[fmt.Sprintf("%c", i)] = abcd[fmt.Sprintf("%c", i)]
		}
		//下半区
		if j >= length/2 {
			if _az == '0' {
				_az = i
			}
			lowerHalf[fmt.Sprintf("%c", i)] = abcd[fmt.Sprintf("%c", i)]
		}
		i++
	}
	_func := func(half map[string][]int64, az int32) (pvpA [][]*pvp) {
		_upperHalfA := map[string][]int64{}
		_upperHalfB := map[string][]int64{}
		for i, j := az, 0; j < len(half); j++ {
			//上半区
			if j < len(half)/2 {
				_upperHalfA[fmt.Sprintf("%c", i)] = half[fmt.Sprintf("%c", i)]
			}
			//下半区
			if j >= len(half)/2 {
				_upperHalfB[fmt.Sprintf("%c", i)] = half[fmt.Sprintf("%c", i)]
			}
			i++
		}
		fmt.Println("============upperHalfA=============", _upperHalfA)
		fmt.Println("============upperHalfB=============", _upperHalfB)
		//对战
		pvpA = make([][]*pvp, 0)
		//字母移位
		incr := len(_upperHalfA)
		for k1, v1 := range _upperHalfA {
			for pvpIndexA, pvpIndexB := 0, len(v1)-1; pvpIndexA < len(v1); pvpIndexA++ {
				item := make([]*pvp, 0)
				item = append(item, &pvp{
					Letter:   k1,
					Position: v1[pvpIndexA],
				})
				_letter := that.az(k1, incr)
				item = append(item, &pvp{
					Letter:   _letter,
					Position: _upperHalfB[_letter][pvpIndexB],
				})
				pvpA = append(pvpA, item)
				pvpIndexB--
			}
			incr++
		}
		return pvpA
	}
	//
	data := NewMP()
	//左区对战
	data["left"] = _func(upperHalf, 'A')
	//右区对战
	data["right"] = _func(lowerHalf, _az)
	return data
}

// 安排字母
func (that *AlgorithmSort) az(a string, incr ...int) string {
	_incr := 1
	if len(incr) > 0 {
		_incr = incr[0]
		Slice := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O",
			"P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
		return Slice[_incr+that.ContainsString(Slice, a)]
	}
	x := []rune(a)
	for index := range x {
		x[index] = x[index] + 1
	}
	return string(x)
}
