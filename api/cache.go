package main

type Context struct {
	role    string
	content string
}

type QueryCache struct {
	Cache map[string][]Context
}

func (qc *QueryCache) Add(sessionId string, userInput string, machineOutput string) {
	var context []Context
	context = qc.Cache[sessionId]

	if len(context)+2 <= cap(context) {
		context = context[2:]
	}
	context = append(context, Context{
		role:    "user",
		content: userInput,
	})
	context = append(context, Context{
		role:    "assistant",
		content: machineOutput,
	})
}
