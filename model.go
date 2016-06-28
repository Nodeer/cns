package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

type Auth struct {
	Id       string `json:"id"`
	Email    string `json:"email" auth:"username"`
	Password string `json:"password" auth:"password"`
	Active   bool   `json:"active" auth:"active"`
	Role     string `json:"role"`
}

type Employee struct {
	Auth
	FirstName string `json:"firstName,omitempty"`
	LastName  string `json:"lastName,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Home      string `json:"home,omitempty"`
	Address
}

type BusinessType int

const (
	SOLE_PROPRIETOR BusinessType = iota
	CORPORATION
	PARTNERSHIP
	LLC
	LLP
	BUSINESS_OTHER
)

type CarrierType int

const (
	PRIVATE CarrierType = iota
	COMMON
	CONTRACT
	CARRIER_OTHER
)

type Company struct {
	Id                string       `json:"id"`
	DOTNum            string       `json:"dotNum,omitempty"`
	Name              string       `json:"name,omitempty"`
	ContactName       string       `json:"contactName,omitempty"`
	DBA               string       `json:"dba,omitempty"`
	ContactTitle      string       `json:"contactTitle,omitempty"`
	ContactSSN        string       `jsni:"contactSSN,omitempty"`
	ContactPhone      string       `json:"contactPhone,omitempty"`
	SecondName        string       `json:"secondName,omitempty"`
	SecondTitle       string       `json:"secondTitle,omitempty"`
	SecondPhone       string       `json:"secondPhone,omitempty"`
	SameAddress       bool         `json:"sameAddress"`
	PhysicalAddress   Address      `json:"pysicalAddress,omitempty"`
	MailingAddress    Address      `json:"mailingAddress,omitempty"`
	BusinessType      BusinessType `json:"businessType,omitempty"`
	BusinessTypeOther string       `json:"businessTypeOther,omitempty"`
	MCNum             string       `json:"mcNum,omitempty"`
	PUCNum            string       `json:"pucNum,omitempty"`
	Fax               string       `json:"fax,omitempty"`
	Email             string       `json:"email,omitempty"`
	EINNum            string       `json:"einNum,omitempty"`
	ARPAccountNum     string       `json:"arpAccountNum,omitempty"`
	CarrierType       CarrierType  `json:"carrierType,omitempty"`
	CarrierTypeOther  string       `json:"carrierTypeOther,omitempty"`
	EntityNum         string       `jaon:"entityNum,omitempty"`
	CreditCard        CreditCard
	NYHutUsername     string `json:"nyHutUsername,omitempty"`
	NYHutPassword     string `json:"nyHutPassword,omitempty"`
	NYOscarUsername   string `json:"nyOrcarUsername,omitempty"`
	NYOscarPassword   string `json:"nyOscarUsername,omitempty"`
	KYUseNum          string `json:"kyUseNum,omitempty"`
	NMHutUsername     string `json:"nmHutUsername,omitempty"`
	NMHutPassword     string `json:"nmHutPassword,omitempty"`
	DOTPin            string `json:"dotPin,omitempty"`
	MCPin             string `json:"mcPin,omitempty"`
	FMCSAUsername     string `json:"fmcsaUsername,omitempty"`
	FMCSAPassword     string `json:"fmcsaPassword,omitempty"`
	IRPNum            string `json:"irpNum,omitempty"`
	//Slug            string  `json:"slug,omitempty"`
}

func (c Company) GetBusinessType() string {
	switch c.BusinessType {
	case SOLE_PROPRIETOR:
		return "Sole Proprietor"
	case CORPORATION:
		return "Corporation"
	case PARTNERSHIP:
		return "Partnership"
	case LLC:
		return "LLC"
	case LLP:
		return "LLP"
	case BUSINESS_OTHER:
		return c.BusinessTypeOther
	}
	return ""
}

func (c Company) GetCarrierType() string {
	switch c.CarrierType {
	case PRIVATE:
		return "Private"
	case COMMON:
		return "Common"
	case CONTRACT:
		return "Contract"
	case CARRIER_OTHER:
		return c.CarrierTypeOther
	}
	return ""
}

func GetCompanyConsts() map[string]interface{} {
	m := map[string]interface{}{
		"SOLE_PROPRIETOR": SOLE_PROPRIETOR,
		"CORPORATION":     CORPORATION,
		"PARTNERSHIP":     PARTNERSHIP,
		"LLC":             LLC,
		"LLP":             LLP,
		"BUSINESS_OTHER":  BUSINESS_OTHER,
		"PRIVATE":         PRIVATE,
		"COMMON":          COMMON,
		"CONTRACT":        CONTRACT,
		"CARRIER_OTHER":   CARRIER_OTHER,
	}

	return m
}

/*func (c *Company) CreateSlug() {
	// slug = title.replaceAll("[;/?:@&=+\\\$,\\{\\}\\|\\\\^\\[\\]`]", "").trim().replace(' ', '_').toLowerCase()
	r, err := regexp.Compile("[;/?:@&=+$,\\{\\}\\|^\\[\\]`]")
	if err != nil {
		log.Printf("model.go -> Company.creaeSlug() -> regexp.Compile() -> %v\n", err)
		return
	}
	c.Slug = strings.ToLower(strings.Replace(strings.Trim(r.ReplaceAllString(c.Name, ""), " "), " ", "-", -1))
}*/

type Driver struct {
	Auth
	Address
	FirstName             string `json:"firstName,omitempty"`
	LastName              string `json:"lastName,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	EmergencyContactName  string `json:"emergencyContactName,omitempty"`
	EmergencyContactPhone string `json:"emergencyContactPhone,omitempty"`
	LicenseNum            string `json:"licenseNum,omitempty"`
	LicenseState          string `json:"licenseState,omitempty"`
	LicenseExpire         string `json:"licenseExpire,omitempty"`
	DOB                   string `json:"dob,omitempty"`
	MedCardExpiry         string `json:"medCardExpiry,omitempty"`
	MVRExpiry             string `json:"mVRExpiry,omitempty"`
	ReviewExpiry          string `json:"reviewExpiry,omitempty"`
	OneEightyExpiry       string `json:"oneEightyExpiry,omitempty"`
	HireDate              string `json:"hireDate,omitempty"`
	TermDate              string `json:"termDate,omitempty"`
	CompanyId             string `json:"companyId,omitempty"`
}

func (d Driver) FormatDOB() string {
	ds := strings.Split(d.DOB, "-")
	if len(ds) != 3 {
		return ""
	}
	if ds[1][0] == '0' {
		ds[1] = ds[1][1:]
	}
	return fmt.Sprintf("%s/%s/%s", ds[1], ds[2], ds[0])
}

func (d Driver) GetAge() int32 {
	dobT, err := time.Parse("01/02/2006", d.DOB)
	if err != nil {
		return 0
	}
	dob := dobT.UnixNano()
	diff := time.Now().UnixNano() - dob
	return int32(math.Floor((float64(diff) / float64(1000) / float64(1000) / float64(1000) / float64(60) / float64(60) / float64(24) / float64(365.25))))
}

type Address struct {
	Street string `json:"street,omitempty"`
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
	Zip    string `json:"zip,omitempty"`
	County string `json:"county,omitempty"`
}

func (a Address) AddrHTML() string {
	return a.Street + "<br>" + a.City + ", " + a.County + ", " + a.State + " " + a.Zip
}

type Document struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	DocumentId string `json:"documentId,omitempty"`
	Complete   bool   `json:"complete"`
	Data       string `json:"data,omitempty"`
	CompanyId  string `json:"companyId,omitempty"`
	DriverId   string `json:"driverId,omitempty"`
}

var DQFS = [][]string{
	[]string{"100", "Driver's Application"},
	[]string{"180", "Certification of Violations"},
	[]string{"200", "Annual Inquery & Review"},
	[]string{"250", "Road Test Certication"},
	[]string{"300", "Previous Driver Inquires"},
	[]string{"400", "Drug & Alcohol Records Request"},
	[]string{"450", "Drug & Alcohol Certified Receipt"},
	[]string{"500", "Certification Compliance"},
	[]string{"600", "Confictions for a Driver Violation"},
	[]string{"700", "New Hire Stmt On Duty Hours"},
	[]string{"750", "Other Ompensated Work"},
	[]string{"775", "Fair Credit Reporting Act"},
}

var CompanyForms = [][]string{
	[]string{"MV 550", ""},
	[]string{"MV 550A", ""},
	[]string{"MV 551", ""},
	[]string{"MV-552A", ""},
	[]string{"MV 558", ""},
	[]string{"MV41", ""},
	[]string{"TMT-39", ""},
	[]string{"PUC App", ""},
	[]string{"MCS 150", ""},
}

type Event struct {
	Id        string    `json:"id,omitempty"`
	Name      string    `json:"name"`
	Title     string    `json:"title,omitempty"`
	AllDay    bool      `json:"allDay,omitempty"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end,omitempty"`
	URL       string    `json:"url,omitempty"`
	ClassName string    `json:"className,omitempty"`
	Editable  bool      `json:"editable"`
	Overlap   bool      `json:"overlap,omitempty"`
}

type Comment struct {
	Id     string `json:"id"`
	Body   string `json:"body"`
	Url    string `json:"url"`
	Page   string `json:"page"`
	Closed bool   `json:"closed"`
}

type BodyType int

const (
	TT BodyType = iota
	TK
	TRL
	BUS
	SW
	BODY_OTHER
)

type Vehicle struct {
	Id            string   `json:"id"`
	CompanyId     string   `json:"companyId,omitempty"`
	VehicleType   string   `json:"vehicleType,omitempty"`
	UnitNumber    string   `json:"unitNumber,omitempty"`
	Make          string   `json:"make,omitempty"`
	VIN           string   `json:"vin,omitempty"`
	Title         string   `json:"title,omitempty"`
	GVW           string   `json:"gvw,omitempty"`
	GCR           string   `json:"gcr,omitempty"`
	UnladenWeight string   `json:"unladenWeight,omitempty"`
	PurchasePrice float32  `json:"purchacePrice,omitempty"`
	PurchaseDate  string   `json:"purchaseDate,omitempty"`
	CurrentValue  float32  `json:"currentValue,omitempty"`
	AxleAmmount   int      `json:"axleAmmount,omitempty"`
	FuelType      string   `json:"fuelType,omitempty"`
	Active        bool     `json:"active"`
	Owner         string   `json:"owner,omitempty"`
	Year          string   `json:"year,omitempty"`
	PlateNum      string   `json:"plateNum,omitempty"`
	PlateExpire   string   `json:"plateExpire,omitempty"`
	BodyType      BodyType `json:"bodyType,omitempty"`
	BodyTypeOther string   `json:"bodyTypeOther,omitempty"`
}

func (v Vehicle) GetBodyType() string {
	switch v.BodyType {
	case TT:
		return "TT"
	case TK:
		return "TK"
	case TRL:
		return "TRL"
	case BUS:
		return "Bus"
	case BODY_OTHER:
		return v.BodyTypeOther
	}
	return ""
}

func GetVehicleConsts() map[string]interface{} {
	m := map[string]interface{}{
		"TT":         TT,
		"TK":         TK,
		"TRL":        TRL,
		"BUS":        BUS,
		"SW":         SW,
		"BODY_OTHER": BODY_OTHER,
	}
	return m
}

type Note struct {
	Id              string `json:"id,omitempty"`
	CompanyId       string `json:"companyId,omitempty"`
	EmployeeId      string `json:"employeeId,omitempty"`
	Communication   string `json:"communication,omitempty"`
	Purpose         string `json:"purpose,omitempty"`
	StartTime       int64  `json:"startTime,omitempty"`
	StartTimePretty string `json:"startTimePretty,omitempty"`
	EndTime         int64  `json:"endTime,omitempty"`
	EndTimePretty   string `json:"endTimePretty,omitempty"`
	Representative  string `json:"representative,omitempty"`
	CallBack        string `json:"callBack,omitempty"`
	EmailEmployee   bool   `json:"emailEmployee,omitempty"`
	Billable        bool   `json:"billable,omitempty"`
	Body            string `json:"body,omitempty"`
}

type NoteRevSort []Note

func (n NoteRevSort) Len() int {
	return len(n)
}

func (n NoteRevSort) Less(i, j int) bool {
	return n[i].StartTime > n[j].StartTime
}

func (n NoteRevSort) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

type QuickNote struct {
	Name string
	Body string
}

type CreditCard struct {
	Number         string `json:"number,omitempty"`
	ExpirationDate string `json:"expirationDate,omitempty"`
	SecurityCode   string `json:"securityCode,omitempty"`
}
