package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func drainChan(in In) {
	for range in { //nolint: revive
	}
}

func doneWrapper(in In, done In) Out {
	out := make(Bi)

	go func() {
		defer drainChan(in)
		defer close(out)
		for {
			select {
			case <-done:
				return
			case data, ok := <-in:
				if !ok {
					return
				}
				out <- data
			}
		}
	}()

	return out
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		wrapped := doneWrapper(in, done)
		in = stage(wrapped)
	}
	return in
}
