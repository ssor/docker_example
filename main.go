package main

import (
    "io/ioutil"
    "fmt"
    "gopkg.in/yaml.v2"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
    "time"
    "net/http"
    "flag"
    "sync"
    "strconv"
)

var (
    port = flag.Int("port", 8001, "listening port")

    status = sync.Map{}

    pairs = []HostPortPair{}
)

func main() {
    flag.Parse()
    if flag.Parsed() == false {
        flag.PrintDefaults()
        return
    }
    if *port == 0 {
        logrus.Warn("must set a port")
        return
    }
    inputPairs, err := readConfig("/config/config.yaml")
    if err != nil {
        logrus.Error("config file error: ", err)
    } else {
        for _, pair := range inputPairs {
            addNewPair(pair)
        }
    }

    go talkToBrothers()
    r := gin.Default()
    r.GET("/ping", func(context *gin.Context) {
        context.Data(http.StatusOK, "plain/text", []byte("pong"))
    })
    r.GET("status", brothersStatus)
    r.GET("add", addBrothers)
    r.Run(fmt.Sprintf(":%d", *port))
}

func addBrothers(context *gin.Context) {
    hostRaw, ok := context.GetQuery("ip")
    if ok == false {
        logrus.Info("ip empty")
        return
    }
    //ip := net.ParseIP(hostRaw)
    //if ip == nil {
    //    logrus.Info("ip invalid: ", hostRaw)
    //    return
    //}

    portRaw, ok := context.GetQuery("port")
    if ok == false {
        logrus.Info("port  empty")
        return
    }
    port, err := strconv.Atoi(portRaw)
    if err != nil {
        logrus.Info("port invalid: ", portRaw)
        return
    }
    pair := HostPortPair{
        Host: hostRaw,
        Port: port,
    }
    pairs = append(pairs, pair)
    if addNewPair(pair) {
        logrus.Info("add brother: ", pair)
    } else {
        logrus.Infof("brother %s already exits", pair)
    }
}
func brothersStatus(context *gin.Context) {
    obj := make(map[string]int)
    status.Range(func(key, value interface{}) bool {
        logrus.Infof("%s -> %d", key, value)
        obj[key.(string)] = value.(int)
        return true
    })
    context.JSON(http.StatusOK, obj)
}

func talkToBrothers() {
    ticker := time.NewTicker(5 * time.Second)
    for {
        <-ticker.C
        for _, pair := range pairs {
            callBrother(pair)
        }
    }
}

func callBrother(pair HostPortPair) {
    key := fmt.Sprintf("%s:%d", pair.Host, pair.Port)
    url := fmt.Sprintf("http://%s/ping", key)
    res, err := http.Get(url)
    if err != nil {
        status.Store(key, 1)
        logrus.Error("request URL ", url, " error: ", err)
        return
    }
    if res.StatusCode != http.StatusOK {
        status.Store(key, 2)
        logrus.Error("request URL ", url, " status : ", res.StatusCode)
        return
    }
    logrus.Info("request url ", url, " success ", time.Now().Format(time.RFC3339))
    status.Store(key, 0)
}

type Config struct {
    Pairs []HostPortPair
}

type HostPortPair struct {
    Host string `yaml:"host"`
    Port int    `yaml:"port"`
}

func (pair HostPortPair) String() string {
    return fmt.Sprintf("%s:%d", pair.Host, pair.Port)
}

func addNewPair(newPair HostPortPair) bool {
    alreadyExists := false
    for _, pair := range pairs {
        if pair.Port == newPair.Port && pair.Host == newPair.Host {
            alreadyExists = true
            break
        }
    }
    if alreadyExists == false {
        pairs = append(pairs, newPair)
        return true
    } else {
        return false
    }
}

func readConfig(filePath string) ([]HostPortPair, error) {
    bs, err := ioutil.ReadFile(filePath)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    //var config Config
    var pairs []HostPortPair
    err = yaml.Unmarshal(bs, &pairs)
    if err != nil {
        fmt.Println(err)
        return nil, err
    }
    return pairs, nil
}
