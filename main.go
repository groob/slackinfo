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
		flUserInfo = flag.Bool("userinfo", false, "print user information")
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

	var w userInfoOutputter
	if *flCSV {
		w = &csvOut{w: csv.NewWriter(os.Stdout)}
	} else {
		w = &basicOut{w: tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)}
	}

	api := slack.New(*flAPIToken)
	client := &Client{client: api}
	if *flUserInfo {
		client.UserInfo(w)
		return
	}
	client.ChannelInfo(w)
}

type Client struct {
	client *slack.Client
}

func (c *Client) UserInfo(w userInfoOutputter) {
	users, err := c.client.GetUsers()
	if err != nil {
		log.Fatal(err)
	}
	w.UserHeader()
	defer w.Footer()

	for _, u := range users {
		line := &UserLine{
			Name:              u.Name,
			RealName:          u.Profile.RealName,
			Title:             u.Profile.Title,
			Email:             u.Profile.Email,
			IsAdmin:           u.IsAdmin,
			IsRestricted:      u.IsRestricted,
			IsUltraRestricted: u.IsUltraRestricted,
			Has2FA:            u.Has2FA,
		}
		w.WriteUserLine(line)
	}
}

func (c *Client) ChannelInfo(w outputter) {
	channels, err := c.client.GetChannels(true)
	if err != nil {
		log.Fatal(err)
	}

	w.Header()
	defer w.Footer()

	for _, ch := range channels {
		user, err := c.client.GetUserInfo(ch.Creator)
		if err != nil {
			log.Println(err)
			continue
		}
		line := &Line{
			ChannelName: ch.Name,
			UserName:    user.Name,
			CreatedDate: ch.Created.Time().UTC().String(),
			NumMembers:  ch.NumMembers,
			Purpose:     ch.Purpose.Value,
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

type UserLine struct {
	Name              string
	RealName          string
	Title             string
	Email             string
	IsAdmin           bool
	IsRestricted      bool
	IsUltraRestricted bool
	Has2FA            bool
}

type userInfoOutputter interface {
	outputter
	UserHeader()
	WriteUserLine(*UserLine)
}

type csvOut struct {
	w *csv.Writer
}

func (out *csvOut) UserHeader() {
	out.w.Write([]string{"Name", "RealName", "Title", "Email", "IsAdmin", "IsRestricted", "IsUltraRestricted", "Has2FA"})
}

func (out *csvOut) Header() {
	out.w.Write([]string{"Name", "Creator", "CreatedDate", "NumMembers", "Purpose"})
}

func (out *csvOut) Footer() {
	out.w.Flush()
}

func (out *csvOut) WriteUserLine(l *UserLine) {
	out.w.Write([]string{
		l.Name,
		l.RealName,
		l.Title,
		l.Email,
		fmt.Sprint(l.IsAdmin),
		fmt.Sprint(l.IsRestricted),
		fmt.Sprint(l.IsUltraRestricted),
		fmt.Sprint(l.Has2FA),
	})
}

func (out *csvOut) WriteLine(l *Line) {
	members := strconv.Itoa(l.NumMembers)
	out.w.Write([]string{l.ChannelName, l.UserName, l.CreatedDate, members, l.Purpose})
}

type basicOut struct {
	w *tabwriter.Writer
}

func (out *basicOut) UserHeader() {
	fmt.Fprintf(out.w, "Name\tRealName\tTitle\tEmail\tIsAdmin\tIsRestricted\tIsUltraRestricted\tHas2FA\n")
}

func (out *basicOut) WriteUserLine(l *UserLine) {
	fmt.Fprintf(out.w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		l.Name,
		l.RealName,
		l.Title,
		l.Email,
		fmt.Sprint(l.IsAdmin),
		fmt.Sprint(l.IsRestricted),
		fmt.Sprint(l.IsUltraRestricted),
		fmt.Sprint(l.Has2FA),
	)
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
