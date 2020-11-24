package kiambe

import (
	"errors"

	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/kiam/acm"
	"github.com/ariefdarmawan/kmsg"
	"github.com/eaciit/toolkit"
)

type RegisterOptions struct {
	SendNotifTopic string
	FnPostRegister func(usr *acm.User, parm toolkit.M)
}

type registerEngine struct {
	Name           string `kf-pos:"1,1" label:"Nama lengkap" required:"true"`
	LoginID        string `kf-pos:"2,1" label:"Login ID" required:"true"`
	Email          string `kf-pos:"3,1" label:"Email" required:"true"`
	Phone          string `kf-pos:"4,1" label:"Mobile phone" required:"true"`
	Password       string `kf-pos:"5,1" kf-control:"password" required:"true"`
	ConfimPassword string `kf-pos:"6,1" kf-control:"password" label:"Confirm password" required:"true"`

	opts RegisterOptions
}

func NewRegisterEngine(o RegisterOptions) *registerEngine {
	r := new(registerEngine)
	r.opts = o
	return r
}

func (re *registerEngine) Register(ctx *kaos.Context, parm toolkit.M) (string, error) {
	h, _ := ctx.DefaultHub()
	if h == nil {
		return "", errors.New("invalid data hub")
	}

	req := new(registerEngine)
	err := toolkit.Serde(parm, req, "")
	if err != nil {
		return "", errors.New("fail to register user: " + err.Error())
	}

	mgr := acm.NewACM(h)
	user, err := mgr.CreateUser(req.LoginID, req.Name, req.Email, req.Phone, req.Password)
	if err != nil {
		return "", errors.New("fail to register user: " + err.Error())
	}

	msg, e := kmsg.NewMessage(h, "UserRegistration", "SMTP", "", user.Email,
		"User anda berhasil diregistrasi",
		"Hai "+user.Name+"\n"+
			"User anda telah berhasil didaftarkan dengan Login ID: "+user.LoginID)
	if e != nil {
		return "", errors.New("system error preparing registration welcome message: " + e.Error())
	}

	if re.opts.FnPostRegister != nil {
		re.opts.FnPostRegister(user, parm)
	}

	ev, e := ctx.DefaultEvent()
	if ev == nil {
		return "", errors.New("invalid eventhub")
	}

	reply := ""
	if err := ev.Publish(re.opts.SendNotifTopic, msg.ID, &reply); err != nil {
		return "", errors.New("system error when sending token: " + err.Error())
	}

	return "OK", nil
}
