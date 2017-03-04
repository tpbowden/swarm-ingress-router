package service

type Service struct {
	Name        string
	DNSNames    []string
	Port        int
	Certificate string
	Key         string
	Secure      bool
	ForceTLS    bool
}
