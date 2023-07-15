package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"separa/common/log"
	"separa/core/report"
	"sync"

	"github.com/gin-gonic/gin"
)

type Nod struct {
	url    string
	status int
}

type NodeMap struct {
	*sync.Mutex
	Nodes map[string]Nod
}

var Nodes = &NodeMap{}

func NodeList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func NodeAdd(c *gin.Context) {
	url, _ := c.GetQuery("url")
	resp, err := http.Get("http://" + url + "/ping")
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	status := data["message"].(int)

	Nodes.Nodes[url] = Nod{url, status}

	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func NodeUpdate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func NodeDelete(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func TaskDistribute(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	num := len(Nodes.Nodes)
	var groups [][]string
	size := (len(task.Targets) + num - 1) / num
	for i := 0; i < len(task.Targets); i += size {
		end := i + size
		if end > len(task.Targets) {
			end = len(task.Targets)
		}
		groups = append(groups, task.Targets[i:end])
	}

	cnt := 0
	for _, v := range Nodes.Nodes {
		task := Task{
			Targets: groups[cnt],
			Port:    task.Port,
			Delay:   task.Delay,
			Output:  task.Output,
		}
		data, _ := json.Marshal(task)
		resp, err := http.Post("http://"+v.url+"/task/run", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Err("%s post error: %s", v.url, err.Error())
		}
		cnt++
		defer resp.Body.Close()
	}
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func TaskStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

func TaskResult(c *gin.Context) {
	out, _ := c.GetQuery("out")

	var result map[string]report.ResultUnit = make(map[string]report.ResultUnit)

	for _, v := range Nodes.Nodes {
		resp, err := http.Get("http://" + v.url + "/task/collect?out=" + out)
		if err != nil {
			log.Err("%s get error: %s", v.url, err.Error())
			continue
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Err("%s body error: %s", v.url, err.Error())
			continue
		}
		var data map[string]report.ResultUnit
		if err := json.Unmarshal(body, &data); err != nil {
			log.Err("%s unmarshal error: %s", v.url, err.Error())
			continue
		}

		for kk, vv := range data {
			result[kk] = vv
		}
	}

	c.JSON(200, result)
}
