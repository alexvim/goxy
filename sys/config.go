package sys

var Config config

type config struct {
	ConnectionAddress string
	ExternalAddress   string
}

func (c *config) Read {

}

func (c *config) ReadFromJson {

}

func (c *config) ReadFromEnv {

}
