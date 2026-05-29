package domain

// MergeUserWithoutAuditDetail ports OwnerInfo.addUserWithoutAuditDetail.
func (o *OwnerInfo) MergeUserWithoutAuditDetail(from OwnerInfo) {
	o.UUID = from.UUID
	o.ID = from.ID
	o.UserName = from.UserName
	o.Password = from.Password
	o.Salutation = from.Salutation
	o.Name = from.Name
	o.Gender = from.Gender
	o.MobileNumber = from.MobileNumber
	o.EmailID = from.EmailID
	o.AltContactNumber = from.AltContactNumber
	o.Pan = from.Pan
	o.AadhaarNumber = from.AadhaarNumber
	o.PermanentAddress = from.PermanentAddress
	o.PermanentCity = from.PermanentCity
	o.PermanentPinCode = from.PermanentPinCode
	o.CorrespondenceAddress = from.CorrespondenceAddress
	o.CorrespondenceCity = from.CorrespondenceCity
	o.CorrespondencePinCode = from.CorrespondencePinCode
	o.Active = from.Active
	o.Dob = from.Dob
	o.PwdExpiryDate = from.PwdExpiryDate
	o.Locale = from.Locale
	o.Type = from.Type
	o.AccountLocked = from.AccountLocked
	o.Roles = from.Roles
	o.FatherOrHusbandName = from.FatherOrHusbandName
	o.BloodGroup = from.BloodGroup
	o.IdentificationMark = from.IdentificationMark
	o.Photo = from.Photo
	o.TenantID = from.TenantID
}

// CompareWithExistingUser ports OwnerInfo.compareWithExistingUser (emailId parity).
func (o OwnerInfo) CompareWithExistingUser(other OwnerInfo) bool {
	return o.Name == other.Name &&
		o.MobileNumber == other.MobileNumber &&
		o.Gender == other.Gender &&
		o.EmailID == other.EmailID &&
		o.FatherOrHusbandName == other.FatherOrHusbandName &&
		o.CorrespondenceAddress == other.CorrespondenceAddress
}
