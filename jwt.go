package kiam

import (
	"errors"
	"fmt"
	"reflect"

	"git.kanosolution.net/kano/kaos"
	"github.com/dgrijalva/jwt-go"
	"github.com/eaciit/toolkit"
)

var (
	secret    string
	claimType reflect.Type
)

func JWTMiddleware(secret string, ref interface{}) func(*kaos.Context) error {
	claimType := reflect.Indirect(reflect.ValueOf(ref)).Type()
	return func(ctx *kaos.Context) error {
		var claim jwt.Claims
		claim = reflect.New(claimType).Interface().(jwt.Claims)

		token := ctx.Data().Get("jwt-token", "").(string)
		if token == "" {
			return nil
		}

		t, e := jwt.ParseWithClaims(token, claim, func(tk *jwt.Token) (interface{}, error) {
			//return secret, nil
			return []byte(secret), nil
		})
		if e != nil {
			return fmt.Errorf("Token error: %s", e.Error())
		}
		if !t.Valid {
			return errors.New("Token is not valid")
		}

		var kc Claim
		kc = claim.(Claim)

		// get session from IAM service
		ev, e := ctx.DefaultEvent()
		if e != nil {
			return nil
			//return fmt.Errorf("unable to get event manager: %s", ev.Error())
		}

		// reference ID
		sess := new(Session)
		if e = ev.Publish("/v1/iam/get", toolkit.M{}.Set("id", kc.SessionID()), sess); e != nil {
			//return nil
			return fmt.Errorf("unable to get session: %s", e.Error())
		}
		ctx.Data().Set("jwt-sessionid", kc.SessionID())
		ctx.Data().Set("jwt-referenceid", sess.ReferenceID)
		return nil
	}
}
