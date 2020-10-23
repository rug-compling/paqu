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
			processLock.Lock()
			owner := task.owner
			// zoek eerste plek na hoogste positie voor owner
			idx := 0
			for i := top - 1; i >= bottom; i-- {
				if workmap[i].owner == owner {
					idx = i + 1
					break
				}
			}
			// zoek laagste posities van alle andere owners vanaf idx
			seen := make(map[string]bool)
			for i := idx; i < top; i++ {
				if !seen[workmap[i].owner] {
					idx = i + 1
					seen[workmap[i].owner] = true
				}
			}
			// invoegen
			for i := top; i > idx; i-- {
				workmap[i] = workmap[i-1]
			}
			workmap[idx] = task
			top++
			// hernummeren
			for i := bottom; i < top; i++ {
				workmap[i].nr = i - bottom
			}
			logf("SCHEDULED: %v", task.id)
			processLock.Unlock()
		case id := <-chDelete:
			processLock.Lock()
			found := false
			for i := bottom; i < top; i++ {
				if found {
					workmap[i] = workmap[i+1]
					workmap[i].nr = i - bottom
				} else {
					if workmap[i].id == id {
						i--
						top--
						found = true
					}
				}
			}
			if found {
				logf("UNSCHEDULED: %v", id)
				delete(workmap, top)
			}
			processLock.Unlock()
		case sema <- struct{}{}:
			processLock.Lock()
			if bottom < top {
				wg.Add(1)
				go func(p *Process) {
					logf("FROM SCHEDULE: %v", p.id)
					work(p)
					wg.Done()
					<-semaphore
				}(workmap[bottom])
				delete(workmap, bottom)
				bottom++
				// hernummeren
				for i := bottom; i < top; i++ {
					workmap[i].nr = i - bottom
				}
			}
			if len(workmap) == 0 {
				// TODO: mag dit?
				top = 0
				bottom = 0
			}
			processLock.Unlock()
		}
	}
}
