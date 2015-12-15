package main

func scheduler() {
	workmap := make(map[int]*Process)
	var top, bottom int
	for {
		var sema chan struct{}
		if top != bottom {
			sema = semaphore
		}
		select {
		case <-chGlobalExit:
			return
		case task := <-chWork:
			workmap[top] = task
			top++
		case sema <- struct{}{}:
			wg.Add(1)
			go func(p *Process) {
				work(p)
				wg.Done()
				<-semaphore
			}(workmap[bottom])
			delete(workmap, bottom)
			bottom++
		}
	}
}
