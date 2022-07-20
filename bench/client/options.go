package client

import (
	"github.com/isucon/isucandar/agent"

	"github.com/recruit-tech/RISUCON2022Summer/bench/constant"
)

var clientOptsGeneratorMap = map[ClientType](func() []agent.AgentOption){
	InitializerType: func() []agent.AgentOption {
		return []agent.AgentOption{
			agent.WithBaseURL(targetUrl),
			agent.WithNoCookie(),
			agent.WithNoCache(),
			agent.WithCloneTransport(agent.DefaultTransport),
			agent.WithTimeout(constant.InitializeTimeout),
			agent.WithUserAgent(InitializerType),
		}
	},
	CompatibilityCheckerType: func() []agent.AgentOption {
		return []agent.AgentOption{
			agent.WithBaseURL(targetUrl),
			agent.WithNoCache(),
			agent.WithCloneTransport(agent.DefaultTransport),
			agent.WithTimeout(constant.CompatibilityCheckRequestTimeout),
			agent.WithUserAgent(CompatibilityCheckerType),
		}
	},
	LoaderType: func() []agent.AgentOption {
		return []agent.AgentOption{
			agent.WithBaseURL(targetUrl),
			agent.WithNoCache(),
			agent.WithCloneTransport(agent.DefaultTransport),
			agent.WithTimeout(constant.LoadRequestTimeout),
			agent.WithUserAgent(LoaderType),
		}
	},
}
