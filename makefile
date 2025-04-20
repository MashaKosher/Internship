# Makefile

AUTH_SERVICE = ./authservice
ADMIN_SERVICE = ./adminservice


start-app:
	air $(AUTH_SERVICE)
	air $(ADMIN_SERVICE)