package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/util/cert"
	"net/http"
)

var (
	caPEM = `
-----BEGIN CERTIFICATE-----
MIICrjCCAZYCCQDmdKc1MJXJ2TANBgkqhkiG9w0BAQsFADAZMRcwFQYDVQQDDA5j
YS5kZWZhdWx0LnN2YzAeFw0yMjA4MDYxNTI3MzVaFw00OTEyMjIxNTI3MzVaMBkx
FzAVBgNVBAMMDmNhLmRlZmF1bHQuc3ZjMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8A
MIIBCgKCAQEArMObHoGhub4gobZQil4D8OYLvxqFZtfFgrUCuR/DTA5a9fNK9ysE
ZFPiltPedmFBqQdX+STg5x5ONpiYZ6mVzAWLE8/i2Lbl6ID4mkBETWePKq4xUfiw
J6XYlmbWHbKXsAs4pLTq4uNd1y0C6azNo1wbIopqNmCjjfoophDRNUA0zyX2X8SN
Aw+n9rJYSRsHR6fvIBjBP2QoKsdF7AToXatFeZmolleMdAH0ja6PQd4G6hDiMhFY
D+XzPGdIUKdOvsNNAPAqnn2u96jA0+676vryxqOoK9WJ1XwkDJGXe3hPlNIkNFmd
vqRJvrVtxGw+oOZ+K6PMZ3oSkv/EmCQvfwIDAQABMA0GCSqGSIb3DQEBCwUAA4IB
AQBniu4YfPP32Nl/C3Mg/lD94hqFCK6Cyi3db38MlSE8zdUmZ4AZUYoK0uqUmsRF
ztVFaciaw69H/C3FG2AWAiYZMsZByEGImgTI8LsgXz8ixPkzeb8EfB1mN/uRlaDb
aENAfFWjOrnfZgfbmlvizlxCjIGbpLOcU6uiWGA2OmXtHrniPIYdaIvhMNQvIbYb
o4Ny16oiezX+Mdt6bLfix1f0/fZ9QFW37p9gDGbNwdkXxEK1TOO8HvO6CfUZul5J
6b4F95M7lwkeDl44Z4CVmfXNxyN7a2girDxbYfIHnXzYRC0YRrIlFlhTJ8T/lPLw
FuINiqLppxlwHJbmGkUHM2qS
-----END CERTIFICATE-----


`
	serverPEM = `
-----BEGIN CERTIFICATE-----
MIIDYzCCAkugAwIBAgIJAL/0iqyKjqnMMA0GCSqGSIb3DQEBBQUAMBkxFzAVBgNV
BAMMDmNhLmRlZmF1bHQuc3ZjMB4XDTIyMDgwNjE1MjczOFoXDTQ5MTIyMjE1Mjcz
OFowczELMAkGA1UEBhMCY24xEDAOBgNVBAgMB3NoYWFueGkxDTALBgNVBAcMBHhp
YW4xEDAOBgNVBAoMB2RlZmF1bHQxETAPBgNVBAsMCHdlYnNldmVyMR4wHAYDVQQD
DBV3ZWJzZXJ2ZXIuZGVmYXVsdC5zdmMwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAw
ggEKAoIBAQDeECMoROU9TUC9+Y1ZLfhWXTIOxJlSZq3bjcFPR5KlM6hd5FnaYQDL
WqEYDTSiS2ifFUqMGnFfdz69JI/Q0Gmqp7ly/WR/bHDoG02JBMO8VsLdvHouLctj
WjtNFVo2hWZeCtRJAc+8+HQH/zELa5T6Af0PFUVDHte4T+NYcHfhbp7FsEBFkip1
SVsobur+FRjRAmoziW4pchLRGjZUVdS8eGg3QoR2u9YjaWdT2L/LURia6dO+l11A
RjfZcv6uqYmJuav+5jw0Tnf8Ekmcq2xhHkyoookZbfk5tVVYltE8dXDjdtwTzMl5
mgcTFe+onXJeat3azeYiHzw2+w+10jw3AgMBAAGjVDBSMAkGA1UdEwQCMAAwCwYD
VR0PBAQDAgQwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAZBgNVHREE
EjAQhwTAqDIBggh3ZWJzZXZlcjANBgkqhkiG9w0BAQUFAAOCAQEAbcJcOyNMcvC6
zsEganZdSc9XfVfaq6gLGBAA1kPmtqQ/Ee0mg7hLtKF8VIYON18PapUXGxA0n+FL
2alXLE1zysS+SxQqZLPnZ8Sciuc2rb2guq25KvKqsqL9hGitShX7QfjzMDJE8PtS
UefhRODuILp1RJ3t//UOWKlVu1yTMdXvIts7hIEmwglAJ2Z/Zm1SbTaF7cbNfc6s
Maag5wgMamrhDLkGd4h8/ozwyUsxA7LNsxAEYPH+3RtbphHGeBVcy74U5hDohVyR
c0ailttQ3OqXLWoeRH1cmuYbKAzgo1hDrF5lyiljAA4k7eWRPz1Zm0m2SM7/uFwP
veVoSxQlZA==
-----END CERTIFICATE-----

`
)

func main() {
	//cryptoDemo()
	//httpsDemo()
	clientGoDemo()
}

func clientGoDemo() {
	pem, _ := cert.ParseCertsPEM([]byte("LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNyakNDQVpZQ0NRRG1kS2MxTUpYSjJUQU5CZ2txaGtpRzl3MEJBUXNGQURBWk1SY3dGUVlEVlFRRERBNWoKWVM1a1pXWmhkV3gwTG5OMll6QWVGdzB5TWpBNE1EWXhOVEkzTXpWYUZ3MDBPVEV5TWpJeE5USTNNelZhTUJreApGekFWQmdOVkJBTU1EbU5oTG1SbFptRjFiSFF1YzNaak1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBCk1JSUJDZ0tDQVFFQXJNT2JIb0dodWI0Z29iWlFpbDREOE9ZTHZ4cUZadGZGZ3JVQ3VSL0RUQTVhOWZOSzl5c0UKWkZQaWx0UGVkbUZCcVFkWCtTVGc1eDVPTnBpWVo2bVZ6QVdMRTgvaTJMYmw2SUQ0bWtCRVRXZVBLcTR4VWZpdwpKNlhZbG1iV0hiS1hzQXM0cExUcTR1TmQxeTBDNmF6Tm8xd2JJb3BxTm1Dampmb29waERSTlVBMHp5WDJYOFNOCkF3K245ckpZU1JzSFI2ZnZJQmpCUDJRb0tzZEY3QVRvWGF0RmVabW9sbGVNZEFIMGphNlBRZDRHNmhEaU1oRlkKRCtYelBHZElVS2RPdnNOTkFQQXFubjJ1OTZqQTArNjc2dnJ5eHFPb0s5V0oxWHdrREpHWGUzaFBsTklrTkZtZAp2cVJKdnJWdHhHdytvT1orSzZQTVozb1Nrdi9FbUNRdmZ3SURBUUFCTUEwR0NTcUdTSWIzRFFFQkN3VUFBNElCCkFRQm5pdTRZZlBQMzJObC9DM01nL2xEOTRocUZDSzZDeWkzZGIzOE1sU0U4emRVbVo0QVpVWW9LMHVxVW1zUkYKenRWRmFjaWF3NjlIL0MzRkcyQVdBaVlaTXNaQnlFR0ltZ1RJOExzZ1h6OGl4UGt6ZWI4RWZCMW1OL3VSbGFEYgphRU5BZkZXak9ybmZaZ2ZibWx2aXpseENqSUdicExPY1U2dWlXR0EyT21YdEhybmlQSVlkYUl2aE1OUXZJYlliCm80TnkxNm9pZXpYK01kdDZiTGZpeDFmMC9mWjlRRlczN3A5Z0RHYk53ZGtYeEVLMVRPTzhIdk82Q2ZVWnVsNUoKNmI0Rjk1TTdsd2tlRGw0NFo0Q1ZtZlhOeHlON2EyZ2lyRHhiWWZJSG5YellSQzBZUnJJbEZsaFRKOFQvbFBMdwpGdUlOaXFMcHB4bHdISmJtR2tVSE0ycVMKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="))

	c, err := rest.RESTClientFor(&rest.Config{
		Host: "https://192.168.50.1:9443/",
		ContentConfig: rest.ContentConfig{
			GroupVersion:         &schema.GroupVersion{Group: "", Version: "v1"},
			NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
		},
		TLSClientConfig: rest.TLSClientConfig{
			Insecure: false,
			//ServerName: "xxx",
			CAData: pem,
		},
	})
	if err != nil {
		panic(err)
	}
	request := c.Get()
	do := request.Do(context.Background())
	get, err := do.Get()
	if err != nil {
		panic(err)
	}
	fmt.Println(get)
}

func httpsDemo() {
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(serverPEM))
	tr := &http.Transport{
		////把从服务器传过来的非叶子证书，添加到中间证书的池中，使用设置的根证书和中间证书对叶子证书进行验证。

		TLSClientConfig: &tls.Config{RootCAs: roots},

		//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //InsecureSkipVerify用来控制客户端是否证书和服务器主机名。如果设置为true,//       //则不会校验证书以及证书中的主机名和服务器主机名是否一致。

	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("https://192.168.50.1:9443")
	if err != nil {
		fmt.Println("Get error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func cryptoDemo() {

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(caPEM))
	if !ok {
		panic("failed to parse root certificate")
	}
	_, err := tls.Dial("tcp", "192.168.50.1:9443", &tls.Config{
		RootCAs: roots,
	})

	if err != nil {
		panic("Server doesn't support SSL certificate err: " + err.Error())
	}
	fmt.Println("suc")
}
