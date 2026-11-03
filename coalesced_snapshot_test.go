package main

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

// Pins theorems/Theorems/CoalesceAlias.lean: waiters never alias the primary's response buffer.
func TestSnapshotCoalescedResponse_OwnsBody(t *testing.T) {
	app := fiber.New()
	fctx := &fasthttp.RequestCtx{}
	c := app.AcquireCtx(fctx)
	defer app.ReleaseCtx(c)

	c.Response().SetStatusCode(201)
	c.Response().SetBody([]byte(`{"user":"alice"}`))

	snap := snapshotCoalescedResponse(c)

	c.Response().SetBody([]byte(`{"user":"bobby"}`)) // buffer reused by a later request

	assert.Equal(t, 201, snap.StatusCode)
	assert.Equal(t, `{"user":"alice"}`, string(snap.Body))
}
