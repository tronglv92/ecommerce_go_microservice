package consul

import (
	"flag"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"github.com/tronglv92/ecommerce_go_common/logger"
)

type consulClient struct {
	prefix string

	logger       logger.Logger
	consulClient *api.Client
	TTL          time.Duration
	CheckID      string
	ID           string
	NameCluster  string
	Tags         []string
	Address      string
	Port         int
}

func NewConsulClient(prefix string, tag string) *consulClient {
	return &consulClient{
		prefix: prefix,
		Tags:   []string{tag},
	}
}

func (cs *consulClient) GetPrefix() string {
	return cs.prefix
}

func (cs *consulClient) Get() interface{} {
	return cs
}

func (cs *consulClient) Name() string {
	return cs.prefix
}

func (cs *consulClient) InitFlags() {
	flag.DurationVar(&cs.TTL, cs.GetPrefix()+"-ttl", time.Second*8, "ttl and deregister server after second ")
	flag.StringVar(&cs.CheckID, cs.GetPrefix()+"-checkid", "checkalive", "Check ID")
	flag.StringVar(&cs.ID, cs.GetPrefix()+"-id", "login_service", "ID")
	flag.StringVar(&cs.NameCluster, cs.GetPrefix()+"-name-cluster", "mycluster", "NameC luster")

	flag.StringVar(&cs.Address, cs.GetPrefix()+"-address", "127.0.0.1", "Address")
	flag.IntVar(&cs.Port, cs.GetPrefix()+"-port", 3000, "Port")
	flag.Parse()
}

func (cs *consulClient) Configure() error {
	cs.logger = logger.GetCurrent().GetLogger(cs.prefix)

	cs.logger.Debugf("ttl consul: %v", cs.TTL)

	cs.logger.Debugf("tags consul: %v", cs.Tags)
	config := api.DefaultConfig()
	client, err := api.NewClient(config)
	if err != nil {

		return err
	}
	cs.consulClient = client

	if err := cs.registerService(); err != nil {
		return err
	}
	go cs.updateHealthCheck()
	return nil
}

func (cs *consulClient) Run() error {
	return cs.Configure()
}

func (cs *consulClient) Stop() <-chan bool {
	c := make(chan bool)

	go func() {

		cs.consulClient.Agent().ServiceDeregister(cs.ID)
		c <- true
	}()
	return c
}

func (cs *consulClient) registerService() error {
	check := &api.AgentServiceCheck{
		DeregisterCriticalServiceAfter: cs.TTL.String(),
		TLSSkipVerify:                  true,
		TTL:                            cs.TTL.String(),
		CheckID:                        cs.CheckID,
	}

	serviceName := cs.ID
	uniqueIdentifier := uuid.New().String()
	serviceID := fmt.Sprintf("%s-%s", serviceName, uniqueIdentifier)
	cs.ID = serviceID

	register := &api.AgentServiceRegistration{
		ID:      cs.ID,
		Name:    cs.NameCluster,
		Tags:    cs.Tags,
		Address: cs.Address,
		Port:    cs.Port,
		Check:   check,
	}
	if err := cs.consulClient.Agent().ServiceRegister(register); err != nil {
		return err
	}

	return nil
}
func (cs *consulClient) updateHealthCheck() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		err := cs.consulClient.Agent().UpdateTTL(cs.CheckID, "online", api.HealthPassing)
		if err != nil {
			cs.logger.Errorf("update health check failed: %v", err)
		}
		<-ticker.C
	}
}
