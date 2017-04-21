package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/nlopes/slack"
)

var version = "development" // set by build script

func main() {
	var (
		flAPIToken = flag.String("api.token", "", "slack api token")
		flCSV      = flag.Bool("csv", false, "print csv output")
		flVersion  = flag.Bool("version", false, "print version and exit")
	)
	flag.Parse()

	if *flVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	if *flAPIToken == "" {
		fmt.Println("must provide an API token")
		flag.Usage()
		os.Exit(1)
	}

	var w outputter
	if *flCSV {
		w = &csvOut{w: csv.NewWriter(os.Stdout)}
	} else {
		w = &basicOut{w: tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)}
	}

	api := slack.New(*flAPIToken)
	channels, err := api.GetChannels(true)
	if err != nil {
		log.Fatal(err)
	}

	w.Header()
	defer w.Footer()

	for _, c := range channels {
		user, err := api.GetUserInfo(c.Creator)
		if err != nil {
			log.Println(err)
			continue
		}
		line := &Line{
			ChannelName: c.Name,
			UserName:    user.Name,
			CreatedDate: c.Created.Time().UTC().String(),
			NumMembers:  c.NumMembers,
			Purpose:     c.Purpose.Value,
		}
		w.WriteLine(line)
	}
}

type Line struct {
	ChannelName string
	UserName    string
	CreatedDate string
	NumMembers  int
	Purpose     string
}

type outputter interface {
	Header()
	WriteLine(*Line)
	Footer()
}

type csvOut struct {
	w *csv.Writer
}

func (out *csvOut) Header() {
	out.w.Write([]string{"Name", "Creator", "CreatedDate", "NumMembers", "Purpose"})
}
func (out *csvOut) Footer() {
	out.w.Flush()
}

func (out *csvOut) WriteLine(l *Line) {
	members := strconv.Itoa(l.NumMembers)
	out.w.Write([]string{l.ChannelName, l.UserName, l.CreatedDate, members, l.Purpose})
}

type basicOut struct {
	w *tabwriter.Writer
}

func (out *basicOut) Header() {
	fmt.Fprintf(out.w, "Name\tCreator\tCreatedDate\tNumMembers\tPurpose\n")
}

func (out *basicOut) WriteLine(l *Line) {
	fmt.Fprintf(out.w, "%s\t%s\t%s\t%d\t%s\n", l.ChannelName, l.UserName, l.CreatedDate, l.NumMembers, l.Purpose)
}

func (out *basicOut) Footer() {
	out.w.Flush()
}
