<!-- private key -->
openssl genrsa -out jwt-private.pem 2048




<!-- public keu -->
openssl rsa -in jwt-private.pem -pubout -out jwt-public.pem



