package cache

import "blog/help"

var (
	maliciousIpMap = make(map[string]bool)
)

func IfIpInMaliciousMap(ip string) bool {
	return maliciousIpMap[ip]
}

func AddMaliciousIpMap(ip string) {
	maliciousIpMap[ip] = true
}

func getMaliciousMap() error {
	if err := help.GetMaliciousIp(maliciousIpMap); err != nil {
		return err
	}
	return nil
}
