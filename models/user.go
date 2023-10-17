package models

type User struct {
	Id       int64  `json:"id"`
	IIN      string `json:"IIN"`
	Name     string `json:"name" `
	Password string `json:"password"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}

// // Entity
// type Entity struct {
// 	IIN          string    `json:"IIN"`
// 	Name         string    `json:"name"`
// 	Type         int       `json:"type"`
// 	Mail         string    `json:"mail"`
// 	Info         string    `json:"info"`
// 	PlatonusID   *int64    `json:"platonusId"`
// 	Photo        string    `json:"photo"`
// 	Position     string    `json:"position"`
// 	Postal       *string   `json:"postal"`
// 	Course       string    `json:"course"`
// 	Degree       string    `json:"degree"`
// 	Faculty      string    `json:"faculty"`
// 	Group        string    `json:"group"`
// 	FullName     string    `json:"fullName"`
// 	UserID       int64     `json:"userID"`
// 	CreatedDate  time.Time `json:"createdDate"`
// 	ModiDate     time.Time `json:"modiDate"`
// 	WhenChanged  string    `json:"whenChanged"`
// 	PhoneNumber  *string   `json:"phoneNumber"`
// 	FirstName    string    `json:"firstName"` // есімі
// 	LastName     string    `json:"lastName"`  // тегі - фамилия
// 	ThirdName    string    `json:"thirdName"` // әкесі - отчество
// 	FirstnameEn  *string   `json:"firstnameEn"`
// 	MiddlenameEn *string   `json:"lastnameEn"`
// 	LastnameEn   *string   `json:"thirdnameEn"`
// 	Email        *string   `json:"email"`
// 	Address      *string   `json:"address"`
// 	Resident     *int      `json:"resident"`
// 	LDAP         int       `json:"fromLDAP"`
// 	//Roles           []Role            `json:"roles"`
// 	// Positions       []Position        `json:"positions"`
// 	// AcademicDegree  DictionaryElement `json:"academicDegree"`
// 	// AcademicTitle   DictionaryElement `json:"academicTitle"`
// 	// Organization    Organization      `json:"organization"`
// 	// Department      Department        `json:"department"`
// 	// MainPosition    Position          `json:"mainPosition"`
// 	BirthDay *time.Time `json:"birthday"`
// 	Gender   *int       `json:"gender"`
// 	//Locality        Locality          `json:"locality"`
// 	AddressRus  *string `json:"addressrus"`
// 	BankAccount *string `json:"bankaccount"`
// 	//Bank            Bank              `json:"bank"`
// 	State           *int       `json:"state"`
// 	IDNumber        *string    `json:"idnumber"`
// 	IDDate          *time.Time `json:"iddate"`
// 	IDGivenOrg      *int       `json:"idsource"`
// 	IDExpire        *time.Time `json:"idexpire"`
// 	Password        *string    `json:"password"`
// 	HasAlreadyReg   bool       `json:"hasAlreadyReg"`
// 	SertType        *string    `json:"sertType"`
// 	HasUpdated      bool       `json:"hasUpdated"`
// 	HomePhone       *string
// 	CountryCode     *int
// 	BuildingKZ      *string
// 	BuildingRU      *string
// 	BuildingEN      *string
// 	FirstNameEN     *string
// 	LastNameEN      *string
// 	DepartmentKZ    *string
// 	DepartmentRU    *string
// 	DepartmentEN    *string
// 	PostNameKZ      *string
// 	PostNameRU      *string
// 	PostNameEN      *string
// 	PostNameTutorKZ *string
// 	PostNameTutorRU *string
// 	PostNameTutorEN *string
// 	UnitNameTutorKZ *string
// 	UnitNameTutorRU *string
// 	UnitNameTutorEN *string
// 	DoesExternal    int  `json:"doesExternal"`
// 	IsContragent    bool `json:"isContragent"`
// }
