# Example 1
```
package main

import "time"
import "fmt"
import "github.com/k7575/rtlimit"

func main() {
  a := rtlimit.New(5, 1*time.Second, 60*time.Second)
  if a.Check("192.168.1.1") {
    fmt.Println("Ok")
  }
}
```
# Example 2
```
package main

import "time"
import "fmt"
import "github.com/k7575/rtlimit"

func main() {
  go func() {
    rtlimit.Run("127.0.0.1:8080", 5, 1*time.Second, 60*time.Second)
  }()
  time.Sleep(1 * time.Second)

  if rtlimit.Client("http://127.0.0.1:8080/?id=192.168.1.1") {
    fmt.Println("Ok")
  }
}
```
