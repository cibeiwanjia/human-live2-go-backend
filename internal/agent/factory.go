package agent

// NewAgent creates an agent by name
func NewAgent(name, desc string) Agent {
	switch name {
	case "RepeaterAgent":
		return NewRepeaterAgent()
	default:
		return nil
	}
}
