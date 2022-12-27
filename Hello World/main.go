package main

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"context"
	"log"
	"net"

	"github.com/nats-io/nats.go"
)

type Player struct {
	On    bool
	Ammo  int
	Power int
}

func (p *Player) Shoot() (r bool) {

	r = p.On && p.Ammo > 0

	if r {
		p.Ammo--
	}

	return
}

func (p *Player) RideBike() (r bool) {

	r = p.On && p.Power > 0

	if r {
		p.Power--
	}

	return
}

func mul2X2Matrix(a [2][2]big.Int, b [2][2]big.Int) [2][2]big.Int {
	var c [2][2]big.Int
	var x1, x2 big.Int

	x1.Mul(&a[0][0], &b[0][0])
	x2.Mul(&a[0][1], &b[1][0])
	c[0][0].Add(&x1, &x2)

	x1.Mul(&a[0][0], &b[0][1])
	x2.Mul(&a[0][1], &b[1][1])
	c[0][1].Add(&x1, &x2)

	x1.Mul(&a[1][0], &b[0][0])
	x2.Mul(&a[1][1], &b[1][0])
	c[1][0].Add(&x1, &x2)

	x1.Mul(&a[1][0], &b[0][1])
	x2.Mul(&a[1][1], &b[1][1])
	c[1][1].Add(&x1, &x2)

	return c
}

func pow2X2Matrix(a [2][2]big.Int, n int) [2][2]big.Int {

	var res [2][2]big.Int

	copy(res[:], a[:])

	for n > 0 {
		if n&1 != 0 {
			res = mul2X2Matrix(res, a)
		}

		a = mul2X2Matrix(a, a)
		n >>= 1
	}

	return res
}

func fibonacci(n int) big.Int {
	F := [2][2]big.Int{
		{*big.NewInt(1), *big.NewInt(1)},
		{*big.NewInt(1), *big.NewInt(0)},
	}

	F = pow2X2Matrix(F, n)

	return F[1][1]
}

func mul2X2Matrix2(a [2][2]int, b [2][2]int) [2][2]int {
	var c [2][2]int
	c[0][0] = a[0][0]*b[0][0] + a[0][1]*b[1][0]
	c[0][1] = a[0][0]*b[0][1] + a[0][1]*b[1][1]
	c[1][0] = a[1][0]*b[0][0] + a[1][1]*b[1][0]
	c[1][1] = a[1][0]*b[0][1] + a[1][1]*b[1][1]

	return c
}

func pow2X2Matrix2(a [2][2]int, n int) [2][2]int {
	var res [2][2]int

	copy(res[:], a[:])

	for n > 0 {
		if n&1 != 0 {
			res = mul2X2Matrix2(res, a)
		}

		a = mul2X2Matrix2(a, a)
		n >>= 1
	}
	return res
}

func test(x1 *int, x2 *int) {

	*x1, *x2 = *x2, *x1
	fmt.Println(*x1, *x2)
}

func fibonacci2(n int) int {
	F := [2][2]int{
		{1, 1},
		{1, 0},
	}

	return pow2X2Matrix2(F, n)[1][1]
}

func sumInt(a ...int) (int, int) {
	s := 0

	for _, e := range a {
		s += e
	}

	return len(a), s
}

func testStruct(p *Player) {
	p.Ammo = 10
	p.On = true
	fmt.Print(p)
}

func isPolindrom(r []rune) bool {
	for i := 0; i < len(r)/2; i++ {
		if r[i] != r[len(r)-1-i] {
			return false
		}
	}
	return true
}

func checkPassword(r []rune) bool {

	if len(r) < 5 {
		return false
	}

	for _, ch := range r {
		if !unicode.IsDigit(ch) && !unicode.Is(unicode.Latin, ch) {
			return false
		}
	}

	return true
}

var k, p, v float64 = 1296, 6, 6

func T() float64 {
	return 6 / W()
}

func W() float64 {
	return math.Sqrt(k / M())
}

func M() float64 {
	return p * v
}

// func contains(s []string, a string) bool {
// 	for _, e := range s {
// 		if e == a {
// 			return true
// 		}
// 	}
// 	return false
// }

func findFirstInt(s string) int64 {

	runes := []rune(s)
	num := []rune("")

	i := 0

	for ; i < len(runes) && !unicode.IsDigit(runes[i]); i++ {
	}

	for ; i < len(runes) && unicode.IsDigit(runes[i]); i++ {
		num = append(num, runes[i])
	}

	v, err := strconv.ParseInt(string(num), 10, 64)

	if err != nil {
		panic("error parse")
	}

	return v
}

func adding(s1 string, s2 string) int64 {
	return findFirstInt(s1) + findFirstInt(s2)
}

func parseFloats(a ...interface{}) (parsed []float64, err error) {
	for _, i := range a {
		if v, ok := i.(float64); ok {
			parsed = append(parsed, v)
		} else {
			err = errors.New("value=" + fmt.Sprintf("%v", i) + ": " + reflect.TypeOf(i).String())
			return
		}
	}
	return
}

func calc(a, b float64, c string) (float64, error) {
	ops := map[string]func(float64, float64) float64{
		"+": func(f1, f2 float64) float64 { return f1 + f2 },
		"-": func(f1, f2 float64) float64 { return f1 - f2 },
		"*": func(f1, f2 float64) float64 { return f1 * f2 },
		"/": func(f1, f2 float64) float64 { return f1 / f2 },
	}

	if f, ok := ops[c]; ok {
		return f(a, b), nil
	}

	return -1, errors.New("неизвестная операция")
}

type battery uint

func (b battery) String() string {
	return fmt.Sprintf("[%10s]", strings.Repeat("X", int(b)))
}

func listFiles(file *zip.File) error {
	// fileread, err := file.Open()
	// if err != nil {
	// 	msg := "Failed to open zip %s for reading: %s"
	// 	return fmt.Errorf(msg, file.Name, err)
	// }
	// defer fileread.Close()

	if !file.FileInfo().IsDir() {

		f, err := file.Open()

		if err != nil {
			return nil
		}

		r := csv.NewReader(f)

		record, _ := r.ReadAll()
		if len(record) > 1 {
			fmt.Println(record[4][2])
		}
		// ext := file.Name[len(file.Name)-4:]

		// if file.FileInfo().Size() > 300 {
		// 	fmt.Println(file.Name)
		// }
		// if ext == ".csv" {
		// 	fmt.Fprintf(os.Stdout, "%s:", file.Name)
		// }
	}

	// if err != nil {
	// 	msg := "Failed to read zip %s for reading: %s"
	// 	return fmt.Errorf(msg, file.Name, err)
	// }
	return nil
}

// func walk(path *zip.File) {
// 	for _, file := range path.File {
// 		if err := listFiles(file); err != nil {
// 			log.Fatalf("Failed to read %s from zip: %s", file.Name, err)
// 		}
// 	}
// }

func walk(path string) {
	read, err := zip.OpenReader(path)
	if err != nil {
		msg := "Failed to open: %s"
		log.Fatalf(msg, err)
	}
	defer read.Close()

	for _, file := range read.File {
		if err := listFiles(file); err != nil {
			log.Fatalf("Failed to read %s from zip: %s", file.Name, err)
		}
	}
}

type (
	Group struct {
		Students []Student
	}

	Student struct {
		Rating []int
	}

	Average struct {
		Average float64
	}
)

type Object struct {
	GlobalId uint64 `json:"global_id"`
}

func add(a *int) {
	*a++
}

const now = 1589570165

func task2(chanel chan string, s string) {
	for i := 0; i < 5; i++ {
		chanel <- s + " "
	}
}

func removeDuplicates(inputStream chan string, outputStream chan string) {
	var prev string

	for s := range inputStream {
		if prev != s {
			outputStream <- s
			prev = s
		}
	}
	close(outputStream)
}

func task(chanel chan int, n int) {
	chanel <- n + 1
}

func merge2Channels(fn func(int) int, in1 <-chan int, in2 <-chan int, out chan<- int, n int) {
	a := make([]int, n)
	mu := new(sync.Mutex)
	wg := new(sync.WaitGroup)

	calc := func(x int, i int) {
		defer wg.Done()

		y := fn(x)

		mu.Lock()
		a[i] += y
		mu.Unlock()
	}

	go func() {
		defer close(out)

		wg.Add(n * 2)
		for i := 0; i < n; i++ {
			go calc(<-in1, i)
			go calc(<-in2, i)
		}
		wg.Wait()

		fmt.Println(a)
		for i := 0; i < n; i++ {
			out <- a[i]
		}
	}()
}

const N = 3

// func main() {
// 	// walk(`C:\Users\sergi\Downloads\task.zip`)

// 	// reader := bufio.NewReader(os.Stdin)
// 	// s, _ := reader.ReadString('\n')

// 	// p := "2006-01-02 15:04:05"

// 	// fmt.Println(strings.TrimRight(s, "\n"))

// 	// t, err := time.Parse(p, strings.TrimRight(s, "\n"))

// 	// if err != nil {
// 	// 	fmt.Println("ошибка парсинга даты:", err)
// 	// }

// 	// if t.Hour() > 13 {
// 	// 	t = t.AddDate(0, 0, 1)
// 	// }

// 	// fmt.Println(t.Format(p))

// 	// reader := bufio.NewReader(os.Stdin)

// 	// str, _ := reader.ReadString('\n')

// 	// dates := strings.Split(strings.TrimSpace(str), ",")

// 	fn := func(x int) int {
// 		time.Sleep(time.Duration(rand.Int31n(N)) * time.Second)
// 		return x * 2
// 	}
// 	in1 := make(chan int, N)
// 	in2 := make(chan int, N)
// 	out := make(chan int, N)

// 	start := time.Now()
// 	merge2Channels(fn, in1, in2, out, N+1)
// 	for i := 0; i < N+1; i++ {
// 		in1 <- i
// 		in2 <- i
// 	}

// 	orderFail := false
// 	EvenFail := false
// 	for i, prev := 0, 0; i < N; i++ {
// 		c := <-out
// 		if c%2 != 0 {
// 			EvenFail = true
// 		}
// 		if prev >= c && i != 0 {
// 			orderFail = true
// 		}
// 		prev = c
// 		fmt.Println(c)
// 	}
// 	if orderFail {
// 		fmt.Println("порядок нарушен")
// 	}
// 	if EvenFail {
// 		fmt.Println("Есть не четные")
// 	}
// 	duration := time.Since(start)
// 	if duration.Seconds() > N {
// 		fmt.Println("Время превышено")
// 	}
// 	fmt.Println("Время выполнения: ", duration)
// }

func Sqrt(x float64) float64 {
	z := 100.0
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Println(z)
	}
	return z
}

func WordCount(s string) map[string]int {

	counter := map[string]int{}

	for _, word := range strings.Split(s, " ") {
		counter[word]++
	}

	return counter
}

func fibonaccif() func() int {
	a, b := 1, 0

	return func() int {
		a, b = b, a+b
		return a
	}
}

type IPAddr [4]byte

func (ipAddr IPAddr) String() string {
	return fmt.Sprintf("%v.%v.%v.%v", ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

func main() {
	// hosts := map[string]IPAddr{
	// 	"loopback":  {127, 0, 0, 1},
	// 	"googleDNS": {8, 8, 8, 8},
	// }
	// for name, ip := range hosts {
	// 	fmt.Printf("%v: %v\n", name, ip)
	// }

	fmt.Println('A' - 'z')

	// cap()
	// nc, _ := nats.Connect(nats.DefaultURL)
	// defer nc.Close()

	// nc.Subscribe("foo.*.baz", func(m *nats.Msg) {
	// 	fmt.Printf("Msg received on [%s] : %s\n", m.Subject, string(m.Data))
	// })

	// nc.Subscribe("foo.bar.*", func(m *nats.Msg) {
	// 	fmt.Printf("Msg received on [%s] : %s\n", m.Subject, string(m.Data))
	// })

	// // ">" matches any length of the tail of a subject, and can only be the last token
	// // E.g. 'foo.>' will match 'foo.bar', 'foo.bar.baz', 'foo.foo.bar.bax.22'
	// nc.Subscribe("foo.>", func(m *nats.Msg) {
	// 	fmt.Printf("Msg received on [%s] : %s\n", m.Subject, string(m.Data))
	// })

	// // Matches all of the above
	// nc.Publish("foo.bar.baz", []byte("Hello World"))

	// time.Sleep(3000 * time.Microsecond)

	// recvCh := make(chan *person)
	// ec.BindRecvChan("hello", recvCh)

	// sendCh := make(chan *person)
	// ec.BindSendChan("hello", sendCh)

	// me := &person{Name: "derek", Age: 22, Address: "140 New Montgomery Street"}

	// // Send via Go channels
	// sendCh <- me

	// // Receive via Go channels
	// who := <-recvCh

	// fmt.Println(*who)
}

type customDialer struct {
	ctx             context.Context
	nc              *nats.Conn
	connectTimeout  time.Duration
	connectTimeWait time.Duration
}

func (cd *customDialer) Dial(network, address string) (net.Conn, error) {
	ctx, cancel := context.WithTimeout(cd.ctx, cd.connectTimeout)
	defer cancel()

	for {
		log.Println("Attempting to connect to", address)
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		select {
		case <-cd.ctx.Done():
			return nil, cd.ctx.Err()
		default:
			d := &net.Dialer{}
			if conn, err := d.DialContext(ctx, network, address); err == nil {
				log.Println("Connected to NATS successfully")
				return conn, nil
			} else {
				time.Sleep(cd.connectTimeWait)
			}
		}
	}
}

// Merge2Channels below

func calculator(arguments <-chan int, done <-chan struct{}) <-chan int {
	c := make(chan int)

	go func() {

		defer close(c)
		sum := 0

		for {
			select {
			case a := <-arguments:
				sum += a
			case <-done:
				c <- sum
				return
			}
		}
	}()

	return c
}

func worker(id int, limit <-chan time.Time, wg *sync.WaitGroup) {
	defer wg.Done()
	<-limit

	fmt.Printf("worker %d выполнил работу\n", id)

	// По идее значение счетчика должно быть 1000, но крайне вероятно, что этого не произойдет

	// var m, s int64
	// fmt.Scanf("%d мин. %d сек.", &m, &s)
	// fmt.Println(time.Unix(now+m*60+s, 0).UTC().Format(time.UnixDate))
}

func convert(a int64) uint16 {
	return uint16(a)
}

func divide(x, y float64) float64 {
	if y == 0 {
		panic("division by zero!")
	}
	return x / y
}
