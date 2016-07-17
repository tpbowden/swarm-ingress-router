package service

import (
	"testing"

	"github.com/docker/engine-api/types/swarm"
)

var key = `
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAqc7E3vRddKrOgsssOVr51/43cBmwWm3I0et0Xhwl3IALA5eD
f2mzEK8/Q8znr4mKqnuFVK57XHwv8gC0oXVHmGqZ2E9k7Qz1hjWE+yaGeZdmvu21
QtN3m554YwSAyTUPwKXn3mLu8ZXutrs8HW9UU7M4/3eumKLLAeXoIuc4KzNCI7Ax
9EZYnrxDgsaMyGR7rtvMJUFUqeK3/uc/VrMMWeINOnisIsNyXP4kXn2JcZvihAkB
38fjG5I8aQS8OZSI+foz+HZa4t4Cf6x+lIVQY4ClEACVR4/jIikStVNS25qT3Cli
CZocU6/M6aqT/PSSjb7xCG7LRxniLok8fAkopwIDAQABAoIBACBCGYBJZ9+b+zM5
qaMTt1bjAmEXYDstk/LdCfQ+4Kps4KrEA8UWrV/ECv60VGcyd9c/M+sMYn8FUI5m
t+ncP8eiqKtBGek2hwYB8MtNXEqJConG2qJhTd2QU0psNpMCs4VjfxLuNHpMg1bP
ZvIojhqAd5mAgPiaJHqt1FWdGzyJZxFE651xHFTW7NOyUuuOR025jGMI5QnSCkJq
GA9VoJ1yXdG2U7EHmY5+dEJBglDsb5JJGthMBtY+o0goT/HPUqAfV86uGzffgWUq
kQPO9YPwwClwOr17fPE+8Qw7VIeGAySNewYXjmjtfgoyDazuAd6KzorRjplvPHBl
HvZMMAECgYEA1NpVX0vdTgNP6l5VUvTbXlUas/9mZ6XA3Z1ngma0d0rlkhzw1LgS
Ww54qJwJoR5C4gjzgK1hqCPF1CRv0aVFQ0j/dfCHRfqytcKUdCPldMZa4W1xHnVa
jAQqFaEi6Cj+427S/eMn1NDUEWL1VvFs8WVrHnpsFddTDAb8oL6ZTKcCgYEAzDqs
b9QHSwzcIacSNzJimr03h9R8DNTwtSdh/sogDwL/tCdGLYtxAcvzb1NMSvlvHd62
+jKLwPIs2p5n7I0BKvjIkmW9nGmdOstn8SQCkhCUGWAta2+BP1UF5DctpVLb8/uo
EruCwouzvVeTzpB96F+1Tm3uckKLC0uASMNAxAECgYEAvZG2dykZ8GECy7k4RKnO
tjUVkznj+mulWbrWdU0DbTtHOtqLouhNcMtyqrhN2zEYYDeYpwHD9/vkNQw+inin
N0XMPz35PFoKz9Z8YPOXaGlAh4TxOi9KdWlAEgNxE1Nvrx8EyxmEYAWc2d9IoiZi
4Jtyy7I8kTc0v4F5fbBC2AECgYALUKjjHT02NEUx/B6vPjRmXFtqRCSHVXjsoHz2
b95s1n6yTYa+2T3umo0nOtc8RCua3Q8IN6q0ivZfOm2Jlppc9iGussJZmyRh8IkW
vCcETrTV3xVFIY1oo95KsZ/uy/NxxhyexRLOkoznzaVbyXegW0UhTkfqvrMTciBu
Z5r8AQKBgH0+HKWvjWt6Adk62FSLgTVfNOOyETSXdso1S0fh6/JRlnvUhBhgYB59
zxkL4jmqp1YtSBUrbPXlsWPiV1CLU0DUopDq9bebXT4q7lETek5IWAO467hV4S2v
82ypi+3PZhgFh9hzE0nPcxkPCxtX2E8pwzCfhbpvtdlZNIjQge0S
-----END RSA PRIVATE KEY-----
`

var certificate = `
-----BEGIN CERTIFICATE-----
MIIE5TCCAs2gAwIBAgIBADANBgkqhkiG9w0BAQsFADCBszELMAkGA1UEBhMCR0Ix
CzAJBgNVBAgTAlVLMR8wHQYDVQQHExZBZGFzdHJhbCBQYXJrLCBJcHN3aWNoMScw
JQYDVQQKEx5Ccml0aXNoIFRlbGVjb21tdW5pY2F0aW9ucyBwbGMxGDAWBgNVBAsT
D0lQIEFwcGxpY2F0aW9uczESMBAGA1UEAxMJRGl2ZWJvYXJkMR8wHQYJKoZIhvcN
AQkBFhBkaXZlYm9hcmRAYnQuY29tMB4XDTE2MDcxNTE1NDAwOFoXDTE5MDcxNTE1
NDAwOFowgbcxCzAJBgNVBAYTAkdCMQswCQYDVQQIEwJVSzEfMB0GA1UEBxMWQWRh
c3RyYWwgUGFyaywgSXBzd2ljaDEnMCUGA1UECgweQnJpdGlzaCBUZWxlY29tbXVu
aWNhdGlvbnMgcGxjMRgwFgYDVQQLDA9JUCBBcHBsaWNhdGlvbnMxFjAUBgNVBAMM
DWV4YW1wbGUubG9jYWwxHzAdBgkqhkiG9w0BCQEMEGRpdmVib2FyZEBidC5jb20w
ggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCpzsTe9F10qs6Cyyw5WvnX
/jdwGbBabcjR63ReHCXcgAsDl4N/abMQrz9DzOeviYqqe4VUrntcfC/yALShdUeY
apnYT2TtDPWGNYT7JoZ5l2a+7bVC03ebnnhjBIDJNQ/ApefeYu7xle62uzwdb1RT
szj/d66YossB5egi5zgrM0IjsDH0RlievEOCxozIZHuu28wlQVSp4rf+5z9WswxZ
4g06eKwiw3Jc/iRefYlxm+KECQHfx+MbkjxpBLw5lIj5+jP4dlri3gJ/rH6UhVBj
gKUQAJVHj+MiKRK1U1LbmpPcKWIJmhxTr8zpqpP89JKNvvEIbstHGeIuiTx8CSin
AgMBAAEwDQYJKoZIhvcNAQELBQADggIBADrJ1/0VxXeslbBrFvAYxAuEJx/NTpcc
ImqIyhJ9Kg7C1xHYX6C8wR2Lw16aidcpWnuQZ86EVCs6onb5HJVdxP34Fmd9VIcQ
NwjfGuaUAM+nyd/bfLjgJymQ2ehEF9RikfjQMGc9fvtaP0qFBDUrof5Z7922NOHM
8Nsa3hInpVUNz2a9ZyXbf/liqwmDy2GXbc0ed9+4MT+bgOcAUI4XZ79WUMv9uHcx
VlCYDzyPZQh93u8w7q5TWJhfAEeOUzr6KBSyIhMNg2jTXjFBl9KZFu0mzxy1OdnW
kVH7eqnfQwkObS1NxfqpSOkf5bz2InxNybeMh/+x3i5WNNDD3mn3W0SjUQsSFBoh
vPi+YYqM01Y2L+MNSI4nazZ9OoczNfdMDiT2dldYjDvPYYatxbNzDQaXH5Kbdf51
p85qwrSoIEkPveAEsyiDIQJMxWC6xkjxDnu26qO/j25C2kae6fXKc35dO8zve4QI
1xZQ4le3H1OzVv3/foiNZZitRMBgUrlOkvZAkR1xiwhFSH1NdurKeKAvZ04y8twA
XJnHVxmqBp7/o69UoOtAwUMGbks7vAeG21KwUF2Fx6d4gdQjR1EDNDSYNYeV+g/g
gDmrXsJilmWdVCehduHoYorwcICbSr1TcERMrlTUXfW+wuMdzF1WUoScYatvLXaX
sKi+bSEUZYKb
-----END CERTIFICATE-----`

type FakeClient struct {
	services []swarm.Service
}

func (f FakeClient) GetServices(filters map[string]string) []swarm.Service {
	return f.services
}

func TestLoadingServices(t *testing.T) {
	labels := map[string]string{
		"ingress.targetport": "8080",
		"ingress.dnsname":    "example.com",
	}

	fakeService := swarm.Service{
		ID: "123",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "myservice", Labels: labels},
		},
	}

	result := LoadAll(FakeClient{services: []swarm.Service{fakeService}})

	parsedService := result[0]

	expectedName := "example.com"
	actualName := parsedService.DNSName()

	if expectedName != actualName {
		t.Errorf("Expected DNS name of %s, got %s", expectedName, actualName)
	}

	expectedURL := "http://myservice:8080"
	actualURL := parsedService.URL()

	if expectedURL != actualURL {
		t.Errorf("Expected URL of %s, got %s", expectedURL, actualURL)
	}

	if _, ok := parsedService.Certificate(); ok {
		t.Error("Expected the insecure service not to have a certificate")
	}
}

func TestLoadingTLSServices(t *testing.T) {

	labels := map[string]string{
		"ingress.targetport": "8443",
		"ingress.dnsname":    "example.local",
		"ingress.tls":        "true",
		"ingress.cert":       certificate,
		"ingress.key":        key,
	}

	fakeService := swarm.Service{
		ID: "654",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "tlsservice", Labels: labels},
		},
	}

	result := LoadAll(FakeClient{services: []swarm.Service{fakeService}})

	parsedService := result[0]

	expectedName := "example.local"
	actualName := parsedService.DNSName()

	if expectedName != actualName {
		t.Errorf("Expected DNS name of %s, got %s", expectedName, actualName)
	}

	expectedURL := "http://tlsservice:8443"
	actualURL := parsedService.URL()

	if expectedURL != actualURL {
		t.Errorf("Expected URL of %s, got %s", expectedURL, actualURL)
	}

	if _, ok := parsedService.Certificate(); !ok {
		t.Error("Expected the TLS service to have a certificate")
	}
}

func TestLoadingWithoutCertOrKey(t *testing.T) {
	noCertLabels := map[string]string{
		"ingress":            "true",
		"ingress.targetport": "200",
		"ingress.dnsname":    "example.local",
		"ingress.tls":        "true",
		"ingress.key":        key,
	}

	ignoredService1 := swarm.Service{
		ID: "567",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "ignored", Labels: noCertLabels},
		},
	}

	noKeyLabels := map[string]string{
		"ingress":            "true",
		"ingress.targetport": "200",
		"ingress.dnsname":    "example.local",
		"ingress.tls":        "true",
		"ingress.cert":       certificate,
	}

	ignoredService2 := swarm.Service{
		ID: "567",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "ignored", Labels: noKeyLabels},
		},
	}

	result := LoadAll(FakeClient{services: []swarm.Service{
		ignoredService1,
		ignoredService2}})
	if len(result) != 0 {
		t.Errorf("Expected no services to be created, got %d", len(result))
	}
}

func TestLoadingInvalidService(t *testing.T) {
	labels := map[string]string{
		"ingress.targetport": "abc",
		"ingress.dnsname":    "example.local",
	}

	ignoredService := swarm.Service{
		ID: "567",
		Spec: swarm.ServiceSpec{
			Annotations: swarm.Annotations{Name: "ignored", Labels: labels},
		},
	}

	result := LoadAll(FakeClient{services: []swarm.Service{ignoredService}})
	if len(result) != 0 {
		t.Errorf("Expected no services to be created, got %d", len(result))
	}
}
