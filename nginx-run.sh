#openssl s_client -connect 127.0.0.1:6443 \
# -cert certs/cluster-2.dev.power.sd.istio.space/server.crt \
#-key certs/cluster-2.dev.power.sd.istio.space/server.key \
#-CAfile certs/cluster-2.dev.power.sd.istio.space/ca.crt
#

#openssl s_client -connect 192.168.182.124:6443 -CAfile certs/cluster-2.dev.power.sd.istio.space/ca.crt
#
#openssl x509 -in certs/cluster-2.dev.power.sd.istio.space/server.crt -noout -text
#openssl x509 -in certs/cluster-2.dev.power.sd.istio.space/server.crt -text -noout | grep -A 1 "X509v3 Extended Key Usage"
#openssl x509 -in certs/cluster-2.dev.power.sd.istio.space/server.crt -noout -text | grep -A1 "Subject Alternative Name"
#openssl x509 -noout -modulus -in certs/cluster-2.dev.power.sd.istio.space/server.crt | openssl md5
#openssl rsa -noout -modulus -in certs/cluster-2.dev.power.sd.istio.space/server.key | openssl md5
#
#openssl x509 -in certs/apiserver.crt -text -noout | grep -A 1 "X509v3 Extended Key Usage"
#openssl x509 -in certs/apiserver.crt -noout -text | grep -A1 "Subject Alternative Name"

#
docker run -it  --rm -p 9443:9443 \
   -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf \
   -v $(pwd)/certs:/etc/nginx/certs \
   -v $(pwd)/logs:/var/log/nginx \
   nginx