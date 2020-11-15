package kiamui

type registerEngine struct {
	Name           string `kf-pos:"1,1" label:"Nama lengkap" required:"true"`
	LoginID        string `kf-pos:"2,1" label:"Login ID" required:"true"`
	Email          string `kf-pos:"3,1" label:"Email" required:"true"`
	Phone          string `kf-pos:"3,2" label:"Mobile phone" required:"true"`
	Password       string `kf-pos:"4,1" kf-control:"password" required:"true"`
	ConfimPassword string `kf-pos:"4,2" kf-control:"password" label:"Confirm password" required:"true"`
}

func NewRegisterEngine() *registerEngine {
	r := new(registerEngine)
	return r
}
