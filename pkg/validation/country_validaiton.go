package validation

import (
	"errors"
	"strings"
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
	country = strings.ReplaceAll(country, " ", "")
	if _, exists := validCountries[country]; !exists {
		return false, errors.New("invalid country name")
	}
	return true, nil
}
