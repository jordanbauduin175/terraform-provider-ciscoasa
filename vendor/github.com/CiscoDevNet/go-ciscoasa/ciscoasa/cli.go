package ciscoasa

type CLIRequest struct {
	Commands []string `json:"commands"`
}

func (c *Client) PostCLI(commands []string) error {
	reqBody := CLIRequest{
		Commands: commands,
	}

	req, err := c.newRequest("POST", "/api/cli", reqBody)
	if err != nil {
		return err
	}

	_, err = c.do(req, nil)
	return err
}
