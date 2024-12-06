package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/segmentio/encoding/json"

	"github.com/valyala/fasthttp"

	"github.com/maid-zone/soundcloak-instances/data"
	"github.com/maid-zone/soundcloak-instances/templates"
)

// The script will write data to file + ".json" and file + ".html"
//const file = "/home/maid.zone/website/assets/soundcloak/instances"

const file = "test"

type JSONSettings struct {
	ProxyImages  bool
	ProxyStreams bool
	Restream     bool
}

type JSONStatus struct {
	Error          string
	SkippedResolve bool
}

type JSONInstance struct {
	URL        string
	Country    string
	Maintainer string
	Note       string

	Settings JSONSettings
	Status   JSONStatus
}

func DoWithRetry(req *fasthttp.Request, resp *fasthttp.Response) (err error) {
	for i := 1; i <= 5; i++ {
		err = fasthttp.Do(req, resp)
		if err == nil {
			return nil
		}

		// god damn
		time.Sleep(time.Duration(i) * time.Second)
	}

	return
}

func resolve(instance string) (i data.InstanceInfo, err error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(instance + "/_/info")
	req.Header.SetUserAgent("github.com/maid-zone/soundcloak-instances")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err = DoWithRetry(req, resp)
	if err != nil {
		return
	}

	body := resp.Body()

	if string(body) == "resolve: got status code 404" || resp.StatusCode() == 404 {
		err = data.ErrNotFound
		return
	}

	err = json.Unmarshal(body, &i)
	return
}

func main() {
	var instances = make(map[data.Instance]data.InstanceInfo, len(data.Instances))

	for _, instance := range data.Instances {
		var i data.InstanceInfo
		if instance.SkipResolve {
			i.SkippedResolve = true
			instances[instance] = i
			continue
		}

		i, err := resolve(instance.URL.URL)
		if err != nil {
			i.Error = err
		}

		instances[instance] = i
	}

	fd, err := os.OpenFile(file+".html", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0766)
	if err != nil {
		log.Fatalln("failed to open", file+".html:", err)
	}

	err = templates.Instances(instances).Render(context.Background(), fd)
	if err != nil {
		log.Fatalln("failed to render webpage:", err)
	}

	fd.Close()

	var formattedInstances = []JSONInstance{}
	for _, instance := range data.Instances {
		info := instances[instance]
		i := JSONInstance{
			URL:        instance.URL.URL,
			Country:    instance.Country,
			Maintainer: instance.Maintainer.URL,
			Note:       instance.Note,

			Settings: JSONSettings{
				info.ProxyImages,
				info.ProxyStreams,
				info.Restream,
			},

			Status: JSONStatus{
				SkippedResolve: info.SkippedResolve,
			},
		}

		if info.Error != nil {
			i.Status.Error = info.Error.Error()
		}

		formattedInstances = append(formattedInstances, i)
	}

	data, err := json.Marshal(formattedInstances)
	if err != nil {
		log.Fatalln("failed to marshal instance list:", err)
	}

	os.WriteFile(file+".json", data, 0766)
}
