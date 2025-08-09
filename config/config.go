package config

import (
	"os"
	"strings"

)

// to get the map of allowed domains, key is string and value type is empty struct
func GetAllowedDomains() map[string]struct{} {
	// to get the domainStr from the .env, domains separated by commas
	domainStr := os.Getenv("ALLOWED_DOMAINS")
	//splitting domainStr to get all the domains
	domains := strings.Split(domainStr, ",")
	//creating maps to store all allowed domains of length domains
	allowed := make(map[string]struct{}, len(domains))


	for _, d := range domains{
		//to remove the whitespaces from the domains
		d = strings.ToLower(strings.TrimSpace(d))
		if d != "" {
			// struct{}{} is used as a placeholder since we only care about keys, that is domains
			// hindi mae bolu to, bas domain exist krta hai, mtlb sirf key exist krti hai
			// struct isliye use kiya kyonki ye 0 byte of storage leta hai
			allowed[d] = struct{}{}
		}
	}

	return allowed

}




