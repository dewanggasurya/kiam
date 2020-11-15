package kiamui

import (
	"errors"

	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/kiam"
	"github.com/ariefdarmawan/kiam/acm"
	"github.com/eaciit/toolkit"
)

type loginEngine struct {
	LoginID    string `kf-pos:"1,1" required:"true" label:"Login ID, email atau mobile no"`
	Password   string `kf-pos:"2,1" required:"true" label:"Password" kf-control:"password"`
	RememberMe bool   `kf-pos:"3,1"`

	im *kiam.Manager
}

func NewLoginEngine(im *kiam.Manager) *loginEngine {
	l := new(loginEngine)
	l.im = im
	return l
}

/* need to do:
taken care of allowed number of multi session, currently has no control on it
*/
func (l *loginEngine) Authenticate(ctx *kaos.Context, req *loginEngine) (toolkit.M, error) {
	h, e := ctx.DefaultHub()
	if e != nil {
		return toolkit.M{}, errors.New("invalid hub")
	}

	// auth
	mgr := acm.NewACM(h)
	uid, e := mgr.Authenticate(req.LoginID, req.Password)
	if e != nil {
		return toolkit.M{}, e
	}

	usr, _ := mgr.GetUser("ID", uid)

	// setup jwt
	token, err := l.im.FindOrCreate(ctx, toolkit.M{}.Set("ID", uid).Set("Duration", 0))
	if err != nil {
		return toolkit.M{}, err
	}

	return toolkit.M{}.Set("Token", token.SessionID).Set("Name", usr.Name), nil
}

func (l *loginEngine) Logout(ctx *kaos.Context, req toolkit.M) (string, error) {
	tokenid := ctx.Data().Get("jwt-token-id", "").(string)
	if tokenid == "" {
		return "", nil
	}

	l.im.Remove(ctx, toolkit.M{}.Set("ID", tokenid))
	return "", nil
}
