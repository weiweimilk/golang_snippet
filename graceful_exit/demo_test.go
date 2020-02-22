package graceful_exit

import (
    "bytes"
    "fmt"
    "os"
    "os/signal"
    "regexp"
    "syscall"
    "testing"
    "time"
)

func TestManager(t *testing.T) {
    queueSize := 10
    manager, err := NewManager(queueSize)
    if err != nil {
        t.Errorf("NewManager return error[%v]", err)
        return
    }

    for i := 0; i < queueSize; i++ {
        err := manager.Insert(fmt.Sprintf("data%d", i))
        if err != nil {
            t.Errorf("manager.Insert:%d result should be nil, but get [%v]", i, err)
            return
        }
    }
    err = manager.Insert("data")

    if err != FullError {
        t.Errorf("manager.Insert result should be error, but get[%v]", err)
        return
    }

    for i := 0; i < queueSize - 1; i++ {
        manager.processUtil(0)
    }
    t.Logf("start to stop")
    manager.Stop()
}

func TestExitBySignal(t *testing.T) {
    signals := make(chan os.Signal, 1)
    done := make(chan bool, 1)

    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        sig := <-signals
        fmt.Println()
        fmt.Println(sig)
        done <- true
    }()

    // 模拟ctrl+c 行为
    go func() {
       time.Sleep(time.Second)
        signals <- syscall.SIGINT
    }()

    fmt.Println("awaiting signal")
    <-done
    fmt.Println("exiting")
}


func worker(id int, jobs <-chan int, results chan<- int) {
    for j := range jobs {
        fmt.Println("worker", id, "started  job", j)
        time.Sleep(time.Microsecond)
        fmt.Println("worker", id, "finished job", j)
        results <- j * 2
    }
}
func TestWorkerPool(t *testing.T) {
    const numJobs = 5
    jobs := make(chan int, numJobs)
    results := make(chan int, numJobs)

    for w := 1; w <= 3; w++ {
        go worker(w, jobs, results)
    }

    for j := 1; j <= numJobs; j++ {
        jobs <- j
    }
    close(jobs)

    for a := 1; a <= numJobs; a++ {
        <-results
    }
}

func TestClosingChannel(t *testing.T) {
    jobs := make(chan int, 5)
    done := make(chan bool)

    go func() {
        for {
            j, more := <-jobs
            if more {
                fmt.Println("received job", j)
            } else {
                fmt.Println("received all jobs")
                done <- true
                return
            }
        }
    }()

    for j := 1; j <= 3; j++ {
        jobs <- j
        fmt.Println("sent job", j)
    }
    close(jobs)
    fmt.Println("sent all jobs")

    <-done
}

type point struct {
    x, y int
}
func TestFormat(t *testing.T) {
    p := point{1, 2}
    fmt.Printf("%v\n", p)

    fmt.Printf("%+v\n", p)

    fmt.Printf("%#v\n", p)

    fmt.Printf("%T\n", p)

    fmt.Printf("%t\n", true)

    fmt.Printf("%d\n", 123)

    fmt.Printf("%b\n", 14)

    fmt.Printf("%c\n", 33)

    fmt.Printf("%x\n", 456)

    fmt.Printf("%f\n", 78.9)

    fmt.Printf("%e\n", 123400000.0)
    fmt.Printf("%E\n", 123400000.0)

    fmt.Printf("%s\n", "\"string\"")

    fmt.Printf("%q\n", "\"string\"")

    fmt.Printf("%x\n", "hex this")

    fmt.Printf("%p\n", &p)

    fmt.Printf("|%6d|%6d|\n", 12, 345)

    fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

    fmt.Printf("|%-6.2f|%-6.2f|\n", 1.2, 3.45)

    fmt.Printf("|%6s|%6s|\n", "foo", "b")

    fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

    s := fmt.Sprintf("a %s", "string")
    fmt.Println(s)

    fmt.Fprintf(os.Stderr, "an %s\n", "error")
}

func TestRegex(t *testing.T) {

    match, _ := regexp.MatchString("p([a-z]+)ch", "peach")
    fmt.Println(match)

    r, _ := regexp.Compile("p([a-z]+)ch")

    fmt.Println(r.MatchString("peach"))

    fmt.Println(r.FindString("peach punch"))

    fmt.Println(r.FindStringIndex("peach punch"))

    fmt.Println(r.FindStringSubmatch("peach punch"))

    fmt.Println(r.FindStringSubmatchIndex("peach punch"))

    fmt.Println(r.FindAllString("peach punch pinch", -1))

    fmt.Println(r.FindAllStringSubmatchIndex(
        "peach punch pinch", -1))

    fmt.Println(r.FindAllString("peach punch pinch", 2))

    fmt.Println(r.Match([]byte("peach")))

    r = regexp.MustCompile("p([a-z]+)ch")
    fmt.Println(r)

    fmt.Println(r.ReplaceAllString("a peach", "<fruit>"))

    in := []byte("a peach")
    out := r.ReplaceAllFunc(in, bytes.ToUpper)
    fmt.Println(string(out))
}

func TestFuntionTimeout(t *testing.T) {

}