package lfu

import cache "5_cache"

type entry struct {
	key    string      //queue 是一个 entry 指针切片；
	value  interface{} //	entry 和 FIFO 中的区别是多了两个字段：weight 和 index；
	weight int         //	weight 表示该 entry 在 queue 中权重（优先级），访问次数越多，权重越高；
	index  int         //index 代表该 entry 在堆（heap）中的索引；
}

func (e *entry) Len() int {
	//还要加上key的长度， weight int 和 index int：如果你是在 64 位机器，int 是 8 字节
	return len(e.key) + cache.CalcLen(e.value) + 8 + 8
}

type queue []*entry

func (q queue) Len() int {
	return len(q)
}

func (q queue) Less(i, j int) bool {
	return q[i].weight < q[j].weight
}
func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].index = i
	q[j].index = j
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	e := x.(*entry)
	e.index = n
	*q = append(*q, e)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	e := old[n-1]
	old[n-1] = nil // avoid memory leak
	e.index = -1   // for safety
	*q = old[:n-1]
	return e
}
