package mail

type Client struct{}

func (c Client) Send(_ string, _ string, _ string) error { return nil }
