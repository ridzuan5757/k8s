PKI = openssl
DAYS = 36525
C = MY
ST = SELANGOR
L = BANGI
O = SILENTMODE
OU = ENGINEERING
CERTS = admin node1

all: root-ca.pem $(addsuffix .pem, $(CERTS))

# ---------------------------
# Root CA Generation
# ---------------------------

root-ca-key.pem:
	$(PKI) genrsa -out root-ca-key.pem 2048

root-ca.pem: root-ca-key.pem
	$(PKI) req -new -x509 -sha256 \
		-key root-ca-key.pem \
		-subj "/C=$(C)/ST=$(ST)/L=$(L)/O=$(O)/OU=$(OU)/CN=root.dns.a-record" \
		-days $(DAYS) \
		-out root-ca.pem

# ---------------------------
# Admin Certificate Generation
# ---------------------------

admin-key-temp.pem:
	$(PKI) genrsa -out admin-key-temp.pem 2048

admin-key.pem: admin-key-temp.pem:
	$(PKI) pkcs8 -topk8 -nocrypt -v1 PBE-SHA1-3DES \
		-inform PEM \
		-outform PEM \
		-in admin-key-temp.pem \
 		-out admin-key.pem \
	
admin.csr: admin-key.pem
	$(PKI) req -new \
		-key admin-key.pem \
		-subj "/C=CA/ST=ONTARIO/L=TORONTO/O=ORG/OU=UNIT/CN=A" \
		-out admin.csr

