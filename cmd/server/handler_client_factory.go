package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory/herodot"
	"github.com/ory/hydra/client"
	"github.com/ory/hydra/config"
	"github.com/vpugar/hydra-boltdb-backend/backend"
)

func newClientManager(c *config.Config) client.Manager {
	ctx := c.Context()

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		return &client.MemoryManager{
			Clients: map[string]client.Client{},
			Hasher:  ctx.Hasher,
		}
	case *config.SQLConnection:
		return &client.SQLManager{
			DB:     con.GetDatabase(),
			Hasher: ctx.Hasher,
		}
	case *config.PluginConnection:
		if m, err := con.NewClientManager(); err != nil {
			c.GetLogger().Fatalf("Could not load client manager plugin %s", err)
		} else {
			return m
		}
		break
	case *backend.BoltdbConnection:
		if m, err := con.NewClientManager(ctx.Hasher); err != nil {
			c.GetLogger().Fatalf("Could not load client manager plugin %s", err)
		} else {
			return m
		}
		break
	default:
		panic("Unknown connection type.")
	}
	return nil
}

func newClientHandler(c *config.Config, router *httprouter.Router, manager client.Manager) *client.Handler {
	ctx := c.Context()
	h := &client.Handler{
		H: herodot.NewJSONWriter(c.GetLogger()),
		W: ctx.Warden, Manager: manager,
	}

	h.SetRoutes(router)
	return h
}
