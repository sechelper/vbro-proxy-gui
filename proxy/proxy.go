package proxy

import (
	"crypto/tls"
	"github.com/google/martian/v3/header"
	"github.com/rs/xid"
	"github.com/rs/zerolog/log"
	"github.com/sechelper/vbro/proxy"
	"github.com/sechelper/vbro/proxy/api"
	"github.com/sechelper/vbro/proxy/cache"
	"github.com/sechelper/vbro/proxy/filter"
	"github.com/sechelper/vbro/proxy/intercepter"
	"github.com/sechelper/vbro/utils"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run() {

	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(viper.GetInt64("transport.timeout")) * time.Second,
			KeepAlive: time.Duration(viper.GetInt64("transport.Keep_alive")) * time.Second,
		}).DialContext,
		TLSHandshakeTimeout:   time.Duration(viper.GetInt64("transport.tls_handshake_timeout")) * time.Second,
		ExpectContinueTimeout: time.Duration(viper.GetInt64("transport.expect_continue_timeout")) * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: viper.GetBool("transport.insecure_skip_verify"),
		},
	}

	vbroProxy := proxy.NewVbroProxy(
		proxy.Address{
			Host: viper.GetString("listen.host"),
			Port: viper.GetInt("listen.port"),
		},
		transport)

	if viper.GetString("proxy.url") != "" {
		if err := vbroProxy.UseProxy(viper.GetString("proxy.url")); err != nil {
			panic(err)
		}
	}

	// 配置支持https
	if err := vbroProxy.MitmConfig(viper.GetString("certificate.org"), viper.GetString("certificate.cert"),
		viper.GetString("certificate.key"),
		viper.GetDuration("certificate.validity"), viper.GetBool("certificate.tlsVerifySkip")); err != nil {
		log.Error().Msg(err.Error())
		return
	}

	//fg := fifo.NewGroup()
	//
	//fg.AddRequestModifier(header.NewBadFramingModifier())
	//fg.AddRequestModifier(header.NewHopByHopModifier())
	//fg.AddResponseModifier(header.NewHopByHopModifier())

	// 存储历史记录
	//f0 := filter.NewFilter()
	//var mg = filter.MatcherSpecifyGroup{filter.NewMatcherSpecify(xid.New().String(),
	//	filter.MatchNone, filter.MatchHttpMethod, // 执行
	//	filter.DoesNoMatch, `CONNECT`)}
	//f0.Matcher.AddRequestSpecifyGroup(mg...)
	//f0.Matcher.AddResponseSpecifyGroup(mg...)
	//record := history.NewHistory(cache.NewRedis(fmt.Sprintf("%s:%d", viper.GetString("cache.redis.address"),
	//	viper.GetInt("cache.redis.port")), viper.GetString("cache.redis.password"), viper.GetInt("cache.redis.db")))
	//f0.SetRequestModifier(record)
	//f0.SetResponseModifier(record)
	//fg.AddRequestModifier(f0)
	//fg.AddResponseModifier(f0)

	// 拦截器
	//f1 := filter.NewFilter()
	//f1.MatcherIntercept.AddRequestSpecifyGroup(filter.NewMatcherSpecify(xid.New().String(),
	//	filter.MatchNone, filter.MatchDomainName, // 执行
	//	filter.Match, `.*\.sechelper\.com`),
	//	filter.NewMatcherSpecify(xid.New().String(), filter.MatchAnd, filter.MatchFileExtension, // 执行
	//		filter.Match, `\.(gif|jpg|png|css|js|ico|svg|eot|woff|woff2|ttf)$`),
	//	filter.NewMatcherSpecify(xid.New().String(), filter.MatchAnd, filter.MatchHeaderKey, // 不执行
	//		filter.Match, `User-Agent`),
	//	filter.NewMatcherSpecify(xid.New().String(), filter.MatchOr, filter.MatchHttpMethod, // 执行
	//		filter.Match, `GET`))
	//interceptor := intercept.NewIntercept()
	//f1.SetRequestModifier(interceptor)
	//f1.SetResponseModifier(interceptor)
	//cache.FIFOGroup.AddRequestModifier(f1)
	//cache.FIFOGroup.AddResponseModifier(f1)

	//vbroProxy.SetRequestModifier(fg)
	//vbroProxy.SetResponseModifier(fg)

	// 拦截器组
	cache.IntercepterGroup = intercepter.NewGroup()
	cache.IntercepterGroup.AddRequestModifier(header.NewBadFramingModifier())
	cache.IntercepterGroup.AddRequestModifier(header.NewHopByHopModifier())
	cache.IntercepterGroup.AddResponseModifier(header.NewHopByHopModifier())
	cache.IntercepterGroup.SetAggregateErrors(true)
	f1 := filter.NewFilter()
	f1.Matcher.AddRequestSpecifyGroup(filter.NewMatcherSpecify(xid.New().String(),
		filter.MatchAnd, filter.MatchDomainName, // 执行
		filter.Match, viper.GetString("filter.domain")))
	f1.SetRequestModifier(cache.IntercepterGroup)
	f1.SetResponseModifier(cache.IntercepterGroup)
	vbroProxy.SetRequestModifier(f1)
	vbroProxy.SetResponseModifier(f1)

	go func() {
		if err := vbroProxy.Run(); err != nil {
			panic(err)
		}
	}()

	// 使用 unix sock 方式建立api服务
	l, _ := utils.NewGinWithSock("vbro-proxy.sock")
	go api.HttpServer.RunListener(l)

	//go api.HttpServer.Run(fmt.Sprintf("%s:%d", viper.GetString("api.listen.host"),
	//	viper.GetInt("api.listen.port")))
	wait()
}

// exit: ctr + c
func wait() {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)

	<-s

	os.Exit(0)
}

func init() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
