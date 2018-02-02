#!/bin/sh

#icinga2 API cert - regenerate new private key and certificate when running in a new container
if [ ! -f "/etc/icinga2/pki/$(hostname).key" ]; then
	icinga2 node setup --master
fi

exec icinga2 daemon