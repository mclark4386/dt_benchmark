package actions

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
	"github.com/mclark4386/dt_benchmark/middleware"
	"github.com/mclark4386/dt_benchmark/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// LoginRequest holds the two fields used in a login form
// allows for simpler/cleaner binding to a login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	req := &LoginRequest{}

	if err := c.Bind(req); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	u := &models.User{}

	// setup error function for later use
	bad := func(msg string) error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password: "+msg)
		c.Set("errors", verrs)
		return c.Error(http.StatusUnprocessableEntity, verrs)
	}

	// check for empty password before anything else
	pwd := req.Password
	if len(pwd) == 0 {
		return bad("bad password")
	}

	if err := tx.Where("email = ?", strings.TrimSpace(strings.ToLower(req.Email))).First(u); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return bad("no user with that email")
		}
		return errors.WithStack(err)
	}

	// compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		return bad("bad password")
	}

	if err := BuildSetToken(c, u); err != nil {
		return c.Error(http.StatusBadRequest, err)
	}
	return c.Render(200, r.JSON(u))
}

// BuildSetToken builds and sets a new JWT
// after checking the refresh token for the user
func BuildSetToken(c buffalo.Context, u *models.User) error {
	// check for refresh token
	tx := c.Value("tx").(*pop.Connection)
	if u.RefreshToken.UUID == uuid.Nil {
		u.RefreshToken = nulls.NewUUID(uuid.Must(uuid.NewV4()))
		verrs, err := u.Update(tx)
		if err != nil {
			return errors.WithStack(err)
		}
		if verrs.HasAny() {
			return c.Render(400, r.JSON(verrs))
		}
	}

	token, err := middleware.BuildJWT(u, 0)
	if err != nil {
		return err
	}
	middleware.SetToken(c, token)
	c.Set("current_user", u)
	return nil
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	u := c.Value("current_user").(*models.User)
	u.RefreshToken.UUID = uuid.Nil
	u.RefreshToken.Valid = false

	tx := c.Value("tx").(*pop.Connection)
	tx.Save(u)

	middleware.SetToken(c, "")
	return c.Render(200, r.JSON("success"))
}
