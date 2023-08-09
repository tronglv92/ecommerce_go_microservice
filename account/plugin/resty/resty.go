package resty

import (
	"github.com/go-resty/resty/v2"
	"github.com/tronglv92/accounts/common"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

func NewRestService() *restService { return &restService{} }

type restService struct {
	client *resty.Client
	// serviceURL string
	logger logger.Logger
}

func (*restService) GetPrefix() string {
	return common.PluginRestService
}

func (s *restService) Get() interface{} {
	return s.client
}

func (restService) Name() string {
	return common.PluginRestService
}

func (s *restService) InitFlags() {
	// flag.StringVar(&s.serviceURL, s.GetPrefix()+"-url", "", "URL of user service (Ex: http://user-service:8080)")
}

func (s *restService) Configure() error {
	s.client = resty.New()
	s.logger = logger.GetCurrent().GetLogger(s.GetPrefix())

	// if s.serviceURL == "" {
	// 	s.logger.Errorln("Missing service URL")
	// 	return errors.New("missing service URL")
	// }

	return nil
}

func (s *restService) Run() error {
	return s.Configure()
}

func (s *restService) Stop() <-chan bool {
	c := make(chan bool)

	go func() {
		c <- true
		s.logger.Infoln("Stopped")
	}()
	return c
}
