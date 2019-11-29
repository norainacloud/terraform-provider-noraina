package noraina

import (
	"context"
	"fmt"
	"net/http"
)

func (c *Client) deleteResource(ctx context.Context, id string, resourceRoute string) error {
	if id == "" {
		return fmt.Errorf("[ERROR] deleteResource: Resource id cannot be blank")
	}

	if resourceRoute == "" {
		return fmt.Errorf("[ERROR] deleteResource: Resource route cannot be blank")
	}

	path := fmt.Sprintf("%s/%s", resourceRoute, id)
	req, err := c.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}

	return c.Do(ctx, req, nil)
}
