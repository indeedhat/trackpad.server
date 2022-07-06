# Trackpad server
Server componnet for indeedhat/trackpad.client

allow your android device to be used as a wireless keyboadr and trackpad

## TODO
- [ ] add make/mage file for easy deployment
- [ ] make the multicast port configurable (app will also need a setting for this)

## Config
all done via a .env file placed in the server directory
| envar | default | description |
| --- | --- | --- |
| SERVER_PORT | 8181 | Local port to listen for connections on |
| CONNECT_PASS | N/A | When set clients will be asked for a password on connect |
| DISCOVERY_INTERVAL | 5 | seconds interval at which the server will send out a udp multicast message for server discovery in app |

## Credits
(github.com/micmonay/keybd_event)[github.com/micmonay/keybd_event]
(github.com/joho/godotenv)[github.com/joho/godotenv]
(github.com/go-vgo/robotgo)[github.com/go-vgo/robotgo]
(github.com/gorilla/websocket)[github.com/gorilla/websocket]

