package kiambe

import (
	"errors"
	"net/http"

	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/kiam"
	"github.com/ariefdarmawan/kiam/acm"
	"github.com/dgrijalva/jwt-go"
	"github.com/eaciit/toolkit"
)

type loginEngine struct {
	LoginID    string `kf-pos:"1,1" required:"true" label:"Login ID, email atau mobile no"`
	Password   string `kf-pos:"2,1" required:"true" label:"Password" kf-control:"password"`
	RememberMe bool   `kf-pos:"3,1"`

	im                 *kiam.Manager
	fnPostAuthenticate func(u *acm.User, m *toolkit.M)
}

func NewLoginEngine(im *kiam.Manager) *loginEngine {
	l := new(loginEngine)
	l.im = im
	return l
}

func (l *loginEngine) SetFnPostAuthenticate(fn func(u *acm.User, m *toolkit.M)) *loginEngine {
	l.fnPostAuthenticate = fn
	return l
}

/* need to do:
taken care of allowed number of multi session, currently has no control on it
*/
func (l *loginEngine) Authenticate(ctx *kaos.Context, req *loginEngine) (toolkit.M, error) {
	signMtd := l.im.Options().SignMethod
	signKey := []byte(l.im.Options().SignSecret)

	h, e := ctx.DefaultHub()
	if e != nil {
		return toolkit.M{}, errors.New("invalid hub")
	}

	// auth
	var (
		uid string
	)
	mgr := acm.NewACM(h)
	if req == nil {
		hr := ctx.Data().Get("http-request", new(http.Request)).(*http.Request)
		loginID, password, ok := hr.BasicAuth()
		if !ok {
			return toolkit.M{}, errors.New("unable to get auth info")
		}
		uid, e = mgr.Authenticate(loginID, password)
		if e != nil {
			return toolkit.M{}, e
		}
	} else {
		uid, e = mgr.Authenticate(req.LoginID, req.Password)
		if e != nil {
			return toolkit.M{}, e
		}
	}

	usr, _ := mgr.GetUser("ID", uid)
	tokenData := toolkit.M{}
	if l.fnPostAuthenticate != nil {
		l.fnPostAuthenticate(usr, &tokenData)
	}

	// setup jwt
	token, err := l.im.FindOrCreate(ctx, toolkit.M{}.Set("ID", uid).Set("Duration", 0), tokenData)
	if err != nil {
		return toolkit.M{}, err
	}
	// jwt-ed
	bc := new(jwt.StandardClaims)
	bc.Id = token.SessionID
	jtkn := jwt.NewWithClaims(signMtd, bc)
	tknString, err := jtkn.SignedString(signKey)
	if err != nil {
		return toolkit.M{}, err
	}

	res := toolkit.M{}
	for k, v := range tokenData {
		res.Set(k, v)
	}
	res.Set("Token", tknString).Set("Name", usr.Name)
	return res, nil
}

func (l *loginEngine) Logout(ctx *kaos.Context, req toolkit.M) (string, error) {
	tokenid := ctx.Data().Get("jwt-token-id", "").(string)
	if tokenid == "" {
		return "", nil
	}

	l.im.Remove(ctx, toolkit.M{}.Set("ID", tokenid))
	return "", nil
}
