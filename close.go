package redis_replica

// setCloseStatus устанавливает статус соединения "close"
func (c *Client) setCloseStatus() {
	c.Lock()
	defer c.Unlock()

	c.closeStatus = true
	for _, ch := range c.closeNotifyChan {
		ch <- struct{}{}
	}
	c.closeNotifyChan = nil
}

// WaitClose ждёт закрытия соединения с мастером
func (c *Client) WaitClose() (<-chan struct{}) {
	ch := make(chan struct{}, 1)

	c.Lock()
	defer c.Unlock()

	if c.closeStatus {
		ch <- struct{}{}
	} else {
		c.closeNotifyChan = append(c.closeNotifyChan, ch)
	}
	return ch
}
