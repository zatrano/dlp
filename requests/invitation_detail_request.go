package requests

type CreateInvitationDetailRequest struct {
	Title              string `form:"title"`
	Person             string `form:"person"`
	
	// Person's Parents
	IsMotherLive       bool   `form:"is_mother_live"`
	MotherName         string `form:"mother_name"`
	MotherSurname      string `form:"mother_surname"`
	IsFatherLive       bool   `form:"is_father_live"`
	FatherName         string `form:"father_name"`
	FatherSurname      string `form:"father_surname"`

	// Bride Details
	BrideName          string `form:"bride_name"`
	BrideSurname       string `form:"bride_surname"`
	IsBrideMotherLive  bool   `form:"is_bride_mother_live"`
	BrideMotherName    string `form:"bride_mother_name"`
	BrideMotherSurname string `form:"bride_mother_surname"`
	IsBrideFatherLive  bool   `form:"is_bride_father_live"`
	BrideFatherName    string `form:"bride_father_name"`
	BrideFatherSurname string `form:"bride_father_surname"`

	// Groom Details
	GroomName          string `form:"groom_name"`
	GroomSurname       string `form:"groom_surname"`
	IsGroomMotherLive  bool   `form:"is_groom_mother_live"`
	GroomMotherName    string `form:"groom_mother_name"`
	GroomMotherSurname string `form:"groom_mother_surname"`
	IsGroomFatherLive  bool   `form:"is_groom_father_live"`
	GroomFatherName    string `form:"groom_father_name"`
	GroomFatherSurname string `form:"groom_father_surname"`
}

type UpdateInvitationDetailRequest struct {
	Title              string `form:"title"`
	Person             string `form:"person"`
	
	// Person's Parents
	IsMotherLive       bool   `form:"is_mother_live"`
	MotherName         string `form:"mother_name"`
	MotherSurname      string `form:"mother_surname"`
	IsFatherLive       bool   `form:"is_father_live"`
	FatherName         string `form:"father_name"`
	FatherSurname      string `form:"father_surname"`

	// Bride Details
	BrideName          string `form:"bride_name"`
	BrideSurname       string `form:"bride_surname"`
	IsBrideMotherLive  bool   `form:"is_bride_mother_live"`
	BrideMotherName    string `form:"bride_mother_name"`
	BrideMotherSurname string `form:"bride_mother_surname"`
	IsBrideFatherLive  bool   `form:"is_bride_father_live"`
	BrideFatherName    string `form:"bride_father_name"`
	BrideFatherSurname string `form:"bride_father_surname"`

	// Groom Details
	GroomName          string `form:"groom_name"`
	GroomSurname       string `form:"groom_surname"`
	IsGroomMotherLive  bool   `form:"is_groom_mother_live"`
	GroomMotherName    string `form:"groom_mother_name"`
	GroomMotherSurname string `form:"groom_mother_surname"`
	IsGroomFatherLive  bool   `form:"is_groom_father_live"`
	GroomFatherName    string `form:"groom_father_name"`
	GroomFatherSurname string `form:"groom_father_surname"`
}

func ValidateCreateInvitationDetailRequest(c *fiber.Ctx) error {
	var req CreateInvitationDetailRequest
	errorMessages := map[string]string{}

	if err := validateRequest(c, &req, errorMessages, "/panel/invitations"); err != nil {
		return err
	}

	c.Locals("createInvitationDetailRequest", req)
	return c.Next()
}

func ValidateUpdateInvitationDetailRequest(c *fiber.Ctx) error {
	var req UpdateInvitationDetailRequest
	errorMessages := map[string]string{}

	if err := validateRequest(c, &req, errorMessages, "/panel/invitations"); err != nil {
		return err
	}

	c.Locals("updateInvitationDetailRequest", req)
	return c.Next()
}
