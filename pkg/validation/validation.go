package validation

import (
	"cli-project/pkg/utils"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var validCountries = map[string]struct{}{
	"afghanistan": {}, "albania": {}, "algeria": {}, "andorra": {}, "angola": {}, "antiguaandbarbuda": {},
	"argentina": {}, "armenia": {}, "australia": {}, "austria": {}, "azerbaijan": {}, "bahamas": {},
	"bahrain": {}, "bangladesh": {}, "barbados": {}, "belarus": {}, "belgium": {}, "belize": {},
	"benin": {}, "bhutan": {}, "bolivia": {}, "bosniaandherzegovina": {}, "botswana": {}, "brazil": {},
	"brunei": {}, "bulgaria": {}, "burkinafaso": {}, "burundi": {}, "caboverde": {}, "cambodia": {},
	"cameroon": {}, "canada": {}, "centralafricanrepublic": {}, "chad": {}, "chile": {}, "china": {},
	"colombia": {}, "comoros": {}, "congodemocraticrepublicofthe": {}, "congorepublicofthe": {}, "costarica": {},
	"croatia": {}, "cuba": {}, "cyprus": {}, "czechrepublic": {}, "denmark": {}, "djibouti": {},
	"dominica": {}, "dominicanrepublic": {}, "easttimor": {}, "ecuador": {}, "egypt": {}, "elsalvador": {},
	"equatorialguinea": {}, "eritrea": {}, "estonia": {}, "eswatini": {}, "ethiopia": {}, "fiji": {},
	"finland": {}, "france": {}, "gabon": {}, "gambia": {}, "georgia": {}, "germany": {}, "ghana": {},
	"greece": {}, "grenada": {}, "guatemala": {}, "guinea": {}, "guineabissau": {}, "guyana": {}, "haiti": {},
	"honduras": {}, "hungary": {}, "iceland": {}, "india": {}, "indonesia": {}, "iran": {}, "iraq": {},
	"ireland": {}, "israel": {}, "italy": {}, "ivorycoast": {}, "jamaica": {}, "japan": {}, "jordan": {},
	"kazakhstan": {}, "kenya": {}, "kiribati": {}, "koreanorth": {}, "koreasouth": {}, "kosovo": {}, "kuwait": {},
	"kyrgyzstan": {}, "laos": {}, "latvia": {}, "lebanon": {}, "lesotho": {}, "liberia": {}, "libya": {},
	"liechtenstein": {}, "lithuania": {}, "luxembourg": {}, "madagascar": {}, "malawi": {}, "malaysia": {},
	"maldives": {}, "mali": {}, "malta": {}, "marshallislands": {}, "mauritania": {}, "mauritius": {},
	"mexico": {}, "micronesia": {}, "moldova": {}, "monaco": {}, "mongolia": {}, "montenegro": {},
	"morocco": {}, "mozambique": {}, "myanmar": {}, "namibia": {}, "nauru": {}, "nepal": {}, "netherlands": {},
	"newzealand": {}, "nicaragua": {}, "niger": {}, "nigeria": {}, "northmacedonia": {}, "norway": {}, "oman": {},
	"pakistan": {}, "palau": {}, "panama": {}, "papuanewguinea": {}, "paraguay": {}, "peru": {}, "philippines": {},
	"poland": {}, "portugal": {}, "qatar": {}, "romania": {}, "russia": {}, "rwanda": {}, "saintkittsandnevis": {},
	"saintlucia": {}, "saintvincentandthegrenadines": {}, "samoa": {}, "sanmarino": {}, "saotomeandprincipe": {},
	"saudiarabia": {}, "senegal": {}, "serbia": {}, "seychelles": {}, "sierraleone": {}, "singapore": {},
	"slovakia": {}, "slovenia": {}, "solomonislands": {}, "somalia": {}, "southafrica": {}, "southsudan": {},
	"spain": {}, "srilanka": {}, "sudan": {}, "suriname": {}, "sweden": {}, "switzerland": {}, "syria": {},
	"taiwan": {}, "tajikistan": {}, "tanzania": {}, "thailand": {}, "togo": {}, "tonga": {},
	"trinidadandtobago": {}, "tunisia": {}, "turkey": {}, "turkmenistan": {}, "tuvalu": {},
	"uganda": {}, "ukraine": {}, "unitedarabemirates": {}, "unitedkingdom": {}, "unitedstates": {}, "uruguay": {},
	"uzbekistan": {}, "vanuatu": {}, "vaticancity": {}, "venezuela": {}, "vietnam": {}, "yemen": {},
	"zambia": {}, "zimbabwe": {},
}

func ValidateCountryName(country string) (bool, error) {
	country = strings.ToLower(country)
	country = strings.ReplaceAll(country, " ", "")
	if _, exists := validCountries[country]; !exists {
		return false, errors.New("invalid country name")
	}
	return true, nil
}

func ValidateEmail(email string) (bool, bool) {

	// Extract the domain from the email
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, false
	}

	// Define a list of reputable email domains
	reputableDomains := []string{"gmail.com", "outlook.com", "yahoo.com", "watchguard.com", "hotmail.com", "icloud.com"}

	// Regular expression to match a valid email format
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(emailRegex, email)
	if !match {
		return false, false
	}

	domain := parts[1]

	// Check if the domain is in the list of reputable domains
	for _, reputableDomain := range reputableDomains {
		if domain == reputableDomain {
			return true, true
		}
	}
	return true, false
}

// ValidateName checks if the name is valid
func ValidateName(name string) bool {
	// Name must be non-empty and only contain alphabetic characters and spaces
	if len(name) <= 2 || len(name) >= 31 {
		return false
	}

	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == ' ') {
			return false
		}
	}
	return true
}

func ValidateOrganizationName(orgName string) (bool, error) {

	if len(orgName) <= 1 || len(orgName) > 40 {
		return false, errors.New("invalid organization name : name must be between 2 and 40 characters")
	}
	// Regex to allow only letters and spaces
	const orgNameRegex = `^[a-zA-Z\s]+$`
	match, _ := regexp.MatchString(orgNameRegex, orgName)

	if !match {
		return false, errors.New("invalid organization name : only letters and spaces are allowed")
	}

	return true, nil
}

const minPasswordLength = 8

func ValidatePassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}

	hasUpper := false
	for _, r := range password {
		if unicode.IsUpper(r) {
			hasUpper = true
			break
		}
	}
	if !hasUpper {
		return false
	}

	hasLower := false
	for _, r := range password {
		if unicode.IsLower(r) {
			hasLower = true
			break
		}
	}
	if !hasLower {
		return false
	}

	hasDigit := false
	for _, r := range password {
		if unicode.IsDigit(r) {
			hasDigit = true
			break
		}
	}
	if !hasDigit {
		return false
	}

	hasSpecial := false
	specialChars := []byte("!@#$%^&*()-+?_=,<>/{}[]|`~;")

	// Convert the rune to a byte for comparison
	for _, r := range password {
		if unicode.IsPunct(rune(r)) && specialChars[0] != byte(r) && specialChars[len(specialChars)-1] != byte(r) {
			hasSpecial = true
			break
		}
	}
	if !hasSpecial {
		return false
	}

	return true
}

func ValidateTitleSlug(titleSlug string) (bool, error) {
	// Add validation logic for titleSlug, such as format checks
	titleSlug = strings.TrimSpace(titleSlug)
	if len(titleSlug) == 0 {
		return false, fmt.Errorf("title slug cannot be empty")
	}
	// Add more validations as needed
	return true, nil
}

// ValidateUsername checks if the username is valid
func ValidateUsername(username string) bool {

	if len(username) <= 3 || len(username) >= 21 {
		return false
	}

	hasLetter := false
	hasDigitAfterLetter := false

	for _, r := range username {
		if unicode.IsLetter(r) {
			hasLetter = true
			// Digits after a letter are allowed
			hasDigitAfterLetter = true
		} else if unicode.IsDigit(r) {
			if hasLetter {
				hasDigitAfterLetter = true
			} else {
				// Digits are not allowed if there has been no letter before
				return false
			}
		} else {
			// Invalid character found
			return false
		}
	}

	return hasLetter && hasDigitAfterLetter
}

func ValidateQuestionDifficulty(difficulty string) (string, error) {
	lowerDifficulty := utils.CleanString(difficulty)
	validDifficulties := map[string]bool{"easy": true, "medium": true, "hard": true}

	if !validDifficulties[lowerDifficulty] {
		return "", errors.New("invalid difficulty level: must be 'easy', 'medium', or 'hard'")
	}
	return lowerDifficulty, nil
}

func ValidateQuestionLink(link string) (string, error) {
	lowerLink := utils.CleanString(link)
	parsedURL, err := url.Parse(lowerLink)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" || !strings.Contains(parsedURL.Host, "leetcode.com") {
		return "", errors.New("invalid question link: must be a valid Leetcode link")
	}
	return lowerLink, nil
}

func ValidateQuestionID(questionID string) (bool, error) {
	qid, err := strconv.Atoi(questionID)
	if err != nil || qid <= 0 {
		return false, errors.New("invalid question ID : must be a positive number")
	}
	return true, nil
}
