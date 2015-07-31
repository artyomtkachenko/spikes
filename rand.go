package main
import (
  "fmt"
  "math/rand"
  "time"
  "sort"
  "encoding/json"
  "strconv"
  "runtime"
)


func getMax(max int, res map[int]int) int {
  for k, v := range res {
    fmt.Println(k,v)
    if v == max {
      return k
    }
  }
  return 0
}

func init() {
  rand.Seed(time.Now().UnixNano())
}

func getRandNumber(n int) int {
  res := make(map[int]int) //will contain 36 numbers
  for i := 0; i < 65535; i++ {
    res[rand.Intn(n)] += 1
  }
  max := 0
  maxKey := -1
  for k, v := range res {
    if v > max {
      max = v
      maxKey = k
    }
  }
  return maxKey
}

func genTicket(n int, m int) map[string]string {
  res := make(map[int]int) //will contain 6 numbers
  for i := 1; i <= n ; i++ {
    again:
    max := getRandNumber(m)
    val, ok := res[max]
    if ok {
      goto again // we do have it already
    } else if val == 0 && ok {
      goto again
    } else if max == 0 {
      goto again
    } else {
      res[max] = i
    }
  }
  return convertToString(res)
}

func convertToString(m map[int]int) map[string]string  {
  res := make(map[string]string)
  for k, v := range m {
    res[strconv.Itoa(k)] = strconv.Itoa(v)
  }
  return res
}

func sortMapByValues(m map[int]int) {
  arr := make([]int, 0)
  for _, v := range m {
    arr = append(arr, v)
  }
  sort.Ints(arr)
  for _, v := range arr {
    fmt.Println(v)
  }
}

func produce(n int,queue chan map[string]string) { // writes into the channel
  for i := 0; i <= n; i++ {
    queue <- genTicket(7, 45)
  }
}

func reduce(tasks int, queue chan map[string]string) { //reads from the channel
  arr := make([]map[string]string, 0)
  // for r := range queue {
  for i := 0; i <= tasks; i++ {
    res := <- queue
    arr = append(arr, res)
  }

  out, err := json.Marshal(arr)
  if err != nil {
    panic(err)
  }
  fmt.Println(string(out))
}

func main() {
  cores := runtime.NumCPU()

  fmt.Printf("This machine has %d CPU cores. \n", cores)
  runtime.GOMAXPROCS(cores)

  ticketsToGenerate := 2000
  tasksPerCore := ticketsToGenerate / cores
  queue := make(chan map[string]string, ticketsToGenerate)
  for i := 0; i != cores; i++ {
    go produce(tasksPerCore, queue)
  }
  reduce(ticketsToGenerate, queue)
}
