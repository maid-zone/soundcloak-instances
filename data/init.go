package data

import "errors"

type Link struct {
	URL  string
	Text string
}

type Instance struct {
	URL         Link
	SkipResolve bool   // whether the script should skip checking if instance is up and getting it's settings
	Country     string // Flag of the country, where this instance is hosted
	Maintainer  Link   // What person or group is maintaining this instance
	Note        string // just a note with some important things about the instance
}

const Germany = "🇩🇪"
const Netherlands = "🇳🇱"
const Russia = "🇷🇺"

var Instances = []Instance{
	{
		URL:        Link{"https://sc.maid.zone", "sc.maid.zone"},
		Country:    Germany,
		Maintainer: Link{"https://maid.zone", "maid.zone"},
	},

	{
		URL:        Link{"https://tunes.floppa.nl", "tunes.floppa.nl"},
		Country:    Germany,
		Maintainer: Link{"https://jonas.tf", "jonas"},
	},

	{
		URL:        Link{"https://sc.bloat.cat", "sc.bloat.cat"},
		Country:    Germany,
		Maintainer: Link{"https://bloat.cat", "bloat.cat"},
	},

	{
		URL:        Link{"https://soundcloak.fly.dev", "soundcloak.fly.dev"},
		Country:    Netherlands,
		Maintainer: Link{"https://laptopc.at", "laptopcat"},
		Note:       "This is the staging/flagship instance. I usually push changes to here even before I commit them to the repository to test stuff out.",
	},

	{
		URL:        Link{"https://soundcloak.laincorp.tech", "soundcloak.laincorp.tech"},
		Country:    Russia,
		Maintainer: Link{"https://laincorp.tech", "lain"},
	},
}

// We're only interested in these 3 properties for now
type InstanceInfo struct {
	Commit       string
	Repo         string
	ProxyImages  bool
	ProxyStreams bool
	Restream     bool

	// This isn't returned by /_/info, but used to indicate that we skipped resolving this instance's info OR that something went wrong
	SkippedResolve bool
	Error          error
}

var ErrNotFound = errors.New("got status code 404")
