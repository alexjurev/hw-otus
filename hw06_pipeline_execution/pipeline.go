package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	inChan := in

	for _, s := range stages {
		outChan := make(Bi)

		go func(s Stage, out Bi, in In) {
			defer close(out)
			ch := s(in)

			for {
				select {
				case data, ok := <-ch:
					if !ok {
						return
					}
					out <- data
				case <-done:
					return
				}
			}
		}(s, outChan, inChan)
		inChan = outChan
	}

	return inChan
}
