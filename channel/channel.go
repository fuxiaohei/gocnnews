package channel

import "github.com/gocnnews/model"

type ChannelFunc func() []*model.Article

var Channels []ChannelFunc

func init() {
	Channels = append(Channels, Hexacosa)
	Channels = append(Channels, Golangnews)
	Channels = append(Channels, Juejin)
	Channels = append(Channels, Oschina)
	Channels = append(Channels, Segmentfault)
	Channels = append(Channels, Toutiao)
	Channels = append(Channels, Reddit)
	Channels = append(Channels, Hacknews)
}
