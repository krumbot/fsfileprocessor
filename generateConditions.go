package fsfileprocessor

func (config Crawler) generateConditionFunction(errChannel chan error) (func(infoChannel <-chan WalkInfo), <-chan WalkInfo) {
	validChannel := make(chan WalkInfo)

	recursionCheck := config.generateRecursionCheck(errChannel)
	fileExtCheck := config.generateFileExtCheck(errChannel)
	earliestTimeModifiedCheck := config.genereateEarliestTimeModifiedCheck(errChannel)

	conditionals := append(config.Conditionals, recursionCheck, fileExtCheck, earliestTimeModifiedCheck)

	conditionFunc := func(infoChannel <-chan WalkInfo) {
		defer close(validChannel)

		for info := range infoChannel {
			shouldPass := true
			shouldPassChannel := make(chan bool, len(conditionals))

			for _, fnc := range conditionals {
				go fnc(shouldPassChannel, info)
			}

			for i := 0; i < len(conditionals); i++ {
				shouldPass = shouldPass && <-shouldPassChannel
			}

			if shouldPass {
				validChannel <- info
			}
		}
	}

	return conditionFunc, validChannel

}

func (config Crawler) generateRecursionCheck(errChannel chan error) func(shouldPassChannel chan<- bool, w WalkInfo) {
	if config.Controller.Recursive {
		return func(s chan<- bool, w WalkInfo) { s <- true }
	}

	return func(s chan<- bool, w WalkInfo) {
		if w.Info.IsDir() {
			s <- false
		} else {
			s <- true
		}
	}
}

func (config Crawler) generateFileExtCheck(errChannel chan error) func(shouldPassChannel chan<- bool, w WalkInfo) {
	if config.Controller.FileExt == nil {
		return func(s chan<- bool, w WalkInfo) { s <- true }
	}

	return func(s chan<- bool, w WalkInfo) {
		s <- config.Controller.FileExt.MatchString(w.Path)
	}
}

func (config Crawler) genereateEarliestTimeModifiedCheck(errChannel chan error) func(shouldPassChannel chan<- bool, w WalkInfo) {
	if config.Controller.EarliestTimeModified.IsZero() {
		return func(s chan<- bool, w WalkInfo) { s <- true }
	}

	return func(s chan<- bool, w WalkInfo) {
		s <- w.Info.ModTime().After(config.Controller.EarliestTimeModified)
	}
}
