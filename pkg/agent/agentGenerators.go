package agent

type AgentGenerator[T IAgent[T]] func(IExposedServerFunctions[T]) T

type AgentGeneratorCountPair[T IAgent[T]] struct {
	Generator AgentGenerator[T]
	Count     int
}

func MakeAgentGeneratorCountPair[T IAgent[T]](generatorFunction AgentGenerator[T], count int) AgentGeneratorCountPair[T] {
	return AgentGeneratorCountPair[T]{
		Generator: generatorFunction,
		Count:     count,
	}
}
