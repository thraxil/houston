package main

import (
	"flag"
	"fmt"
	"github.com/stvp/go-toml-config"
	"html/template"
	"log"
	"net/http"
	"strings"
)

type Graph struct {
	Base    string
	Attrs   map[string]string
	Targets []string
}

func NewGraph(base string) *Graph {
	var g = &Graph{Base: base, Attrs: map[string]string{}, Targets: []string{}}
	return g
}

func (g *Graph) Param(k string, v string) *Graph {
	g.Attrs[k] = v
	return g
}

func (g *Graph) Target(t string) *Graph {
	g.Targets = append(g.Targets, t)
	return g
}

func (g Graph) Render() template.URL {
	pairs := make([]string, 0)
	for _, t := range g.Targets {
		pairs = append(pairs,
			"target="+template.URLQueryEscaper(t))
	}
	for k, v := range g.Attrs {
		pairs = append(pairs,
			k+"="+template.URLQueryEscaper(v))
	}
	return template.URL(g.Base + "render/?" + strings.Join(pairs, "&"))
}

type ApplicationInfo struct {
	Name         string
	GraphiteBase string
	Width        int
	Height       int
}

type ServerInfo struct {
	Name         string
	GraphiteBase string
	Width        int
	Height       int
}

func (s ServerInfo) GraphBase() *Graph {
	return NewGraph(s.GraphiteBase).Param(
		"_salt", "1369503684.499").Param(
		"width", fmt.Sprintf("%d", s.Width)).Param(
		"height", fmt.Sprintf("%d", s.Height)).Param(
		"bgcolor", "FFFFFF").Param(
		"colorList", "#999999,#006699").Param(
		"hideGrid", "true").Param(
		"hideLegend", "true").Param(
		"hideAxes", "true").Param(
		"graphOnly", "true")
}

func (s ServerInfo) CPUGraphBase() *Graph {
	return s.GraphBase().Param(
		"yMin", "0").Param(
		"areaMode", "first").Target(
		"server." + s.Name + ".cpu.load_average.1_minute").Target(
		"constantLine(1)")
}

func (s ServerInfo) CPUGraphUrl() template.URL {
	return s.CPUGraphBase().Render()
}

func (s ServerInfo) CPUGraphUrlWeekly() template.URL {
	return s.CPUGraphBase().Weekly().Render()
}

func (s ServerInfo) MemGraphBase() *Graph {
	return s.GraphBase().Param(
		"yMin", "0").Param(
		"yMax", "100").Param(
		"lineMode", "staircase").Param(
		"areaMode", "first").Target(
		"server." + s.Name + ".memory.MemFree.percent").Target(
		"constantLine(10)")
}

func (s ServerInfo) MemGraphUrl() template.URL {
	return s.MemGraphBase().Render()
}

func (s ServerInfo) MemGraphUrlWeekly() template.URL {
	return s.MemGraphBase().Weekly().Render()
}

func (s ServerInfo) DiskUsageGraphBase() *Graph {
	return s.GraphBase().Param(
		"yMin", "0").Param(
		"yMax", "100").Param(
		"lineMode", "staircase").Param(
		"areaMode", "first").Target(
		"asPercent(server." + s.Name + ".disk.usage.available," +
			"server." + s.Name + ".disk.usage.total)").Target(
		"constantLine(10)")
}

func (s ServerInfo) DiskUsageGraphUrl() template.URL {
	return s.DiskUsageGraphBase().Render()
}

func (s ServerInfo) DiskUsageGraphUrlWeekly() template.URL {
	return s.DiskUsageGraphBase().Weekly().Render()
}

func (s ServerInfo) NetworkGraphBase() *Graph {
	return s.GraphBase().Target(
		"nonNegativeDerivative(server." +
			s.Name + ".network.eth0.receive.byte_count)").Target(
		"nonNegativeDerivative(server." +
			s.Name + ".network.eth0.transmit.byte_count)")
}

func (s ServerInfo) NetworkGraphUrl() template.URL {
	return s.NetworkGraphBase().Render()
}

func (s ServerInfo) NetworkGraphUrlWeekly() template.URL {
	return s.NetworkGraphBase().Weekly().Render()
}

func (g *Graph) Weekly() *Graph {
	return g.Param(
		"from", "-7days").Param(
		"height", "50").Param(
		"colorList", "#cccccc,#6699cc").Param(
		"bgcolor", "#eeeeee")
}

func (a ApplicationInfo) GraphBase() *Graph {
	return NewGraph(a.GraphiteBase).Param(
		"_salt", "1369503684.499").Param(
		"width", fmt.Sprintf("%d", a.Width)).Param(
		"height", fmt.Sprintf("%d", a.Height)).Param(
		"bgcolor", "FFFFFF").Param(
		"colorList", "#999999,#006699").Param(
		"hideGrid", "true").Param(
		"hideLegend", "true").Param(
		"hideAxes", "true").Param(
		"graphOnly", "true")
}

func (a ApplicationInfo) ReqsGraphBase() *Graph {
	return a.GraphBase().Param(
		"yMin", "0").Param(
		"lineMode", "connected").Target(
		"hitcount(stats_counts." + a.Name + ".response.200,\"10s\")")
}

func (a ApplicationInfo) ReqsGraphUrl() template.URL {
	return a.ReqsGraphBase().Render()
}

func (a ApplicationInfo) ReqsGraphUrlWeekly() template.URL {
	return a.ReqsGraphBase().Weekly().Render()
}

func (a ApplicationInfo) FiveHundredsGraphBase() *Graph {
	return a.GraphBase().Param(
		"yMin", "0").Param(
		"drawNullAsZero", "True").Param(
		"lineMode", "connected").Target(
		"stats_counts." + a.Name + ".response.500")
}

func (a ApplicationInfo) FiveHundredsGraphUrl() template.URL {
	return a.FiveHundredsGraphBase().Render()
}

func (a ApplicationInfo) FiveHundredsGraphUrlWeekly() template.URL {
	return a.FiveHundredsGraphBase().Weekly().Render()
}

func (a ApplicationInfo) TimesGraphBase() *Graph {
	return a.GraphBase().Param(
		"yMin", "0").Param(
		"lineMode", "connected").Target(
		"keepLastValue(stats.timers." + a.Name + ".view.GET.mean)").Target(
		"keepLastValue(stats.timers." + a.Name + ".view.GET.upper_90)").Target(
		"constantLine(100)")
}

func (a ApplicationInfo) TimesGraphUrl() template.URL {
	return a.TimesGraphBase().Render()
}

func (a ApplicationInfo) TimesGraphUrlWeekly() template.URL {
	return a.TimesGraphBase().Weekly().Render()
}

func RiakGraphBase(graphite_base string, width int, height int) *Graph {
	return NewGraph(graphite_base).Param(
		"_salt", "1369503684.499").Param(
		"width", fmt.Sprintf("%d", width)).Param(
		"height", fmt.Sprintf("%d", height)).Param(
		"bgcolor", "FFFFFF").Param(
		"colorList", "#999999,#006699").Param(
		"hideGrid", "true").Param(
		"hideLegend", "true").Param(
		"hideAxes", "true").Param(
		"graphOnly", "true").Target(
		"riak_stats.*.node_gets.count").Target(
		"riak_stats.*.node_puts.count").Param(
		"lineMode", "connected").Param(
		"drawNullAsZero", "True")
}

func (g *Graph) RiakGraphUrl() template.URL {
	return g.Render()
}

func NginxGraphBase(graphite_base string, width int, height int) *Graph {
	return NewGraph(graphite_base).Param(
		"_salt", "1369503684.499").Param(
		"width", fmt.Sprintf("%d", width)).Param(
		"height", fmt.Sprintf("%d", height)).Param(
		"bgcolor", "FFFFFF").Param(
		"colorList", "#999999,#006699").Param(
		"hideGrid", "true").Param(
		"hideLegend", "true").Param(
		"hideAxes", "true").Param(
		"graphOnly", "true").Target(
		"nonNegativeDerivative(keepLastValue(nginx.*.requests))").Param(
		"lineMode", "connected").Param(
		"drawNullAsZero", "True")
}

func (g *Graph) NginxGraphUrl() template.URL {
	return g.Render()
}

func (g *Graph) NginxGraphUrlWeekly() template.URL {
	return g.Weekly().Render()
}

type PageResponse struct {
	Servers       []ServerInfo
	Applications  []ApplicationInfo
	RiakGraph     *Graph
	NginxGraph    *Graph
}

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "./dev.conf", "TOML config file")
	flag.Parse()
	var (
		port          = config.String("port", "8888")
		media_dir     = config.String("media_dir", "static")
		graphite_base = config.String("graphite_base", "")
		servers       = config.String("servers", "")
		apps          = config.String("apps", "")
		template_file = config.String("template", "index.html")
	)
	config.Parse(configFile)

	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			var serverinfo = []ServerInfo{}
			var appinfo = []ApplicationInfo{}
			for _, s := range strings.Split(*servers, ",") {
				serverinfo = append(serverinfo,
					ServerInfo{
						Name:         s,
						GraphiteBase: *graphite_base,
						Width:        300,
						Height:       100,
					})
			}
			for _, a := range strings.Split(*apps, ",") {
				appinfo = append(appinfo,
					ApplicationInfo{
						Name:         a,
						GraphiteBase: *graphite_base,
						Width:        400,
						Height:       100,
					})
			}
			pr := PageResponse{
				Servers:       serverinfo,
				Applications:  appinfo,
				RiakGraph:     RiakGraphBase(*graphite_base, 1200, 100),
				NginxGraph:    NginxGraphBase(*graphite_base, 1200, 100),
			}
			t, err := template.ParseFiles(*template_file)
			if err != nil {
				fmt.Println(fmt.Sprintf("%v", err))
			}
			t.Execute(w, pr)
		})
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir(*media_dir))))
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
