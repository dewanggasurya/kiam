module github.com/dewanggasurya/kiam

go 1.16

replace github.com/ariefdarmawan/kiam => ./

require (
	git.kanosolution.net/kano/appkit v0.0.1
	git.kanosolution.net/kano/dbflex v1.0.15
	git.kanosolution.net/kano/kaos v0.1.1
	github.com/ariefdarmawan/datahub v0.2.0
	github.com/ariefdarmawan/kiam v0.0.4
	github.com/ariefdarmawan/kmsg v0.0.2
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eaciit/toolkit v0.0.0-20210610161449-593d5fadf78e
	github.com/google/uuid v1.3.0
	go.mongodb.org/mongo-driver v1.7.2
)
