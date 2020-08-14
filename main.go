package main

import (
	"github.com/ASalimov/jbuilder/cmd"
)

func main() {
	//s:="привет"s
	//fmt.Println(string([]rune(s)[:3]))
	//s:="جا"
	//s1:="в"
	//fmt.Println(len(s))
	//fmt.Println(len(s1))
	//cmd.InitView()
	//fmt.Println("eee")
	//os.Exit(1)
	cmd.Execute()
	//time.Sleep(time.Second)
	//n := 20
	//b := bar.NewWithOpts(
	//	bar.WithDimensions(20, 20),
	//	bar.WithFormat(
	//		fmt.Sprintf(
	//			" %sloading...%s :percent :bar %s:rate ops/s%s ",
	//			chalk.Blue,
	//			chalk.Reset,
	//			chalk.Green,
	//			chalk.Reset)))
	//
	//for i := 0; i < n; i++ {
	//	b.Interrupt("asdfasdf")
	//	b.Tick()
	//	time.Sleep(500 * time.Millisecond)
	//}
	//
	//b.Done()
}
