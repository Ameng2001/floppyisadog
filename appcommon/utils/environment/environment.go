package environment

import (
	"fmt"

	"github.com/TarsCloud/TarsGo/tars/util/conf"
)

type Environment struct {
	Name            string
	IsDebug         bool
	ExternalApex    string // Apex domain off of which services operate externally
	InternalApex    string // Apex domain off of which services operate internally
	LogLevel        string // Verbosity of logging
	Scheme          string // default URL scheme - http or https
	JWTTokenSecret  string // secret for generating JWT token
	CSRFTokenSecret string // secret for generating CSRF token
}

var (
	currentEnv string
	envs       map[string]*Environment
)

func InitEnvironment(c *conf.Conf, domain string) bool {
	currentEnv = c.GetStringWithDef(domain+"<CurrentEnvironment>", "development")

	subdomains := c.GetDomain(domain)
	envs = make(map[string]*Environment)
	for _, env := range subdomains {
		envInfo := &Environment{
			Name:            c.GetString(domain + env + "<Name>"),
			IsDebug:         c.GetBoolWithDef(domain+env+"<IsDebug>", false),
			ExternalApex:    c.GetString(domain + env + "<ExternalApex>"),
			InternalApex:    c.GetString(domain + env + "<InternalApex>"),
			LogLevel:        c.GetString(domain + env + "<LogLevel>"),
			Scheme:          c.GetString(domain + env + "<Scheme>"),
			JWTTokenSecret:  c.GetString(domain + env + "<JWTTokenSecret>"),
			CSRFTokenSecret: c.GetString(domain + env + "<CSRFTokenSecret>"),
		}
		envs[env] = envInfo
		fmt.Printf("Envirionment (%s) : (%v)\n", env, envInfo)
	}

	return true
}

func GetEnvs() map[string]*Environment {
	return envs
}

func GetCurrEnv() *Environment {
	env, ok := envs[currentEnv]
	if !ok {
		return nil
	}

	return env
}
