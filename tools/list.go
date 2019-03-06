//这里是一个将Slice切片封装成List结构体的轮子
//为什么没有使用官方包的"container/list"？主要是测试后发现性能的差异，所以应用场景的使用上也要有所差异
//声明：测试结果为个人电脑的测试结果，仅供参考。
//===========================================================
//t := time.Now()
//sli := make([]int, 10)
//for i := 0; i < 1*100000*1000; i++ {
//	sli = append(sli, 1)
//}
//fmt.Println("slice 创建速度：" + time.Now().Sub(t).String())

//t = time.Now()
//l := list.New()
//for i := 0; i < 1*100000*1000; i++ {
//	l.PushBack(1)
//}
//fmt.Println("list 创建速度: " + time.Now().Sub(t).String())

//slice 创建速度：3.1410928s
//list 创建速度: 41.4315308s
//对于1亿条数据，slice 的创建和添加元素的速度约是list的13倍。
//============================================================
//sli := make([]int, 10)
//for i := 0; i < 1*100000*1000; i++ {
//	sli = append(sli, 1)
//}

//l := list.New()
//for i := 0; i < 1*100000*1000; i++ {
//	l.PushBack(1)
//}
//// 比较遍历
//t := time.Now()
//for _, _ = range sli {
//	//fmt.Printf("values[%d]=%d\n", i, item)
//}
//fmt.Println("遍历slice的速度:" + time.Now().Sub(t).String())
//t = time.Now()
//for e := l.Front(); e != nil; e = e.Next() {
//	//fmt.Println(e.Value)
//}
//fmt.Println("遍历list的速度:" + time.Now().Sub(t).String())

//遍历slice的速度:65.1759ms
//遍历list的速度:28.7595276s
//这差距。。
//============================================================
//    sli:=make([]int ,10)
//    for i := 0; i<1*100000*1000;i++  {
//    sli=append(sli, 1)
//    }

//    l:=list.New()
//    for i := 0; i<1*100000*1000;i++  {
//        l.PushBack(1)
//    }
//    //比较插入
//    t := time.Now()
//    slif:=sli[:100000*500]
//    slib:=sli[100000*500:]
//    slif=append(slif, 10)
//    slif=append(slif, slib...)
//    fmt.Println("slice的插入速度" + time.Now().Sub(t).String())

//    var em *list.Element
//    len:=l.Len()
//    var i int
//    for e := l.Front(); e != nil; e = e.Next() {
//        i++
//        if i ==len/2 {
//            em=e
//            break
//        }
//    }
//    //忽略掉找中间元素的速度。
//    t = time.Now()
//    ef:=l.PushBack(2)
//    l.MoveBefore(ef,em)
//    fmt.Println("list的插入速度: " + time.Now().Sub(t).String())

//slice的插入速度:1.9442905s
//list的插入速度:2.0326ms
//list的优势在快速的插入数据
//====================================================================

